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

// Yes, edit it
const (
	playerPatternWithoutNickname     = "%v"             // ID
	playerPrefixPattern              = "%v: %v"         // ID, Nick
	spectatorPrefixPattern           = "(spectator) %v" // Nick
	deadPrefixPatternWithoutNickname = "(dead) %v"      // ID
	deadPrefixPattern                = "(dead) %v: %v"  // ID, Nick
)

var (
	getNewPlayerNickname                    = func(ID int, oldNick string) string { return fmt.Sprintf(playerPrefixPattern, ID, oldNick) }
	getNewPlayerNicknameWithoutNick         = func(ID int) string { return fmt.Sprintf(playerPatternWithoutNickname, ID) }
	getNewSpectatorNickname                 = func(oldNick string) string { return fmt.Sprintf(spectatorPrefixPattern, oldNick) }
	getNewPlayerDeadNicknameWithoutNickname = func(ID int) string { return fmt.Sprintf(deadPrefixPatternWithoutNickname, ID) }
	getNewPlayerDeadNickname                = func(ID int, oldNick string) string { return fmt.Sprintf(deadPrefixPattern, ID, oldNick) }

	logIsEmptyProvider = func(serverUserID string, oldNick string, newNick string, channelIID string) {
		log.Printf("renameProvider is not provided. User with ServerID %v %v will "+
			"not be renamed to %v in %v channel.",
			serverUserID, oldNick, newNick, channelIID)
	}
)

var (
	InvalidID          = errors.New("invalid ID")
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
		newNick = getNewPlayerNicknameWithoutNick(p.ID)
	} else {
		newNick = getNewPlayerNickname(p.ID, p.OldNick)
	}
	p.Nick = newNick
	if provider == nil {
		logIsEmptyProvider(p.Tag, p.OldNick, newNick, channelIID)
		return nil
	}
	return provider.RenameUser(channelIID, p.Tag, newNick)
}

func (p *Player) RenameToSpectator(provider RenameUserProviderInterface, channelIID string) error {
	if p.LifeStatus != Spectating {
		return UserIsNotSpectator
	}
	var newNick string
	if len(p.OldNick) == 0 {
		newNick = getNewPlayerNicknameWithoutNick(p.ID)
	} else {
		newNick = getNewSpectatorNickname(p.OldNick)
	}
	p.Nick = newNick
	if provider == nil {
		logIsEmptyProvider(p.Tag, p.OldNick, newNick, channelIID)
		return nil
	}
	return provider.RenameUser(channelIID, p.Tag, newNick)
}

func (p *Player) RenameToDeadPlayer(provider RenameUserProviderInterface, channelIID string) error {
	if p.ID <= 0 {
		return InvalidID
	}
	if p.LifeStatus != Dead {
		return UserIsNotDead
	}
	var newNick string
	if len(p.OldNick) == 0 {
		newNick = getNewPlayerDeadNicknameWithoutNickname(p.ID)
	} else {
		newNick = getNewPlayerDeadNickname(p.ID, p.OldNick)
	}
	p.Nick = newNick
	if provider == nil {
		logIsEmptyProvider(p.Tag, p.OldNick, newNick, channelIID)
		return nil
	}

	return provider.RenameUser(channelIID, p.Tag, newNick)
}

func (p *Player) RenameUserAfterGame(provider RenameUserProviderInterface, channelIID string) error {
	newNick := p.OldNick
	oldNick := getNewPlayerNickname(p.ID, p.OldNick)
	p.Nick = newNick
	p.OldNick = newNick
	if provider == nil {
		logIsEmptyProvider(p.Tag, oldNick, newNick, channelIID)
		return nil
	}
	return provider.RenameUser(channelIID, p.Tag, newNick)
}
