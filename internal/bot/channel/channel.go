package channel

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
)

// BotRoleChannel provided a channel, with interaction with role.
// Core RoleChannel implementation.

type BotRoleChannel struct {
	Chat       *discordgo.Channel
	s          *discordgo.Session
	ChannelIID string
	Role       *roles.Role
}

func NewBotRoleChannel(s *discordgo.Session, channelIID string, roleName string) (*BotRoleChannel, error) {
	channel, err := s.Channel(channelIID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting channel: %v", err))
	}
	role, exists := roles.GetRoleByName(roleName)
	if !exists {
		return nil, errors.New(fmt.Sprintf("empty role %v", roleName))
	}

	return &BotRoleChannel{
		s:          s,
		Chat:       channel,
		ChannelIID: channelIID,
		Role:       role,
	}, nil
}

func (ch BotRoleChannel) AddPlayer(serverUserID string) error {
	overridePermission := discordgo.PermissionOverwriteTypeMember
	var allChannelPermission int64 = discordgo.PermissionAllChannel

	err := ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, overridePermission, allChannelPermission, 0)
	return err
}

func (ch BotRoleChannel) AddSpectator(serverUserID string) error {
	overridePermission := discordgo.PermissionOverwriteTypeMember
	var viewPermutation int64 = discordgo.PermissionViewChannel

	err := ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, overridePermission, viewPermutation, 0)
	return err
}

func (ch BotRoleChannel) RemoveUser(serverUserID string) error {
	err := ch.s.ChannelPermissionDelete(ch.ChannelIID, serverUserID)
	return err
}

func (ch BotRoleChannel) GetServerID() string {
	return ch.ChannelIID
}

func (ch BotRoleChannel) GetRole() *roles.Role {
	return ch.Role
}

func (ch BotRoleChannel) Write(b []byte) (n int, err error) {
	return SendMessage(ch.s, ch.ChannelIID, b)
}

// MainChannel provided a main interaction channel with players
// Core MainChannel implementation.

type BotMainChannel struct {
	ChannelIID string             `json:"channel_iid"`
	Chat       *discordgo.Channel `json:"channel"`
	s          *discordgo.Session
}

func NewBotMainChannel(s *discordgo.Session, channelIID string, roleName string) (*BotRoleChannel, error) {
	channel, err := s.Channel(channelIID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting channel: %v", err))
	}
	role, exists := roles.GetRoleByName(roleName)
	if !exists {
		return nil, errors.New(fmt.Sprintf("empty role %v", roleName))
	}

	return &BotRoleChannel{
		s:          s,
		Chat:       channel,
		ChannelIID: channelIID,
		Role:       role,
	}, nil
}

func (ch BotMainChannel) AddPlayer(serverUserID string) error {
	overridePermission := discordgo.PermissionOverwriteTypeMember
	var allChannelPermission int64 = discordgo.PermissionAllChannel

	err := ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, overridePermission, allChannelPermission, 0)
	return err
}

func (ch BotMainChannel) AddSpectator(serverUserID string) error {
	overridePermission := discordgo.PermissionOverwriteTypeMember
	var viewPermutation int64 = discordgo.PermissionViewChannel

	err := ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, overridePermission, viewPermutation, 0)
	return err
}

func (ch BotMainChannel) RemoveUser(serverUserID string) error {
	err := ch.s.ChannelPermissionDelete(ch.ChannelIID, serverUserID)
	return err
}

func (ch BotMainChannel) GetServerID() string {
	return ch.ChannelIID
}

func (ch BotMainChannel) Write(b []byte) (n int, err error) {
	return SendMessage(ch.s, ch.ChannelIID, b)
}

// SendMessage send message provided channelIID channel, using discordgo.Session
func SendMessage(s *discordgo.Session, channelIID string, b []byte) (n int, err error) {
	if s == nil {
		return 0, errors.New("nil session")
	}
	if s.TryLock() {
		defer s.Unlock()
	}
	_, err = s.ChannelMessageSend(channelIID, string(b))
	if err != nil {
		return 0, err
	}
	return len(b), nil
}
