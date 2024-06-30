package player

import (
	"github.com/https-whoyan/MafiaBot/core/roles"
	"slices"
	"sort"
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

// DeadPlayersToPlayers convert slice of DeadPlayer's To Player's
func DeadPlayersToPlayers(deadPlayers []*DeadPlayer) []*Player {
	var players []*Player
	for _, deadPlayer := range deadPlayers {
		players = append(players, deadPlayer.P)
	}
	return players
}

// SortPlayersByTeamAndDead used for messaging after game.
func SortPlayersByTeamAndDead(players []*Player) []*Player {
	sort.Slice(players, func(i, j int) bool {
		playerI, playerJ := players[i], players[j]
		if playerI.Role.Team != playerJ.Role.Team {
			return playerI.Role.Team < playerJ.Role.Team
		}

		if playerI.LifeStatus != playerJ.LifeStatus {
			if playerI.LifeStatus == Alive {
				return true
			}
			return false
		}

		return playerI.ID < playerJ.ID
	})

	return players
}
