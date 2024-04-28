package game

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type RegisterGameCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewRegisterGameCommand() *RegisterGameCommand {
	name := "register_game"
	return &RegisterGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Register new Game",
		},
		name: name,
	}
}

func (c *RegisterGameCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *RegisterGameCommand) GetName() string {
	return c.name
}

func (c *RegisterGameCommand) GetExecuteFunc() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return c.Execute
}

func (c *RegisterGameCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Print("RegisterGameCommand TODO!")
	// TODO!
}
