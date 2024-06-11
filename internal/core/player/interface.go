package player

import (
	"errors"
	"fmt"
)

// RenameUserProviderInterface need to rename people at game startup by adding prefixes to them,
// like “<ID in game>: <old nickname>”
//
// See playerPrefixPattern and spectatorPrefixPattern (const)
type RenameUserProviderInterface interface {
	RenameUser(userServerID string, newNick string) error
}

const (
	playerPrefixPattern    = "%v: %v"
	spectatorPrefixPattern = "(spectator) %v"
)

var (
	getNewPlayerNickname = func(ID int, oldNick string) string {
		return fmt.Sprintf(playerPrefixPattern, ID, oldNick)
	}

	getNewSpectatorNickname = func(oldNick string) string {
		return fmt.Sprintf(spectatorPrefixPattern, oldNick)
	}
)

var (
	InvalidID = errors.New("invalid ID")
)

// UserRenamingFunction

func (p *Player) RenameAfterGettingID() error {
	if p.ID <= 0 {
		return InvalidID
	}
	newNick := getNewPlayerNickname(p.ID, p.OldNick)
	return p.RenameProvider.RenameUser(p.Tag, newNick)
}

func (p *Player) renameFunctionAfterGame() error {
	return p.RenameProvider.RenameUser(p.Tag, p.OldNick)
}

func (p *Player) RenameToSpectator() error {
	if p.LifeStatus != Spectating {
		return errors.New("user is not a spectator")
	}
	newNick := getNewSpectatorNickname(p.OldNick)
	return p.RenameProvider.RenameUser(p.Tag, newNick)
}
