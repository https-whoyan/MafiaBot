package game

import (
	"github.com/bwmarrin/discordgo"
)

func GetIDVoiceChannelGame(s *discordgo.Session, memberID string) string {
	return ""
	// TODO!
}

func GetActiveConnectionsMembers(s *discordgo.Session, voiceID string) map[string]bool {
	membersMap := make(map[string]bool)
	guild := s.State.Guilds[0]
	voiceConnections := guild.VoiceStates
	for _, voiceConnection := range voiceConnections {
		membersMap[voiceConnection.UserID] = true
	}
	return membersMap
}
