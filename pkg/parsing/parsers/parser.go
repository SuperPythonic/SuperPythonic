package parsers

import (
	"github.com/SuperPythonic/SuperPythonic/internal/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

type Keyword = parsers.Keyword

var (
	NewState = parsers.NewState

	Start = parsers.Start
	End   = parsers.End

	Range = parsers.Range

	Lowercase = parsers.Lowercase
	Uppercase = parsers.Uppercase
	CamelBack = parsers.CamelBack
	CamelCase = parsers.CamelCase

	Int = parsers.Int

	Seq    = parsers.Seq
	Choice = parsers.Choice
	Many   = parsers.Many
	Option = parsers.Option
)

func Parse(parser parsing.Parser, text string) parsing.State {
	return parser.Parse(NewState(text))
}
