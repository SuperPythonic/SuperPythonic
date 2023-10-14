package conc

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/theories"
)

func Parse(text string) (theories.Prog, parsing.State) {
	p := new(prog)
	return p, parsers.Parse(p.Parse, text)
}

func (p *prog) Parse(s parsing.State) parsing.State {
	return parsers.Seq(parsers.Start, parsers.Many(parsers.Choice(p.parseFn)), parsers.End)(s)
}

func (p *prog) parseFn(s parsing.State) parsing.State {
	f := new(fn)
	if s = f.Parse(s); !s.IsError() {
		p.defs = append(p.defs, f)
	}
	return s
}

func (p *fn) Parse(s parsing.State) parsing.State {
	// TODO: Function body.
	return parsers.Seq(parsers.Word("def"), p.parseName, p.parseParams, parsers.Word(":"))(s)
}

func (p *fn) parseName(s parsing.State) parsing.State {
	if s = parsers.Lowercase(s); !s.IsError() {
		p.name = &Var{s.Text(s.Span())}
	}
	return s
}

func (p *fn) parseParams(s parsing.State) parsing.State {
	return parsers.Choice(
		parsers.Seq(parsers.Word("("), parsers.Word(")")),
		parsers.Seq(
			parsers.Word("("),
			p.parseParam,
			parsers.Many(parsers.Seq(parsers.Word(","), p.parseParam)),
			parsers.Word(")"),
		),
	)(s)
}

func (p *fn) parseParam(s parsing.State) parsing.State {
	if s = parsers.Lowercase(s); !s.IsError() {
		p.params = append(p.params, &Var{s.Text(s.Span())})
	}
	return s
}
