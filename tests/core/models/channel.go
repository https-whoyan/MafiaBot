package models

import (
	"github.com/https-whoyan/MafiaBot/core/roles"
	"strconv"
)

type TestChannel struct {
	Messages   []string
	ChannelIID string
}

func NewTestChannel(channelIID string) *TestChannel {
	return &TestChannel{
		Messages:   make([]string, 0),
		ChannelIID: channelIID,
	}
}

func (c *TestChannel) Write(b []byte) (n int, err error) {
	c.Messages = append(c.Messages, string(b))
	return len(b), nil
}

func (c *TestChannel) AddPlayer(_ string) error    { return nil }
func (c *TestChannel) RemoveUser(_ string) error   { return nil }
func (c *TestChannel) AddSpectator(_ string) error { return nil }
func (c *TestChannel) GetServerID() string         { return c.ChannelIID }

type TestRoleChannel struct {
	TestChannel
	Role *roles.Role
}

func NewTestRoleChannel(channelIID string, role *roles.Role) *TestRoleChannel {
	return &TestRoleChannel{
		TestChannel: *NewTestChannel(channelIID),
		Role:        role,
	}
}

func (c *TestRoleChannel) GetRole() *roles.Role { return c.Role }

type TestMainChannel struct {
	TestChannel
}

func NewTestMainChannel(channelIID string) *TestMainChannel {
	return &TestMainChannel{
		TestChannel: *NewTestChannel(channelIID),
	}
}

func NewTestChannels() []*TestRoleChannel {
	allIntersectionRolesNames := roles.GetAllNightInteractionRolesNames()

	var channels []*TestRoleChannel
	for i, roleName := range allIntersectionRolesNames {
		channelIID := strconv.Itoa(i)
		intersectionRole, _ := roles.GetRoleByName(roleName)
		channels = append(channels, NewTestRoleChannel(channelIID, intersectionRole))
	}

	return channels
}

const (
	TestMainChannelIID = "20"
)

func NewTestMainChannels() *TestMainChannel {
	return NewTestMainChannel(TestMainChannelIID)
}
