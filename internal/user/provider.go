package user

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/wrap"
)

// BotUserRenameProvider is
// Core RenameUserProviderInterface realization
type BotUserRenameProvider struct {
	guildID string
	s       *discordgo.Session
}

func NewBotUserRenameProvider(s *discordgo.Session, guildID string) *BotUserRenameProvider {
	return &BotUserRenameProvider{
		guildID: guildID,
		s:       s,
	}
}

func (p BotUserRenameProvider) RenameUser(_ string, userServerID string, newNick string) error {
	if p.s == nil || p.guildID == "" {
		return errors.New("bot User Rename Error, empty fields")
	}
	err := p.s.GuildMemberNickname(p.guildID, userServerID, newNick)
	return wrap.UnwrapDiscordRESTErr(err, discordgo.ErrCodeMissingPermissions)
}
