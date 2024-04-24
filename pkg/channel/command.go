package channel

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type AddChannelRoleCommand struct {
	cmd *discordgo.ApplicationCommand
}

func NewAddChannelRole() *AddChannelRoleCommand {
	return &AddChannelRoleCommand{
		cmd: &discordgo.ApplicationCommand{
			Name: "add-channel-role",
			Description: "Allows you to define the chat that " +
				"will be used for the role to send messages about its selection.",
		},
	}
}

func (c *AddChannelRoleCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("sldksldk")
	if m.Author.ID == s.State.User.ID {
		return
	}
	args := strings.Split(m.Content, " ")[2:]
	message, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v", args))
	fmt.Println(message, err)
}

func (c *AddChannelRoleCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}
