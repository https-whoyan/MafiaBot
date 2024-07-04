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

func GetTestPlayer(i int) *player.NonPlayingPlayer {
	tag := testTagPrefix + ":" + strconv.Itoa(i)
	username := testUsernamePrefix + ":" + strconv.Itoa(i)
	serverUsername := testServerUsernamePrefix + ":" + strconv.Itoa(i)
	return player.NewNonPlayingPlayer(tag, username, serverUsername)
}

func GetTestPlayers(cnt int) *player.NonPlayingPlayers {
	var players = &player.NonPlayingPlayers{}
	for i := 1; i <= cnt; i++ {
		players.Append(GetTestPlayer(i))
	}
	return players
}
