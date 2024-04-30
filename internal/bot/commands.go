package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"log"
	"time"
)

// YanLohCommand command
type YanLohCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewYanLohCommand() *YanLohCommand {
	name := "yan_loh"
	return &YanLohCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Call Yan with this command!",
		},
		name: name,
	}
}

func (c *YanLohCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *YanLohCommand) GetName() string {
	return c.name
}

func (c *YanLohCommand) GetExecuteFunc() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return c.Execute
}

func (c *YanLohCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Yan Loh.")
	messageContent := "Возможно, что ян и лох. И древлян. Но что бы его же ботом его обзывать..."
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: messageContent,
		},
	})
	if err != nil {
		log.Print(err)
	}
	// async kick requester
	guildID := s.State.Guilds[0].ID
	go func(sessId, kickedUserID string) {
		var kickPing time.Duration = 3
		time.Sleep(time.Second * kickPing)

		err = s.GuildMemberMove(sessId, kickedUserID, nil)
		if err != nil {
			log.Printf("failed kick user, err: %v", err)
		}
	}(guildID, i.Interaction.Member.User.ID)

}

func (c *YanLohCommand) GameInteraction(g *game.Game) *game.Game {
	return g
}
