package mongo

import (
	"context"
	"github.com/https-whoyan/MafiaCore/game"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *mongoDB) InitNewGame(ctx context.Context, g game.DeepCloneGame) error {
	coll, err := s.getColl(dbName, gameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, newMongoGameLog(g))
	return err
}

func (s *mongoDB) SaveNightLog(ctx context.Context, g game.DeepCloneGame, l game.NightLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := bson.D{{
		"$push", bson.D{{
			"nightLogs", newMongoGameNight(l),
		}},
	}}
	coll, err := s.getColl(dbName, gameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}

func (s *mongoDB) SaveDayLog(ctx context.Context, g game.DeepCloneGame, l game.DayLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := bson.D{{
		"$push", bson.D{{
			"day_log", newMongoGameDay(l),
		}},
	}}
	coll, err := s.getColl(dbName, gameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}

func (s *mongoDB) SaveFinishLog(ctx context.Context, g game.DeepCloneGame, l game.FinishLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := getUpdateByByNightLog(l)
	coll, err := s.getColl(dbName, gameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}

func (s *mongoDB) NameAGame(ctx context.Context, g game.DeepCloneGame, gameName string) error {
	coll, err := s.getColl(dbName, gameStorageCollection)
	if err != nil {
		return err
	}
	filter := bson.M{
		"guildID":    g.GuildID,
		"time_start": g.TimeStart,
	}
	update := bson.D{{
		"$set", bson.D{{
			"name", gameName,
		}},
	}}
	_, err = coll.UpdateOne(ctx, filter, update)
	return err
}
