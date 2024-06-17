package bot

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	coreConfigPack "github.com/https-whoyan/MafiaBot/core/config"
	coreGamePack "github.com/https-whoyan/MafiaBot/core/game"
	coreMessagePack "github.com/https-whoyan/MafiaBot/core/message"
	corePlayerPack "github.com/https-whoyan/MafiaBot/core/player"
	coreRolesPack "github.com/https-whoyan/MafiaBot/core/roles"
	botCnvPack "github.com/https-whoyan/MafiaBot/internal/converter"
	botFMTPack "github.com/https-whoyan/MafiaBot/internal/fmt"
	botGameCfgPack "github.com/https-whoyan/MafiaBot/internal/game"
	botMsgPack "github.com/https-whoyan/MafiaBot/internal/message"
	botTimeConstsPack "github.com/https-whoyan/MafiaBot/internal/time"
	botUserPack "github.com/https-whoyan/MafiaBot/internal/user"

	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"
	"github.com/https-whoyan/MafiaBot/pkg/repository/redis"

	"github.com/bwmarrin/discordgo"
)

// _________________________________
// ********************************
// This contains all bot commands.
//********************************
// _________________________________

const (
	AddChannelRoleCommandName   = "add_channel_role"
	AddMainChannelCommandName   = "add_main_channel"
	RegisterGameCommandName     = "register_game"
	ChoiceGameConfigCommandName = "choose_game_config"
	YanLohCommandName           = "yan_loh"
	AboutRolesCommandName       = "about_roles"
	StartGameCommandName        = "start_game"
)

// _______________________
// Channels
// _______________________

// AddChannelRoleCommand command logic
type AddChannelRoleCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewAddChannelRoleCommand() *AddChannelRoleCommand {
	generateOption := func(roleName string) *discordgo.ApplicationCommandOption {
		return &discordgo.ApplicationCommandOption{
			Name:        roleName,
			Description: fmt.Sprintf("Add %s interationChat", roleName),
			Type:        discordgo.ApplicationCommandOptionString,
		}
	}

	generateOptions := func() []*discordgo.ApplicationCommandOption {
		allNamesOfRoles := coreRolesPack.GetInteractionRoleNamesWhoHasOwnChat()
		var options []*discordgo.ApplicationCommandOption
		for _, roleName := range allNamesOfRoles {
			options = append(options, generateOption(strings.ToLower(roleName)))
		}
		return options
	}

	return &AddChannelRoleCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        AddChannelRoleCommandName,
			Description: "Define a chat room where the interaction between the bot and the role will take place.",
			Options:     generateOptions(),
		},
		isUsedForGame: false,
		name:          AddChannelRoleCommandName,
	}
}

func (c AddChannelRoleCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c AddChannelRoleCommand) GetName() string                       { return c.name }
func (c AddChannelRoleCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

func (c AddChannelRoleCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	if len(i.ApplicationCommandData().Options) == 0 {
		content := "Must be option!"
		Response(s, i, content)
		return
	}

	// Get variables by options
	roleName := i.ApplicationCommandData().Options[0].Name
	requestedChatID := i.ApplicationCommandData().Options[0].Value.(string)

	currDB, containsDB := mongo.GetCurrMongoDB()
	if !containsDB {
		content := "Internal Server Error!"
		Response(s, i, content)
		return
	}

	isFreeChatID, _ := currDB.IsFreeChannelIID(i.GuildID, requestedChatID)
	if !isFreeChatID {
		content := f.B("Provided channel is already used.") + f.NL() + "Please, provide another, available chatID."
		Response(s, i, content)
		return
	}

	if !isCorrectChatID(s, requestedChatID) {
		content := "Invalid Chat ID!"
		Response(s, i, content)
		return
	}

	err := currDB.SetRoleChannel(i.GuildID, requestedChatID, roleName)
	if err != nil {
		content := "Internal Server Error!"
		Response(s, i, content)
		return
	}

	noticeChatContent := fmt.Sprintf("Now chat is used for %v role.", roleName)
	_, err = sendMessages(s, requestedChatID, noticeChatContent)
	if err != nil {
		content := "Bot can't write to given chat. " + f.NL() + f.Bold("Reuse command.")
		Response(s, i, content)
		return
	}
	messageContent := fmt.Sprintf("Done, now is %v chat is used for %v role.", requestedChatID, roleName) +
		f.NL() + f.BU("Thanks!")
	Response(s, i, messageContent)
}

// AddMainChannelCommand command logic
type AddMainChannelCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewAddMainChannelCommand() *AddMainChannelCommand {
	return &AddMainChannelCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        AddMainChannelCommandName,
			Description: "Define a chat room where the interaction between the bot and all game participants.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "chat_id",
					Description: "Add main game chat.",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		isUsedForGame: false,
		name:          AddMainChannelCommandName,
	}
}

func (c AddMainChannelCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c AddMainChannelCommand) GetName() string                       { return c.name }
func (c AddMainChannelCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

func (c AddMainChannelCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	if len(i.ApplicationCommandData().Options) == 0 {
		content := "Must be option!"
		Response(s, i, content)
		return
	}

	// Get variables by options
	newChatID := i.ApplicationCommandData().Options[0].Value.(string)

	currDB, containsDB := mongo.GetCurrMongoDB()
	if !containsDB {
		content := "Internal Server Error! Main Chat not add"
		Response(s, i, content)
		return
	}

	isFreeChatID, _ := currDB.IsFreeChannelIID(i.GuildID, newChatID)
	if !isFreeChatID {
		content := f.B("Provided channel is already used.") + f.NL() + "Please, provide another, available chatID."
		Response(s, i, content)
		return
	}

	if !isCorrectChatID(s, newChatID) {
		content := "Invalid Chat ID!"
		Response(s, i, content)
		return
	}

	err := currDB.SetMainChannel(i.GuildID, newChatID)
	if err != nil {
		content := "Internal Server Error! Main Chat not add"
		Response(s, i, content)
		return
	}

	noticeChatContent := fmt.Sprintf("Now chat will be used as main for games.")
	_, err = sendMessages(s, newChatID, noticeChatContent)
	if err != nil {
		content := "Bot can't write to given chat. " + f.NL() + f.B("Reuse command.")
		Response(s, i, content)
		return
	}
	messageContent := "Done, you configure main interaction chat."
	Response(s, i, messageContent)
}

// _____________________
// Game
// _____________________

// RegisterGameCommand command logic
type RegisterGameCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewRegisterGameCommand() *RegisterGameCommand {
	return &RegisterGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        RegisterGameCommandName,
			Description: "Register new Game",
		},
		isUsedForGame: true,
		name:          RegisterGameCommandName,
	}
}

func (c RegisterGameCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c RegisterGameCommand) GetName() string                       { return c.name }
func (c RegisterGameCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

func (c RegisterGameCommand) trySetChannels(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) (correct bool) {
	emptyRoles, err := setRolesChannels(s, i.GuildID, g)
	if err != nil {
		content := "Internal Server Error"
		Response(s, i, content)
		return
	}

	var content string
	if len(emptyRoles) != 0 {
		content += "You don't configure all channels for bot interaction. " +
			"Please, use " + f.BU(AddChannelRoleCommandName) + " to fix " + strings.Join(emptyRoles, ", ") +
			" roles."
	}

	if !existsMainChannel(i.GuildID) {
		if len(content) != 0 {
			content += f.NL() + f.NL()
		}

		content += f.B("You don't configure main channel for bot interaction.") + f.NL() +
			f.I("All messages regarding the game will be sent there.") + f.NL() +
			"To add a main chat channel for interaction, use the command " +
			f.BU(AddMainChannelCommandName)
	}
	setMainChannel(s, i.GuildID, g)
	if len(content) == 0 {
		return true
	}
	Response(s, i, content)
	return false
}

func (c RegisterGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	// Validation
	if !c.trySetChannels(s, i, g, f) {

		return
	}

	// Send message.
	Response(s, i, "Ok. Message below.")

	// Send additional message and save it ID
	deadlineStr := strconv.Itoa(botTimeConstsPack.RegistrationDeadlineMinutes)
	responseMessageText := "Registration has begun." + f.NL() +
		f.B("Post "+botFMTPack.RegistrationPlayerSticker+" reaction below.") + f.I(" If you want to be a spectator, "+
		"put the reaction "+botFMTPack.RegistrationSpectatorSticker+".") + f.NL() + f.NL() + f.B(
		f.U("Deadline: "+deadlineStr+" minutes"))

	channelID := i.ChannelID
	message, err := s.ChannelMessageSend(channelID, responseMessageText)
	if err != nil {
		log.Print(err)
	}
	g.SwitchState()
	messageID := message.ID
	// Safe messageID to redis
	currDB, isContains := redis.GetCurrRedisDB()
	if !isContains {
		log.Print("DB isn't initialed, couldn't get currDB")
		return
	}
	err = currDB.SetInitialGameMessageID(
		i.GuildID, messageID,
		botTimeConstsPack.RegistrationDeadlineSeconds*time.Second)
	if err != nil {
		log.Print(err)
	}
}

// ChoiceGameConfigCommand command logic
type ChoiceGameConfigCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewChoiceGameConfigCommand() *ChoiceGameConfigCommand {
	return &ChoiceGameConfigCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        ChoiceGameConfigCommandName,
			Description: "This output a list of game configs for voting.",
		},
		isUsedForGame: true,
		name:          ChoiceGameConfigCommandName,
	}
}

func (c ChoiceGameConfigCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c ChoiceGameConfigCommand) GetName() string                       { return c.name }
func (c ChoiceGameConfigCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

func (c ChoiceGameConfigCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	currRedisDB, isContains := redis.GetCurrRedisDB()
	if !isContains {
		log.Println("redis is not exists, command: startGameCommand")
		content := "Internal Server Error!"
		Response(s, i, content)
	}

	registrationMessageID, err := currRedisDB.GetInitialGameMessageID(i.GuildID)
	if (err != nil || registrationMessageID == "") && g.State == coreGamePack.NonDefinedState {
		messageContent := f.U("Registration Deadline passed!") + f.NL() + "Please, " +
			f.B("use the "+RegisterGameCommandName+" command") + " to register a new game."
		g.SetState(coreGamePack.NonDefinedState)
		Response(s, i, messageContent)
		return
	}

	// If playersCount not in range [minAvailableCount, maxAvailableCount],
	// Send message that it's impossible to choice config.
	registrationStickerUnicode := botFMTPack.RegistrationPlayerSticker
	_, playersCount := botMsgPack.GetUsersByEmojiID(s, i.ChannelID, registrationMessageID, registrationStickerUnicode)

	allConfigs, nearest, err := coreConfigPack.GetConfigsByPlayersCount(playersCount)
	switch {
	case errors.Is(err, coreConfigPack.SmallCountOfPeopleToConfig):
		content := f.B("The number of players is too small to start the game.") + f.NL() +
			"Number of registered players: " + f.Bl(strconv.Itoa(playersCount)) +
			f.NL() + "Minimum number to vote on game config choices: " + f.Bl(strconv.Itoa(nearest))
		Response(s, i, content)
		return
	case errors.Is(err, coreConfigPack.BigCountOfPeopleToConfig):
		content := f.B("The number of players is too large to start the game.") + f.NL() +
			"Number of registered players: " + f.Bl(strconv.Itoa(playersCount)) + f.NL() +
			"Maximum number to vote on game config choices: " + f.Bl(strconv.Itoa(nearest))
		Response(s, i, content)
		return
	}

	// If playersCount is ok,
	// set empty players to game (to safe it.)
	startBotPlayers, _ := botMsgPack.GetUsersByEmojiID(
		s, i.ChannelID, registrationMessageID, botFMTPack.RegistrationPlayerSticker)
	startGamePlayers := botCnvPack.DiscordUsersToEmptyPlayers(s, i.GuildID,
		startBotPlayers, false)
	g.SetStartPlayers(startGamePlayers)

	// And spectators.
	registerSpectatorStickerUnicode := botFMTPack.RegistrationSpectatorSticker
	botSpectators, _ := botMsgPack.GetUsersByEmojiID(s, i.ChannelID,
		registrationMessageID, registerSpectatorStickerUnicode)
	// We need no duplicates in active and spectators.
	// Then, I get unique Spectators, which not include in active players.
	uniqueBotSpectators, _ := botUserPack.GetUsersNotInclude(botSpectators, startBotPlayers)
	spectators := botCnvPack.DiscordUsersToEmptyPlayers(s, i.GuildID, uniqueBotSpectators, true)
	g.SetSpectators(spectators)

	content := "Below is a list of available game configurations. " + f.NL() +
		f.B("If you like the configuration") + ",  please " +
		f.BU("put a reaction ") + botFMTPack.ConfigChoiceSticker + f.BU(" on the post.") + f.NL() +
		"You can give a reaction to any number of configurations. " + f.NL() + f.NL()
	content += f.I("The bot will choose the configuration with the most reactions. "+
		"If there are several configurations, a random one will be selected.") + f.NL() + f.NL()
	deadlineStr := strconv.Itoa(botTimeConstsPack.VotingGameConfigDeadlineMinutes)
	content += "Deadline: " + f.BU(deadlineStr+" minutes") + "."
	Response(s, i, content)

	db, isContains := redis.GetCurrRedisDB()
	if !isContains {
		log.Println("redis is not exists, command: startGameCommand")
		content = "Internal Server Error!"
		_, _ = sendMessages(s, i.ChannelID, content)
		return
	}

	// Create a ConfigMessages structure to safe all message to redis
	cfgMessages := botGameCfgPack.NewConfigMessages(i.GuildID, playersCount, len(allConfigs))
	for index, config := range allConfigs {
		_, _ = sendMessages(s, i.ChannelID, f.InfoSplitter())
		messageContent := config.GetMessageAboutConfig(f)
		mpMessages, _ := sendMessages(s, i.ChannelID, messageContent)
		cfgMessages.AddNewMessage(index, mpMessages[messageContent].ID)
	}
	err = db.SetConfigGameVotingMessages(cfgMessages, botTimeConstsPack.VotingGameConfigDeadlineSeconds*time.Second)
	if err != nil {
		log.Println(err)
	}

	return
}

type StartGameCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewStartGameCommand() *StartGameCommand {
	return &StartGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        StartGameCommandName,
			Description: "Init game after game config choosing",
		},
		isUsedForGame: true,
		name:          StartGameCommandName,
	}
}

func (c StartGameCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c StartGameCommand) GetName() string                       { return c.name }
func (c StartGameCommand) IsUsedForGame() bool                   { return c.isUsedForGame }
func (c StartGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	r, isContains := redis.GetCurrRedisDB()
	if !isContains {
		Response(s, i, "Internal Server Error!")
		return
	}

	cfgMessages, err := r.GetConfigGameVotingMessageID(i.GuildID)
	if err != nil {
		Response(s, i, "Internal Server Error!")
		return
	}

	g.SwitchState()

	// Set Reactions count to cfgMessages
	for _, message := range cfgMessages.Messages {
		messageIID := message.MessageID
		emoji := botFMTPack.ConfigChoiceSticker
		usersWhoLiked, _ := botMsgPack.GetUsersByEmojiID(s, i.ChannelID, messageIID, emoji)

		// Validate, we only want users who participate in the game.
		_, playersMessageCount := botUserPack.GetUsersOnlyIncludeInTags(
			usersWhoLiked,
			corePlayerPack.GetTagsByPlayers(g.StartPlayers))

		// Set It.
		message.SetReactionCount(playersMessageCount)
	}

	winner, playerCount, isRandom := cfgMessages.GetWinner()
	if isRandom {
		Response(s, i, f.B("Choosing Random configuration..."))
	} else {
		Response(s, i, f.B("Setting the game's configuration..."))
	}
	// Simulation :))))
	time.Sleep(time.Millisecond * time.Duration(1400))
	winnerConfig := coreConfigPack.GetConfigByPlayersCountAndIndex(playerCount, winner)
	//TODO!!!!
	_, _ = sendMessages(s, i.ChannelID, f.InfoSplitter())
	err = g.Init(winnerConfig)
	if err != nil {
		_, _ = sendMessages(s, i.ChannelID, f.IU("Can't start game, internal server error!"))
		log.Println(err)
		g.SetState(coreGamePack.RegisterState)
		return
	}
	for _, player := range g.StartPlayers {
		err = SendToUser(s, player.Tag, coreMessagePack.GetStartPlayerDefinition(player, f))
		if err != nil {
			log.Println(err)
		}
	}
	log.Printf("Init Game in %v Guild", i.GuildID)
	return
}

// ______________
// Voting
// ______________

//TODO!

// ___________
// Other
// ___________

// YanLohCommand command
type YanLohCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewYanLohCommand() *YanLohCommand {
	return &YanLohCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        YanLohCommandName,
			Description: "Call Yan with this command!",
		},
		isUsedForGame: false,
		name:          YanLohCommandName,
	}
}

func (c YanLohCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c YanLohCommand) GetName() string                       { return c.name }
func (c YanLohCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

func (c YanLohCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game, _ *botFMTPack.DiscordFMTer) {
	messageContent := "Не, лол, что бы его же ботом его обзывать..."
	Response(s, i, messageContent)

	// async kick requester
	guildID := i.GuildID
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(sessId, kickedUserID string) {
		defer wg.Done()
		var kickPing time.Duration = 3
		time.Sleep(time.Second * kickPing)

		err := s.GuildMemberMove(sessId, kickedUserID, nil)
		if err != nil {
			log.Printf("failed kick user, err: %v", err)
		}
	}(guildID, i.Member.User.ID)
	wg.Wait()
}

// AboutRolesCommand command logic
type AboutRolesCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewAboutRolesCommand() *AboutRolesCommand {
	return &AboutRolesCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        AboutRolesCommandName,
			Description: "Send description about roles",
		},
		isUsedForGame: false,
		name:          AboutRolesCommandName,
	}
}

func (c AboutRolesCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c AboutRolesCommand) GetName() string                       { return c.name }
func (c AboutRolesCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

func (c AboutRolesCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game,
	f *botFMTPack.DiscordFMTer) {
	messageContent := f.Bold("Below information about all roles:") + f.NL()
	Response(s, i, messageContent)

	sendMessage := func(s *discordgo.Session, i *discordgo.Interaction, message string) {
		_, err := s.ChannelMessageSend(i.ChannelID, message)
		if err != nil {
			log.Print(err)
		}
	}

	messages := coreRolesPack.GetDefinitionsOfAllRoles(f, 2000)
	for _, message := range messages {
		sendMessage(s, i, coreRolesPack.FixDescription(message))
	}
}
