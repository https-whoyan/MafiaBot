package game

import (
	"github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"
	"strings"
)

// This is where all the code regarding reincarnation and role reversal is contained.
// This may be necessary if a role is specified to become a different role in certain scenarios.

func (g *Game) reincarnation(ch chan<- Signal, p *player.Player) {
	switch p.Role {
	case roles.Don:
		g.donReincarnation(ch, p)
	}
	return
}

func (g *Game) donReincarnation(ch chan<- Signal, p *player.Player) {
	// I find out he's the only one left on the mafia team.
	g.RLock()
	mafiaTeamCounter := 0
	for _, activePlayer := range *g.Active {
		if activePlayer.Role.Team == roles.MafiaTeam {
			mafiaTeamCounter++
		}
	}

	if mafiaTeamCounter > 1 {
		g.RUnlock()
		return
	}
	p.Role = roles.Mafia
	safeSendErrSignal(ch, g.RoleChannels[strings.ToLower(roles.Don.Name)].RemoveUser(p.Tag))
	safeSendErrSignal(ch, g.RoleChannels[strings.ToLower(roles.Mafia.Name)].AddPlayer(p.Tag))

	f := g.messenger.f
	g.RUnlock()
	var message string
	message = f.Bold("Hello, dear ") + f.Mention(p.ServerNick) + "." + f.LineSplitter()
	message += "You are the last player left alive from the mafia team, so you become mafia." + f.LineSplitter()
	message += f.Underline("Don't reveal yourself.")
	_, err := g.RoleChannels[strings.ToLower(roles.Mafia.Name)].Write([]byte(message))
	safeSendErrSignal(ch, err)
}
