package player

import (
	"github.com/https-whoyan/MafiaBot/core/roles"
)

// _______________________
// Types and constants
// _______________________

type AliveStatus int

const (
	Alive AliveStatus = iota
	Dead
	Spectating
)

type VoteStatus int

const (
	Passed VoteStatus = iota
	Muted
)

type IDType int

// _______________
// Structs
// _______________

// NonPlayingPlayer
// Used for peoples, are used for people who do not have their ID in the game, i.e. are not participating in it.
type NonPlayingPlayer struct {
	// Tag Represent account IDType on the presentation platform
	Tag string `json:"tag"`
	// Represent Server nick in your implementation
	// Using ONLY for Mentions.
	//
	// In my case, it is a Tag/ServerID of player
	ServerNick string `json:"server_nick"`
	// OldNick before renaming.
	// My implementation of the game assumes bot will change player nicknames to
	// their IDs for easier recognition (1, 2, 3...)
	OldNick string `json:"oldNick"`
	// Nick after renaming.
	Nick string `json:"nick"`
}

// For interfacing all structs

func (n NonPlayingPlayer) GetTag() string        { return n.Tag }
func (n NonPlayingPlayer) GetServerNick() string { return n.ServerNick }
func (n NonPlayingPlayer) GetOldNick() string    { return n.OldNick }
func (n NonPlayingPlayer) GetNick() string       { return n.Nick }

type Player struct {
	NonPlayingPlayer
	ID   IDType      `json:"id"`
	Role *roles.Role `json:"role"`
	// Votes stores all the night votes the player casts.
	//
	// NOTE - if the number of players for a role is >2, there will be added voices (is player not be muted)
	// for these voices will be identical.
	//
	// NOTE - For the detective, or roles, who used two votes system, this array will be empty after each night.
	// During the night 2 voices will be recorded here, but after the night they will be deleted immediately.
	Votes []int `json:"votes"`
	// DayVote stores the player's vote in the day vote.
	DayVote    int         `json:"dayVote"`
	LifeStatus AliveStatus `json:"lifeStatus"`
	// InteractionStatus What a player should be doing right now
	InteractionStatus VoteStatus `json:"interactionStatus"`
}

// DeadPlayer superstructure on top of the player.
// Shows more additional fields.
type DeadPlayer struct {
	Player
	DeadReason DeadReason `json:"deadReason"`
	LivedDays  int        `json:"livedDays"`
}

type DeadReason string

const (
	KilledAtNight     DeadReason = "KilledAtNight"
	KilledByDayVoting DeadReason = "KilledByDayVoting"
)

func NewDeadPlayer(p *Player, reason DeadReason, dayLived int) *DeadPlayer {
	return &DeadPlayer{
		Player:     *p,
		DeadReason: reason,
		LivedDays:  dayLived,
	}
}

// ________________________________________________
// Functions to get new players (or Spectating)
// ________________________________________________

func NewNonPlayingPlayer(tag string, username string, serverUsername string) *NonPlayingPlayer {
	return &NonPlayingPlayer{
		Tag:        tag,
		Nick:       username,
		OldNick:    username,
		ServerNick: serverUsername,
	}
}

func NewPlayer(id IDType, tag string, username string, serverUsername string, role *roles.Role) *Player {
	return &Player{
		NonPlayingPlayer: NonPlayingPlayer{
			Tag:        tag,
			ServerNick: serverUsername,
			OldNick:    username,
			Nick:       username,
		},
		ID:                id,
		Role:              role,
		LifeStatus:        Alive,
		InteractionStatus: Passed,
	}
}
