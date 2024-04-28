package channel

import (
	"github.com/bwmarrin/discordgo"
	"log"
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
	messageContent := "Yan Loh"
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
