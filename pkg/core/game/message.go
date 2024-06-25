package game

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/https-whoyan/MafiaBot/core/fmt"
)

// !!!!!!!!!!!!!!!!!!
// The use of all functions here assumes that you g.Init() has already been
// called and successfully executed without errors.
// !!!!!!!!!!!!!!!!!!

var (
	playersCalling = []string{"poopsies", "players", "ladies and gentlemen"}
	playerCalling  = []string{"poops", "ancient", "modern", "member"}
)

var (
	sInt = func(s int) string { return strconv.Itoa(s) }
	sCap = func(s string) string {
		if len(s) <= 1 {
			return s
		}
		return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
	}
	getRandomPlayersCalling = func() string { return playersCalling[rand.Intn(len(playersCalling))] }
	getRandomPlayerCalling  = func() string { return playerCalling[rand.Intn(len(playerCalling))] }
)

// ____________
// Init
// ____________

func (g *Game) GetStartMessage() string {
	var message string

	f := g.fmtEr

	nl := f.LineSplitter()
	dNl := nl + nl
	iL := f.InfoSplitter()

	message = f.Bold("Have a good day, " + getRandomPlayersCalling() + "!")
	message += dNl
	message += fmt.BoldUnderline(f, "Today, our players:") + nl

	var aboutIDMessages []string
	activePlayers := g.Active
	sort.Slice(activePlayers, func(i, j int) bool {
		return activePlayers[i].ID < activePlayers[j].ID
	})
	for _, player := range activePlayers {
		messageAboutPlayerID := f.Tab() + f.Bold(sCap(getRandomPlayerCalling())) + " " + f.Mention(player.ServerNick)
		messageAboutPlayerID += " with ID in game " + f.Block(sInt(player.ID))

		aboutIDMessages = append(aboutIDMessages, messageAboutPlayerID)
	}
	message += strings.Join(aboutIDMessages, nl)

	if len(g.Spectators) != 0 {
		message += dNl
		message += "From behind the scenes to support us: "
		var spectatorMentions []string

		for _, spectator := range g.Spectators {
			spectatorMentions = append(spectatorMentions, f.Mention(spectator.ServerNick))
		}
		message += strings.Join(spectatorMentions, ", ")
	}

	message += nl + iL + nl
	message += fmt.ItalicUnderline(f, "Selected game configuration:") + nl
	message += g.RolesConfig.GetMessageAboutConfig(f)
	message += nl + iL + nl

	// Redo it if it false!!!!
	message += "A private message has been sent to each of you, you can find your ID and role in it."
	message += nl
	message += "Also, " + f.Bold("if you have an active night role, you have been added to special channels, where "+
		"you can send commands to the bot anonymously")
	if len(g.Spectators) != 0 {
		message += f.Italic(" (but there's no hiding from observers))))")
	}
	message += "." + nl
	if g.renameMode != NotRenameModeMode {
		message += nl
		message += "Also, all participants have been prefixed with their IDs to make it more convenient for you."
	}
	message += dNl + iL + nl
	message += f.Bold("Welcome, welcome, welcome... Happy hunger games and the odds be ever in your favor! ") +
		f.Italic("(Or just have a good game!) 🍀")

	return message
}
