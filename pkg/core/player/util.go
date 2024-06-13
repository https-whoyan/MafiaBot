package player

import "strconv"

func SearchPlayerByServerID(players []*Player, ID string) *Player {
	for _, player := range players {
		if player.Tag == ID {
			return player
		}
	}

	return nil
}

func SearchPlayerByGameID(players []*Player, ID string) *Player {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return nil
	}

	for _, player := range players {
		if player.ID == intID {
			return player
		}
	}
	return nil
}

func SearchPlayerByID(players []*Player, ID string, isServerID bool) *Player {
	if isServerID {
		return SearchPlayerByServerID(players, ID)
	}
	return SearchPlayerByGameID(players, ID)
}
