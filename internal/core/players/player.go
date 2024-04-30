package players

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
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
	ID                int
	OldNick           string
	Tag               *discordgo.User
	Role              *roles.Role
	LifeStatus        int
	InteractionStatus int
}

// test
func getPlayer() *Player {
	return &Player{
		ID:         1,
		Tag:        nil,
		Role:       nil,
		LifeStatus: alive,
	}
}
