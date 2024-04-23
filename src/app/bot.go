package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

func getBot() *discordgo.Session {
	token := os.Getenv("BOT_TOKEN")
	botStr := "Bot " + token
	bot, err := discordgo.New(botStr)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}
