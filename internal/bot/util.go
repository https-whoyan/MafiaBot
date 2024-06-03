package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/bot/channel"
	"github.com/https-whoyan/MafiaBot/internal/core/config"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
	"github.com/https-whoyan/MafiaBot/pkg/db/mongo"
	"log"
)

// _____________
// Text Style
// _____________

func Bold(s string) string {
	return "**" + s + "**"
}

func Italic(s string) string {
	return "_" + s + "_"
}

func Emphasized(s string) string {
	return "__" + s + "__"
}

func CodeBlock(language, text string) string {
	return "```" + language + text + "```"
}

// _____________

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

// finds out if a message has been sent to private messages
func isPrivateMessage(i *discordgo.InteractionCreate) bool {
	return i.GuildID == ""
}

func noticePrivateChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := Bold("All commands are used on the server.\n") + "If you have difficulties in using the bot, " +
		"please refer to the repository documentation: https://github.com/https-whoyan/MafiaBot"
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		log.Println(errors.Join(
			errors.New("there was an error when sending a private message, err: "), err),
		)
	}
}

// If game not exists
func noticeIsEmptyGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You can't interact with the game because you haven't registered it\n" +
				Bold("Write the "+Emphasized("/register_game")+" command") + " to start the game.",
		},
	})
	if err != nil {
		log.Println(err)
	}
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
	mappedRoles := make(map[string]*channel.RoleChannel)
	for _, roleName := range allRolesNames {
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

	fmt.Println(emptyRolesMp)
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
		emptyRolesArr = append(emptyRolesArr, emptyRole)
	}
	// Return
	return emptyRolesArr, nil
}

func CreateConfigMessage(cfg *config.RolesConfig) string {
	return ""
}
