package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	loadDotEnv()
	bot := getBot()
	err := bot.Open()
	if err != nil {
		log.Fatal(err)
	}
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
