package converter

import (
	botChannel "github.com/https-whoyan/MafiaBot/internal/channel"
	IChannel "github.com/https-whoyan/MafiaCore/channel"
)

func ConvertRoleChannelsSliceToIChannelSlice(sl []*botChannel.BotRoleChannel) []IChannel.RoleChannel {
	ans := make([]IChannel.RoleChannel, len(sl))
	for i, roleChannel := range sl {
		ans[i] = roleChannel
	}
	return ans
}
