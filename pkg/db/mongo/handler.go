package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const (
	DbName = "discord-go_mafia_bot"
)

const (
	RolesChannelCollection = "roles_channels"
)

func (db *MongoDB) getColl(dbName, collection string) (*mongo.Collection, error) {
	coll := db.db.Database(dbName).Collection(collection)
	if coll == nil {
		return nil, errors.New("empty mongo collection")
	}
	return coll, nil
}

// ________________
// ChannelRole
// ________________

type ChannelRoleStruct struct {
	GuildID    string `bson:"guild_id"`
	ChannelIID string `bson:"channel_iid"`
	Role       string `bson:"role"`
}

func (db *MongoDB) DeleteRoleChannel(guildID string, role string) error {
	coll, err := db.getColl(DbName, RolesChannelCollection)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = coll.DeleteOne(ctx, bson.D{
		{
			"guild_id", guildID,
		},
		{
			"role", role,
		},
	})
	return err
}

func (db *MongoDB) SetRoleChannel(guildID string, channelIID string, role string) error {
	// Delete if contains
	err := db.DeleteRoleChannel(guildID, role)
	coll, err := db.getColl(DbName, RolesChannelCollection)
	if err != nil {
		return err
	}

	ctx := context.Background()
	insertedRow := ChannelRoleStruct{
		GuildID:    guildID,
		ChannelIID: channelIID,
		Role:       role,
	}
	_, err = coll.InsertOne(ctx, insertedRow)
	if err != nil {
		return err
	}
	log.Printf("add define for %v channel: guild: %v, role: %v", channelIID, guildID, role)
	return err
}

func (db *MongoDB) GetRoleByChannelIID(guildID string, channelIID string) (string, error) {
	coll, err := db.getColl(DbName, RolesChannelCollection)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	var result *ChannelRoleStruct
	err = coll.FindOne(ctx, bson.D{
		{
			"guild_id", guildID,
		},
		{
			"channel_iid", channelIID,
		},
	}).Decode(result)
	if err != nil {
		return "", err
	}

	return result.Role, nil
}

func (db *MongoDB) GetChannelIIDByRole(guildID string, role string) (string, error) {
	coll, err := db.getColl(DbName, RolesChannelCollection)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	var result *ChannelRoleStruct
	err = coll.FindOne(ctx, bson.D{
		{
			"guild_id", guildID,
		},
		{
			"role", role,
		},
	}).Decode(result)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", errors.New("role not found")
	}

	return result.Role, nil
}
