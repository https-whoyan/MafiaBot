package channel

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
)

type RoleChannel struct {
	Chat *discordgo.Channel `json:"name"`
	Role *roles.Role        `json:"role"`
}


func LoadAllRolesChat(s *discordgo.Session, guildID string)  {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return //TODO!
	}
	for
}
