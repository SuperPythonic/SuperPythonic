package parsers

import (
	"unicode"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

type opts struct{}

func (*opts) IsSpace(r rune) bool     { return unicode.IsSpace(r) }
func (*opts) IsNewline(r rune) bool   { return r == '\n' }
func (*opts) IndentWord() string      { return "\t" }
func (*opts) OnNewline(parsing.State) {}
