package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
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
	fmt.Println("Register.")
	messageID := i.Message.ID
	responseMessageText := "Registration has begun. \n" +
		"Post any reactions below. And if you want to be a spectator, put the reaction :smiling_imp:"
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessageText,
		},
	})
	if err != nil {
		log.Print(err)
	}
	log.Println("MessageID: ", messageID)
}

func (c *RegisterGameCommand) GameInteraction(g *game.Game) (newGameState *game.Game) {
	return game.NewGame(0)
}
