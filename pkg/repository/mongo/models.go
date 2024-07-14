package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"

	gameConfig "github.com/https-whoyan/MafiaBot/core/config"
	"github.com/https-whoyan/MafiaBot/core/game"
	"github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"
)

// Document struct

type mongoGameLog struct {
	ID      string `json:"id,omitempty" bson:"_id,omitempty"`
	GuildID string `json:"guildID,omitempty" bson:"guildID,omitempty"`
	LogName string `json:"logName,omitempty" bson:"logName,omitempty"`

	StartTime time.Time `json:"startTime,omitempty" bson:"startTime,omitempty"`
	Config    *gameConfig.RolesConfig

	Players   []mongoGamePlayer `json:"players,omitempty" bson:"players,omitempty"`
	NightLogs []mongoGameNight  `json:"nightLogs,omitempty" bson:"nightLogs,omitempty"`
	DayLogs   []mongoGameDay    `json:"dayLogs,omitempty" bson:"dayLogs,omitempty"`

	NightsPlayed int        `json:"nightsPlayed,omitempty" bson:"nightsPlayed,omitempty"`
	FinishTeam   roles.Team `json:"finishTeam,omitempty" bson:"finishTeam,omitempty"`
	WinByFool    bool       `json:"winByFool,omitempty" bson:"winByFool,omitempty"`
	IsSuspended  bool       `json:"isSuspended,omitempty" bson:"isSuspended,omitempty"`
}

func newMongoGameLog(g *game.Game) mongoGameLog {
	var playersModel []mongoGamePlayer
	for _, p := range *g.Active {
		mongoP := newMongoGamePlayer(p)
		playersModel = append(playersModel, mongoP)
	}
	model := mongoGameLog{
		GuildID:     g.GuildID,
		StartTime:   g.TimeStart,
		Players:     playersModel,
		Config:      g.RolesConfig,
		IsSuspended: true,
	}
	return model
}

// Player struct

type mongoGamePlayer struct {
	IDInGame       int    `json:"idInGame,omitempty" bson:"idInGame,omitempty"`
	ServerID       string `json:"serverID,omitempty" bson:"serverID,omitempty"`
	ServerUsername string `json:"serverUsername" bson:"serverUsername"`
	Role           string `json:"role,omitempty" bson:"role,omitempty"`
}

func newMongoGamePlayer(p *player.Player) mongoGamePlayer {
	return mongoGamePlayer{
		IDInGame:       int(p.ID),
		ServerID:       p.Tag,
		ServerUsername: p.ServerNick,
		Role:           p.Role.Name,
	}
}

// NightLog struct

type mongoGameNight struct {
	Number      int                    `json:"number,omitempty" bson:"number,omitempty"`
	Votes       map[int]mongoGameVotes `json:"votes,omitempty" bson:"votes,omitempty"`
	DeadPlayers []int                  `json:"deadPlayers,omitempty" bson:"deadPlayers,omitempty"`
}

func newMongoGameNight(l game.NightLog) mongoGameNight {
	votes := make(map[int]mongoGameVotes)
	for voter, voterVotes := range l.NightVotes {
		votes[voter] = mongoGameVotes{voterVotes}
	}
	return mongoGameNight{
		Number:      l.NightNumber,
		Votes:       votes,
		DeadPlayers: l.Dead,
	}
}

type mongoGameVotes struct {
	Votes []int `json:"votes,omitempty" bson:"votes,omitempty"`
}

// DayLog struct

type mongoGameDay struct {
	Number int         `json:"number,omitempty" bson:"number,omitempty"`
	Votes  map[int]int `json:"votes,omitempty" bson:"votes,omitempty"`
	Kicked *int        `json:"kicked,omitempty" bson:"kicked,omitempty"`
	IsSkip bool        `json:"isSkip,omitempty" bson:"isSkip,omitempty"`
}

func newMongoGameDay(l game.DayLog) mongoGameDay {
	return mongoGameDay{
		Number: l.DayNumber,
		Votes:  l.DayVotes,
		IsSkip: l.IsSkip,
		Kicked: l.Kicked,
	}
}

// Finish get an update filter
func getUpdateByByNightLog(l game.FinishLog) bson.M {
	if l.IsFool {
		return bson.M{
			"$set": bson.M{
				"isSuspended":  false,
				"nightsPlayed": l.TotalNights,
				"winByFool":    true,
			},
		}
	}
	return bson.M{
		"$set": bson.M{
			"isSuspended":  false,
			"nightsPlayed": l.TotalNights,
			"finishTeam":   l.WinnerTeam,
		},
	}

}
