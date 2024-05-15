package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"log"
	"time"
)

// YanLohCommand command
type YanLohCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewYanLohCommand() *YanLohCommand {
	name := "yan_loh"
	return &YanLohCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Call Yan with this command!",
		},
		name: name,
	}
}

func (c *YanLohCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *YanLohCommand) GetName() string {
	return c.name
}

func (c *YanLohCommand) Execute(s *discordgo.Session, i *discordgo.Interaction) {
	messageContent := "Возможно, что ян и лох. И древлян. Но что бы его же ботом его обзывать..."
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: messageContent,
		},
	})
	if err != nil {
		log.Print(err)
	}
	// async kick requester
	guildID := i.GuildID
	go func(sessId, kickedUserID string) {
		var kickPing time.Duration = 3
		time.Sleep(time.Second * kickPing)

		err = s.GuildMemberMove(sessId, kickedUserID, nil)
		if err != nil {
			log.Printf("failed kick user, err: %v", err)
		}
	}(guildID, i.Member.User.ID)
}

func (c *YanLohCommand) GameInteraction(g *game.Game) {
	//...
}

// RegisterGameCommand command
type RegisterGameCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewRegisterGameCommand() *RegisterGameCommand {
	name := "register_game"
	return &RegisterGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Register new Game",
		},
		name: name,
	}
}

func (c *RegisterGameCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *RegisterGameCommand) GetName() string {
	return c.name
}

func (c *RegisterGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction) {
	responseMessageText := "Registration has begun. \n" +
		"Post any reactions below. And if you want to be a spectator, put the reaction :smiling_imp:"
	channelID := i.ChannelID
	message, err := s.ChannelMessageSend(channelID, responseMessageText)
	if err != nil {
		log.Print(err)
	}
	log.Println("MessageID: ", message.ID)
}

func (c *RegisterGameCommand) GameInteraction(g *game.Game) {
	g = game.NewUndefinedGame()
}

// AddChannelRoleCommand command logic
type AddChannelRoleCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewAddChannelRole() *AddChannelRoleCommand {
	name := "add_channel_role"
	return &AddChannelRoleCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Define a chat room where the interaction between the bot and the role will take place.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "mafia",
					Description: "Add Mafia interaction chat",
					Type:        discordgo.ApplicationCommandOptionString,
				},
				{
					Name:        "doctor",
					Description: "Add Doctor interaction chat",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		name: name,
	}
}

func (c *AddChannelRoleCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *AddChannelRoleCommand) GetName() string {
	return c.name
}

func (c *AddChannelRoleCommand) Execute(s *discordgo.Session, i *discordgo.Interaction) {
	if len(i.ApplicationCommandData().Options) == 0 {
		err := s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Must be option!",
			},
		})
		if err != nil {
			log.Print(err)
		}
		return
	}

	roleName := i.ApplicationCommandData().Options[0].Name
	requestedChatID := i.ApplicationCommandData().Options[0].Value.(string)

	if ok := noticeChat(s, roleName, requestedChatID); ok != nil {
		err := s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid Chat ID!",
			},
		})
		if err != nil {
			log.Print(err)
		}
		return
	}

	messageContent := fmt.Sprintf("Done, now is %v chat is %v chat.", requestedChatID, roleName)
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: messageContent,
		},
	})
	if err != nil {
		log.Print(err)
	}
}

func noticeChat(s *discordgo.Session, chatType, chatID string) error {
	messageContent := fmt.Sprintf("Now this chat for %v role.", chatType)
	message, err := s.ChannelMessageSend(chatID, messageContent)
	message.ID = "1"
	// TODO!
	return err
}

func (c *AddChannelRoleCommand) GameInteraction(g *game.Game) {
	//..
}
