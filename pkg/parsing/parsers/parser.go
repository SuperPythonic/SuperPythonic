package parsers

import (
	"github.com/SuperPythonic/SuperPythonic/internal/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

var (
	NewState = parsers.NewState

	Start = parsers.Start
	End   = parsers.End

	Keyword   = parsers.Keyword
	Lowercase = parsers.Lowercase
	Uppercase = parsers.Uppercase
	CamelBack = parsers.CamelBack
	CamelCase = parsers.CamelCase

	Seq    = parsers.Seq
	Choice = parsers.Choice
	Many   = parsers.Many
)

func Parse(parser parsing.Parser, text string) parsing.State {
	return parser.Parse(NewState(text))
}
