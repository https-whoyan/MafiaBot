package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/pkg/channel"
	"log"
)

type Command interface {
	Execute(s *discordgo.Session, m *discordgo.MessageCreate)
	GetCmd() *discordgo.ApplicationCommand
}

var (
	commands []Command
)

func (b *Bot) InitBotCommands() {
	// Add role channel command
	addChannelRoleCommand := channel.NewAddChannelRole()
	commands = append(commands, addChannelRoleCommand)
	commandName := addChannelRoleCommand.GetCmd().Name
	b.Commands[commandName] = addChannelRoleCommand
}

func (b *Bot) RegisterHandlers() {
	for _, cmd := range b.Commands {
		b.Session.AddHandler(cmd.Execute)
	}
}

func (b *Bot) RegisterCommands() {
	stateId := b.Session.State.User.ID
	for _, cmd := range b.Commands {
		_, err := b.Session.ApplicationCommandCreate(stateId, "", cmd.GetCmd())
		if err != nil {
			log.Fatal(err)
		}
	}
}
