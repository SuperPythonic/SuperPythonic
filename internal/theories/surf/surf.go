package surf

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
)

func Parse(text string) parsing.State {
	return parsers.Parse(Prog(), text)
}

func Prog() parsing.Parser {
	return parsers.Seq(parsers.Start(), parsers.Many(parsers.Choice(Fn())), parsers.End())
}

type fn struct {
	name *fnName
}

func (f *fn) Parse(s parsing.State) parsing.State {
	return parsers.
		Seq(
			parsers.Keyword("def"),
			f.name,
			parsers.Choice(
				parsers.Seq(parsers.Keyword("("), parsers.Keyword(")")),
				parsers.Seq(
					parsers.Keyword("("),
					parsers.Lowercase(),
					parsers.Many(parsers.Seq(parsers.Keyword(","), parsers.Lowercase())),
					parsers.Keyword(")"),
				),
			),
			parsers.Keyword(":"),
			// TODO: Function body.
		).
		Parse(s)
}

type fnName struct{ name string }

func (r *fnName) Parse(s parsing.State) parsing.State {
	if s = parsers.Lowercase().Parse(s); !parsing.IsError(s) {
		r.name = s.Text(s.Cur())
	}
	return s
}

func Fn() parsing.Parser { return &fn{new(fnName)} }
