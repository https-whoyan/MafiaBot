package game

import (
	"errors"

	configPack "github.com/https-whoyan/MafiaBot/core/config"
	rolesPack "github.com/https-whoyan/MafiaBot/core/roles"
)

// ___________________________
// Game.Start validator
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

// Start validation Errors.
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

func (g *Game) validationStart(cfg *configPack.RolesConfig) error {
	g.RLock()
	defer g.RUnlock()

	var err error
	if cfg == nil {
		return EmptyConfigErr
	}
	if cfg.PlayersCount != len(g.Active) {
		err = errors.Join(err, MismatchPlayersCountAndGamePlayersCountErr)
	}
	if len(g.RoleChannels) != len(rolesPack.GetAllNightInteractionRolesNames()) {
		err = errors.Join(err, NotFullRoleChannelInfoErr)
	}
	if g.MainChannel == nil {
		err = errors.Join(err, NotMainChannelInfoErr)
	}
	if g.Ch == nil {
		err = errors.Join(err, EmptyChanErr)
	}
	if g.fmtEr == nil {
		err = errors.Join(err, EmptyFMTerErr)
	}
	if g.RenameMode == NotRenameModeMode {
		return err
	}
	if g.renameProvider == nil {
		err = errors.Join(err, EmptyRenameProviderErr)
	}
	switch g.RenameMode {
	case RenameInGuildMode:
		return err
	case RenameOnlyInMainChannelMode:
		return err
	case RenameInAllChannelsMode:
		return err
	}
	err = errors.Join(err, EmptyRenameModeErr)
	return err
}
