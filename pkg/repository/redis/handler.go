package redis

import (
	"context"
	"time"

	coreTime "github.com/https-whoyan/MafiaBot/core/time"
)

const (
	initialGameTB      = "initialGames"
	configVotingGameTB = "configVotingGameTB"
)

func (r *RedisDB) SetInitialGameMessageID(guildID string, messageID string) error {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := initialGameTB + ":" + guildID
	lifeDuration := coreTime.RegistrationDeadlineSeconds * time.Second
	err := r.db.Set(ctx, key, messageID, lifeDuration).Err()
	return err
}

func (r *RedisDB) GetInitialGameMessageID(guildID string) (string, error) {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := initialGameTB + ":" + guildID
	cmd := r.db.Get(ctx, key)
	return cmd.Result()
}

func (r *RedisDB) SetConfigGameVotingMessageID(guildID string, messageID string) error {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := configVotingGameTB + ":" + guildID
	lifeDuration := coreTime.VotingGameConfigDeadlineSeconds * time.Second
	err := r.db.Set(ctx, key, messageID, lifeDuration).Err()
	return err
}

func (r *RedisDB) GetConfigGameVotingMessageID(guildID string) (string, error) {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := configVotingGameTB + ":" + guildID
	cmd := r.db.Get(ctx, key)
	return cmd.Result()
}
