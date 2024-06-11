package player

import (
	"errors"
	"math/rand"

	"github.com/https-whoyan/MafiaBot/internal/core/config"
)

// ___________________________________
// Use to start a game starting.
// Role reversal, to put it simply.
// ___________________________________

func generateRandomOrderToIDs(n int) []int {
	var IDs []int
	for i := 1; i <= n; i++ {
		IDs = append(IDs, i)
	}
	rand.Shuffle(n, func(i, j int) {
		IDs[i], IDs[j] = IDs[j], IDs[i]
	})

	return IDs
}

func GeneratePlayers(tags []string, oldUsernames []string, cfg *config.RolesConfig) ([]*Player, error) {
	if len(tags) != cfg.PlayersCount {
		return []*Player{}, errors.New("unexpected mismatch of playing participants and configs")
	}
	if len(tags) != len(oldUsernames) {
		return []*Player{}, errors.New("unexpected mismatch of playing participants and nicknames")
	}

	n := len(tags)
	IDs := generateRandomOrderToIDs(n)
	rolesArr := config.GetShuffledRolesConfig(cfg)

	players := make([]*Player, n)

	for i := 0; i <= n-1; i++ {
		players[i] = &Player{
			ID:                IDs[i],
			OldNick:           oldUsernames[i],
			Tag:               tags[i],
			Role:              rolesArr[i],
			Vote:              -1,
			LifeStatus:        Alive,
			InteractionStatus: Passed,
		}
	}

	return players, nil
}

func GetUsernamesByPlayers(players []*Player) []string {
	names := make([]string, len(players))
	for i, player := range players {
		names[i] = player.OldNick
	}
	return names
}

func GenerateEmptyPlayers(tags []string, usernames []string) []*Player {
	players := make([]*Player, 0)
	for i, tag := range tags {
		players = append(players, &Player{
			Tag:     tag,
			OldNick: usernames[i],
		})
	}
	return players
}

func GetTagsByPlayers(players []*Player) []string {
	var tags []string
	for _, player := range players {
		tags = append(tags, player.Tag)
	}
	return tags
}
