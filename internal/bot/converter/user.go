package converter

import (
	"github.com/bwmarrin/discordgo"
	corePlayerPack "github.com/https-whoyan/MafiaBot/internal/core/player"
)

func DiscordUsersToEmptyPlayers(users []*discordgo.User, isSpectators bool) []*corePlayerPack.Player {
	// First realization
	var (
		tags      []string
		usernames []string
	)

	for _, user := range users {
		tags = append(tags, user.ID)
		usernames = append(usernames, user.Username)
	}

	return corePlayerPack.GenerateEmptyPlayersByTagsAndUsernames(tags, usernames, isSpectators)

	/*
		Second Realization
		getTagAndUsernameFunc := func(u any, index int) (string, string) {
			iUser := u.([]*discordgo.User)[index]
			return iUser.ID, iUser.Username
		}
		return corePlayerPack.GenerateEmptyPlayersByFunc(users, getTagAndUsernameFunc, len(users), isSpectators)
	*/
}
