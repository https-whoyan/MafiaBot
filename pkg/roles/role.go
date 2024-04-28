package roles

const (
	PeacefulTeam = iota
	MafiaTeam
)

type Role struct {
	Name        string `json:"name"`
	Team        int    `json:"team"`
	Description string `json:"description"`
}
