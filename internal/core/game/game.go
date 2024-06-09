package game

import (
	"sync"

	"github.com/https-whoyan/MafiaBot/internal/bot/channel"
	"github.com/https-whoyan/MafiaBot/internal/core/config"
	"github.com/https-whoyan/MafiaBot/internal/core/players"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
)

type State int

const (
	RegisterState   State = 1
	StartingState   State = 2
	NightState      State = 3
	DayState        State = 4
	VotingState     State = 5
	PausedState     State = 6
	FinishState     State = 7
	NonDefinedState State = 8
)

type Game struct {
	sync.RWMutex
	GuildID      string              `json:"guildID"`
	RolesConfig  *config.RolesConfig `json:"rolesConfig"`
	StartPlayers []*players.Player   `json:"startPlayers"`
	Dead         []*players.Player   `json:"dead"`
	Active       []*players.Player   `json:"active"`
	Spectators   []*players.Player   `json:"spectators"`
	// keeps what role is voting right now.
	NightVoting *roles.Role `json:"nightVoting"`
	// presents to the bot which discord chat is used for which role.
	// key: str - role name
	// It is necessary, that would not load mongoDB too much, and that would quickly validate the vote team
	InteractionChannels map[string]*channel.RoleChannel `json:"interactionChannels"`
	ch                  chan int
	PreviousState       State `json:"previousState"`
	State               State `json:"state"`
}

func (g *Game) GetNextState() State {
	switch g.State {
	case RegisterState:
		return StartingState
	case StartingState:
		return NightState
	case NightState:
		return DayState
	case DayState:
		return VotingState
	case VotingState:
		return NightState
	}

	return g.PreviousState
}

func (g *Game) SetState(state State) {
	g.State = state
}

func (g *Game) ChangeStateToPause() {
	g.PreviousState = g.State
	g.State = PausedState
}
