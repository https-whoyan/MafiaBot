package message

import (
	"strconv"

	"github.com/https-whoyan/MafiaBot/core/fmt"
	playerPack "github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"
)

// GetStartPlayerDefinition is used to receive a private message at the beginning of the player.
// Add and modify to your taste.
func GetStartPlayerDefinition(p *playerPack.Player, f fmt.FmtInterface) string {
	message := f.Bold("Hello, "+f.Mention(p.ServerNick)+"!") + f.LineSplitter()
	message += "Today, in game you play in " + f.Bold(roles.StringTeam[p.Role.Team]) +
		" , your role is " + f.Block(p.Role.Name) + " and your ID is " + f.Block(strconv.Itoa(int(p.ID)))
	message += f.LineSplitter() + f.InfoSplitter() + f.LineSplitter()
	message += f.Italic("Let me remind you of your role description.") + f.LineSplitter()
	message += roles.FixDescription(p.Role.Description)
	message += f.LineSplitter() + f.LineSplitter() + f.Bold("Have a good game!")

	return message
}
