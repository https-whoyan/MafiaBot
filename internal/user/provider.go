package user

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
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
	if err == nil {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("Error renaming user to ", newNick))
		}

	}()

	var discordGoErr *discordgo.RESTError
	if errors.As(err, &discordGoErr) {
		if discordGoErr.Message.Code == discordgo.ErrCodeMissingPermissions {
			log.Println("User Rename Error, permission denied, but, ok.")
			return nil
		}
	}
	return err
}
