package converter

import (
	corePlayerPack "github.com/https-whoyan/MafiaBot/core/player"

	"github.com/bwmarrin/discordgo"
)

func getUserUsernameInGuild(s *discordgo.Session, guildID string, memberID string) string {
	member, err := s.GuildMember(guildID, memberID)
	if err != nil || member == nil {
		return ""
	}
	nick := member.Nick
	if nick == "" {
		nick = member.User.Username
	}
	return nick
}

func SafeGetUserUsernameInGuild(s *discordgo.Session, guildID string, user *discordgo.User) string {
	nick := getUserUsernameInGuild(s, guildID, user.ID)
	if nick == "" {
		return user.Username
	}
	return nick
}

func DiscordUsersToEmptyPlayers(s *discordgo.Session, guildID string,
	users []*discordgo.User, isSpectators bool) []*corePlayerPack.Player {
	// First realization
	var (
		tags            []string
		usernames       []string
		serverUsernames []string
	)

	for _, user := range users {
		tags = append(tags, user.ID)
		usernames = append(usernames, SafeGetUserUsernameInGuild(s, guildID, user))
		serverUsernames = append(serverUsernames, user.ID)
	}

	return corePlayerPack.GenerateEmptyPlayersByTagsAndUsernames(tags, usernames, serverUsernames, isSpectators)

	//Second Realization
	/*
		getTagAndUsernameFunc := func(u any, index int) (string, string, string) {
			iUser := u.([]*discordgo.User)[index]
			return iUser.ID, SafeGetUserUsernameInGuild(s, guildID, iUser), iUser.ID
		}
		return corePlayerPack.GenerateEmptyPlayersByFunc(users, getTagAndUsernameFunc, len(users), isSpectators)
	*/
}
