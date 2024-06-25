package game

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"

	"github.com/https-whoyan/MafiaBot/core/config"
	"github.com/https-whoyan/MafiaBot/core/game"

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

	t.Parallel()

	t.Run("Empty votes, config 5;1", func(t *testing.T) {
		cfg := config.GetConfigByPlayersCountAndIndex(5, 1)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 6;0", func(t *testing.T) {
		cfg := config.GetConfigByPlayersCountAndIndex(6, 0)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 7;2", func(t *testing.T) {
		cfg := config.GetConfigByPlayersCountAndIndex(7, 2)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 9;0", func(t *testing.T) {
		cfg := config.GetConfigByPlayersCountAndIndex(9, 0)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
	t.Run("Empty votes, config 10;0", func(t *testing.T) {
		cfg := config.GetConfigByPlayersCountAndIndex(10, 0)
		deadlockCheckerHelper(cfg)
		assert.True(t, true, "No deadlock")
	})
}
