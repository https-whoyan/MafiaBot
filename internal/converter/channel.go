package converter

import (
	IChannel "github.com/https-whoyan/MafiaBot/core/channel"
	botChannel "github.com/https-whoyan/MafiaBot/internal/channel"
)

func ConvertRoleChannelsSliceToIChannelSlice(sl []*botChannel.BotRoleChannel) []IChannel.RoleChannel {
	ans := make([]IChannel.RoleChannel, len(sl))
	for i, roleChannel := range sl {
		ans[i] = roleChannel
	}
	return ans
}
