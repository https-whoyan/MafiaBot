package config

import (
	"strconv"
	"strings"

	"github.com/https-whoyan/MafiaBot/core/message/fmt"
	"github.com/https-whoyan/MafiaBot/core/roles"
)

func (cfg *RolesConfig) GetMessageAboutConfig(f fmt.FmtInterface) string {
	teamsMp := cfg.GetMapKeyByTeamValuesRoleCfg()
	teamsCount := len(teamsMp)

	var message string

	NL := f.LineSplitter()
	doubleNL := f.LineSplitter() + f.LineSplitter()
	tripleNL := doubleNL + f.LineSplitter()

	playersPlayedMsg := f.Bold("Players") + " count: " + f.Block(strconv.Itoa(cfg.PlayersCount))
	teamsPlayersMsg := f.Bold("Teams") + " count: " + f.Block(strconv.Itoa(teamsCount))
	rolesMsg := f.Bold("Roles") + " count: " + f.Block(strconv.Itoa(len(cfg.RolesMp)))

	hasFool := cfg.HasFool()

	message = playersPlayedMsg + NL + teamsPlayersMsg + NL + rolesMsg
	if hasFool {
		message += NL + f.Italic("(It's worth mentioning that the fool counts as a peaceful player, however "+
			"he plays as a separate team. "+"When checked by the detective, he is considered as a peaceful "+
			"player, but this is not entirely true.)")
	}
	message += tripleNL + f.InfoSplitter() + NL

	teamsMessages := make([]string, teamsCount)
	for _, team := range cfg.GetTeamsByCfg() {
		var teamMessage string
		playersInTeamsCount := cfg.GetPlayersCountByTeam(team)
		teamMessage = f.Bold("In "+roles.StringTeam[team]+" plays ") + f.Block(strconv.Itoa(playersInTeamsCount))
		teamMessage += NL

		rolesMessages := make([]string, len(teamsMp[team]))
		for _, roleCfg := range teamsMp[team] {
			roleMessage := f.Tab() + roleCfg.Role.Name + " " + f.Block(strconv.Itoa(roleCfg.Count))
			rolesMessages = append(rolesMessages, roleMessage)
		}
		teamsMessages = append(teamsMessages, strings.Join(rolesMessages, NL))
	}
	message += strings.Join(teamsMessages, doubleNL)
	return message
}
