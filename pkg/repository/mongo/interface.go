package mongo

import (
	"context"
	coreGamePack "github.com/https-whoyan/MafiaCore/game"
)

type Storage interface {
	coreGamePack.Logger

	IsFreeChannelIID(ctx context.Context, guildID string, ChannelID string) (bool, error)

	SetRoleChannel(ctx context.Context, guildID string, channelIID string, role string) error
	GetRoleByChannelIID(ctx context.Context, guildID string, channelIID string) (string, error)
	GetChannelIIDByRole(ctx context.Context, guildID string, role string) (string, error)
	DeleteRoleChannel(ctx context.Context, guildID string, role string) (isDeleted bool, err error)

	SetMainChannel(ctx context.Context, guildID string, channelIID string) error
	GetMainChannelIID(ctx context.Context, guildID string) (string, error)

	Close(ctx context.Context) error
}
