package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
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
	GuildID    string `bson:"guild_id" json:"guild_id"`
	ChannelIID string `bson:"channel_iid" json:"channel_iid"`
	Role       string `bson:"role" json:"role"`
}

func (db *MongoDB) DeleteRoleChannel(guildID string, role string) error {
	if db.TryLock() {
		defer db.Unlock()
	}
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
	if db.TryLock() {
		defer db.Unlock()
	}
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
	if db.TryLock() {
		defer db.Unlock()
	}
	coll, err := db.getColl(DbName, RolesChannelCollection)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	result := ChannelRoleStruct{}
	filter := bson.D{
		{
			"guild_id", guildID,
		},
		{
			"channel_iid", channelIID,
		},
	}
	err = coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Role, nil
}

func (db *MongoDB) GetChannelIIDByRole(guildID string, role string) (string, error) {
	if db.TryLock() {
		defer db.Unlock()
	}
	coll, err := db.getColl(DbName, RolesChannelCollection)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	result := ChannelRoleStruct{}
	filter := bson.D{
		{
			"guild_id", guildID,
		},
		{
			"role", strings.ToLower(role),
		},
	}
	err = coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.ChannelIID, nil
}
