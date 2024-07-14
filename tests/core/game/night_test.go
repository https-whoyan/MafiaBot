package game

import (
	"github.com/https-whoyan/MafiaBot/core/converter"
	"github.com/https-whoyan/MafiaBot/core/player"
	"sync"
	"testing"

	"github.com/https-whoyan/MafiaBot/core/config"
	"github.com/https-whoyan/MafiaBot/core/game"
	"github.com/https-whoyan/MafiaBot/core/roles"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNightCheckDeadlock_EmptyVotes(t *testing.T) {
	t.Parallel()
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
		mappedPlayers := playersHelper(*g.Active)
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{mappedPlayers[roles.Mafia][0].ID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{mappedPlayers[roles.Peaceful][2].ID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{mappedPlayers[roles.Mafia][0].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})
	t.Run("Excepted No dies, 2", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{nonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{mappedPlayers[roles.Peaceful][2].ID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{mappedPlayers[roles.Peaceful][2].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})
	t.Run("Excepted No dies, 3", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		doctorID := mappedPlayers[roles.Doctor][0].ID
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{nonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{doctorID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{doctorID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})

	t.Run("Excepted Detective die", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		detectiveID := mappedPlayers[roles.Doctor][0].ID
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{mappedPlayers[roles.Doctor][0].ID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{detectiveID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{detectiveID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{
			detectiveID: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})

	t.Run("Excepted Doctor die", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		doctorID := mappedPlayers[roles.Doctor][0].ID
		vCfg := votesCfg{
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{nonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{doctorID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{mappedPlayers[roles.Peaceful][0].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{
			doctorID: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
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
		mappedPlayers := playersHelper(*g.Active)
		detectiveID := mappedPlayers[roles.Detective][0].ID
		maniacID := mappedPlayers[roles.Maniac][0].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{maniacID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{detectiveID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{nonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{detectiveID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{mappedPlayers[roles.Peaceful][0].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
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
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{nonVote},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{nonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{nonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{nonVote},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{nonVote},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})

	t.Run("Excepted one peaceful die", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		maniacID := mappedPlayers[roles.Maniac][0].ID
		doctorID := mappedPlayers[roles.Doctor][0].ID

		randomPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		randomPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{maniacID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{doctorID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{randomPeacefulID2},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{maniacID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{randomPeacefulID1},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{
			randomPeacefulID2: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})

	t.Run("Excepted two peaceful dies", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		maniacID := mappedPlayers[roles.Maniac][0].ID

		randomPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		randomPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		randomPeacefulID3 := mappedPlayers[roles.Peaceful][2].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{randomPeacefulID1},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{nonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{randomPeacefulID2},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{maniacID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{randomPeacefulID3},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{randomPeacefulID3},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{
			randomPeacefulID2: true,
			randomPeacefulID3: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})

	// Citizen test

	t.Run("Citizen test: Excepted No dies: citizen saved", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		detectiveID := mappedPlayers[roles.Detective][0].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{nonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{detectiveID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{detectiveID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{nonVote},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{nonVote},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})
	t.Run("Citizen test: Excepted No dies, 2; all voted to detective))", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)

		detectiveID := mappedPlayers[roles.Detective][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{detectiveID},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{detectiveID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{detectiveID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{detectiveID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{detectiveID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{detectiveID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})
	t.Run("Citizen test: Excepted no dies, 3; doctor saved citizen, citizen saved peaceful", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)

		donID := mappedPlayers[roles.Don][0].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID
		rndPeacefulID := mappedPlayers[roles.Peaceful][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{donID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{rndPeacefulID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{citizenID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{rndPeacefulID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})

	t.Run("Citizen test: Excepted one die: doctor", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)

		doctorID := mappedPlayers[roles.Doctor][0].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{citizenID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{doctorID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{doctorID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{nonVote},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{doctorID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{
			doctorID: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})
	t.Run("Citizen test: Excepted one die: citizen, citizen muted", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)

		citizenID := mappedPlayers[roles.Citizen][0].ID
		rndPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{citizenID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{rndPeacefulID1},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{nonVote},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{nonVote},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{
			citizenID: true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})

	t.Run("Citizen test: Excepted two die: one peaceful and and citizen", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)

		rndPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		rndPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{nonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{rndPeacefulID1},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{rndPeacefulID2},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{rndPeacefulID2},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{
			rndPeacefulID1: true,
			citizenID:      true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})

	t.Run("Citizen test: Excepted three dies: two peaceful dies and citizen (lol)", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)

		rndPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		rndPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID

		detectiveID := mappedPlayers[roles.Detective][0].ID

		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{nonVote, nonVote},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{nonVote},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{rndPeacefulID1},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{detectiveID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{rndPeacefulID2},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[player.IDType]bool{
			rndPeacefulID1: true,
			rndPeacefulID2: true,
			citizenID:      true,
		}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, convertIntMpToIDTypeMp(actualDead))
	})
}

func TestReincarnation(t *testing.T) {
	t.Parallel()
	var testedCfg = config.GetConfigByPlayersCountAndIndex(10, 2)

	t.Run("Excepted Reincarnation", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		detectiveID := mappedPlayers[roles.Detective][0].ID
		citizenID := mappedPlayers[roles.Citizen][0].ID
		mafiaID := mappedPlayers[roles.Mafia][0].ID
		rndPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		don := mappedPlayers[roles.Don][0]
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{citizenID, don.ID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{citizenID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{mafiaID},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{rndPeacefulID1},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{detectiveID},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{detectiveID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{citizenID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		ch := make(chan game.Signal)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			g.AffectNight(nightLog, ch)
			close(ch)
		}()
		go func() {
			defer wg.Done()
			for range ch {
			}
		}()
		wg.Wait()

		exceptedDoneRole := roles.Mafia
		actualDonRole := don.Role

		assert.Equal(t, exceptedDoneRole, actualDonRole)
	})

	t.Run("Excepted no reincarnation", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(*g.Active)
		detectiveID := mappedPlayers[roles.Detective][0].ID
		maniacID := mappedPlayers[roles.Maniac][0].ID
		don := mappedPlayers[roles.Don][0]
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []player.IDType{maniacID, don.ID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []player.IDType{nonVote},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []player.IDType{maniacID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []player.IDType{detectiveID},
			},
			roles.Citizen: {
				role:  roles.Citizen,
				votes: []player.IDType{nonVote},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []player.IDType{detectiveID},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []player.IDType{mappedPlayers[roles.Peaceful][0].ID},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		ch := make(chan game.Signal)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			g.AffectNight(nightLog, ch)
			close(ch)
		}()
		go func() {
			defer wg.Done()
			for range ch {
			}
		}()
		wg.Wait()

		exceptedDoneRole := roles.Don
		actualDonRole := don.Role

		assert.Equal(t, exceptedDoneRole, actualDonRole)
	})
}
