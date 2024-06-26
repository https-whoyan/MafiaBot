package handlers

import (
	"github.com/https-whoyan/MafiaBot/core/game"

	"github.com/https-whoyan/MafiaBot/internal/fmt"
)

// ValidateCommandByGameState validate, is correct command name by current game State.
func ValidateCommandByGameState(commandName string, g *game.Game, fmtEr *fmt.DiscordFMTer) (content string, isOk bool) {
	gameState := g.State

	gameIn := "Game " + g.State.String() + "."
	cantUse := fmtEr.B("Couldn't use /" + commandName + " command")
	content = gameIn + fmtEr.NL() + cantUse

	switch commandName {
	case RegisterGameCommandName:
		switch gameState {
		case game.NonDefinedState:
			return "", true
		case game.FinishState:
			return "", true
		}
	case ChoiceGameConfigCommandName:
		switch gameState {
		case game.RegisterState:
			return "", true
		}
	case StartGameCommandName:
		switch gameState {
		case game.RegisterState:
			return "", true
		}
	case VoteGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	case TwoVoteGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	case FinishGameCommandName:
		return "", true
	}

	return content, false
}
