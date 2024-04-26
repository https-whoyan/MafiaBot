package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/pkg/channel"
	"log"
)

type Command interface {
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
	GetCmd() *discordgo.ApplicationCommand
}

func (b *Bot) initCommand(c Command) {
	commandName := c.GetCmd().Name
	b.Commands[commandName] = c
}

func (b *Bot) InitBotCommands() {
	b.initCommand(channel.NewAddChannelRole())
}

func (b *Bot) RegisterHandlers() {
	log.Print("Register handlers")
	for _, cmd := range b.Commands {
		b.Session.AddHandler(cmd.Execute)
	}
}

func (b *Bot) RemoveRegisteredCommands() {
	log.Println("Remove commands")
	stateId := b.Session.State.User.ID
	for _, registeredCmd := range b.registeredCommands {
		err := b.Session.ApplicationCommandDelete(stateId, "", registeredCmd.ID)
		if err != nil {
			log.Printf("cannot delete command, err: %v", err)
		}
	}
}

func (b *Bot) RegisterCommands() {
	log.Println("Register commands")
	stateId := b.Session.State.User.ID
	for _, cmd := range b.Commands {
		registeredCmd, err := b.Session.ApplicationCommandCreate(stateId, "", cmd.GetCmd())
		if err != nil {
			log.Print(err)
		}
		b.registeredCommands = append(b.registeredCommands, registeredCmd)
	}
}
