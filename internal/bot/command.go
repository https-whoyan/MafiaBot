package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
)

type Command interface {
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
	GetCmd() *discordgo.ApplicationCommand
	GetName() string
	GameInteraction(b *game.Game) *game.Game
}
