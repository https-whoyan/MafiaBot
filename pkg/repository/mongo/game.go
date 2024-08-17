package mongo

import (
	"github.com/https-whoyan/MafiaCore/game"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *MongoDB) InitNewGame(g game.DeepCloneGame) error {
	coll, err := db.getColl(DbName, GameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, newMongoGameLog(g))
	return err
}

func (db *MongoDB) SaveNightLog(g game.DeepCloneGame, l game.NightLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := bson.D{{
		"$push", bson.D{{
			"nightLogs", newMongoGameNight(l),
		}},
	}}
	coll, err := db.getColl(DbName, GameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}

func (db *MongoDB) SaveDayLog(g game.DeepCloneGame, l game.DayLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := bson.D{{
		"$push", bson.D{{
			"day_log", newMongoGameDay(l),
		}},
	}}
	coll, err := db.getColl(DbName, GameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}

func (db *MongoDB) SaveFinishLog(g game.DeepCloneGame, l game.FinishLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := getUpdateByByNightLog(l)
	coll, err := db.getColl(DbName, GameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}
