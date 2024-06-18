package game

import (
	"context"
	"errors"
	"slices"
	"sync"

	channelPack "github.com/https-whoyan/MafiaBot/core/channel"
	configPack "github.com/https-whoyan/MafiaBot/core/config"
	fmtPack "github.com/https-whoyan/MafiaBot/core/fmt"
	playerPack "github.com/https-whoyan/MafiaBot/core/player"
	rolesPack "github.com/https-whoyan/MafiaBot/core/roles"
)

// This file describes the structure of the game, as well as the start and end functions of the game.

// ____________________
// Types and constants
// ____________________

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

	// Keeps what role is voting (in night) right now.
	NightVoting *rolesPack.Role `json:"nightVoting"`
	// Presents to the application which chat is used for which role.
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

// ___________________________
// Game.Init validator
// __________________________
/*
	After RegisterGame I must have all information about
		1) Tags and usernames of players
		2) RoleChannels info
		3) GuildID (Ok, optional)
		4) MainChannel implementation
		5) Spectators
		6) And chan (See GetNewGame)
		7) fmtEr
		8) renameProvider
		9) RenameMode

	Let's validate it.
*/

// Init validation Errors.
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

// validationStart is used to validate the game before it is fully initialized.
func (g *Game) validationStart(cfg *configPack.RolesConfig) error {
	g.RLock()
	defer g.RUnlock()

	joinErr := func(err, addedErr error) {
		err = errors.Join(err, addedErr)
	}

	var err error

	if cfg == nil {
		return EmptyConfigErr
	}

	if cfg.PlayersCount != len(g.StartPlayers) {
		joinErr(err, MismatchPlayersCountAndGamePlayersCountErr)
	}
	if len(g.RoleChannels) != len(rolesPack.GetAllNightInteractionRolesNames()) {
		joinErr(err, NotFullRoleChannelInfoErr)
	}
	if g.MainChannel == nil {
		joinErr(err, NotMainChannelInfoErr)
	}
	if g.VoteChan == nil {
		joinErr(err, EmptyChanErr)
	}
	if g.fmtEr == nil {
		joinErr(err, EmptyFMTerErr)
	}
	if g.RenameMode == NotRenameModeMode {
		return err
	}
	if g.renameProvider == nil {
		joinErr(err, EmptyRenameProviderErr)
	}
	switch g.RenameMode {
	case RenameInGuildMode:
		return err
	case RenameOnlyInMainChannelMode:
		return err
	case RenameInAllChannelsMode:
		return err
	}
	joinErr(err, EmptyRenameModeErr)
	return err
}

// Init
/*
The Init function is used to generate all players, add all players to channels, and rename all players.
It is also used to validate all fields of the game.
This is the penultimate and mandatory function that you must call before starting the game.

Before using it, you must have all options set, all players must have known ServerIDs,
Tags and serverUsernames (all of which must be in StartPlayers), and all channels,
both role-based and non-role-based, must be set.
See the realization of the ValidationStart function (line 139)

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
		for _, spectator := range g.Spectators {
			err = spectator.RenameToSpectator(g.renameProvider, "")
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
	return nil
}

// ********************
// ____________________
// Main Cycle in game.
// ___________________
// ********************
// ********************

var (
	NilContext            = errors.New("nil context")
	ErrGameAlreadyStarted = errors.New("game already started")
)

// Run
/*
Is used to start the game.

Runs the run method in its goroutine.
Used after g.Init()

Also call deferred Finish() (or FinishAnyway(), if game was stopped by context)

It is recommended to use context.Background()

Return receive chan of Signal type
*/
func (g *Game) Run(ctx context.Context) <-chan Signal {
	// Send Message About New Game
	_, _ = g.MainChannel.Write([]byte(g.GetStartMessage()))
	var ch chan Signal

	defer func() {
		defer close(ch)
		switch {
		case ctx == nil:
			ch <- NewFatalSignal(NilContext)
		case g.ctx != nil:
			ch <- NewFatalSignal(ErrGameAlreadyStarted)
		default:
			g.Lock()
			g.ctx = ctx
			g.Unlock()

			var isStoppedByCtx bool

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				isStoppedByCtx = g.run(ch)
			}()
			wg.Wait()

			switch isStoppedByCtx {
			case true:
				g.FinishAnyway(ch)
			case false:
				g.Finish(ch)
			}
		}
	}()

	ch = make(chan Signal)
	return ch
}

func (g *Game) run(ch chan<- Signal) (isStoppedByCtx bool) {
	// FinishState will be set when the winner is already clear.
	// This will be determined after the night and after the day's voting.
	for g.State != FinishState {
		select {
		case <-g.ctx.Done():
			isStoppedByCtx = true
		default:
			g.SetState(NightState)
			g.night(ch)
			//log := g.GetNightLog()
			//log
			//winnerTeam, err := log.StateWinner()
			//if err != nil {
			//	g.SetState(FinishState)
			//	return
			//}
		}
	}
	return
}

// ********************
// ____________________
// Finishing functions
// ___________________
// ********************
// ********************

// FinishAnyway is used to end the running game anyway.
//
// Not recommended for use.
func (g *Game) FinishAnyway(ch chan<- Signal) {
	content := "The game was suspended."
	_, err := g.MainChannel.Write([]byte(g.fmtEr.Bold(content)))
	if err != nil {
		ch <- NewErrSignal(err)
	}
	g.SetState(FinishState)
	g.Lock()
	if g.ctx == nil {
		g.ctx = context.Background()
	}
	newCtx, cancel := context.WithCancel(g.ctx)
	g.ctx = newCtx
	g.Unlock()
	cancel()
	g.Finish(ch)
}
func (g *Game) Finish(ch chan<- Signal) {
	if g.State != FinishState {
		ch <- NewCloseSignal(errors.New("game is not finished"))
		return
	}

	// Delete from channels
	for _, player := range g.StartPlayers {
		if player.Role.NightVoteOrder == -1 {
			continue
		}

		playerChannel := g.RoleChannels[player.Role.Name]
		safeSendErrSignal(ch, playerChannel.RemoveUser(player.Tag))
	}

	// Then remove spectators from game
	for _, spectator := range g.Spectators {
		for _, interactionChannel := range g.RoleChannels {
			safeSendErrSignal(ch, interactionChannel.RemoveUser(spectator.Tag))
		}
	}

	// Then, remove all players of main chat.
	for _, player := range g.Active {
		safeSendErrSignal(ch, g.MainChannel.RemoveUser(player.Tag))
	}
	// And spectators.
	for _, spectator := range g.Spectators {
		safeSendErrSignal(ch, g.MainChannel.RemoveUser(spectator.Tag))
	}

	// _______________
	// Renaming.
	// _______________
	switch g.RenameMode {
	case NotRenameModeMode: // No actions
	case RenameInGuildMode:
		for _, player := range append(g.StartPlayers, g.Spectators...) {
			safeSendErrSignal(ch, player.RenameUserAfterGame(g.renameProvider, ""))
		}
	case RenameOnlyInMainChannelMode:
		mainChannelServerID := g.MainChannel.GetServerID()

		for _, player := range g.StartPlayers {
			err := player.RenameUserAfterGame(g.renameProvider, mainChannelServerID)
			safeSendErrSignal(ch, err)
		}
	case RenameInAllChannelsMode:
		// Rename from Role Channels.
		for _, player := range g.StartPlayers {
			if player.Role.NightVoteOrder == -1 {
				continue
			}

			playerRoleName := player.Role.Name
			playerInteractionChannel := g.RoleChannels[playerRoleName]
			playerInteractionChannelIID := playerInteractionChannel.GetServerID()
			err := player.RenameUserAfterGame(g.renameProvider, playerInteractionChannelIID)
			safeSendErrSignal(ch, err)
		}

		// Rename from main channel
		mainChannelServerID := g.MainChannel.GetServerID()

		for _, player := range g.StartPlayers {
			err := player.RenameUserAfterGame(g.renameProvider, mainChannelServerID)
			safeSendErrSignal(ch, err)
		}
	default:
		ch <- NewFatalSignal(errors.New("invalid rename mode"))
	}
}
