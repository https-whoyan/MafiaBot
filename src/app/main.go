package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	botPack "github.com/https-whoyan/MafiaBot/src/bot"
)

func main() {
	loadDotEnv()
	bot := botPack.InitBot()
	bot.InitBotCommands()
	bot.RegisterHandlers()
	err := bot.Session.Open()
	bot.RegisterCommands()
	if err != nil {
		log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	defer func(b *botPack.Bot) {
		b.Close()
		b.RemoveRegisteredCommands()
	}(bot)
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
