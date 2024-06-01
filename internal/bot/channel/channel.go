package channel

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
)

type RoleChannel struct {
	Chat       *discordgo.Channel `json:"name"`
	ChannelIID string             `json:"channel_iid"`
	Role       *roles.Role        `json:"role"`
}

func LoadRoleChannel(s *discordgo.Session, channelIID string, roleName string) (*RoleChannel, error) {
	channel, err := s.Channel(channelIID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting channel: %v", err))
	}
	role, exists := roles.GetRoleByName(roleName)
	if !exists {
		return nil, errors.New(fmt.Sprintf("empty role %v", roleName))
	}

	return &RoleChannel{
		Chat:       channel,
		ChannelIID: channelIID,
		Role:       role,
	}, nil
}

func (ch *RoleChannel) AddPlayerInChannel(s *discordgo.Session, u *discordgo.User) error {
	overridePermission := discordgo.PermissionOverwriteTypeMember
	var allChannelPermission int64 = discordgo.PermissionAllChannel

	err := s.ChannelPermissionSet(ch.ChannelIID, u.ID, overridePermission, allChannelPermission, 0)
	return err
}

func (ch *RoleChannel) AddSpectatorInChannel(s *discordgo.Session, u *discordgo.User) error {
	overridePermission := discordgo.PermissionOverwriteTypeMember
	var viewPermutation int64 = discordgo.PermissionViewChannel

	err := s.ChannelPermissionSet(ch.ChannelIID, u.ID, overridePermission, viewPermutation, 0)
	return err
}

func (ch *RoleChannel) RemoveUserOfChannel(s *discordgo.Session, u *discordgo.User) error {
	err := s.ChannelPermissionDelete(ch.ChannelIID, u.ID)
	return err
}
