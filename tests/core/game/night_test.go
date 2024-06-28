package game

import (
	"github.com/https-whoyan/MafiaBot/core/converter"
	"sync"
	"testing"

	"github.com/https-whoyan/MafiaBot/core/config"
	"github.com/https-whoyan/MafiaBot/core/game"
	"github.com/https-whoyan/MafiaBot/core/roles"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	NonVote = game.EmptyVoteInt
)

func TestNightCheckDeadlock_EmptyVotes(t *testing.T) {
	deadlockCheckerHelper := func(cfg *config.RolesConfig) {
		g, err := initHelper(cfg)
		if err != nil {
			require.Fail(t, err.Error())
		}

		ch := make(chan game.Signal)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			g.Night(ch)
			close(ch)
		}()
		go func() {
			defer wg.Done()
			select {
			case <-ch:
				return
			default:
				for range ch {
				}
			}
		}()
		wg.Wait()
	}

	t.Run("Empty votes, config 5;1", func(t *testing.T) {
		t.Parallel()
		cfg := config.GetConfigByPlayersCountAndIndex(5, 1)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 6;0", func(t *testing.T) {
		t.Parallel()
		cfg := config.GetConfigByPlayersCountAndIndex(6, 0)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 7;2", func(t *testing.T) {
		t.Parallel()
		cfg := config.GetConfigByPlayersCountAndIndex(7, 2)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 9;0", func(t *testing.T) {
		t.Parallel()
		cfg := config.GetConfigByPlayersCountAndIndex(9, 0)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 10;0", func(t *testing.T) {
		t.Parallel()
		cfg := config.GetConfigByPlayersCountAndIndex(10, 0)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 11;1", func(t *testing.T) {
		t.Parallel()
		cfg := config.GetConfigByPlayersCountAndIndex(11, 1)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 13;2", func(t *testing.T) {
		t.Parallel()
		cfg := config.GetConfigByPlayersCountAndIndex(13, 2)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 14;0", func(t *testing.T) {
		t.Parallel()
		cfg := config.GetConfigByPlayersCountAndIndex(14, 0)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
}

/*
PlayersCount: 7,
RolesMp: map[string]*RoleConfig{
	"Peaceful": {
		Role:  roles.Peaceful,
		Count: 3,
	},
	"Doctor": {
		Role:  roles.Doctor,
		Count: 1,
	},
	"Whore": {
		Role:  roles.Whore,
		Count: 1,
	},
	"Mafia": {
		Role:  roles.Mafia,
		Count: 1,
	},
	"Don": {
		Role:  roles.Don,
		Count: 1,
	},
},
*/

func TestNightConfig7_3(t *testing.T) {
	t.Parallel()
	var testedCfg = config.GetConfigByPlayersCountAndIndex(7, 3)

	t.Run("Excepted No dies", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{mappedPlayers[roles.Mafia][0].ID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{mappedPlayers[roles.Peaceful][2].ID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{mappedPlayers[roles.Mafia][0].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})
	t.Run("Excepted No dies, 2", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{NonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{mappedPlayers[roles.Peaceful][2].ID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{mappedPlayers[roles.Peaceful][2].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})
	t.Run("Excepted No dies, 3", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		doctorID := mappedPlayers[roles.Doctor][0].ID
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{NonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{doctorID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{doctorID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	t.Run("Excepted Detective die", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		detectiveID := mappedPlayers[roles.Doctor][0].ID
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{mappedPlayers[roles.Doctor][0].ID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{detectiveID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{detectiveID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{
			detectiveID: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	t.Run("Excepted Doctor die", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		doctorID := mappedPlayers[roles.Doctor][0].ID
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{NonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{doctorID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{mappedPlayers[roles.Peaceful][0].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{
			doctorID: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})
}

/*
	{
		PlayersCount: 10,
		RolesMp: map[string]*RoleConfig{
			"Peaceful": {
				Role:  roles.Peaceful,
				Count: 3,
			},
			"Citizen": {
				Role:  roles.Citizen,
				Count: 1,
			},
			"Doctor": {
				Role:  roles.Doctor,
				Count: 1,
			},
			"Whore": {
				Role:  roles.Whore,
				Count: 1,
			},
			"Detective": {
				Role:  roles.Detective,
				Count: 1,
			},
			"Mafia": {
				Role:  roles.Mafia,
				Count: 1,
			},
			"Don": {
				Role:  roles.Don,
				Count: 1,
			},
			"Maniac": {
				Role:  roles.Maniac,
				Count: 1,
			},
		},
	}
*/

func TestNightConfig10_2(t *testing.T) {
	t.Parallel()
	var testedCfg = config.GetConfigByPlayersCountAndIndex(10, 2)

	// No citizen test

	t.Run("Excepted No dies", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		detectiveID := mappedPlayers[roles.Detective][0].ID
		maniacID := mappedPlayers[roles.Maniac][0].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{maniacID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{detectiveID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{NonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{detectiveID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{mappedPlayers[roles.Peaceful][0].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})
	t.Run("Excepted No dies, 2", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{NonVote},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{NonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{NonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{NonVote},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{NonVote},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	t.Run("Excepted one peaceful die", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		maniacID := mappedPlayers[roles.Maniac][0].ID
		doctorID := mappedPlayers[roles.Doctor][0].ID

		randomPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		randomPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{maniacID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{doctorID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{randomPeacefulID2},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{maniacID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{randomPeacefulID1},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{
			randomPeacefulID2: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	t.Run("Excepted two peaceful dies", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		maniacID := mappedPlayers[roles.Maniac][0].ID

		randomPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		randomPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		randomPeacefulID3 := mappedPlayers[roles.Peaceful][2].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{randomPeacefulID1},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{NonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{randomPeacefulID2},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{maniacID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{randomPeacefulID3},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{randomPeacefulID3},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{
			randomPeacefulID2: true,
			randomPeacefulID3: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	// Citizen test

	t.Run("Citizen test: Excepted No dies: citizen saved", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		detectiveID := mappedPlayers[roles.Detective][0].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{NonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{detectiveID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{detectiveID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{NonVote},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{NonVote},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})
	t.Run("Citizen test: Excepted No dies, 2; all voted to detective))", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)

		detectiveID := mappedPlayers[roles.Detective][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{detectiveID},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{detectiveID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{detectiveID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{detectiveID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{detectiveID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{detectiveID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})
	t.Run("Citizen test: Excepted no dies, 3; doctor saved citizen, citizen saved peaceful", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)

		donID := mappedPlayers[roles.Don][0].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID
		rndPeacefulID := mappedPlayers[roles.Peaceful][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{donID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{rndPeacefulID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{citizenID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{rndPeacefulID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	t.Run("Citizen test: Excepted one die: doctor", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)

		doctorID := mappedPlayers[roles.Doctor][0].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{citizenID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{doctorID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{doctorID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{NonVote},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{doctorID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{
			doctorID: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})
	t.Run("Citizen test: Excepted one die: citizen, citizen muted", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)

		citizenID := mappedPlayers[roles.Citizen][0].ID
		rndPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{citizenID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{rndPeacefulID1},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{NonVote},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{NonVote},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{
			citizenID: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	t.Run("Citizen test: Excepted two die: one peaceful and and citizen", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)

		rndPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		rndPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{NonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{rndPeacefulID1},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{rndPeacefulID2},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{rndPeacefulID2},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{
			rndPeacefulID1: true,
			citizenID:      true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	t.Run("Citizen test: Excepted three dies: two peaceful dies and citizen (lol)", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)

		rndPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		rndPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID

		detectiveID := mappedPlayers[roles.Detective][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{NonVote, NonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{NonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{NonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []int{rndPeacefulID1},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{detectiveID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{rndPeacefulID2},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{
			rndPeacefulID1: true,
			rndPeacefulID2: true,
			citizenID:      true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})
}
