package channel

import (
	"io"

	"github.com/https-whoyan/MafiaBot/internal/core/roles"
)

type Channel interface {
	// Writer to write message about game.
	io.Writer
	// AddPlayer  , AddSpectator , AddSpectator
	// Optional realization.
	// You can write like this:
	// func (ch ChannelImpl) AddPlayer(serverUserID string) error { return nil }
	AddPlayer(serverUserID string) error
	AddSpectator(serverUserID string) error
	RemoveUser(serverUserID string) error
}

type RoleChannel interface {
	Channel
	GetRole() *roles.Role
}

type MainChannel interface {
	Channel
}
