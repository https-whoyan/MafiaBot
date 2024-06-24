package game

import (
	"fmt"
	"testing"

	"github.com/https-whoyan/MafiaBot/internal/converter"
	"github.com/https-whoyan/MafiaBot/internal/wrap"

	"github.com/https-whoyan/MafiaBot/core/game"

	"github.com/https-whoyan/MafiaBot/tests/core/config"
	"github.com/https-whoyan/MafiaBot/tests/core/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_init(t *testing.T) {
	t.Parallel()

	t.Run("Nothing is set", func(t *testing.T) {
		g := game.GetNewGame(models.TestingGuildID)
		err := g.Init(nil)

		errs := wrap.UnwrapErrToErrSlices(err)
		assert.Equal(t, []error{game.EmptyConfigErr}, errs)
	})

	t.Run("Random Config Set", func(t *testing.T) {
		g := game.GetNewGame(models.TestingGuildID)
		err := g.Init(config.GetRandomConfig())

		errs := wrap.UnwrapErrToErrSlices(err)
		mpErrs := converter.SliceToSet(errs)
		exceptedErrs := map[error]bool{
			game.MismatchPlayersCountAndGamePlayersCountErr: true,
			game.NotFullRoleChannelInfoErr:                  true,
			game.NotMainChannelInfoErr:                      true,
			game.EmptyFMTerErr:                              true,
		}
		assert.Equal(t, exceptedErrs, mpErrs)
	})

	t.Run("All loaded", func(t *testing.T) {
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
			t.Fatal(internalErr)
		}
		for _, roleCh := range allRoleChannels {
			internalErr = g.SetNewRoleChannel(roleCh)
			if internalErr != nil {
				t.Fatal(internalErr)
			}
		}

		cfg := config.GetRandomConfig()
		testPlayers := models.GetTestPlayers(cfg.PlayersCount)
		g.SetStartPlayers(testPlayers)

		err := g.Init(cfg)
		if err != nil {
			require.FailNow(t, fmt.Sprintf("Everything is supplied, but an error is made: %v", err))
		}
	})
}
