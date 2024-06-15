package redis

import (
	"context"
	"reflect"
	"strconv"
	"time"

	botGame "github.com/https-whoyan/MafiaBot/internal/game"
)

const (
	initialGameTB      = "initialGames"
	configVotingGameTB = "configVotingGameTB"
)

// Utils

// It can be r.db.Del...
func (r *RedisDB) deleteKey(key string) { r.db.PExpire(context.Background(), key, time.Millisecond) }

func redisNameByFieldName(t reflect.Type, fieldName string) string {
	field, isFound := t.FieldByName(fieldName)
	if !isFound {
		return "NON_FOUND"
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

// _________________
// ConfigsVoting
// _________________

func (r *RedisDB) SetConfigGameVotingMessages(c *botGame.ConfigMessages, lifeDuration time.Duration) (err error) {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := configVotingGameTB + ":" + c.GuildID

	t := reflect.TypeOf(*c)

	// Set firstValue
	err = r.db.HSet(ctx, key, redisNameByFieldName(t, "PlayersCount"), c.PlayersCount).Err()
	if err != nil {
		return
	}
	// Set lifeDuration
	err = r.db.PExpire(ctx, key, lifeDuration).Err()
	if err != nil {
		return
	}
	// Set other fields
	err = r.db.HSet(ctx, key, redisNameByFieldName(t, "ConfigsCount"), c.ConfigsCount).Err()
	if err != nil {
		return
	}
	fieldPrefix := redisNameByFieldName(t, "Messages") + ":"
	for _, message := range c.Messages {
		fieldName := fieldPrefix + strconv.Itoa(message.ConfigIndex)
		err = r.db.HSet(ctx, key, fieldName, message.MessageID).Err()
		if err != nil {
			return
		}
	}
	return err
}

func (r *RedisDB) GetConfigGameVotingMessageID(guildID string) (*botGame.ConfigMessages, error) {
	r.Lock()
	defer r.Unlock()

	ctx := context.Background()
	key := configVotingGameTB + ":" + guildID

	defer r.deleteKey(key)

	t := reflect.TypeOf(botGame.ConfigMessages{})

	playersCountStr, err := r.db.HGet(ctx, key, redisNameByFieldName(t, "PlayersCount")).Result()
	playersCount, _ := strconv.Atoi(playersCountStr)
	if err != nil {
		return nil, err
	}
	configsCountStr, err := r.db.HGet(ctx, key, redisNameByFieldName(t, "ConfigsCount")).Result()
	configsCount, _ := strconv.Atoi(configsCountStr)
	if err != nil {
		return nil, err
	}

	ans := botGame.NewConfigMessages(guildID, playersCount, configsCount)

	fieldPrefix := redisNameByFieldName(t, "Messages") + ":"
	keys := make([]string, ans.ConfigsCount)
	for i := 0; i <= ans.ConfigsCount-1; i++ {
		keys[i] = fieldPrefix + strconv.Itoa(i)
	}
	values, err := r.db.HMGet(ctx, key, keys...).Result()
	if err != nil {
		return ans, err
	}
	for i, v := range values {
		ans.AddNewMessage(i, v.(string))
	}
	return ans, err
}
