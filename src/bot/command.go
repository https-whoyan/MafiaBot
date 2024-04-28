package bot

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
	GetCmd() *discordgo.ApplicationCommand
	GetName() string
	GetExecuteFunc() func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
