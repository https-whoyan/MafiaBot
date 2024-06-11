package user

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

// It is
// Core RenameUserProviderInterface realization

type BotUserRenameProvider struct {
	s       *discordgo.Session
	guildID string
}

func NewBotUserRenameProvider(s *discordgo.Session, guildID string) *BotUserRenameProvider {
	return &BotUserRenameProvider{
		s:       s,
		guildID: guildID,
	}
}

func (p BotUserRenameProvider) RenameUser(userServerID string, newNick string) error {
	if p.s == nil || p.guildID == "" {
		return errors.New("bot User Rename Error, empty fields")
	}
	err := p.s.GuildMemberNickname(p.guildID, userServerID, newNick)
	return err
}
