package player

import (
	"github.com/https-whoyan/MafiaBot/core/roles"
)

// _______________________
// Types and constants
// _______________________

type AliveStatus int

const (
	Alive      AliveStatus = 1
	Dead       AliveStatus = 2
	Spectating AliveStatus = 3
)

type VoteStatus int

const (
	Chooses VoteStatus = 1
	Passed  VoteStatus = 2
	Muted   VoteStatus = 3
)

// _______________
// Player Struct
// _______________

type Player struct {
	ID int `json:"id"`
	// Tag Represent account ID on the presentation platform
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
	Nick string      `json:"nick"`
	Role *roles.Role `json:"role"`
	// Vote which player ID the player is voting for
	Vote       int         `json:"vote"`
	LifeStatus AliveStatus `json:"lifeStatus"`
	// InteractionStatus What a player should be doing right now
	InteractionStatus VoteStatus `json:"interactionStatus"`
}

// ________________________________________________
// Functions to get new players (or Spectating)
// ________________________________________________

func NewEmptyPlayer(tag string, username string, serverUsername string) *Player {
	return &Player{
		Tag:               tag,
		Nick:              username,
		OldNick:           username,
		ServerNick:        serverUsername,
		LifeStatus:        Alive,
		InteractionStatus: Passed,
	}
}

func NewPlayer(id int, tag string, username string, serverUsername string, role *roles.Role) *Player {
	return &Player{
		ID:                id,
		OldNick:           username,
		Nick:              username,
		Tag:               tag,
		ServerNick:        serverUsername,
		Role:              role,
		LifeStatus:        Alive,
		InteractionStatus: Passed,
	}
}

func NewSpectator(tag string, username string, serverUsername string) *Player {
	return &Player{
		Tag:               tag,
		OldNick:           username,
		Nick:              username,
		ServerNick:        serverUsername,
		LifeStatus:        Spectating,
		InteractionStatus: Muted,
		Vote:              -1,
	}
}
