package game

import (
	"errors"
	"fmt"
	"sync"

	channelPack "github.com/https-whoyan/MafiaBot/internal/core/channel"
	configPack "github.com/https-whoyan/MafiaBot/internal/core/config"
	myFmt "github.com/https-whoyan/MafiaBot/internal/core/message/fmt"
	playerPack "github.com/https-whoyan/MafiaBot/internal/core/player"
	rolesPack "github.com/https-whoyan/MafiaBot/internal/core/roles"
)

type State int

const (
	NonDefinedState State = 1
	RegisterState   State = 2
	StartingState   State = 3
	NightState      State = 4
	DayState        State = 5
	VotingState     State = 6
	PausedState     State = 7
	FinishState     State = 8
)

type Game struct {
	sync.RWMutex
	// Presents the server where the game is running.
	// Or GameID.
	// Depends on the implementation.
	GuildID      string                  `json:"guildID"`
	PlayersCount int                     `json:"playersCount"`
	RolesConfig  *configPack.RolesConfig `json:"rolesConfig"`

	StartPlayers []*playerPack.Player `json:"startPlayers"`
	Dead         []*playerPack.Player `json:"dead"`
	Active       []*playerPack.Player `json:"active"`
	Spectators   []*playerPack.Player `json:"spectators"`

	// keeps what role is voting right now.
	NightVoting *rolesPack.Role `json:"nightVoting"`
	// presents to the bot which chat is used for which role.
	// key: str - role name
	RoleChannels  map[string]channelPack.RoleChannel
	MainChannel   channelPack.MainChannel
	ch            chan int
	PreviousState State `json:"previousState"`
	State         State `json:"state"`
	// For beautiful messages
	fmtEr myFmt.FmtInterface
}

// ___________________
// For Register Game
// ___________________

func GetNewGame(guildID string, fmtEr myFmt.FmtInterface) *Game {
	return &Game{
		GuildID: guildID,
		State:   NonDefinedState,
		// Chan create.
		ch: make(chan int),
		// And slices.
		Active:     make([]*playerPack.Player, 0),
		Dead:       make([]*playerPack.Player, 0),
		Spectators: make([]*playerPack.Player, 0),
		fmtEr:      fmtEr,
	}
}

// Add Channels

func (g *Game) tryAddNewRoleChannel(ch channelPack.RoleChannel) error {
	addedRole := ch.GetRole()
	if addedRole == nil {
		return errors.New("no role in channel")
	}

	roleName := addedRole.Name
	_, alreadyContained := g.RoleChannels[roleName]
	if alreadyContained {
		return errors.New(fmt.Sprintf("roleChannel %v already exists", roleName))
	}
	return nil
}

func (g *Game) SetNewRoleChannel(ch channelPack.RoleChannel) error {
	if err := g.tryAddNewRoleChannel(ch); err != nil {
		return err
	}
	roleName := ch.GetRole().Name
	g.Lock()
	g.RoleChannels[roleName] = ch
	g.Unlock()
	return nil
}

func (g *Game) SetRoleChannels(chs []channelPack.RoleChannel) (err error) {
	for _, ch := range chs {
		if err = g.tryAddNewRoleChannel(ch); err != nil {
			return
		}
	}

	for _, ch := range chs {
		err = g.SetNewRoleChannel(ch)
	}
	return err
}

// _________________
// States functions
// _________________

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
	g.Lock()
	defer g.Unlock()
	currGState := g.State
	g.PreviousState = currGState
	g.State = state
}

func (g *Game) SwitchState() {
	g.Lock()
	defer g.Unlock()
	nextState := g.GetNextState()
	g.SetState(nextState)
}

func (g *Game) ChangeStateToPause() {
	g.Lock()
	defer g.Unlock()
	currGState := g.State
	g.PreviousState = currGState
	g.State = PausedState
}

// ___________________
// Start
// ___________________

// After RegisterGame I have all information about
// (All Valid)
// 1) Tags and usernames of players
// 2) RoleChannels info
// 3) GuildID
// 4) MainChannel implementation
// 5) Spectators
// 6) And chan (See GetNewGame)

func (g *Game) Start(cfg *configPack.RolesConfig) error {
	g.Lock()
	defer g.Unlock()
	// Set state, config and players count
	g.SwitchState()
	g.RolesConfig = cfg
	g.PlayersCount = cfg.PlayersCount

	// Get Players
	tags := playerPack.GetTagsByPlayers(g.StartPlayers)
	oldNicknames := playerPack.GetUsernamesByPlayers(g.StartPlayers)
	players, err := playerPack.GeneratePlayers(tags, oldNicknames, cfg)
	if err != nil {
		return err
	}
	// And state it to active and startPlayers fields
	g.StartPlayers = players
	g.Active = players

	// ________________
	// Add to channels
	// ________________

	// We need to add spectators and players to channel.
	// First, add users to hit channels.
	for _, player := range g.StartPlayers {
		if player.Role.NightVoteOrder == -1 {
			continue
		}

		// Use mafia interaction channel
		if player.Role.Name == "Don" {
			mafiaChannel := g.RoleChannels["mafia"]
			err = mafiaChannel.AddPlayer(player.Tag)
			if err != nil {
				return err
			}
			continue
		}
		playerChannel := g.RoleChannels[player.Role.Name]
		err = playerChannel.AddPlayer(player.Tag)
		if err != nil {
			return err
		}
	}

	// Then add spectators to game
	for _, spectator := range g.Spectators {
		for _, interactionChannel := range g.RoleChannels {
			err = interactionChannel.AddSpectator(spectator.Tag)
			if err != nil {
				return err
			}
		}
	}

	// Then, add all players to main chat.
	for _, player := range g.StartPlayers {
		err = g.MainChannel.AddPlayer(player.Tag)
		if err != nil {
			return err
		}
	}
	// And spectators.
	for _, spectator := range g.Spectators {
		err = g.MainChannel.AddSpectator(spectator.Tag)
		if err != nil {
			return err
		}
	}

	// _______________
	// Renaming.
	// _______________

	for _, player := range g.StartPlayers {
		err = player.RenameAfterGettingID()
		if err != nil {
			return err
		}
	}
	return nil
}
