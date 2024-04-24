package players

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/pkg/roles"
)

const (
	alive = iota + 1
	dead
	spectating
)

const (
	chooses = iota + 1
	passed
	muted
)

type Player struct {
	Num               int
	Tag               *discordgo.User
	Role              *roles.Role
	lifeStatus        int
	interactionStatus int
}

// test
func getPlayer() *Player {
	return &Player{
		Num:        1,
		Tag:        nil,
		Role:       nil,
		lifeStatus: alive,
	}
}
