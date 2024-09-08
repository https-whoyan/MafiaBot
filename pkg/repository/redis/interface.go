package redis

import (
	"context"
	"time"

	botGame "github.com/https-whoyan/MafiaBot/internal/game"
)

type Hasher interface {
	SetInitialGameMessageID(ctx context.Context, guildID string, messageID string, lifeDuration time.Duration) error
	GetInitialGameMessageID(ctx context.Context, guildID string) (string, error)

	SetConfigGameVotingMessages(ctx context.Context, c *botGame.ConfigMessages, lifeDuration time.Duration) (err error)
	GetConfigGameVotingMessageID(ctx context.Context, guildID string) (*botGame.ConfigMessages, error)

	SetChannelStorage(ctx context.Context, guildID string, channelIID string, storageType СhannelStorageType, duration time.Duration) error
	GetChannelStorage(ctx context.Context, guildID string, storageType СhannelStorageType) (channelIID string, err error)

	Close(ctx context.Context) error
}
