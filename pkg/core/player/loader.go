package player

import (
	"errors"
	"github.com/https-whoyan/MafiaBot/core/config"
	"log"
)

// ___________________________________
// Use to start a game starting.
// Role reversal, to put it simply.
// ___________________________________

func generateListToN(n int) []int {
	var IDs []int
	for i := 1; i <= n; i++ {
		IDs = append(IDs, i)
	}

	return IDs
}

func GeneratePlayers(tags []string, oldUsernames []string,
	serverUsernames []string, cfg *config.RolesConfig) ([]*Player, error) {
	if len(tags) != cfg.PlayersCount {
		return []*Player{}, errors.New("unexpected mismatch of playing participants and configs")
	}
	if len(tags) != len(oldUsernames) {
		return []*Player{}, errors.New("unexpected mismatch of playing participants and nicknames")
	}

	n := len(tags)
	IDs := generateListToN(n)
	rolesArr := cfg.GetShuffledRolesConfig()

	players := make([]*Player, n)

	for i := 0; i <= n-1; i++ {
		players[i] = NewPlayer(IDs[i], tags[i], oldUsernames[i], serverUsernames[i], rolesArr[i])
	}

	return players, nil
}

// _____________________________________________________________
// Load Players
// 2 different player loading options for your convenience
// _____________________________________________________________

// First

func GenerateEmptyPlayersByTagsAndUsernames(tags []string, usernames []string, serverUsernames []string,
	isAllSpectators bool) []*Player {
	if len(tags) != len(usernames) {
		log.Println("Unexpected mismatch of playing participants and nicknames")
		return []*Player{}
	}
	var players []*Player
	for i, tag := range tags {
		var newPlayer *Player
		if isAllSpectators {
			newPlayer = NewSpectator(tag, usernames[i], serverUsernames[i])
		} else {
			newPlayer = NewEmptyPlayer(tag, usernames[i], serverUsernames[i])
		}
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

func GetServerNamesByPlayers(players []*Player) []string {
	names := make([]string, len(players))
	for i, player := range players {
		names[i] = player.ServerNick
	}
	return names
}

// Second

func GenerateEmptyPlayersByFunc(
	x any,
	getTagUsernameAndServerUsername func(x any, index int) (string, string, string),
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
		tag, username, serverUsername := getTagUsernameAndServerUsername(x, i)
		var newPlayer *Player
		if isAllSpectators {
			newPlayer = NewSpectator(tag, username, serverUsername)
		} else {
			newPlayer = NewEmptyPlayer(tag, username, serverUsername)
		}
		players[i] = newPlayer
	}

	if isRecovered {
		return []*Player{}
	}
	return players
}
