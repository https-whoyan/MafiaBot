package game

import (
	botFMT "github.com/https-whoyan/MafiaBot/internal/bot/fmt"
	coreUserPack "github.com/https-whoyan/MafiaBot/internal/bot/user"
	coreGamePack "github.com/https-whoyan/MafiaBot/internal/core/game"
)

var (
	FMTer           = botFMT.DiscordFMTInstance // Same once struct
	ConstRenameMode = coreGamePack.RenameInGuildMode
)

func GetNewGameConfig(renameProvider *coreUserPack.BotUserRenameProvider) []coreGamePack.GameOption {
	options := []coreGamePack.GameOption{
		coreGamePack.FMTerOpt(FMTer),
		coreGamePack.RenameModeOpt(ConstRenameMode),
		coreGamePack.RenamePrOpt(renameProvider),
	}
	return options
}
