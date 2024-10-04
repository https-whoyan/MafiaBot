package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/fmt"
	botGamePack "github.com/https-whoyan/MafiaBot/internal/game"
	"github.com/https-whoyan/MafiaBot/internal/handlers/names"
	userPack "github.com/https-whoyan/MafiaBot/internal/user"
	"github.com/https-whoyan/MafiaBot/pkg"
	"github.com/https-whoyan/MafiaCore/game"
)

// ValidateCommandByGameState validate, is correct command name by current game State.
func ValidateCommandByGameState(s *discordgo.Session,
	commandName string, g *game.Game, fmtEr *fmt.DiscordFMTer, databases *pkg.Database) (
	content string, isOk bool) {
	gameState := g.GetState()

	gameIn := "Game " + g.GetState().String() + "."
	cantUse := fmtEr.B("Couldn't use /" + commandName + " command")
	content = gameIn + fmtEr.NL() + cantUse

	switch commandName {
	case names.RegisterGameCommandName:
		switch gameState {
		case game.NonDefinedState:
			return "", true
		case game.FinishState:
			userRenameProvider := userPack.NewBotUserRenameProvider(s, g.GuildID())
			*g = *game.GetNewGame(g.GuildID(), botGamePack.GetNewGameConfig(
				userRenameProvider, databases.Storage)...)
			return "", true
		default:
			break
		}
	case names.ChoiceGameConfigCommandName:
		switch gameState {
		case game.RegisterState:
			return "", true
		default:
			break
		}
	case names.StartGameCommandName:
		switch gameState {
		case game.InitState:
			return "", true
		default:
			break
		}
	case names.VoteGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	case names.TwoVoteGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	case names.DayVoteGameCommandName:
		switch gameState {
		case game.DayState:
			return "", true
		default:
			return "", true
		}
	case names.FinishGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	}

	return content, false
}
