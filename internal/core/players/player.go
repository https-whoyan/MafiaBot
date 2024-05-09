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
	Vote              int `json:"vote"`
	LifeStatus        int
	InteractionStatus int
}
