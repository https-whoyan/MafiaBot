package player

import "github.com/https-whoyan/MafiaBot/core/roles"

// All interactions for roles are declared here to avoid cyclic import.

type Message string

func (p *Player) Interaction() Message {
	switch p.Role {
	case roles.Peaceful:
		return ""
	case roles.Fool:
		return ""
	case roles.Mafia:
		mafiaInteraction(p)
	}
	return ""
}

// Todo!
// ________________
// Mafia
// ________________

func mafiaInteraction(p *Player) {
	// Todo!
}
