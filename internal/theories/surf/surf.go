package surf

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
)

func Parse(text string) parsing.State {
	return parsers.Parse(Prog(), text)
}

func Prog() parsing.Parser {
	return parsers.Seq(
		parsers.Start(),
		parsers.Many(Def()),
		parsers.End(),
	)
}

func Def() parsing.Parser {
	return parsers.Seq(
		parsers.Keyword("def"),
		parsers.Lowercase(),
		parsers.Choice(
			parsers.Seq(
				parsers.Keyword("("),
				parsers.Keyword(")"),
			),
			parsers.Seq(
				parsers.Keyword("("),
				parsers.Lowercase(),
				parsers.Many(
					parsers.Seq(
						parsers.Keyword(","),
						parsers.Lowercase(),
					),
				),
				parsers.Keyword(")"),
			),
		),
		parsers.Keyword(":"),
	)
}
