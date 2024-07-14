package game

import (
	"fmt"
	"math"
	"testing"

	"github.com/https-whoyan/MafiaBot/core/config"
	"github.com/https-whoyan/MafiaBot/core/game"
)

// SeeCalculationDayDeadline Designed to look at giving time for the day under certain circumstances.
// May be useful if you want to change your formula or weights.
func Test_SeeCalculationDayDeadline(t *testing.T) {
	cfg1 := config.GetConfigByPlayersCountAndIndex(5, 1)
	cfg2 := config.GetConfigByPlayersCountAndIndex(7, 0)
	cfg3 := config.GetConfigByPlayersCountAndIndex(8, 2)
	cfg4 := config.GetConfigByPlayersCountAndIndex(10, 1)
	cfg5 := config.GetConfigByPlayersCountAndIndex(12, 0)
	cfg6 := config.GetConfigByPlayersCountAndIndex(13, 1)

	type seeCase struct {
		cfg              *config.RolesConfig
		deadPlayersCount int
		nightCount       int
	}

	tests := []seeCase{
		// cfg1
		{cfg: cfg1, deadPlayersCount: 0, nightCount: 1},
		{cfg: cfg1, deadPlayersCount: 2, nightCount: 3},
		{cfg: cfg1, deadPlayersCount: 3, nightCount: 3},

		// cfg2
		{cfg: cfg2, deadPlayersCount: 0, nightCount: 1},
		{cfg: cfg2, deadPlayersCount: 0, nightCount: 1},
		{cfg: cfg2, deadPlayersCount: 2, nightCount: 3},

		// cfg3
		{cfg: cfg3, deadPlayersCount: 1, nightCount: 1},
		{cfg: cfg3, deadPlayersCount: 1, nightCount: 3},
		{cfg: cfg3, deadPlayersCount: 4, nightCount: 4},

		// cfg4
		{cfg: cfg4, deadPlayersCount: 1, nightCount: 1},
		{cfg: cfg4, deadPlayersCount: 0, nightCount: 2},
		{cfg: cfg4, deadPlayersCount: 2, nightCount: 4},

		// cfg5
		{cfg: cfg5, deadPlayersCount: 1, nightCount: 1},
		{cfg: cfg5, deadPlayersCount: 4, nightCount: 5},
		{cfg: cfg5, deadPlayersCount: 6, nightCount: 5},

		// cfg6
		{cfg: cfg6, deadPlayersCount: 0, nightCount: 1},
		{cfg: cfg6, deadPlayersCount: 2, nightCount: 1},
		{cfg: cfg6, deadPlayersCount: 3, nightCount: 6},
	}

	giveMinutes := func(test seeCase) int {
		deadline := game.CalculateDayDeadline(test.nightCount, test.deadPlayersCount, test.cfg.PlayersCount)
		return int(math.Ceil(deadline.Minutes()))
	}

	for _, test := range tests {
		fmt.Printf("Configs total players: %v, dead players count: %v, night count: %v\n",
			test.cfg.PlayersCount, test.deadPlayersCount, test.nightCount)
		fmt.Printf("Deadline: %v\n+++++++++++++++\n", giveMinutes(test))
	}
}

func Test_Day(t *testing.T) {
	t.Parallel()
}
