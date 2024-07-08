package mongo

import (
	"context"
	"errors"
	"github.com/https-whoyan/MafiaBot/core/game"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
)

const (
	DbName = "discord-go_mafia_bot"
)

const (
	GuildChannelsCollection = "guild_channels"
	GameStorageCollection   = "game_storage"
)

var (
	ctx = context.Background()
)

func (db *MongoDB) getColl(dbName, collection string) (*mongo.Collection, error) {
	coll := db.db.Database(dbName).Collection(collection)
	if coll == nil {
		return nil, errors.New("empty mongo collection")
	}
	return coll, nil
}

// *******************
// _________________
// Channels
// _________________
// *******************

// Structs

type RoleChannel struct {
	Role       string `bson:"role"`
	ChannelIID string `bson:"channelIID"`
}

type MainChannel struct {
	ChannelIID string `bson:"channelIID"`
}

type GuildChannels struct {
	GuildID      string        `bson:"guildID"`
	RoleChannels []RoleChannel `bson:"roleChannels"`
	MainChannel  MainChannel   `bson:"mainChannel"`
}

func (g GuildChannels) searchChannelIIDByRole(role string) string {
	for _, channel := range g.RoleChannels {
		if channel.Role == role {
			return channel.ChannelIID
		}
	}

	return ""
}

func (g GuildChannels) searchRoleByChannelIID(chanelIID string) string {
	for _, channel := range g.RoleChannels {
		if channel.ChannelIID == chanelIID {
			return channel.Role
		}
	}

	return ""
}

func (g GuildChannels) getMainChannelIID() string {
	var channelID string
	defer func() {
		recover()
	}()
	channelID = g.MainChannel.ChannelIID
	return channelID
}

// ** Utils **

func (db *MongoDB) IsFreeChannelIID(guildID string, ChannelID string) (bool, error) {
	if db.TryLock() {
		defer db.Unlock()
	}

	coll, err := db.getColl(DbName, GuildChannelsCollection)
	ctx := context.Background()
	if err != nil {
		return false, err
	}

	filterRoleChannels := bson.D{
		{"guildID", guildID},
		{
			"roleChannels", bson.M{
				"$elemMatch": bson.M{"channelIID": ChannelID},
			},
		},
	}
	err = coll.FindOne(ctx, filterRoleChannels).Err()
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	filterMainChannels := bson.D{
		{
			"GuildID", guildID,
		},
		{
			"MainChannel.channelIID", ChannelID,
		},
	}
	err = coll.FindOne(ctx, filterMainChannels).Err()
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	return true, nil
}

// Push if not exists information about GuildID
func (db *MongoDB) pushIfNotExistsGuildChannels(guildID string) (isInserted bool, err error) {
	if db.TryLock() {
		defer db.Unlock()
	}

	coll, err := db.getColl(DbName, GuildChannelsCollection)
	if err != nil {
		return false, err
	}

	filter := bson.D{{"guildID", guildID}}
	ctx := context.Background()

	// I get an error if there are no documents in mongo
	isContains := !errors.Is(coll.FindOne(ctx, filter).Err(), mongo.ErrNoDocuments)
	// If contains
	// Return false (not add)
	if isContains {
		return false, nil
	}

	// Else push it
	newGuildInfo := GuildChannels{
		GuildID:      guildID,
		RoleChannels: []RoleChannel{},
		MainChannel:  MainChannel{},
	}

	_, err = coll.InsertOne(ctx, newGuildInfo)
	log.Println("Inserted new GuildChannels, guildID:", guildID)
	return true, err
}

// ________________
// ChannelRole
// ________________

func (db *MongoDB) getEntryByGuildID(guildID string) (*GuildChannels, error) {
	if db.TryLock() {
		defer db.Unlock()
	}
	coll, err := db.getColl(DbName, GuildChannelsCollection)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	result := GuildChannels{}
	filter := bson.D{{"guildID", guildID}}

	err = coll.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}
	return &result, nil
}

func (db *MongoDB) DeleteRoleChannel(guildID string, role string) (isDeleted bool, err error) {
	if db.TryLock() {
		defer db.Unlock()
	}

	coll, err := db.getColl(DbName, GuildChannelsCollection)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	filter := bson.D{
		{
			"guildID", guildID,
		},
	}
	update := bson.D{{
		"$pull", bson.D{
			{
				"roleChannels",
				bson.D{{"role", role}},
			},
		},
	}}

	result, err := coll.UpdateMany(ctx, filter, update)
	if result.ModifiedCount == 0 {
		return false, err
	}

	log.Printf(
		"Delete role %v ChannelIID in %v GuildID.",
		role,
		guildID)
	return true, err

}

func (db *MongoDB) SetRoleChannel(guildID string, channelIID string, role string) error {
	role = strings.ToLower(role)
	if db.TryLock() {
		defer db.Unlock()
	}
	// If channelIID used in other role:
	isFree, err := db.IsFreeChannelIID(guildID, channelIID)
	if err != nil || !isFree {
		return errors.New("permission denied")
	}
	if _, err = db.pushIfNotExistsGuildChannels(guildID); err != nil {
		return err
	}

	// Delete if contains
	_, err = db.DeleteRoleChannel(guildID, role)
	if err != nil {
		return err
	}
	ctx := context.Background()
	coll, err := db.getColl(DbName, GuildChannelsCollection)
	if err != nil {
		return err
	}
	filter := bson.D{
		{
			"guildID", guildID,
		},
	}

	updatePush := bson.D{{
		"$push", bson.D{{
			"roleChannels", bson.D{
				{"role", role},
				{"channelIID", channelIID},
			},
		}},
	}}

	_, err = coll.UpdateOne(ctx, filter, updatePush)

	log.Printf("add define for %v channel: guild: %v, role: %v", channelIID, guildID, role)
	return err
}

func (db *MongoDB) GetRoleByChannelIID(guildID string, channelIID string) (string, error) {
	if db.TryLock() {
		defer db.Unlock()
	}
	coll, err := db.getColl(DbName, GuildChannelsCollection)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	result := &GuildChannels{}
	filter := bson.D{
		{"guildID", guildID},
		{"channels", bson.M{
			"$elemMatch": bson.M{"channelIID": channelIID},
		}},
	}
	err = coll.FindOne(ctx, filter).Decode(result)
	if err != nil || result == nil {
		return "", nil
	}

	role := result.searchRoleByChannelIID(channelIID)
	return role, nil
}

func (db *MongoDB) GetChannelIIDByRole(guildID string, role string) (string, error) {
	entry, err := db.getEntryByGuildID(guildID)
	if err != nil {
		return "", err
	}

	role = strings.ToLower(role)
	return entry.searchChannelIIDByRole(role), nil
}

// ________________
// Main Channel
// ________________

// Not to need Delete, just push

func (db *MongoDB) SetMainChannel(guildID string, channelIID string) error {
	if db.TryLock() {
		defer db.Unlock()
	}
	// If channelIID used in other role:
	isFree, err := db.IsFreeChannelIID(guildID, channelIID)
	if err != nil || !isFree {
		return errors.New("permission denied")
	}

	if _, err = db.pushIfNotExistsGuildChannels(guildID); err != nil {
		return err
	}

	coll, err := db.getColl(DbName, GuildChannelsCollection)
	if err != nil {
		return err
	}

	filter := bson.D{{"guildID", guildID}}
	updateSet := bson.D{{
		"$set", bson.D{{
			"mainChannel.channelIID", channelIID,
		}},
	}}
	ctx := context.Background()

	result, err := coll.UpdateOne(ctx, filter, updateSet)
	if err != nil || result.ModifiedCount == 0 {
		return errors.New("no updated documents")
	}

	log.Printf("set main channel %v GuildID: %v", channelIID, guildID)
	return err
}

func (db *MongoDB) GetMainChannelIID(guildID string) (string, error) {
	entry, err := db.getEntryByGuildID(guildID)
	if err != nil {
		return "", err
	}

	return entry.getMainChannelIID(), nil
}

// *******************
// _________________
// Logs
// _________________
// *******************

func (db *MongoDB) InitNewGame(g *game.Game) error {
	db.Lock()
	defer db.Unlock()
	coll, err := db.getColl(DbName, GameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, newMongoGameLog(g))
	return err
}

func (db *MongoDB) SaveNightLog(g *game.Game, l game.NightLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := bson.D{{
		"$push", bson.D{{
			"nightLogs", newMongoGameNight(l),
		}},
	}}
	db.Lock()
	defer db.Unlock()
	coll, err := db.getColl(DbName, GameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}

func (db *MongoDB) SaveDayLog(g *game.Game, l game.DayLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := bson.D{{
		"$push", bson.D{{
			"da", newMongoGameDay(l),
		}},
	}}
	db.Lock()
	defer db.Unlock()
	coll, err := db.getColl(DbName, GameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}

func (db *MongoDB) SaveFinishLog(g *game.Game, l game.FinishLog) error {
	filter := bson.M{
		"guildID":   g.GuildID,
		"startTime": g.TimeStart,
	}
	updatePush := getUpdateByByNightLog(l)
	db.Lock()
	defer db.Unlock()
	coll, err := db.getColl(DbName, GameStorageCollection)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(ctx, filter, updatePush)
	return err
}
