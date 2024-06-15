package user

import "github.com/bwmarrin/discordgo"

func GetUsersOnlyIncludeInTags(users []*discordgo.User, tags []string) ([]*discordgo.User, int) {
	mpTags := make(map[string]bool)
	for _, tag := range tags {
		mpTags[tag] = true
	}
	ans := make([]*discordgo.User, 0)
	for _, user := range users {
		if mpTags[user.ID] {
			ans = append(ans, user)
		}
	}
	return ans, len(ans)
}

func GetUsersNotInclude(users []*discordgo.User, needNotInclude []*discordgo.User) ([]*discordgo.User, int) {
	mpIds := make(map[string]bool)
	for _, user := range needNotInclude {
		mpIds[user.ID] = true
	}
	ans := make([]*discordgo.User, 0)
	for _, user := range users {
		if !mpIds[user.ID] {
			ans = append(ans, user)
		}
	}
	return ans, len(ans)
}
