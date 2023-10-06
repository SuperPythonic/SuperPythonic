package surf

import (
	"github.com/SuperPythonic/SuperPythonic/internal/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

const (
	Def parsing.TokenKind = iota + 100
	Main
)

type prog struct{}

func Prog() parsing.Parser {
	return new(prog)
}

func (p *prog) Run(s parsing.State) parsing.State {
	return parsers.Seq(
		parsers.Start(),
		parsers.Keyword(Def, "def"),
		parsers.Keyword(Main, "main"),
		parsers.End(),
	).
		Run(s)
}
