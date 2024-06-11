package bot

import (
	"errors"
	"github.com/https-whoyan/MafiaBot/internal/bot/fmt"
	"log"
	"strings"

	"github.com/https-whoyan/MafiaBot/internal/bot/channel"
	"github.com/https-whoyan/MafiaBot/internal/core/config"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
	"github.com/https-whoyan/MafiaBot/pkg/db/mongo"

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
func noticeChat(s *discordgo.Session, chatID string, content ...string) (map[string]*discordgo.Message, error) {
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

func NoticePrivateChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := fmt.Bold("All commands are used on the server.\n") + "If you have difficulties in using the bot, " +
		"please refer to the repository documentation: https://github.com/https-whoyan/MafiaBot"
	Response(s, i.Interaction, content)
}

// NoticeIsEmptyGame If game not exists
func NoticeIsEmptyGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := "You can't interact with the game because you haven't registered it\n" +
		fmt.Bold("Write the "+fmt.Emphasized("/register_game")+" command") + " to start the game."
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
	// Try to lock session
	if s.TryLock() {
		defer s.Unlock()
	}
	// emptyRolesMp: save not contains channel roles
	emptyRolesMp := make(map[string]bool)
	// mappedRoles: save contains channels roles
	mappedRoles := make(map[string]*channel.BotRoleChannel)
	for _, roleName := range allRolesNames {
		if strings.ToLower(roleName) == "don" {
			continue
		}
		channelIID, err := currDB.GetChannelIIDByRole(guildID, roleName)
		if channelIID == "" {
			emptyRolesMp[roleName] = true
			continue
		}
		newRoleChannel, err := channel.LoadRoleChannel(s, channelIID, roleName)
		if err != nil {
			emptyRolesMp[roleName] = true
			continue
		}
		mappedRoles[roleName] = newRoleChannel
	}

	// Try lock game
	if g.TryLock() {
		defer g.Unlock()
	}

	// If a have all roles
	if len(emptyRolesMp) == 0 {
		// Save it to g.InteractionChannels
		g.InteractionChannels = mappedRoles
		// And return it
		return []string{}, nil
	}
	// Convert a map to slice
	var emptyRolesArr []string
	for emptyRole, _ := range emptyRolesMp {
		emptyRolesArr = append(emptyRolesArr, strings.ToLower(emptyRole))
	}
	// Return
	return emptyRolesArr, nil
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

func CreateConfigMessage(cfg *config.RolesConfig) string {
	return ""
}
