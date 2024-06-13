package converter

import "github.com/https-whoyan/MafiaBot/internal/core/channel"

func RoleSliceToChannelSlice(roleChannels []channel.RoleChannel) []channel.Channel {
	var channelSlice []channel.Channel
	for _, roleChannel := range roleChannels {
		channelSlice = append(channelSlice, roleChannel)
	}
	return channelSlice
}
