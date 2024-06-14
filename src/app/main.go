package main

import (
	fmt2 "fmt"
	"github.com/https-whoyan/MafiaBot/config"
	"github.com/https-whoyan/MafiaBot/core/roles"
	"github.com/https-whoyan/MafiaBot/internal/bot/fmt"
)

const (
	isTest = false
)

func main() {
	if isTest {
		test()
		return
	}
	cfg := config.LoadConfig()
	cfg.Run()
}

func test() {
	fmt2.Println(roles.GetDefinitionsOfAllRoles(fmt.DiscordFMTInstance, 2000))
}
