package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/core/game"
	"github.com/https-whoyan/MafiaBot/internal/fmt"
	botGamePack "github.com/https-whoyan/MafiaBot/internal/game"
	userPack "github.com/https-whoyan/MafiaBot/internal/user"
)

// ValidateCommandByGameState validate, is correct command name by current game State.
func ValidateCommandByGameState(s *discordgo.Session, commandName string, g *game.Game, fmtEr *fmt.DiscordFMTer) (
	content string, isOk bool) {
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
			userRenameProvider := userPack.NewBotUserRenameProvider(s, g.GuildID)
			*g = *game.GetNewGame(g.GuildID, botGamePack.GetNewGameConfig(userRenameProvider)...)
			return "", true
		default:
			break
		}
	case ChoiceGameConfigCommandName:
		switch gameState {
		case game.RegisterState:
			return "", true
		default:
			break
		}
	case StartGameCommandName:
		switch gameState {
		case game.InitState:
			return "", true
		default:
			break
		}
	case VoteGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	case TwoVoteGameCommandName:
		if g.IsRunning() {
			return "", true
		}
	case DayVoteGameCommandName:
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
