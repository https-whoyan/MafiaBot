package main

import (
	"github.com/https-whoyan/MafiaBot/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	cfg.Run()
}
