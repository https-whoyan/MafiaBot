package mongo

import (
	"context"
	"errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const (
	dbName = "discord-go_mafia_bot"
)

const (
	guildChannelsCollection = "guild_channels"
	gameStorageCollection   = "game_storage"
)

var (
	PermissionDeniedErr  = errors.New("permission denied")
	NoUpdatedDocumentErr = errors.New("no updated documents")
)

func (s *mongoDB) getColl(dbName, collection string) (*mongo.Collection, error) {
	coll := s.db.Database(dbName).Collection(collection)
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
	ans, contains := lo.Find(g.RoleChannels, func(ch RoleChannel) bool {
		return ch.Role == role
	})
	if !contains {
		return ""
	}
	return ans.ChannelIID
}

func (g GuildChannels) searchRoleByChannelIID(chanelIID string) string {
	ans, contains := lo.Find(g.RoleChannels, func(ch RoleChannel) bool {
		return ch.ChannelIID == chanelIID
	})
	if !contains {
		return ""
	}
	return ans.ChannelIID
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

func (s *mongoDB) IsFreeChannelIID(ctx context.Context, guildID string, ChannelID string) (bool, error) {
	coll, err := s.getColl(dbName, guildChannelsCollection)
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
func (s *mongoDB) pushIfNotExistsGuildChannels(ctx context.Context, guildID string) (isInserted bool, err error) {
	coll, err := s.getColl(dbName, guildChannelsCollection)
	if err != nil {
		return false, err
	}

	filter := bson.D{{"guildID", guildID}}

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
	s.lg.Println("Inserted new GuildChannels, guildID:", guildID)
	return true, err
}

// ________________
// ChannelRole
// ________________

func (s *mongoDB) getEntryByGuildID(ctx context.Context, guildID string) (*GuildChannels, error) {
	coll, err := s.getColl(dbName, guildChannelsCollection)
	if err != nil {
		return nil, err
	}

	result := GuildChannels{}
	filter := bson.D{{"guildID", guildID}}

	err = coll.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found")
	}
	return &result, nil
}

func (s *mongoDB) DeleteRoleChannel(ctx context.Context, guildID string, role string) (isDeleted bool, err error) {
	coll, err := s.getColl(dbName, guildChannelsCollection)
	if err != nil {
		return false, err
	}

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

	s.lg.Printf(
		"Delete role %v ChannelIID in %v GuildID.",
		role,
		guildID)
	return true, err

}

func (s *mongoDB) SetRoleChannel(ctx context.Context, guildID string, channelIID string, role string) error {
	role = strings.ToLower(role)
	// If channelIID used in other role:
	isFree, err := s.IsFreeChannelIID(ctx, guildID, channelIID)
	if err != nil {
		return err
	}
	if !isFree {
		return PermissionDeniedErr
	}
	if _, err = s.pushIfNotExistsGuildChannels(ctx, guildID); err != nil {
		return err
	}

	// Delete if contains
	_, err = s.DeleteRoleChannel(ctx, guildID, role)
	if err != nil {
		return err
	}
	coll, err := s.getColl(dbName, guildChannelsCollection)
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

	s.lg.Printf("add define for %v channel: guild: %v, role: %v", channelIID, guildID, role)
	return err
}

func (s *mongoDB) GetRoleByChannelIID(ctx context.Context, guildID string, channelIID string) (string, error) {
	coll, err := s.getColl(dbName, guildChannelsCollection)
	if err != nil {
		return "", err
	}

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

func (s *mongoDB) GetChannelIIDByRole(ctx context.Context, guildID string, role string) (string, error) {
	entry, err := s.getEntryByGuildID(ctx, guildID)
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

func (s *mongoDB) SetMainChannel(ctx context.Context, guildID string, channelIID string) error {
	// If channelIID used in other role:
	isFree, err := s.IsFreeChannelIID(ctx, guildID, channelIID)
	if err != nil {
		return err
	}
	if !isFree {
		return PermissionDeniedErr
	}

	if _, err = s.pushIfNotExistsGuildChannels(ctx, guildID); err != nil {
		return err
	}

	coll, err := s.getColl(dbName, guildChannelsCollection)
	if err != nil {
		return err
	}

	filter := bson.D{{"guildID", guildID}}
	updateSet := bson.D{{
		"$set", bson.D{{
			"mainChannel.channelIID", channelIID,
		}},
	}}

	result, err := coll.UpdateOne(ctx, filter, updateSet)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return NoUpdatedDocumentErr
	}

	s.lg.Printf("set main channel %v GuildID: %v", channelIID, guildID)
	return nil
}

func (s *mongoDB) GetMainChannelIID(ctx context.Context, guildID string) (string, error) {
	entry, err := s.getEntryByGuildID(ctx, guildID)
	if err != nil {
		return "", err
	}

	return entry.getMainChannelIID(), nil
}
