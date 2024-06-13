package roles

// For beauty messaging

var MappedEmoji = map[string]string{
	"Citizen":   "",
	"Detective": "",
	"Doctor":    "",
	"Don":       "",
	"Fool":      "",
	"Mafia":     "",
	"Maniac":    "",
	"Peaceful":  "",
	"Whose":     "",
}

// TODO!!!!!!!!! (Replace discord stickers to unicode)
var StringTeam = map[Team]string{
	PeacefulTeam: ":heart: Peaceful",
	MafiaTeam:    ":black_heart: Mafia Team",
	ManiacTeam:   ":grey_heart: Maniac Team",
}

func GetEmojiByName(name string) string {
	return MappedEmoji[name]
}
