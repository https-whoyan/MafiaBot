package roles

import (
	"github.com/https-whoyan/MafiaBot/internal/core/roles/mafia"
	"github.com/https-whoyan/MafiaBot/internal/core/roles/maniac"
	"github.com/https-whoyan/MafiaBot/internal/core/roles/peaceful"
)

type Team int

const (
	PeacefulTeam = iota + 1
	MafiaTeam
	ManiacTeam
)

type Role struct {
	Name           string `json:"name"`
	Team           Team   `json:"team"`
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
	"Maniac":    maniac.Maniac,
	"Citizen":   peaceful.Citizen,
	"Fool":      peaceful.Fool,
}
