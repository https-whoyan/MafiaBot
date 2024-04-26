package bot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token              string
	Session            *discordgo.Session
	Commands           map[string]Command
	registeredCommands []*discordgo.ApplicationCommand
}

func InitBot() *Bot {
	token := os.Getenv("BOT_TOKEN")
	botStr := "Bot " + token
	s, err := discordgo.New(botStr)
	if err != nil {
		log.Fatal(err)
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	return &Bot{
		token:    token,
		Session:  s,
		Commands: make(map[string]Command),
	}
}

func (b *Bot) Close() {
	err := b.Session.Close()
	if err != nil {
		log.Printf("err in closing bot, err: %v", err)
	}
}

// DeleteAllGloballyRegisteredCommands Delete all registered to bot functions. Globally Registered
func (b *Bot) DeleteAllGloballyRegisteredCommands() {
	userId := b.Session.State.User.ID
	globallyRegisteredCommands, err := b.Session.ApplicationCommands(userId, "")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(globallyRegisteredCommands)
	for _, command := range globallyRegisteredCommands {
		err = b.Session.ApplicationCommandDelete(userId, "", command.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
