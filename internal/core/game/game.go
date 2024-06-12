package game

import (
	"errors"
	"sync"

	channelPack "github.com/https-whoyan/MafiaBot/internal/core/channel"
	configPack "github.com/https-whoyan/MafiaBot/internal/core/config"
	fmtPack "github.com/https-whoyan/MafiaBot/internal/core/message/fmt"
	playerPack "github.com/https-whoyan/MafiaBot/internal/core/player"
	rolesPack "github.com/https-whoyan/MafiaBot/internal/core/roles"
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
// Setting
// ____________________

type Setting struct {
	FMTer          fmtPack.FmtInterface
	RenameProvider playerPack.RenameUserProviderInterface
	RenameMode     RenameMode
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
	ch            chan int
	PreviousState State `json:"previousState"`
	State         State `json:"state"`
	// For beautiful messages
	fmtEr fmtPack.FmtInterface
	// Use to rename user in your interpretation
	RenameProvider playerPack.RenameUserProviderInterface
	RenameMode     RenameMode `json:"renameMode"`
}

func GetNewGame(guildID string, cfg Setting) *Game {
	return &Game{
		GuildID: guildID,
		State:   NonDefinedState,
		// Chan create.
		ch: make(chan int),
		// Slices.
		Active:     make([]*playerPack.Player, 0),
		Dead:       make([]*playerPack.Player, 0),
		Spectators: make([]*playerPack.Player, 0),
		// Set interfaces
		fmtEr:          cfg.FMTer,
		RenameProvider: cfg.RenameProvider,
		RenameMode:     cfg.RenameMode,
	}
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

// ______________
// Start
// ______________

/*
	After RegisterGame I must have all information about
		1) Tags and usernames of players
		2) RoleChannels info
		3) GuildID (Ok, optional)
		4) MainChannel implementation
		5) Spectators
		6) And chan (See GetNewGame)
		7) fmtEr
		8) RenameProvider
		9) RenameMode

	Let's validate it.
*/

// Validation Errors.
var (
	EmptyConfigErr                             = errors.New("empty config")
	MismatchPlayersCountAndGamePlayersCountErr = errors.New("mismatch config playersCount and game players")
	NotFullRoleChannelInfoErr                  = errors.New("not full role channel info")
	NotMainChannelInfoErr                      = errors.New("not main channel info")
	EmptyChanErr                               = errors.New("empty channel")
	EmptyFMTerErr                              = errors.New("empty FMTer")
	EmptyRenameProviderErr                     = errors.New("empty rename provider")
	EmptyRenameModeErr                         = errors.New("empty rename mode")
)

func (g *Game) validation(cfg *configPack.RolesConfig) error {
	g.RLock()
	defer g.RUnlock()

	var err error
	if cfg == nil {
		return EmptyConfigErr
	}
	if cfg.PlayersCount != len(g.Active) {
		err = errors.Join(err, MismatchPlayersCountAndGamePlayersCountErr)
	}
	if len(g.RoleChannels) != len(rolesPack.GetAllNightInteractionRolesNames()) {
		err = errors.Join(err, NotFullRoleChannelInfoErr)
	}
	if g.MainChannel == nil {
		err = errors.Join(err, NotMainChannelInfoErr)
	}
	if g.ch == nil {
		err = errors.Join(err, EmptyChanErr)
	}
	if g.fmtEr == nil {
		err = errors.Join(err, EmptyFMTerErr)
	}
	if g.RenameMode == NotRenameModeMode {
		return err
	}
	if g.RenameProvider == nil {
		err = errors.Join(err, EmptyRenameProviderErr)
	}
	switch g.RenameMode {
	case RenameInGuildMode:
		return err
	case RenameOnlyInMainChannelMode:
		return err
	case RenameInAllChannelsMode:
		return err
	}
	err = errors.Join(err, EmptyRenameModeErr)
	return err
}

func (g *Game) Start(cfg *configPack.RolesConfig) error {
	if err := g.validation(cfg); err != nil {
		return err
	}
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

	switch g.RenameMode {
	case NotRenameModeMode: // No actions
	case RenameInGuildMode:
		for _, player := range g.StartPlayers {
			err = player.RenameAfterGettingID(g.RenameProvider, "")
			if err != nil {
				return err
			}
		}
	case RenameOnlyInMainChannelMode:
		mainChannelServerID := g.MainChannel.GetServerID()

		for _, player := range g.StartPlayers {
			err = player.RenameAfterGettingID(g.RenameProvider, mainChannelServerID)
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
			err = player.RenameAfterGettingID(g.RenameProvider, playerInteractionChannelIID)
			if err != nil {
				return err
			}
		}

		// Add to main channel
		mainChannelServerID := g.MainChannel.GetServerID()

		for _, player := range g.StartPlayers {
			err = player.RenameAfterGettingID(g.RenameProvider, mainChannelServerID)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("invalid rename mode")
	}
	return nil
}
