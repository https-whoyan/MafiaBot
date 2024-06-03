package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
)

type Command interface {
	Execute(s *discordgo.Session, i *discordgo.Interaction, g *game.Game)
	GetCmd() *discordgo.ApplicationCommand
	GetName() string
	IsUsedForGame() bool
}
