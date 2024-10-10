package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/https-whoyan/MafiaCore/game"
	"github.com/redis/go-redis/v9"
	"log"
	"reflect"
	"strconv"
	"time"

	botGame "github.com/https-whoyan/MafiaBot/internal/game"
)

type redisDB struct {
	db *redis.Client
	lg *log.Logger
}

const (
	initialGameTB      = "initialGames"
	configVotingGameTB = "configVotingGame"
	indicatorGameTable = "indicatorGames"
)

const (
	nonFound = "non_found"
)

// Utils

// It can be r.db.Del...
func (r *redisDB) deleteKey(key string) { r.db.PExpire(context.Background(), key, time.Millisecond) }

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

func (r *redisDB) SetInitialGameMessageID(ctx context.Context, guildID string, messageID string, lifeDuration time.Duration) error {
	key := initialGameTB + ":" + guildID

	pipe := r.db.TxPipeline()
	pipe.Set(ctx, key, messageID, lifeDuration)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisDB) GetInitialGameMessageID(ctx context.Context, guildID string) (string, error) {
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

func (r *redisDB) SetConfigGameVotingMessages(ctx context.Context, c *botGame.ConfigMessages, lifeDuration time.Duration) (err error) {
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

func (r *redisDB) GetConfigGameVotingMessageID(ctx context.Context, guildID string) (*botGame.ConfigMessages, error) {
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

// ChannelsIDStorage

type СhannelStorageType string

const (
	channelStoragePrefix = "channelStorage"
)

const (
	SetInitialGameStorage СhannelStorageType = "setInitialGame"
	//...
)

func (r *redisDB) SetChannelStorage(ctx context.Context, guildID string, channelIID string,
	storageType СhannelStorageType, duration time.Duration) error {
	key := channelStoragePrefix + ":" + guildID + ":" + string(storageType)

	pipe := r.db.TxPipeline()
	pipe.Set(ctx, key, channelIID, duration)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisDB) GetChannelStorage(ctx context.Context, guildID string, storageType СhannelStorageType) (channelIID string, err error) {
	key := channelStoragePrefix + ":" + guildID + ":" + string(storageType)

	pipe := r.db.TxPipeline()
	cmd := pipe.Get(ctx, key)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return "", err
	}
	return cmd.Result()
}

// _____________
// Rename Game
// _____________

func (r *redisDB) SaveGameIndicator(ctx context.Context, indicator string, g game.DeepCloneGame) error {
	pipe := r.db.TxPipeline()
	raw, err := g.MarshalJSON()
	if err != nil {
		return err
	}
	pipe.HSet(ctx, indicatorGameTable, indicator, raw)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *redisDB) GetGameByIndicator(ctx context.Context, indicator string) (out game.DeepCloneGame, err error) {
	pipe := r.db.TxPipeline()
	cmd := pipe.HGet(ctx, indicatorGameTable, indicator)
	// Del HKey
	pipe.HDel(ctx, indicatorGameTable, indicator)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return out, err
	}
	strRaw, err := cmd.Result()
	if err != nil {
		return out, err
	}
	var dst game.DeepCloneGame
	err = json.Unmarshal([]byte(strRaw), &dst)
	if err != nil {
		return out, err
	}
	return dst, nil
}

// Close

func (r *redisDB) Close(_ context.Context) error {
	r.lg.Println("Disconnect Redis")
	return r.db.Close()
}
