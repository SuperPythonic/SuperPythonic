package surf

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
)

func Prog() parsing.Parser {
	return parsers.Seq(
		parsers.Start(),
		parsers.Choice(
			Def(),
			Class(),
		),
		parsers.End(),
	)
}

func Def() parsing.Parser { return parsers.Keyword("def") }

func Class() parsing.Parser { return parsers.Keyword("class") }
