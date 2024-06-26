package game

import (
	"log"
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
	var testedCfg = config.GetConfigByPlayersCountAndIndex(7, 3)

	takeANight := func(g *game.Game, c votesCfg) {

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
					log.Println(votedRole)
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

	t.Run("Excepted No dies", func(t *testing.T) {
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
	})
}
