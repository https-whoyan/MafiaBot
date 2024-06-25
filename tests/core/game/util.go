package game

import (
	"github.com/https-whoyan/MafiaBot/core/config"
	"github.com/https-whoyan/MafiaBot/core/game"
	"github.com/https-whoyan/MafiaBot/tests/core/models"
)

func initHelper(cfg *config.RolesConfig) (*game.Game, error) {
	var internalErr error

	opts := []game.GameOption{
		game.FMTerOpt(models.TestFMTInstance),
		game.RenamePrOpt(models.TestRenameUserProviderInstance),
	}
	g := game.GetNewGame(models.TestingGuildID, opts...)

	allRoleChannels := models.NewTestChannels()
	mainChannel := models.NewTestMainChannels()

	internalErr = g.SetMainChannel(mainChannel)
	if internalErr != nil {
		return nil, internalErr
	}
	for _, roleCh := range allRoleChannels {
		internalErr = g.SetNewRoleChannel(roleCh)
		if internalErr != nil {
			return nil, internalErr
		}
	}

	testPlayers := models.GetTestPlayers(cfg.PlayersCount)
	g.SetStartPlayers(testPlayers)
	internalErr = g.Init(cfg)
	if internalErr != nil {
		return nil, internalErr
	}
	return g, nil
}
