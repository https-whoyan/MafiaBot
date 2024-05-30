package players

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
)

type AliveStatus int

const (
	Alive = iota + 1
	Dead
	Spectating
)

type VoteStatus int

const (
	Chooses = iota + 1
	Passed
	Muted
)

type Player struct {
	ID                int             `json:"ID"`
	OldNick           string          `json:"oldNick"`
	Tag               *discordgo.User `json:"tag"`
	Role              *roles.Role     `json:"role"`
	Vote              int             `json:"vote"`
	LifeStatus        AliveStatus     `json:"lifeStatus"`
	InteractionStatus VoteStatus      `json:"interactionStatus"`
}
