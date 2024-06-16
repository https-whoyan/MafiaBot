package channel

import (
	"errors"
	"fmt"

	"github.com/https-whoyan/MafiaBot/core/roles"
	"github.com/https-whoyan/MafiaBot/internal/wrap"

	"github.com/bwmarrin/discordgo"
)

func NewOverridePermission(allow int64, isNew bool, additionDeny ...int64) *discordgo.PermissionOverwrite {
	var deny int64
	if isNew {
		deny = discordgo.PermissionCreateInstantInvite | discordgo.PermissionManageChannels |
			discordgo.PermissionKickMembers
	}
	for _, v := range additionDeny {
		deny = deny | v
	}
	return &discordgo.PermissionOverwrite{
		Type:  discordgo.PermissionOverwriteTypeMember,
		Allow: allow,
		Deny:  deny,
	}
}

// BotRoleChannel provided a channel, with interaction with role.
// Core RoleChannel implementation.

type BotRoleChannel struct {
	Chat       *discordgo.Channel
	s          *discordgo.Session
	ChannelIID string
	Role       *roles.Role
	// mappedPermissions stores past permissions for each participant in the game.
	mappedPermissions map[string]int64
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
		s:                 s,
		Chat:              channel,
		ChannelIID:        channelIID,
		Role:              role,
		mappedPermissions: make(map[string]int64),
	}, nil
}

func (ch *BotRoleChannel) GetServerID() string  { return ch.ChannelIID }
func (ch *BotRoleChannel) GetRole() *roles.Role { return ch.Role }
func (ch *BotRoleChannel) Write(b []byte) (n int, err error) {
	return SendMessage(ch.s, ch.ChannelIID, b)
}

func (ch *BotRoleChannel) AddPlayer(serverUserID string) error {
	var allow int64 = discordgo.PermissionSendMessages | discordgo.PermissionViewChannel
	ovrd := NewOverridePermission(allow, true)

	if _, contains := ch.mappedPermissions[serverUserID]; !contains {
		ch.mappedPermissions[serverUserID] = GetChannelUserPermission(ch.s, ch.Chat, serverUserID)
	}

	err := ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, ovrd.Type, ovrd.Allow, ovrd.Deny)
	return wrap.UnwrapDiscordRESTErr(err, discordgo.ErrCodeMissingPermissions)
}

func (ch *BotRoleChannel) AddSpectator(serverUserID string) error {
	var allow int64 = discordgo.PermissionViewChannel
	ovrd := NewOverridePermission(allow, true, discordgo.PermissionSendMessages)

	if _, contains := ch.mappedPermissions[serverUserID]; !contains {
		ch.mappedPermissions[serverUserID] = GetChannelUserPermission(ch.s, ch.Chat, serverUserID)
	}

	err := ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, ovrd.Type, ovrd.Allow, ovrd.Deny)
	return wrap.UnwrapDiscordRESTErr(err, discordgo.ErrCodeMissingPermissions)
}

func (ch *BotRoleChannel) RemoveUser(serverUserID string) (err error) {
	prev, isContains := ch.mappedPermissions[serverUserID]
	if isContains {
		delete(ch.mappedPermissions, serverUserID)
	}
	ovrd := NewOverridePermission(prev, false)
	if isContains {
		err = ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, ovrd.Type, ovrd.Allow, 0)
	} else {
		err = ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, ovrd.Type, 0, discordgo.PermissionAll)
	}
	return wrap.UnwrapDiscordRESTErr(err, discordgo.ErrCodeMissingPermissions)
}

// MainChannel provided a main interaction channel with players
// Core MainChannel implementation.

type BotMainChannel struct {
	ChannelIID string             `json:"channel_iid"`
	Chat       *discordgo.Channel `json:"channel"`
	s          *discordgo.Session
	// mappedPermissions stores past permissions for each participant in the game.
	mappedPermissions map[string]int64
}

func NewBotMainChannel(s *discordgo.Session, channelIID string) (*BotMainChannel, error) {
	channel, err := s.Channel(channelIID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting channel: %v", err))
	}

	return &BotMainChannel{
		s:                 s,
		Chat:              channel,
		ChannelIID:        channelIID,
		mappedPermissions: make(map[string]int64),
	}, nil
}

func (ch *BotMainChannel) GetServerID() string { return ch.ChannelIID }
func (ch *BotMainChannel) Write(b []byte) (n int, err error) {
	return SendMessage(ch.s, ch.ChannelIID, b)
}

func (ch *BotMainChannel) AddPlayer(serverUserID string) error {
	var allow int64 = discordgo.PermissionSendMessages | discordgo.PermissionViewChannel
	ovrd := NewOverridePermission(allow, true)

	if _, contains := ch.mappedPermissions[serverUserID]; !contains {
		ch.mappedPermissions[serverUserID] = GetChannelUserPermission(ch.s, ch.Chat, serverUserID)
	}

	err := ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, ovrd.Type, ovrd.Allow, ovrd.Deny)
	return wrap.UnwrapDiscordRESTErr(err, discordgo.ErrCodeMissingPermissions)
}

func (ch *BotMainChannel) AddSpectator(serverUserID string) error {
	var allow int64 = discordgo.PermissionViewChannel
	ovrd := NewOverridePermission(allow, true, discordgo.PermissionSendMessages)

	if _, contains := ch.mappedPermissions[serverUserID]; !contains {
		ch.mappedPermissions[serverUserID] = GetChannelUserPermission(ch.s, ch.Chat, serverUserID)
	}

	err := ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, ovrd.Type, ovrd.Allow, ovrd.Deny)
	return wrap.UnwrapDiscordRESTErr(err, discordgo.ErrCodeMissingPermissions)
}

func (ch *BotMainChannel) RemoveUser(serverUserID string) (err error) {
	prev, isContains := ch.mappedPermissions[serverUserID]
	if isContains {
		delete(ch.mappedPermissions, serverUserID)
	}
	ovrd := NewOverridePermission(prev, false)
	if isContains {
		err = ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, ovrd.Type, ovrd.Allow, 0)
	} else {
		err = ch.s.ChannelPermissionSet(ch.ChannelIID, serverUserID, ovrd.Type, 0, discordgo.PermissionAll)
	}
	return wrap.UnwrapDiscordRESTErr(err, discordgo.ErrCodeMissingPermissions)
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

func GetChannelUserPermission(s *discordgo.Session, ch *discordgo.Channel, userID string) int64 {
	permissions, _ := s.State.UserChannelPermissions(userID, ch.ID)
	return permissions
}
