package game

import (
	"github.com/https-whoyan/MafiaBot/core/converter"
	"strconv"
	"sync"
	"testing"

	"github.com/https-whoyan/MafiaBot/core/config"
	"github.com/https-whoyan/MafiaBot/core/game"
	"github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNightDeadlockEmptyVotes(t *testing.T) {
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
}

func signalHandler(s game.Signal) *roles.Role {
	if sSS, ok := s.(game.SwitchStateSignal); ok {
		if v, ok := sSS.Value.(game.SwitchNightVoteRoleSwitchValue); ok {
			return v.CurrentVotedRole
		}
	}
	return nil
}

func playersHelper(players []*player.Player) map[*roles.Role][]*player.Player {
	mp := make(map[*roles.Role][]*player.Player)
	for _, p := range players {
		mp[p.Role] = append(mp[p.Role], p)
	}
	return mp
}

type voteCfg struct {
	role  *roles.Role
	votes []int
}

func (v voteCfg) ToVotePr(players []*player.Player) *game.VoteProvider {
	return &game.VoteProvider{
		VotedPlayerID:  player.SearchAllPlayersWithRole(players, v.role)[0].Tag,
		Vote:           strconv.Itoa(v.votes[0]),
		IsServerUserID: true,
	}
}

func takeANight(g *game.Game, c votesCfg) {
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
			for s := range ch {
				votedRole := signalHandler(s)
				if votedRole == nil {
					continue
				}
				if votedRole.IsTwoVotes {
					vP := c[votedRole].ToTwoVotePr(g.Active)
					g.TwoVoteChan <- vP
					continue
				}
				vP := c[votedRole].ToVotePr(g.Active)
				g.VoteChan <- vP
				continue
			}
		}
	}()
	wg.Wait()
}

func (v voteCfg) ToTwoVotePr(players []*player.Player) *game.TwoVotesProvider {
	return &game.TwoVotesProvider{
		VotedPlayerID:  player.SearchAllPlayersWithRole(players, v.role)[0].Tag,
		Vote1:          strconv.Itoa(v.votes[0]),
		Vote2:          strconv.Itoa(v.votes[1]),
		IsServerUserID: true,
	}
}

type votesCfg map[*roles.Role]voteCfg

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
				votes: []int{game.EmptyVoteInt},
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
				votes: []int{game.EmptyVoteInt},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{mappedPlayers[roles.Peaceful][2].ID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{game.EmptyVoteInt},
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
				votes: []int{game.EmptyVoteInt},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{doctorID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{game.EmptyVoteInt},
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
				votes: []int{game.EmptyVoteInt},
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
				votes: []int{game.EmptyVoteInt},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{doctorID},
			},
			roles.Don: {
				role:  roles.Don,
				votes: []int{game.EmptyVoteInt},
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
				Count: 4,
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
func TestNightConfig10_3(t *testing.T) {
	t.Parallel()
	var testedCfg = config.GetConfigByPlayersCountAndIndex(10, 2)

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
				votes: []int{game.EmptyVoteInt, game.EmptyVoteInt},
			},
			roles.Don: {
				role:  roles.Detective,
				votes: []int{game.EmptyVoteInt},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{maniacID},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{detectiveID},
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
				votes: []int{game.EmptyVoteInt, game.EmptyVoteInt},
			},
			roles.Don: {
				role:  roles.Detective,
				votes: []int{game.EmptyVoteInt},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{game.EmptyVoteInt},
			},
			roles.Mafia: {
				role:  roles.Mafia,
				votes: []int{game.EmptyVoteInt},
			},
			roles.Doctor: {
				role:  roles.Doctor,
				votes: []int{game.EmptyVoteInt},
			},
			roles.Maniac: {
				role:  roles.Maniac,
				votes: []int{game.EmptyVoteInt},
			},
		}
		takeANight(g, vCfg)
		nightLog := g.NewNightLog()

		exceptedDead := map[int]bool{}
		actualDead := converter.SliceToSet(nightLog.Dead)

		assert.Equal(t, exceptedDead, actualDead)
	})

	t.Run("Excepted one peaceful dies", func(t *testing.T) {
		t.Parallel()
		g, err := initHelper(testedCfg)
		if err != nil {
			t.Fatal(err)
		}
		mappedPlayers := playersHelper(g.Active)
		maniacID := mappedPlayers[roles.Maniac][0].ID

		randomPeacefulID1 := mappedPlayers[roles.Peaceful][0].ID
		randomPeacefulID2 := mappedPlayers[roles.Peaceful][1].ID
		vCfg := votesCfg{
			roles.Detective: {
				role:  roles.Detective,
				votes: []int{game.EmptyVoteInt, game.EmptyVoteInt},
			},
			roles.Don: {
				role:  roles.Detective,
				votes: []int{game.EmptyVoteInt},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{maniacID},
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
				votes: []int{game.EmptyVoteInt, game.EmptyVoteInt},
			},
			roles.Don: {
				role:  roles.Detective,
				votes: []int{game.EmptyVoteInt},
			},
			roles.Whore: {
				role:  roles.Whore,
				votes: []int{randomPeacefulID1},
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
}
