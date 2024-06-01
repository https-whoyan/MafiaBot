package players

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
)

type AliveStatus int

const (
	Alive      AliveStatus = 1
	Dead       AliveStatus = 2
	Spectating AliveStatus = 3
)

type VoteStatus int

const (
	Chooses VoteStatus = 1
	Passed  VoteStatus = 2
	Muted   VoteStatus = 3
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
