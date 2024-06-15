package channel

func SearchRoleChannelByID(channels []RoleChannel, channelID string) RoleChannel {
	for _, channel := range channels {
		if channel.GetServerID() == channelID {
			return channel
		}
	}

	return nil
}

func SearchChannelByGameID(channels []Channel, channelID string) Channel {
	for _, channel := range channels {
		if channel.GetServerID() == channelID {
			return channel
		}
	}

	return nil
}

func RoleSliceToChannelSlice(roleChannels []RoleChannel) []Channel {
	var channelSlice []Channel
	for _, roleChannel := range roleChannels {
		channelSlice = append(channelSlice, roleChannel)
	}
	return channelSlice
}
