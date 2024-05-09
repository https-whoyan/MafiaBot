package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"

	botPack "github.com/https-whoyan/MafiaBot/internal/bot"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"github.com/https-whoyan/MafiaBot/internal/db"
)

func main() {
	loadDotEnv()
	cfg, err := db.LoadMongoDBConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = db.InitMongoDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("я тут")

	currDB, isContains := db.GetCurrDB()
	fmt.Println(isContains, currDB)
	log.Println("Discord-go version:", discordgo.VERSION)
	log.Println("Discord-go API version:", discordgo.APIVersion)
	bot, err := botPack.InitBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	defer func(b *botPack.Bot) {
		b.Close()
	}(bot)

	// If you need delete all registered commands, use here: bot.DeleteAllGloballyRegisteredCommands()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func test(bot *botPack.Bot) {
	game.GetActiveConnectionsMembers(bot.Session, "")
}
