package handlers

import (
	botFMTPack "github.com/https-whoyan/MafiaBot/internal/fmt"
	coreGamePack "github.com/https-whoyan/MafiaCore/game"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	// Execute Would say it's a discord command handler.
	// Uses session, interaction, and formatter.
	// Changes the game.
	Execute(s *discordgo.Session, i *discordgo.Interaction, g *coreGamePack.Game, f *botFMTPack.DiscordFMTer)
	// GetCmd All implementations of the command use ApplicationCommand under the hood.
	// It is also necessary for the bot to remember the commands in order to delete them.
	GetCmd() *discordgo.ApplicationCommand
	// GetName provide information about command name.
	// Need to validate discord-go handler.
	//
	// Also seen from https://github.com/bwmarrin/discordgo/tree/master/examples/slash_commands
	GetName() string
	// IsUsedForGame Not all commands are used to interact with the game.
	//
	//In case a command changes the game, I check if the game is in bot.Games and if not, I say there is no game, please create one.
	//If the command is not used for a game, I execute it.
	IsUsedForGame() bool
}
