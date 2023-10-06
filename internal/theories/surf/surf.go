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

func Def() parsing.Parser {
	return parsers.Seq(
		parsers.Keyword("def"),
		parsers.Lowercase(),
	)
}

func Class() parsing.Parser {
	return parsers.Seq(
		parsers.Keyword("class"),
		parsers.Uppercase(),
	)
}
