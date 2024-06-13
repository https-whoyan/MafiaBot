package channel

import (
	"io"

	"github.com/https-whoyan/MafiaBot/core/roles"
)

type Channel interface {
	// Writer to write message about game.
	io.Writer

	/*
		AddPlayer , AddSpectator , RemoveUser
		All function below used for add and renaming players in channels in your implementation.
		See below.
	*/

	// AddPlayer
	/*
		Optional realization. Use to add user in interactionChannel in your application
		If you don't need to user management, you can write like this:
			func (ch ChannelImpl) AddPlayer(serverUserID string) error { return nil }
	*/
	AddPlayer(serverUserID string) error
	// AddSpectator
	/*
		Optional realization. Use to add spectator (only can view) in interactionChannel in your application
		If you don't need to user management, you can write like this:
			func (ch ChannelImpl) AddSpectator(serverUserID string) error { return nil }
	*/
	AddSpectator(serverUserID string) error
	// RemoveUser
	/*
		Optional realization. Use to remove user from interactionChannel in your application
		If you don't need to user management, you can write like this:
			func (ch ChannelImpl) AddSpectator(serverUserID string) error { return nil }
	*/
	RemoveUser(serverUserID string) error
	// GetServerID
	/*
		Optional realization.
		Need, if you want to auto user nicknames management into interaction channels.

		See RenameUserProviderInterface

		If you don't need to user nicknames management, you can write like this
			func (ch ChannelImpl) GetServerIID() { return "" }
	*/
	GetServerID() string
}

type RoleChannel interface {
	Channel
	// GetRole For validation.
	GetRole() *roles.Role
}

type MainChannel interface {
	Channel
}

// FromUserToSpectator Switch User in channel to spectator
func FromUserToSpectator(channel Channel, serverUserID string) error {
	err := channel.RemoveUser(serverUserID)
	if err != nil {
		return err
	}
	return channel.AddSpectator(serverUserID)
}

// FromSpectatorToUser Switch Spectator in channel to user
func FromSpectatorToUser(channel Channel, serverUserID string) error {
	err := channel.RemoveUser(serverUserID)
	if err != nil {
		return err
	}
	return channel.AddPlayer(serverUserID)
}
