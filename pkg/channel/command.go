package channel

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type AddChannelRoleCommand struct {
	cmd *discordgo.ApplicationCommand
}

func NewAddChannelRole() *AddChannelRoleCommand {
	return &AddChannelRoleCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        "test_command",
			Description: "Yan Loh",
		},
	}
}

func (c *AddChannelRoleCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Yan Loh",
		},
	})
	if err != nil {
		log.Print(err)
	}
}

func (c *AddChannelRoleCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}
