package player

import (
	"github.com/https-whoyan/MafiaBot/core/roles"
	"slices"
	"strconv"
)

func SearchPlayerByServerID(players []*Player, ID string) *Player {
	index := slices.IndexFunc(players, func(player *Player) bool {
		return player.Tag == ID
	})
	if index == -1 {
		return nil
	}

	return players[index]
}

func SearchPlayerByGameID(players []*Player, ID string) *Player {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return nil
	}

	index := slices.IndexFunc(players, func(player *Player) bool {
		return player.ID == intID
	})
	if index == -1 {
		return nil
	}

	return players[index]
}

func SearchPlayerByID(players []*Player, ID string, isServerID bool) *Player {
	if isServerID {
		return SearchPlayerByServerID(players, ID)
	}
	return SearchPlayerByGameID(players, ID)
}

func SearchAllPlayersWithRole(players []*Player, role *roles.Role) []*Player {
	var allPlayers []*Player
	for _, player := range players {
		if player.Role == role {
			allPlayers = append(allPlayers, player)
		}
	}

	return allPlayers
}
