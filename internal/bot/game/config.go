package game

import (
	botFMT "github.com/https-whoyan/MafiaBot/internal/bot/fmt"
	"github.com/https-whoyan/MafiaBot/internal/bot/user"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
)

var (
	FMTer           = botFMT.FMTInstance // Same once struct
	ConstRenameMode = game.RenameInGuildMode
)

func GetNewGameConfig(renameProvider *user.BotUserRenameProvider) game.Setting {
	return game.Setting{
		FMTer:          FMTer,
		RenameProvider: renameProvider,
		RenameMode:     ConstRenameMode,
	}
}
