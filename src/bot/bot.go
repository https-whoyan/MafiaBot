package bot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token    string
	Session  *discordgo.Session
	Commands map[string]Command
}

func InitBot() *Bot {
	token := os.Getenv("BOT_TOKEN")
	botStr := "Bot " + token
	s, err := discordgo.New(botStr)
	if err != nil {
		log.Fatal(err)
	}
	return &Bot{
		token:    token,
		Session:  s,
		Commands: make(map[string]Command),
	}
}
