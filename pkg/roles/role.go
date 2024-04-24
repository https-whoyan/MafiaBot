package roles

const (
	Peaceful = iota
	Mafia
)

type Role struct {
	Name        string `json:"name"`
	Team        int    `json:"team"`
	Description string `json:"description"`
}
