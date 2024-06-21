package redis

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"time"

	botGame "github.com/https-whoyan/MafiaBot/internal/game"
)

const (
	initialGameTB      = "initialGames"
	configVotingGameTB = "configVotingGameTB"
)

const (
	nonFound = "non_found"
)

// Utils

// It can be r.db.Del...
func (r *RedisDB) deleteKey(key string) { r.db.PExpire(context.Background(), key, time.Millisecond) }

func redisNameByFieldName(t reflect.Type, fieldName string) string {
	field, isFound := t.FieldByName(fieldName)
	if !isFound {
		return nonFound
	}
	return field.Tag.Get("redis")
}

// _____________________
// InitialGameMessage
// _____________________

func (r *RedisDB) SetInitialGameMessageID(guildID string, messageID string, lifeDuration time.Duration) error {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := initialGameTB + ":" + guildID

	pipe := r.db.TxPipeline()
	pipe.Set(ctx, key, messageID, lifeDuration)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisDB) GetInitialGameMessageID(guildID string) (string, error) {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := initialGameTB + ":" + guildID

	pipe := r.db.TxPipeline()
	cmd := pipe.Get(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return "", err
	}

	return cmd.Result()
}

// _________________
// ConfigsVoting
// _________________

func (r *RedisDB) SetConfigGameVotingMessages(c *botGame.ConfigMessages, lifeDuration time.Duration) (err error) {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := configVotingGameTB + ":" + c.GuildID

	t := reflect.TypeOf(*c)

	pipe := r.db.TxPipeline()

	pipe.HSet(ctx, key, redisNameByFieldName(t, "PlayersCount"), c.PlayersCount)
	pipe.PExpire(ctx, key, lifeDuration)
	pipe.HSet(ctx, key, redisNameByFieldName(t, "ConfigsCount"), c.ConfigsCount)

	fieldPrefix := redisNameByFieldName(t, "Messages") + ":"
	if fieldPrefix == nonFound+":" {
		pipe.Discard()
		return errors.New("reflect field error")
	}
	for _, message := range c.Messages {
		fieldName := fieldPrefix + strconv.Itoa(message.ConfigIndex)
		pipe.HSet(ctx, key, fieldName, message.MessageID)
	}

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisDB) GetConfigGameVotingMessageID(guildID string) (*botGame.ConfigMessages, error) {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := configVotingGameTB + ":" + guildID

	t := reflect.TypeOf(botGame.ConfigMessages{})

	pipe := r.db.TxPipeline()

	playersCountCmd := pipe.HGet(ctx, key, redisNameByFieldName(t, "PlayersCount"))
	configCountCmd := pipe.HGet(ctx, key, redisNameByFieldName(t, "ConfigsCount"))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	playersCountStr, configCountStr := playersCountCmd.Val(), configCountCmd.Val()
	var playersCount, configsCount int
	playersCount, err = strconv.Atoi(playersCountStr)
	if err != nil {
		return nil, err
	}
	configsCount, err = strconv.Atoi(configCountStr)
	if err != nil {
		return nil, err
	}

	ans := botGame.NewConfigMessages(guildID, playersCount, configsCount)
	fieldPrefix := redisNameByFieldName(t, "Messages") + ":"
	keys := make([]string, ans.ConfigsCount)
	for i := 0; i <= ans.ConfigsCount-1; i++ {
		keys[i] = fieldPrefix + strconv.Itoa(i)
	}

	pipe = r.db.TxPipeline()
	valuesCmd := pipe.HMGet(ctx, key, keys...)
	pipe.Del(ctx, key)

	_, err = pipe.Exec(ctx)
	values := valuesCmd.Val()

	if err != nil {
		return nil, err
	}
	for i, v := range values {
		ans.AddNewMessage(i, v.(string))
	}

	return ans, err
}
