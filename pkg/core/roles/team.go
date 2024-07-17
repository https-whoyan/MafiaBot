package roles

type Team int

const (
	PeacefulTeam Team = iota
	MafiaTeam
	ManiacTeam
)

var StringTeam = map[Team]string{
	PeacefulTeam: "❤️ Peaceful Team",
	MafiaTeam:    "🖤 Mafia Team",
	ManiacTeam:   "\U0001FA76 Maniac Team",
}
