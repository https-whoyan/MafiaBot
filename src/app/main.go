package main

import (
	"github.com/https-whoyan/MafiaBot/config"
)

func main() {
	cfg := config.LoadConfig()
	cfg.Run()
}
