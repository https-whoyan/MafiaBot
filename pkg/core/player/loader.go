package player

import (
	"errors"
	"log"
	"math/rand"

	"github.com/https-whoyan/MafiaBot/core/config"
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

// _____________________________________________________________
// Load Players
// 2 different player loading options for your convenience
// _____________________________________________________________

// First

func GenerateEmptyPlayersByTagsAndUsernames(tags []string, usernames []string, isAllSpectators bool) []*Player {
	if len(tags) != len(usernames) {
		log.Println("Unexpected mismatch of playing participants and nicknames")
		return []*Player{}
	}
	players := make([]*Player, len(tags))
	for i, tag := range tags {
		var newPlayer *Player
		if isAllSpectators {
			newPlayer = NewSpectator(tag, usernames[i])
		}
		newPlayer = NewEmptyPlayer(tag, usernames[i])
		players = append(players, newPlayer)
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

func GetUsernamesByPlayers(players []*Player) []string {
	names := make([]string, len(players))
	for i, player := range players {
		names[i] = player.OldNick
	}
	return names
}

// Second

func GenerateEmptyPlayersByFunc(
	x any,
	getTagAndUsername func(x any, index int) (string, string),
	countOfNewPlayers int, isAllSpectators bool) []*Player {

	isRecovered := false
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			log.Println("Invalid usage of GenerateEmptyPlayersByFunc function! Return empty slice")
			isRecovered = true
		}
	}()

	players := make([]*Player, countOfNewPlayers)

	for i := 0; i <= countOfNewPlayers-1; i++ {
		tag, username := getTagAndUsername(x, i)
		var newPlayer *Player
		if isAllSpectators {
			newPlayer = NewSpectator(tag, username)
		} else {
			newPlayer = NewEmptyPlayer(tag, username)
		}
		players[i] = newPlayer
	}

	if isRecovered {
		return []*Player{}
	}
	return players
}
