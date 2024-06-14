package bot

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	coreConfigPack "github.com/https-whoyan/MafiaBot/core/config"
	coreGamePack "github.com/https-whoyan/MafiaBot/core/game"
	coreRolesPack "github.com/https-whoyan/MafiaBot/core/roles"
	botTimePack "github.com/https-whoyan/MafiaBot/core/time"
	botCnvPack "github.com/https-whoyan/MafiaBot/internal/bot/converter"
	botFMTPack "github.com/https-whoyan/MafiaBot/internal/bot/fmt"
	botMsgPack "github.com/https-whoyan/MafiaBot/internal/bot/message"

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
		allNamesOfRoles := coreRolesPack.GetAllNightInteractionRolesNames()
		var options []*discordgo.ApplicationCommandOption
		for _, roleName := range allNamesOfRoles {
			// For done using mafia chat
			if roleName != "Don" {
				options = append(options, generateOption(strings.ToLower(roleName)))
			}
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

func (c AddChannelRoleCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c AddChannelRoleCommand) GetName() string {
	return c.name
}

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
	_, err = noticeChat(s, requestedChatID, noticeChatContent)
	if err != nil {
		content := "Bot can't write to given chat. " + f.NL() + f.Bold("Reuse command.")
		Response(s, i, content)
		return
	}
	messageContent := fmt.Sprintf("Done, now is %v chat is used for %v role.", requestedChatID, roleName) +
		f.NL() + f.BU("Thanks!")
	Response(s, i, messageContent)
}

func (c AddChannelRoleCommand) IsUsedForGame() bool {
	return c.isUsedForGame
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

func (c AddMainChannelCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c AddMainChannelCommand) GetName() string {
	return c.name
}

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
	_, err = noticeChat(s, newChatID, noticeChatContent)
	if err != nil {
		content := "Bot can't write to given chat. " + f.NL() + f.B("Reuse command.")
		Response(s, i, content)
		return
	}
	messageContent := "Done, you configure main interaction chat."
	Response(s, i, messageContent)
}

func (c AddMainChannelCommand) IsUsedForGame() bool {
	return c.isUsedForGame
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

func (c RegisterGameCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c RegisterGameCommand) GetName() string {
	return c.name
}

func (c RegisterGameCommand) trySetRolesChannel(s *discordgo.Session, i *discordgo.Interaction,
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

	if len(content) == 0 {
		return true
	}
	Response(s, i, content)
	return false
}

func (c RegisterGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	// Validation
	if !c.trySetRolesChannel(s, i, g, f) {

		return
	}

	// Send message.
	Response(s, i, "Ok. Message below.")

	// Send additional message and save it ID
	deadlineStr := strconv.Itoa(botTimePack.RegistrationDeadlineMinutes)
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
	err = currDB.SetInitialGameMessageID(i.GuildID, messageID)
	if err != nil {
		log.Print(err)
	}
}

func (c RegisterGameCommand) IsUsedForGame() bool {
	return c.isUsedForGame
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

func (c ChoiceGameConfigCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c ChoiceGameConfigCommand) GetName() string {
	return c.name
}

func (c ChoiceGameConfigCommand) IsUsedForGame() bool {
	return c.isUsedForGame
}

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
	registrationStickerUnicode := botFMTPack.GetUnicodeBySticker(botFMTPack.RegistrationPlayerSticker)
	_, playersCount := botMsgPack.GetUsersByEmojiID(s, i.ChannelID, registrationMessageID, registrationStickerUnicode)

	_, nearest, err := coreConfigPack.GetConfigsByPlayersCount(playersCount)
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
	startGamePlayers := botCnvPack.DiscordUsersToEmptyPlayers(startBotPlayers, false)
	g.SetStartPlayers(startGamePlayers)

	// And spectators.
	registerSpectatorStickerUnicode := botFMTPack.GetUnicodeBySticker(botFMTPack.RegistrationSpectatorSticker)
	botSpectators, _ := botMsgPack.GetUsersByEmojiID(s, i.ChannelID,
		registrationMessageID, registerSpectatorStickerUnicode)
	spectators := botCnvPack.DiscordUsersToEmptyPlayers(botSpectators, true)
	g.SetSpectators(spectators)

	//TODO!
	Response(s, i, "Empty for now...")
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

func (c YanLohCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c YanLohCommand) GetName() string {
	return c.name
}

func (c YanLohCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game, _ *botFMTPack.DiscordFMTer) {
	messageContent := "Возможно, что ян и лох. И древлян. Но что бы его же ботом его обзывать..."
	Response(s, i, messageContent)

	// async kick requester
	guildID := i.GuildID
	go func(sessId, kickedUserID string) {
		var kickPing time.Duration = 3
		time.Sleep(time.Second * kickPing)

		err := s.GuildMemberMove(sessId, kickedUserID, nil)
		if err != nil {
			log.Printf("failed kick user, err: %v", err)
		}
	}(guildID, i.Member.User.ID)
}

func (c YanLohCommand) IsUsedForGame() bool {
	return c.isUsedForGame
}

// AboutRolesCommand command logic
type AboutRolesCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func (c AboutRolesCommand) IsUsedForGame() bool {
	return c.isUsedForGame
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

func (c AboutRolesCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c AboutRolesCommand) GetName() string {
	return c.name
}

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
		sendMessage(s, i, message)
	}
}
