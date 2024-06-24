package models

type TestFMTer struct{}

func (f TestFMTer) Bold(s string) string       { return s }
func (f TestFMTer) Italic(s string) string     { return s }
func (f TestFMTer) Underline(s string) string  { return s }
func (f TestFMTer) Block(s string) string      { return s }
func (f TestFMTer) LineSplitter() string       { return "\n" }
func (f TestFMTer) InfoSplitter() string       { return "===" }
func (f TestFMTer) Tab() string                { return "\t" }
func (f TestFMTer) Mention(nick string) string { return nick }

var TestFMTInstance = TestFMTer{}
