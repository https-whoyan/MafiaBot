package bot

import (
	"fmt"
	message2 "github.com/https-whoyan/MafiaBot/internal/bot/message"
	"github.com/https-whoyan/MafiaBot/internal/core/config"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"github.com/https-whoyan/MafiaBot/internal/core/players"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
	time2 "github.com/https-whoyan/MafiaBot/internal/time"
	"github.com/https-whoyan/MafiaBot/pkg/db/mongo"
	"github.com/https-whoyan/MafiaBot/pkg/db/redis"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// RegisterGameCommand command
type RegisterGameCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewRegisterGameCommand() *RegisterGameCommand {
	name := "register_game"
	return &RegisterGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Register new Game",
		},
		isUsedForGame: true,
		name:          name,
	}
}

func (c *RegisterGameCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *RegisterGameCommand) GetName() string {
	return c.name
}

func (c *RegisterGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction, g *game.Game) {
	// Load all RoleChannels
	emptyRoles, err := setRolesChannels(s, i.GuildID, g)
	if err != nil {
		content := "Internal Server Error"
		Response(s, i, content)
		return
	}
	if len(emptyRoles) != 0 {
		content := Bold("You don't configure all channels for bot interaction. ") +
			"Please, use " + Emphasized("/add_channel_role") + " to fix " + strings.Join(emptyRoles, ", ") +
			" roles."
		Response(s, i, content)
		return
	}
	// If ok, set game to NonDefinedState
	g.SetState(game.NonDefinedState)
	// Send message.
	Response(s, i, "Ok. Message below.")

	// Send additional message and save it ID
	deadlineStr := strconv.Itoa(time2.RegistrationDeadlineMinutes)
	responseMessageText := "Registration has begun. \n" +
		Bold("Post "+RegistrationPlayerSticker+" reaction below.") + Italic(" If you want to be a spectator, "+
		"put the reaction "+RegistrationSpectatorSticker+".") + "\n\n" + Bold(
		Emphasized("Deadline: "+deadlineStr+" minutes"))

	channelID := i.ChannelID
	message, err := s.ChannelMessageSend(channelID, responseMessageText)
	if err != nil {
		log.Print(err)
	}
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

func (c *RegisterGameCommand) IsUsedForGame() bool {
	return c.isUsedForGame
}

// AddChannelRoleCommand command logic
type AddChannelRoleCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewAddChannelRole() *AddChannelRoleCommand {
	name := "add_channel_role"

	generateOption := func(roleName string) *discordgo.ApplicationCommandOption {
		return &discordgo.ApplicationCommandOption{
			Name:        roleName,
			Description: "Add " + roleName + " interaction chat",
			Type:        discordgo.ApplicationCommandOptionString,
		}
	}

	generateOptions := func() []*discordgo.ApplicationCommandOption {
		allNamesOfRoles := roles.GetAllNightInteractionRolesNames()
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
			Name:        name,
			Description: "Define a chat room where the interaction between the bot and the role will take place.",
			Options:     generateOptions(),
		},
		isUsedForGame: false,
		name:          name,
	}
}

func (c *AddChannelRoleCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *AddChannelRoleCommand) GetName() string {
	return c.name
}

func (c *AddChannelRoleCommand) Execute(s *discordgo.Session, i *discordgo.Interaction, _ *game.Game) {
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
	currDB.Lock()
	currDB.Unlock()
	err := currDB.SetRoleChannel(i.GuildID, requestedChatID, roleName)
	if err != nil {
		content := Bold("Provided channel is already used.\n") + "Please, provide another, available chatID."
		Response(s, i, content)
		return
	}

	noticeChatContent := fmt.Sprintf("Now chat is used for %v role.", roleName)
	if _, ok := noticeChat(s, requestedChatID, noticeChatContent); ok != nil {
		content := "Invalid Chat ID!"
		Response(s, i, content)
		return
	}

	messageContent := fmt.Sprintf("Done, now is %v chat is %v chat.", requestedChatID, roleName)
	Response(s, i, messageContent)
}

func (c *AddChannelRoleCommand) IsUsedForGame() bool {
	return c.isUsedForGame
}

// ChoiceGameConfig command logic
type ChoiceGameConfig struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewChoiceGameConfig() *ChoiceGameConfig {
	name := "choice_config"
	return &ChoiceGameConfig{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "This output a list of game configs for voting.",
		},
		isUsedForGame: true,
		name:          name,
	}
}

func (c *ChoiceGameConfig) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *ChoiceGameConfig) GetName() string {
	return c.name
}

func (c *ChoiceGameConfig) IsUsedForGame() bool {
	return c.isUsedForGame
}

func (c *ChoiceGameConfig) Execute(s *discordgo.Session, i *discordgo.Interaction, g *game.Game) {
	currRedisDB, isContains := redis.GetCurrRedisDB()
	if !isContains {
		log.Println("redis is not exists, command: startGameCommand")
		content := "Internal Server Error!"
		Response(s, i, content)
	}

	registrationMessageID, err := currRedisDB.GetInitialGameMessageID(i.GuildID)
	if (err != nil || registrationMessageID == "") && g.State == game.NonDefinedState {
		messageContent := Emphasized("Registration Deadline passed!") + "\n" + "Please, " +
			Bold("use the /register_game command") + " to register a new game."
		Response(s, i, messageContent)
		return
	}

	// Set empty players to game (to save it)
	_, playersCount := message2.GetUsersByEmojiID(
		s, i.ChannelID, registrationMessageID, RegistrationPlayerSticker)

	// If playersCount not in range [minAvailableCount, maxAvailableCount],
	// Send message that it's impossible to choice config.
	minPossiblePlayers := config.GetMinPlayersCount()
	maxPossiblePlayers := config.GetMaxPlayersCount()
	if playersCount < minPossiblePlayers {
		content := Bold("The number of players is too small to start the game.\n") +
			"Number of registered players: " + CodeBlock("", strconv.Itoa(playersCount)) +
			"\nMinimum number to vote on game config choices: " + CodeBlock("", strconv.Itoa(minPossiblePlayers))
		Response(s, i, content)
	} else if playersCount > maxPossiblePlayers {
		content := Bold("The number of players is too large to start the game.\n") +
			"Number of registered players: " + CodeBlock("", strconv.Itoa(playersCount)) +
			"\nMaximum number to vote on game config choices: " + CodeBlock("", strconv.Itoa(minPossiblePlayers))
		Response(s, i, content)
	}

	// If playersCount is ok, set empty players to game (to safe it.)
	startPlayers, _ := message2.GetUsersByEmojiID(
		s, i.ChannelID, registrationMessageID, RegistrationPlayerSticker)
	g.StartPlayers = players.GenerateEmptyPlayers(startPlayers)
	// And spectators.
	spectators, _ := message2.GetUsersByEmojiID(s, i.ChannelID, registrationMessageID, RegistrationSpectatorSticker)
	g.Spectators = players.GenerateEmptyPlayers(spectators)
	// Send a message to config.
	Response(s, i, "Empty for now...")
}

// AboutRolesCommand command logic
type AboutRolesCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func (c *AboutRolesCommand) IsUsedForGame() bool {
	return c.isUsedForGame
}

func NewAboutRolesCommand() *AboutRolesCommand {
	name := "about_roles"
	return &AboutRolesCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Send description about roles",
		},
		isUsedForGame: false,
		name:          name,
	}
}

func (c *AboutRolesCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *AboutRolesCommand) GetName() string {
	return c.name
}

func (c *AboutRolesCommand) Execute(s *discordgo.Session, i *discordgo.Interaction, _ *game.Game) {
	messageContent := Bold("Below information about all roles:\n")
	Response(s, i, messageContent)
	messageContent = ""

	sendMessage := func(s *discordgo.Session, i *discordgo.Interaction, message string) {
		_, err := s.ChannelMessageSend(i.ChannelID, messageContent)
		if err != nil {
			log.Print(err)
		}
	}

	allSortedRoles := roles.GetAllRolesNames()
	for _, roleName := range allSortedRoles {
		messageContent += "================================\n"
		messageContent += roles.GetDefinitionOfRole(roleName)

		// To erase 2000 max length error
		if len(messageContent) >= 1500 {
			sendMessage(s, i, messageContent)
			messageContent = ""
		}
	}
	sendMessage(s, i, messageContent)
}

// YanLohCommand command
type YanLohCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewYanLohCommand() *YanLohCommand {
	name := "yan_loh"
	return &YanLohCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Call Yan with this command!",
		},
		isUsedForGame: false,
		name:          name,
	}
}

func (c *YanLohCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *YanLohCommand) GetName() string {
	return c.name
}

func (c *YanLohCommand) Execute(s *discordgo.Session, i *discordgo.Interaction, _ *game.Game) {
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

func (c *YanLohCommand) IsUsedForGame() bool {
	return c.isUsedForGame
}
