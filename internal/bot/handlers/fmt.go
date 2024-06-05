package bot

// _____________
// Text Style
// _____________

func Bold(s string) string {
	return "**" + s + "**"
}

func Italic(s string) string {
	return "_" + s + "_"
}

func Emphasized(s string) string {
	return "__" + s + "__"
}

func CodeBlock(language, text string) string {
	return "```" + language + text + "```"
}

// __________
// Stickers
// __________

var (
	RegistrationPlayerSticker    = ":grin:"
	RegistrationSpectatorSticker = ":smiling_imp:"
	LuckySticker                 = ":four_leaf_clover:"
)

// MappedStickersUnicode save stickers Unicode:
var MappedStickersUnicode = map[string]string{
	RegistrationPlayerSticker:    "U+1F601",
	RegistrationSpectatorSticker: "U+1F608",
}

func GetUnicodeBySticker(sticker string) string {
	return MappedStickersUnicode[sticker]
}
