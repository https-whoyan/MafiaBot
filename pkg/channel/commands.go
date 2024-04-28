package channel

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

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

func (c *AddChannelRoleCommand) GetExecuteFunc() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return c.Execute
}

func (c *AddChannelRoleCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
	_, err := s.ChannelMessageSend(chatID, messageContent)
	return err
}
