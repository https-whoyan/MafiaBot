package user

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

func RenameUserInServerByNewID(s *discordgo.Session, guildID string, userGameID int, u *discordgo.User) error {
	prefix := strconv.Itoa(userGameID) + ": "
	newNickname := prefix + u.Username
	err := s.GuildMemberNickname(guildID, u.ID, newNickname)
	return err
}

func RenameUserInServerByNick(s *discordgo.Session, guildID string, nick string, u *discordgo.User) error {
	err := s.GuildMemberNickname(guildID, u.ID, nick)
	return err
}
