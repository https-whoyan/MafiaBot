package game

import (
	"github.com/https-whoyan/MafiaBot/internal/bot/channel"
	"github.com/https-whoyan/MafiaBot/internal/core/config"
	"github.com/https-whoyan/MafiaBot/internal/core/players"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
	"sync"
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
	sync.Mutex
	RolesConfig  *config.RolesConfig `json:"rolesConfig"`
	StartPlayers []*players.Player   `json:"startPlayers"`
	Dead         []*players.Player   `json:"dead"`
	Active       []*players.Player   `json:"active"`
	Spectators   []*players.Player   `json:"spectators"`
	// keeps what role is voting right now.
	NightVoting         *roles.Role                     `json:"nightVoting"`
	InteractionChannels map[string]*channel.RoleChannel `json:"interactionChannels"`
	// presents to the bot which discord chat is used for which role.
	// key: str - role name
	// It is necessary, that would not load mongoDB too much, and that would quickly validate the vote team
	State State `json:"state"`
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

func (g *Game) SetNonDefinedState() {
	g.State = NonDefinedState
}
