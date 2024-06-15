package bot

import (
	"errors"
	"github.com/https-whoyan/MafiaBot/core/game"
	"github.com/https-whoyan/MafiaBot/core/roles"
	"log"
	"strings"

	"github.com/https-whoyan/MafiaBot/internal/channel"
	"github.com/https-whoyan/MafiaBot/internal/converter"
	botFMT "github.com/https-whoyan/MafiaBot/internal/fmt"

	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"

	"github.com/bwmarrin/discordgo"
)

func isCorrectChatID(s *discordgo.Session, chatID string) bool {
	if s.TryLock() {
		defer s.Unlock()
	}

	ch, err := s.Channel(chatID)
	return err == nil && ch != nil
}

// Send message to chatID
func sendMessages(s *discordgo.Session, chatID string, content ...string) (map[string]*discordgo.Message, error) {
	// Represent Message by their content.
	messages := make(map[string]*discordgo.Message)
	for _, onceContent := range content {
		message, err := s.ChannelMessageSend(chatID, onceContent)
		if err != nil {
			return nil, err
		}
		messages[onceContent] = message
	}
	return messages, nil
}

// Response reply to interaction by provided content (s.InteractionResponse)
func Response(s *discordgo.Session, i *discordgo.Interaction, content string) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		log.Print(err)
	}
}

// IsPrivateMessage finds out if a message has been sent to private messages
func IsPrivateMessage(i *discordgo.InteractionCreate) bool {
	return i.GuildID == ""
}

func NoticePrivateChat(s *discordgo.Session, i *discordgo.InteractionCreate, fMTer *botFMT.DiscordFMTer) {
	content := fMTer.Bold("All commands are used on the server.") + fMTer.NL() +
		"If you have difficulties in using the bot, " +
		"please refer to the repository documentation: https://github.com/https-whoyan/MafiaBot"
	Response(s, i.Interaction, content)
}

// NoticeIsEmptyGame If game not exists
func NoticeIsEmptyGame(s *discordgo.Session, i *discordgo.InteractionCreate, fMTer *botFMT.DiscordFMTer) {
	content := "You can't interact with the game because you haven't registered it" + fMTer.NL() +
		fMTer.Bold("Write the "+fMTer.Underline(RegisterGameCommandName)+" command") + " to start the game."
	Response(s, i.Interaction, content)
}

// SetRolesChannels to game.
func setRolesChannels(s *discordgo.Session, guildID string, g *game.Game) ([]string, error) {
	// Get night interaction roles names
	allRolesNames := roles.GetAllNightInteractionRolesNames()
	// Get curr MongoDB struct
	currDB, isContains := mongo.GetCurrMongoDB()
	if !isContains {
		return []string{}, errors.New("MongoDB doesn't initialized")
	}
	// emptyRolesMp: save not contains channel roles
	emptyRolesMp := make(map[string]bool)
	// mappedRoles: save contains channels roles
	mappedRoles := make(map[string]*channel.BotRoleChannel)

	addNewChannelIID := func(roleName, channelName string) {
		channelIID, err := currDB.GetChannelIIDByRole(guildID, channelName)
		if channelIID == "" {
			emptyRolesMp[channelName] = true
			return
		}
		newRoleChannel, err := channel.NewBotRoleChannel(s, channelIID, roleName)
		if err != nil {
			emptyRolesMp[channelName] = true
			return
		}
		mappedRoles[roleName] = newRoleChannel
	}

	for _, roleName := range allRolesNames {
		if strings.ToLower(roleName) == strings.ToLower(roles.Don.Name) {
			addNewChannelIID(roleName, roles.Mafia.Name)
			continue
		}
		addNewChannelIID(roleName, roleName)
	}
	// If a have all roles
	if len(emptyRolesMp) == 0 {
		// Convert
		sliceMappedRoles := converter.GetMapValues(mappedRoles)
		InterfaceRoleChannelSlice := converter.ConvertRoleChannelsSliceToIChannelSlice(sliceMappedRoles)

		// Save it to g.RoleChannels.
		err := g.SetRoleChannels(InterfaceRoleChannelSlice)

		return []string{}, err
	}

	return converter.GetMapKeys(emptyRolesMp), nil
}

// Check, if main channel exists or not
func existsMainChannel(guildID string) bool {
	currMongo, exists := mongo.GetCurrMongoDB()
	if !exists {
		return false
	}
	channelIID, err := currMongo.GetMainChannelIID(guildID)
	if err != nil {
		return false
	}
	return channelIID != ""
}

func setMainChannel(s *discordgo.Session, guildID string, g *game.Game) {
	currMongo, _ := mongo.GetCurrMongoDB()
	channelIID, _ := currMongo.GetMainChannelIID(guildID)
	mainChannel, _ := channel.NewBotMainChannel(s, channelIID)
	_ = g.SetMainChannel(mainChannel)
}
