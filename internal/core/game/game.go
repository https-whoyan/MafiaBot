package game

import (
	"github.com/https-whoyan/MafiaBot/internal/core/config"
	"github.com/https-whoyan/MafiaBot/internal/core/players"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
	"sync"
)

type State int

const (
	RegisterState = iota + 1
	StartingState
	NightState
	DayState
	VotingState
	PausedState
	FinishState
	NonDefinedState
)

type Game struct {
	sync.Mutex
	RolesConfig  *config.RolesConfig `json:"rolesConfig"`
	StartPlayers []*players.Player   `json:"startPlayers"`
	Dead         []*players.Player   `json:"dead"`
	Active       []*players.Player   `json:"active"`
	Spectators   []*players.Player   `json:"spectators"`
	NightVoting  *roles.Role         `json:"nightVoting"`
	State        State               `json:"state"`
}

func NewGame(playersCount int) *Game {
	return &Game{
		StartPlayers: make([]*players.Player, 0, playersCount),
		Dead:         make([]*players.Player, 0, playersCount),
		Active:       make([]*players.Player, 0, playersCount),
		Spectators:   make([]*players.Player, 0, playersCount),
		State:        RegisterState,
	}
}

func NewUndefinedGame() *Game {
	return &Game{
		State: NonDefinedState,
	}
}
