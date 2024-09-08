package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/fmt"
	botGamePack "github.com/https-whoyan/MafiaBot/internal/game"
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
	case RegisterGameCommandName:
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
	case choiceGameConfigCommandName:
		switch gameState {
		case game.RegisterState:
			return "", true
		default:
			break
		}
	case startGameCommandName:
		switch gameState {
		case game.InitState:
			return "", true
		default:
			break
		}
	case voteGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	case twoVoteGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	case dayVoteGameCommandName:
		switch gameState {
		case game.DayState:
			return "", true
		default:
			return "", true
		}
	case FinishGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	}

	return content, false
}
