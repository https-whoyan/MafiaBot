package game

import (
	"context"
	"errors"
	channelPack "github.com/https-whoyan/MafiaBot/core/channel"
	configPack "github.com/https-whoyan/MafiaBot/core/config"
	fmtPack "github.com/https-whoyan/MafiaBot/core/fmt"
	playerPack "github.com/https-whoyan/MafiaBot/core/player"
	rolesPack "github.com/https-whoyan/MafiaBot/core/roles"
	"slices"
	"sync"
)

// ____________________
// Types and constants
// ____________________

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

type RenameMode int

const (
	// NotRenameModeMode used if you not want to rename users in your implementations
	NotRenameModeMode RenameMode = 0
	// RenameInGuildMode used if you want to rename user everything in your guild
	RenameInGuildMode RenameMode = 1
	// RenameOnlyInMainChannelMode used if you want to rename user only in MainChannel
	RenameOnlyInMainChannelMode RenameMode = 2
	// RenameInAllChannelsMode used if you want to rename user in every channel (Roles and Main)
	RenameInAllChannelsMode RenameMode = 3
)

// ____________________
// Options
// ____________________

type GameOption func(g *Game)

func FMTerOpt(fmtEr fmtPack.FmtInterface) GameOption {
	return func(g *Game) { g.fmtEr = fmtEr }
}
func RenamePrOpt(rP playerPack.RenameUserProviderInterface) GameOption {
	return func(g *Game) { g.renameProvider = rP }
}
func RenameModeOpt(mode RenameMode) GameOption {
	return func(g *Game) { g.RenameMode = mode }
}

// __________________
// Game struct
// __________________

type Game struct {
	sync.RWMutex
	// Presents the server where the game is running.
	// Or GameID.
	// Depends on the implementation.
	//
	// Possibly, may be empty.
	GuildID      string                  `json:"guildID"`
	PlayersCount int                     `json:"playersCount"`
	RolesConfig  *configPack.RolesConfig `json:"rolesConfig"`

	StartPlayers []*playerPack.Player `json:"startPlayers"`
	Dead         []*playerPack.Player `json:"dead"`
	Active       []*playerPack.Player `json:"active"`
	Spectators   []*playerPack.Player `json:"spectators"`

	// keeps what role is voting right now.
	NightVoting *rolesPack.Role `json:"nightVoting"`
	// presents to the application which chat is used for which role.
	// key: str - role name
	RoleChannels  map[string]channelPack.RoleChannel
	MainChannel   channelPack.MainChannel
	VoteChan      chan VoteProviderInterface
	PreviousState State `json:"previousState"`
	State         State `json:"state"`
	// For beautiful messages
	fmtEr fmtPack.FmtInterface
	// Use to rename user in your interpretation
	renameProvider playerPack.RenameUserProviderInterface
	RenameMode     RenameMode `json:"RenameMode"`
	ctx            context.Context
}

func GetNewGame(guildID string, opts ...GameOption) *Game {
	newGame := &Game{
		GuildID: guildID,
		State:   NonDefinedState,
		// Chan s create.
		VoteChan: make(chan VoteProviderInterface),
		// Slices.
		Active:     make([]*playerPack.Player, 0),
		Dead:       make([]*playerPack.Player, 0),
		Spectators: make([]*playerPack.Player, 0),
		// Create a map
		RoleChannels: make(map[string]channelPack.RoleChannel),
	}
	// Set options
	for _, opt := range opts {
		opt(newGame)
	}
	return newGame
}

// _________________
// States functions
// _________________

func (g *Game) GetNextState() State {
	g.RLock()
	defer g.RUnlock()
	switch g.State {
	case NonDefinedState:
		return RegisterState
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
	currGState := g.State
	defer g.Unlock()
	g.PreviousState = currGState
	g.State = state
}

func (g *Game) SwitchState() {
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

var StateDefinition = map[State]string{
	NonDefinedState: "is full raw (nothing is known)",
	RegisterState:   "is waited for registration",
	StartingState:   "is prepared for starting",
	NightState:      "is in night state",
	DayState:        "is in day state",
	VotingState:     "is in day voting state",
	PausedState:     "is in paused state",
	FinishState:     "is finished",
}

func GetStateDefinition(state State) string {
	definition, contains := StateDefinition[state]
	if !contains {
		return "is unknown for server"
	}
	return definition
}

// ______________
// Init
// ______________

// Init
/*
The Init function is used to generate all players, add all players to channels, and rename all players.
It is also used to validate all fields of the game.
This is the penultimate and mandatory function that you must call before starting the game.

Before using it, you must have all options set, all players must have known ServerIDs,
Tags and serverUsernames (all of which must be in StartPlayers), and all channels,
both role-based and non-role-based, must be set.
See the implementation of the validationStart function :
https://github.com/https-whoyan/MafiaBot/blob/main/pkg/core/game/validator.go line 39.

Also see the file loaders.go in the same package https://github.com/https-whoyan/MafiaBot/blob/main/pkg/core/game/loaders.go.


More references:
https://github.com/https-whoyan/MafiaBot/blob/main/pkg/core/player/loader.go line 50

(DISCORD ONLY): https://github.com/https-whoyan/MafiaBot/blob/main/internal/converter/user.go
*/
func (g *Game) Init(cfg *configPack.RolesConfig) (err error) {
	if err = g.validationStart(cfg); err != nil {
		return err
	}
	// Set config and players count
	g.SetState(StartingState)
	g.Lock()
	g.RolesConfig = cfg
	g.PlayersCount = cfg.PlayersCount
	g.Unlock()

	// Get Players
	tags := playerPack.GetTagsByPlayers(g.StartPlayers)
	oldNicknames := playerPack.GetUsernamesByPlayers(g.StartPlayers)
	serverUsernames := playerPack.GetServerNamesByPlayers(g.StartPlayers)
	players, err := playerPack.GeneratePlayers(tags, oldNicknames, serverUsernames, cfg)
	if err != nil {
		return err
	}
	// And state it to active and startPlayers fields
	g.Lock()
	g.StartPlayers = slices.Clone(players)
	g.Active = slices.Clone(players)
	g.Unlock()

	g.RLock()
	defer g.RUnlock()
	// ________________
	// Add to channels
	// ________________

	// We need to add spectators and players to channel.
	// First, add users to role channels.
	for _, player := range g.StartPlayers {
		if player.Role.NightVoteOrder == -1 {
			continue
		}

		playerChannel := g.RoleChannels[player.Role.Name]
		err = playerChannel.AddPlayer(player.Tag)
		if err != nil {
			return
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
	switch g.RenameMode {
	case NotRenameModeMode: // No actions
	case RenameInGuildMode:
		for _, player := range g.StartPlayers {
			err = player.RenameAfterGettingID(g.renameProvider, "")
			if err != nil {
				return err
			}
		}
	case RenameOnlyInMainChannelMode:
		mainChannelServerID := g.MainChannel.GetServerID()

		for _, player := range g.StartPlayers {
			err = player.RenameAfterGettingID(g.renameProvider, mainChannelServerID)
			if err != nil {
				return err
			}
		}
	case RenameInAllChannelsMode:
		// Add to Role Channels.
		for _, player := range g.StartPlayers {
			if player.Role.NightVoteOrder == -1 {
				continue
			}

			playerRoleName := player.Role.Name
			playerInteractionChannel := g.RoleChannels[playerRoleName]
			playerInteractionChannelIID := playerInteractionChannel.GetServerID()
			err = player.RenameAfterGettingID(g.renameProvider, playerInteractionChannelIID)
			if err != nil {
				return err
			}
		}

		// Add to main channel
		mainChannelServerID := g.MainChannel.GetServerID()

		for _, player := range g.StartPlayers {
			err = player.RenameAfterGettingID(g.renameProvider, mainChannelServerID)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("invalid rename mode")
	}
	// Send Message About New Game
	_, _ = g.MainChannel.Write([]byte(g.GetStartMessage()))
	return nil
}

// ********************
// ____________________
// Main Cycle in game.
// ___________________
// ********************

// Run
/*
Is used to start the game.

Runs the run method in its goroutine.
Used after g.Init()

Also call deferred Finish() (or FinishAnyway(), if game was stopped by context)

It is recommended to use context.Background()

Returns an error that occurred after the game finished.
*/

var (
	ErrGameAlreadyStarted = errors.New("game already started")
)

func (g *Game) Run(ctx context.Context) error {
	if g.ctx != nil {
		return ErrGameAlreadyStarted
	}
	g.Lock()
	g.ctx = ctx
	g.Unlock()
	var isStoppedByCtx bool
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		isStoppedByCtx = g.run()
	}()
	wg.Wait()
	defer func() {
		if isStoppedByCtx {
			g.FinishAnyway(ctx)
			err = nil
			return
		}
		err = g.Finish()
	}()
	return err
}

func (g *Game) run() (isStoppedByCtx bool) {
	// FinishState will be set when the winner is already clear.
	// This will be determined after the night and after the day's voting.
	for g.State != FinishState {
		select {
		case <-g.ctx.Done():
			return true
		default:
			// Below is the raw and unfinished code.
			// TODO!
			g.SetState(NightState)
			//g.Night()
			//log := g.getNightLog()
			//winnerTeam, err := log.StateWinner()
			//if err != nil {
			//	g.SetState(FinishState)
			//	return
			//}
		}
	}
	return
}

// _______________________
// Finishing functions
// _______________________

// FinishAnyway is used to end the running game anyway.
//
// Not recommended for use.
// Used in Run function, if the game was interrupted using context.
func (g *Game) FinishAnyway(ctx context.Context) {
	content := "The game was suspended."
	_, _ = g.MainChannel.Write([]byte(g.fmtEr.Bold(content)))
	g.Lock()
	g.State = FinishState
	g.Unlock()
	if g.ctx != nil {
		ctx.Done()
	}
	// Else g.Run

	_ = g.Finish()
}

func (g *Game) Finish() error {
	if g.State != FinishState {
		return errors.New("game is not finished")
	}
	// Delete From channels
}
