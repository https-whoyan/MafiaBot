package game

import (
	"context"
	"errors"
	"sync"
	"time"

	channelPack "github.com/https-whoyan/MafiaBot/core/channel"
	configPack "github.com/https-whoyan/MafiaBot/core/config"
	fmtPack "github.com/https-whoyan/MafiaBot/core/fmt"
	playerPack "github.com/https-whoyan/MafiaBot/core/player"
	rolesPack "github.com/https-whoyan/MafiaBot/core/roles"
	timePack "github.com/https-whoyan/MafiaBot/core/time"
)

// This file describes the structure of the game, as well as the start and end functions of the game.

// ____________________
// Types and constants
// ____________________

type RenameMode int8

const (
	// NotRenameMode used if you not want to rename users in your implementations
	NotRenameMode RenameMode = iota
	// RenameInGuildMode used if you want to rename user everything in your guild
	RenameInGuildMode
	// RenameOnlyInMainChannelMode used if you want to rename user only in MainChannel
	RenameOnlyInMainChannelMode
	// RenameInAllChannelsMode used if you want to rename user in every channel (Roles and Main)
	RenameInAllChannelsMode
)

// ____________________
// Options
// ____________________

type GameOption func(g *Game)

func FMTerOpt(fmtEr fmtPack.FmtInterface) GameOption {
	return func(g *Game) {
		messenger := NewGameMessanger(fmtEr, g)
		g.Messenger = messenger
	}
}
func RenamePrOpt(rP playerPack.RenameUserProviderInterface) GameOption {
	return func(g *Game) { g.renameProvider = rP }
}
func RenameModeOpt(mode RenameMode) GameOption {
	return func(g *Game) { g.renameMode = mode }
}
func VotePingOpt(votePing int) GameOption {
	return func(g *Game) { g.VotePing = votePing }
}
func LoggerOpt(logger Logger) GameOption {
	return func(g *Game) { g.logger = logger }
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
	NightCounter int                     `json:"nightCounter"`
	TimeStart    time.Time               `json:"timeStart"`

	StartPlayers *playerPack.NonPlayingPlayers `json:"startPlayers"`
	Active       *playerPack.Players           `json:"active"`
	Dead         *playerPack.DeadPlayers       `json:"dead"`
	Spectators   *playerPack.NonPlayingPlayers `json:"spectators"`

	// Presents to the application which chat is used for which role.
	// key: str - role name
	RoleChannels map[string]channelPack.RoleChannel
	MainChannel  channelPack.MainChannel

	// Keeps what role is voting (in night) right now.
	NightVoting *rolesPack.Role `json:"nightVoting"`
	VoteChan    chan VoteProviderInterface
	TwoVoteChan chan TwoVoteProviderInterface
	// Can the player choose himself
	VoteForYourself bool `json:"voteForYourself"`
	// VotePing presents a delay number for voting for the same player again.
	//
	// Example: A player has voted for players with IDs 5, 4, 3, and VotePing is 2.
	// So the player will not be able to Vote for players 4 and 3 the next night.
	//
	// Default value: 1.
	//
	// Adjustable by option. Set 0, If you want to keep the mechanic that a player can Vote for the same
	// player every night, put -1 or a very large number if you want all players to have completely different votes.
	VotePing int `json:"votePing"`

	timerDone chan struct{}
	timerStop chan struct{}

	PreviousState State `json:"previousState"`
	State         State `json:"state"`
	Messenger     *Messenger
	// Use to rename user in your interpretation
	renameProvider playerPack.RenameUserProviderInterface
	renameMode     RenameMode
	infoSender     chan<- Signal
	logger         Logger
	ctx            context.Context
}

func GetNewGame(guildID string, opts ...GameOption) *Game {
	active := make(playerPack.Players)
	dead := make(playerPack.DeadPlayers)
	spectators := playerPack.NonPlayingPlayers{}
	newGame := &Game{
		GuildID: guildID,
		State:   NonDefinedState,
		// Chan s create.
		VoteChan:    make(chan VoteProviderInterface),
		TwoVoteChan: make(chan TwoVoteProviderInterface),
		timerDone:   make(chan struct{}),
		timerStop:   make(chan struct{}),
		// Slices.
		Active:     &active,
		Dead:       &dead,
		Spectators: &spectators,
		// Create a map
		RoleChannels: make(map[string]channelPack.RoleChannel),
		VotePing:     1,
		ctx:          context.Background(),
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
		6) And chan s (See GetNewGame)
		7) fmtEr
		8) renameProvider
		9) renameMode

	Let's validate it.
*/

// Init validation Errors.
var (
	EmptyConfigErr                             = errors.New("empty config")
	MismatchPlayersCountAndGamePlayersCountErr = errors.New("mismatch config playersCount and game players")
	NotFullRoleChannelInfoErr                  = errors.New("not full role channel info")
	NotMainChannelInfoErr                      = errors.New("not main channel info")
	EmptyFMTerErr                              = errors.New("empty FMTer")
	EmptyRenameProviderErr                     = errors.New("empty rename provider")
	EmptyRenameModeErr                         = errors.New("empty rename mode")
)

// validationStart is used to validate the game before it is fully initialized.
func (g *Game) validationStart(cfg *configPack.RolesConfig) error {
	g.RLock()
	defer g.RUnlock()

	var err error

	if cfg == nil {
		return EmptyConfigErr
	}

	if cfg.PlayersCount != len(*(g.StartPlayers)) {
		err = errors.Join(err, MismatchPlayersCountAndGamePlayersCountErr)
	}
	if len(g.RoleChannels) != len(rolesPack.GetAllNightInteractionRolesNames()) {
		err = errors.Join(err, NotFullRoleChannelInfoErr)
	}
	if g.MainChannel == nil {
		err = errors.Join(err, NotMainChannelInfoErr)
	}
	if g.Messenger == nil {
		err = errors.Join(err, EmptyFMTerErr)
	}
	if g.renameMode == NotRenameMode {
		return err
	}
	if g.renameProvider == nil {
		err = errors.Join(err, EmptyRenameProviderErr)
	}
	switch g.renameMode {
	case RenameInGuildMode:
		return err
	case RenameOnlyInMainChannelMode:
		return err
	case RenameInAllChannelsMode:
		return err
	default:
		err = errors.Join(err, EmptyRenameModeErr)
	}
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
	// Set fmtEr
	// Set config and players count
	g.SetState(StartingState)
	g.Lock()
	g.RolesConfig = cfg
	g.PlayersCount = cfg.PlayersCount
	g.TimeStart = time.Now()
	g.Unlock()

	// Get Players
	tags := g.StartPlayers.GetTags()
	oldNicknames := g.StartPlayers.GetUsernames()
	serverUsernames := g.StartPlayers.GetServerNicknames()
	players, err := playerPack.GeneratePlayers(tags, oldNicknames, serverUsernames, cfg)
	if err != nil {
		return err
	}
	// And state it to active field
	g.Lock()
	g.Active = &players
	g.Unlock()

	g.RLock()
	defer g.RUnlock()
	// ________________
	// Add to channels
	// ________________

	// We need to add spectators and players to channel.
	// First, add users to role channels.
	for _, player := range *g.Active {
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
	for _, spectator := range *g.Spectators {
		for _, interactionChannel := range g.RoleChannels {
			err = interactionChannel.AddSpectator(spectator.Tag)
			if err != nil {
				return err
			}
		}
	}

	// Then, add all players to main chat.
	for _, player := range *g.StartPlayers {
		err = g.MainChannel.AddPlayer(player.Tag)
		if err != nil {
			return err
		}
	}
	// And spectators.
	for _, spectator := range *g.Spectators {
		err = g.MainChannel.AddSpectator(spectator.Tag)
		if err != nil {
			return err
		}
	}

	// _______________
	// Renaming.
	// _______________
	switch g.renameMode {
	case NotRenameMode: // No actions
	case RenameInGuildMode:
		for _, player := range *g.Active {
			err = player.RenameAfterGettingID(g.renameProvider, "")
			if err != nil {
				return err
			}
		}
		for _, spectator := range *g.Spectators {
			err = spectator.RenameToSpectator(g.renameProvider, "")
			if err != nil {
				return err
			}
		}
	case RenameOnlyInMainChannelMode:
		mainChannelServerID := g.MainChannel.GetServerID()

		for _, player := range *g.Active {
			err = player.RenameAfterGettingID(g.renameProvider, mainChannelServerID)
			if err != nil {
				return err
			}
		}
	case RenameInAllChannelsMode:
		// Add to Role Channels.
		for _, player := range *g.Active {
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

		for _, player := range *g.Active {
			err = player.RenameAfterGettingID(g.renameProvider, mainChannelServerID)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("invalid rename mode")
	}
	if g.logger != nil {
		g.RUnlock()
		g.Lock()
		err = g.logger.InitNewGame(g)
		g.Unlock()
		g.RLock()
		return err
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

Also call deferred finish() (or FinishAnyway(), if game was stopped by context)

It is recommended to use context.Background()

Return receive chan of Signal type, that informed you Signal s
*/
func (g *Game) Run(ctx context.Context) <-chan Signal {
	ch := make(chan Signal)
	g.infoSender = ch

	go func() {
		// Send InteractionMessage About New Game
		err := g.Messenger.Init.SendStartMessage(g.MainChannel)
		// Used for participants to familiarize themselves with their roles, and so on.
		time.Sleep(timePack.RoleInfoCount * time.Second)
		safeSendErrSignal(g.infoSender, err)
		switch {
		case ctx == nil:
			sendFatalSignal(g.infoSender, NilContext)
		case g.IsRunning():
			sendFatalSignal(g.infoSender, ErrGameAlreadyStarted)
		default:
			g.Lock()
			g.ctx = ctx
			g.Unlock()

			var isStoppedByCtx bool

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				isStoppedByCtx = g.run()
			}()
			wg.Wait()

			switch isStoppedByCtx {
			case true:
				g.FinishAnyway()
			case false:
				g.finish()
			}
		}
	}()

	return ch
}

func (g *Game) run() (isStoppedByCtx bool) {
	// FinishState will be set when the winner is already clear.
	// This will be determined after the night and after the day's voting.
	for g.State != FinishState {
		select {
		case <-g.ctx.Done():
			isStoppedByCtx = true
		default:
			g.Lock()
			g.NightCounter++
			g.Unlock()

			// Night

			nightLog := g.Night()
			g.AffectNight(nightLog)
			if g.logger != nil {
				err := g.logger.SaveNightLog(g, nightLog)
				safeSendErrSignal(g.infoSender, err)
			}

			// Validate is final?

			winnerTeam := g.UnderstandWinnerTeam()
			if winnerTeam != nil {
				finishLog := g.NewFinishLog(winnerTeam, false)
				g.FinishByFinishLog(finishLog)
				return
			}

			// Day

			dayLog := g.Day()
			g.AffectDay(dayLog)
			if g.logger != nil {
				err := g.logger.SaveDayLog(g, dayLog)
				safeSendErrSignal(g.infoSender, err)
			}

			// Validate is final?
			winnerTeam = g.UnderstandWinnerTeam()
			fool := (*g.Active.SearchAllPlayersWithRole(rolesPack.Fool))[0]

			if dayLog.Kicked != nil && *dayLog.Kicked == int(fool.ID) {
				finishLog := g.NewFinishLog(winnerTeam, false)
				g.FinishByFinishLog(finishLog)
				return
			} else if winnerTeam != nil {
				finishLog := g.NewFinishLog(winnerTeam, false)
				g.FinishByFinishLog(finishLog)
				return
			}
			g.ClearDayVotes()
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

var (
	finishingFuncOnce = sync.Once{}
	finishOnce        = sync.Once{}
)

func (g *Game) FinishByFinishLog(l FinishLog) {
	finishingFuncOnce.Do(func() {
		err := g.Messenger.Finish.SendMessagesAboutEndOfGame(l, g.MainChannel)
		safeSendErrSignal(g.infoSender, err)
		g.SetState(FinishState)
		g.infoSender <- g.newSwitchStateSignal()
		g.replaceCtx()
		g.finish()
	})
}

func (g *Game) replaceCtx() {
	g.Lock()
	if g.ctx == nil {
		g.ctx = context.Background()
	}
	newCtx, cancel := context.WithCancel(g.ctx)
	g.ctx = newCtx
	g.Unlock()
	cancel()
}

// FinishAnyway is used to end the running game anyway.
func (g *Game) FinishAnyway() {
	finishingFuncOnce.Do(func() {
		if !g.IsRunning() {
			return
		}
		content := "The game was suspended."
		_, err := g.MainChannel.Write([]byte(g.Messenger.Finish.f.Bold(content)))
		safeSendErrSignal(g.infoSender, err)
		g.SetState(FinishState)
		g.infoSender <- g.newSwitchStateSignal()
		g.replaceCtx()
		g.finish()
	})
}

func (g *Game) finish() {

	finishOnce.Do(func() {
		if !g.IsFinished() {
			sendFatalSignal(g.infoSender, errors.New("game is not finished"))
			return
		}

		// Delete from channels
		for _, player := range *g.Active {
			if player.Role.NightVoteOrder == -1 {
				continue
			}

			playerChannel := g.RoleChannels[player.Role.Name]
			safeSendErrSignal(g.infoSender, playerChannel.RemoveUser(player.Tag))
		}

		// Then remove spectators from game
		for _, tag := range playerPack.GetTags(g.Dead, g.Spectators) {
			for _, interactionChannel := range g.RoleChannels {
				safeSendErrSignal(g.infoSender, interactionChannel.RemoveUser(tag))
			}
		}

		// Then, remove all players of main chat.
		for _, player := range *g.StartPlayers {
			safeSendErrSignal(g.infoSender, g.MainChannel.RemoveUser(player.Tag))
		}
		// And spectators.
		for _, spectator := range *g.Spectators {
			safeSendErrSignal(g.infoSender, g.MainChannel.RemoveUser(spectator.Tag))
		}

		// _______________
		// Renaming.
		// _______________
		activePlayersAndSpectators := append(*g.StartPlayers, *g.Spectators...)
		switch g.renameMode {
		case NotRenameMode: // No actions
		case RenameInGuildMode:
			for _, player := range activePlayersAndSpectators {
				safeSendErrSignal(g.infoSender, player.RenameUserAfterGame(g.renameProvider, ""))
			}
		case RenameOnlyInMainChannelMode:
			mainChannelServerID := g.MainChannel.GetServerID()

			for _, player := range activePlayersAndSpectators {
				err := player.RenameUserAfterGame(g.renameProvider, mainChannelServerID)
				safeSendErrSignal(g.infoSender, err)
			}
		case RenameInAllChannelsMode:
			// Rename from Role Channels.
			for _, player := range activePlayersAndSpectators {
				for _, interactionChannel := range g.RoleChannels {
					interactionChannelID := interactionChannel.GetServerID()

					err := player.RenameUserAfterGame(g.renameProvider, interactionChannelID)
					safeSendErrSignal(g.infoSender, err)
				}
			}

			// Rename from main channel
			mainChannelServerID := g.MainChannel.GetServerID()

			for _, player := range activePlayersAndSpectators {
				err := player.RenameUserAfterGame(g.renameProvider, mainChannelServerID)
				safeSendErrSignal(g.infoSender, err)
			}
		default:
			sendFatalSignal(g.infoSender, errors.New("invalid rename mode"))
			return
		}
		sendCloseSignal(g.infoSender, "the game has been successfully completed.")
	})
}
