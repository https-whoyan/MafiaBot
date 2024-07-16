package player

import (
	"errors"
	"fmt"
	"log"
)

// ___________________________________
// RenameUserProviderInterface
// ___________________________________

// RenameUserProviderInterface
//
// Need to rename people at game startup by adding prefixes to them,
// like “<ID in game>: <old nickname>”
//
// See below
type RenameUserProviderInterface interface {
	RenameUser(channelIID string, userServerID string, newNick string) error
}

// Edit it.
const (
	playerPatternWithoutNickname = "%v"     // ID
	playerPrefixPattern          = "%v: %v" // ID, Nick

	spectatorPrefixPattern   = "(spectator) %v" // Nick
	spectatorWithoutNickname = "(spectator)"

	deadPrefixPatternWithoutNickname = "(dead) %v"     // ID
	deadPrefixPattern                = "(dead) %v: %v" // ID, Nick
)

var (
	getNewPlayerNickname            = func(ID int, oldNick string) string { return fmt.Sprintf(playerPrefixPattern, ID, oldNick) }
	getNewPlayerNicknameWithoutNick = func(ID int) string { return fmt.Sprintf(playerPatternWithoutNickname, ID) }

	getNewSpectatorNickname            = func(oldNick string) string { return fmt.Sprintf(spectatorPrefixPattern, oldNick) }
	getNewSpectatorNicknameWithoutNick = func() string { return fmt.Sprintf(spectatorWithoutNickname) }

	getNewPlayerDeadNicknameWithoutNickname = func(ID int) string { return fmt.Sprintf(deadPrefixPatternWithoutNickname, ID) }
	getNewPlayerDeadNickname                = func(ID int, oldNick string) string { return fmt.Sprintf(deadPrefixPattern, ID, oldNick) }

	logIsEmptyProvider = func(serverUserID string, oldNick string, newNick string, channelIID string) {
		log.Printf("renameProvider is not provided. User with ServerID %v %v will "+
			"not be renamed to %v in %v channel.",
			serverUserID, oldNick, newNick, channelIID)
	}
)

var (
	InvalidID          = errors.New("invalid ")
	UserIsNotActive    = errors.New("user is not a active user")
	UserIsNotDead      = errors.New("user is not a dead user")
	UserIsNotSpectator = errors.New("user is not spectator")
)

// _______________________
// UserRenamingFunctions
// _______________________

func (p *Player) RenameAfterGettingID(provider RenameUserProviderInterface, channelIID string) error {
	if p.ID <= 0 {
		return InvalidID
	}
	if p.LifeStatus != Alive {
		return UserIsNotActive
	}
	var newNick string
	if len(p.OldNick) == 0 {
		newNick = getNewPlayerNicknameWithoutNick(int(p.ID))
	} else {
		newNick = getNewPlayerNickname(int(p.ID), p.OldNick)
	}
	p.Nick = newNick
	if provider == nil {
		logIsEmptyProvider(p.Tag, p.OldNick, newNick, channelIID)
		return nil
	}
	return provider.RenameUser(channelIID, p.Tag, newNick)
}

func (n *NonPlayingPlayer) RenameToSpectator(provider RenameUserProviderInterface, channelIID string) error {
	var newNick string
	if len(n.OldNick) == 0 {
		newNick = getNewSpectatorNicknameWithoutNick()
	} else {
		newNick = getNewSpectatorNickname(n.OldNick)
	}
	n.Nick = newNick
	if provider == nil {
		logIsEmptyProvider(n.Tag, n.OldNick, newNick, channelIID)
		return nil
	}
	return provider.RenameUser(channelIID, n.Tag, newNick)
}

func (p *DeadPlayer) RenameToDeadPlayer(provider RenameUserProviderInterface, channelIID string) error {
	if p.ID <= 0 {
		return InvalidID
	}
	if p.LifeStatus != Dead {
		return UserIsNotDead
	}
	var newNick string
	if len(p.OldNick) == 0 {
		newNick = getNewPlayerDeadNicknameWithoutNickname(int(p.ID))
	} else {
		newNick = getNewPlayerDeadNickname(int(p.ID), p.OldNick)
	}
	p.Nick = newNick
	if provider == nil {
		logIsEmptyProvider(p.Tag, p.OldNick, newNick, channelIID)
		return nil
	}

	return provider.RenameUser(channelIID, p.Tag, newNick)
}

func (n *NonPlayingPlayer) RenameUserAfterGame(provider RenameUserProviderInterface, channelIID string) error {
	newNick := n.OldNick
	if provider == nil {
		logIsEmptyProvider(n.Tag, n.OldNick, newNick, channelIID)
		return nil
	}
	return provider.RenameUser(channelIID, n.Tag, newNick)
}
