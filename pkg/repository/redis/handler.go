package redis

import (
	"context"
	time2 "github.com/https-whoyan/MafiaBot/internal/core/time"
	"time"
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
	lifeDuration := time2.RegistrationDeadlineSeconds * time.Second
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
	lifeDuration := time2.VotingGameConfigDeadlineSeconds * time.Second
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
