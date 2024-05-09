package roles

const (
	PeacefulTeam = iota + 1
	MafiaTeam
)

type Role struct {
	Name        string `json:"name"`
	Team        int    `json:"team"`
	Description string `json:"description"`
}
