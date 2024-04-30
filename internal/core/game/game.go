package game

import (
	"github.com/https-whoyan/MafiaBot/internal/core/players"
	"sync"
)

type State int

const (
	RegisterState = iota + 1
	StartingState
	NightState
	DayState
	VotingState
	FinishState
	PausedState
	NonDefinedState
)

type Game struct {
	sync.Mutex
	StartPlayers []*players.Player
	Dead         []*players.Player
	Active       []*players.Player
	Spectators   []*players.Player
	State        State
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
