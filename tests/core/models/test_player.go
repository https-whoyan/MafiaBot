package models

import (
	"github.com/https-whoyan/MafiaBot/core/player"
	"strconv"
)

const (
	testTagPrefix            = "TAG"
	testUsernamePrefix       = "USERNAME"
	testServerUsernamePrefix = "SERVER_USERNAME"
)

func GetTestPlayer(i int) *player.Player {
	tag := testTagPrefix + ":" + strconv.Itoa(i)
	username := testUsernamePrefix + ":" + strconv.Itoa(i)
	serverUsername := testServerUsernamePrefix + ":" + strconv.Itoa(i)
	return player.NewEmptyPlayer(tag, username, serverUsername)
}

func GetTestPlayers(cnt int) []*player.Player {
	players := make([]*player.Player, cnt)
	for i := 1; i <= cnt; i++ {
		players[i-1] = GetTestPlayer(i)
	}
	return players
}
