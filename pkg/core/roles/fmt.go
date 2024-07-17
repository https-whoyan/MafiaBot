package roles

import (
	"github.com/https-whoyan/MafiaBot/core/fmt"
	"strings"
)

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

func GetEmojiByName(name string) string {
	return MappedEmoji[name]
}

func FixDescription(description string) string {
	words := strings.Split(description, " ")
	return strings.Join(words, " ")
}

// _____________________________________________________________________
// Beautiful presentations of roles to display information about them.
// _____________________________________________________________________

func GetDefinitionOfRole(f fmt.FmtInterface, roleName string) string {

	role := MappedRoles[roleName]
	var message string

	name := f.Block(role.Name)
	team := f.Bold("Team: ") + StringTeam[role.Team]
	description := FixDescription(role.Description)
	message = name + f.LineSplitter() + f.LineSplitter() + team + f.LineSplitter() + description
	return message
}

func GetDefinitionsOfAllRoles(f fmt.FmtInterface, maxBytesLenInMessage int) (messages []string) {
	allRoles := GetAllSortedRoles()
	var allDescriptions []string

	bytesCounter := 0
	rolesCounter := 0

	infoSptr := f.LineSplitter() + f.InfoSplitter() + f.LineSplitter()

	for _, role := range allRoles {
		roleDescription := GetDefinitionOfRole(f, role.Name)
		// To avoid for looping
		nextTrueBytesMessage := len(roleDescription) + bytesCounter + len(infoSptr)*(rolesCounter-1)
		if nextTrueBytesMessage >= maxBytesLenInMessage {
			messages = append(messages, strings.Join(allDescriptions, infoSptr))
			allDescriptions = []string{}
			bytesCounter = 0
			rolesCounter = 0
		}
		bytesCounter += len(roleDescription)
		allDescriptions = append(allDescriptions, roleDescription)
	}
	if len(allDescriptions) > 0 {
		messages = append(messages, strings.Join(allDescriptions, infoSptr))
	}
	return
}
