package message

import (
	"github.com/bwmarrin/discordgo"
)

func GetUsersByEmojiID(s *discordgo.Session, channelIID, messageID, emojiID string) ([]*discordgo.User, int) {
	users, _ := s.MessageReactions(channelIID, messageID, emojiID, 100, "", "")
	return users, len(users)
}
