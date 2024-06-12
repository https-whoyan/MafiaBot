package converter

import (
	botChannel "github.com/https-whoyan/MafiaBot/internal/bot/channel"
	IChannel "github.com/https-whoyan/MafiaBot/internal/core/channel"
)

func ConvertRoleChannelsSliceToIChannelSlice(sl []*botChannel.BotRoleChannel) []IChannel.RoleChannel {
	ans := make([]IChannel.RoleChannel, len(sl))
	for i, roleChannel := range sl {
		ans[i] = roleChannel
	}
	return ans
}
