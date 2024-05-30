package roles

import (
	"github.com/https-whoyan/MafiaBot/internal/core/roles/mafia"
	"github.com/https-whoyan/MafiaBot/internal/core/roles/peaceful"
)

const (
	PeacefulTeam = iota + 1
	MafiaTeam
)

type Role struct {
	Name           string `json:"name"`
	Team           int    `json:"team"`
	NightVoteOrder int    `json:"nightVoteOrder"`
	Description    string `json:"description"`
}

var MappedRoles = map[string]*Role{
	"Peaceful":  peaceful.Peaceful,
	"Mafia":     mafia.Mafia,
	"Doctor":    peaceful.Doctor,
	"Whore":     peaceful.Whore,
	"Detective": peaceful.Detective,
	"Don":       mafia.Don,
}
