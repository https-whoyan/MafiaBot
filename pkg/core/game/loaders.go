package game

import (
	"errors"
	"fmt"

	channelPack "github.com/https-whoyan/MafiaBot/core/channel"
	playerPack "github.com/https-whoyan/MafiaBot/core/player"
)

// _________________
// Channels
// _________________

func (g *Game) tryAddNewRoleChannel(ch channelPack.RoleChannel) error {
	addedRole := ch.GetRole()
	if addedRole == nil {
		return errors.New("no role in channel")
	}

	roleName := addedRole.Name
	_, alreadyContained := g.RoleChannels[roleName]
	if alreadyContained {
		return errors.New(fmt.Sprintf("roleChannel %v already exists", roleName))
	}
	return nil
}

func (g *Game) SetNewRoleChannel(ch channelPack.RoleChannel) error {
	if err := g.tryAddNewRoleChannel(ch); err != nil {
		return err
	}
	roleName := ch.GetRole().Name
	g.Lock()
	g.RoleChannels[roleName] = ch
	g.Unlock()
	return nil
}

func (g *Game) SetRoleChannels(chs ...channelPack.RoleChannel) (err error) {
	for _, ch := range chs {
		if err = g.tryAddNewRoleChannel(ch); err != nil {
			return
		}
	}

	for _, ch := range chs {
		err = g.SetNewRoleChannel(ch)
	}
	return err
}

func (g *Game) SetRoleChannelByMap(chsMp map[any]channelPack.RoleChannel) (err error) {
	for _, ch := range chsMp {
		if err = g.tryAddNewRoleChannel(ch); err != nil {
			return
		}
	}

	for _, ch := range chsMp {
		err = g.SetNewRoleChannel(ch)
	}
	return err
}

func (g *Game) SetMainChannel(ch channelPack.MainChannel) error {
	g.Lock()
	defer g.Unlock()
	if ch == nil {
		return errors.New("no main channel")
	}
	g.MainChannel = ch
	return nil
}

// ___________________
// Players
// ___________________

func (g *Game) SetStartPlayers(players *playerPack.NonPlayingPlayers) {
	g.Lock()
	defer g.Unlock()
	g.StartPlayers = players
}

func (g *Game) SetSpectators(spectators *playerPack.NonPlayingPlayers) {
	g.Lock()
	defer g.Unlock()
	g.Spectators = spectators
}
