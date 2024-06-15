package game

import (
	"github.com/https-whoyan/MafiaBot/core/channel"
	"github.com/https-whoyan/MafiaBot/core/converter"
)

func (g *Game) searchRoleChannelByIID(channelIID string) channel.RoleChannel {
	allRoleChannels := converter.GetMapValues(g.RoleChannels)
	channelVotedFrom := channel.SearchRoleChannelByID(allRoleChannels, channelIID)

	return channelVotedFrom
}

func (g *Game) searchChannelByID(channelIID string) channel.Channel {
	var allChannels []channel.Channel
	allRoleChannels := converter.GetMapValues(g.RoleChannels)
	allChannels = append(allChannels, channel.RoleSliceToChannelSlice(allRoleChannels)...)

	allChannels = append(allChannels, g.MainChannel)

	channelVotedFrom := channel.SearchChannelByGameID(allChannels, channelIID)
	return channelVotedFrom
}
