package conc

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/theories"
)

func Parse(text string) (theories.Prog, parsing.State) {
	p := new(Prog)
	return p, parsers.Parse(p.Parse, text)
}

func (p *Prog) Parse(s parsing.State) parsing.State {
	return parsers.Seq(parsers.Start, parsers.Many(parsers.Choice(p.parseFn)), parsers.End)(s)
}

func (p *Prog) parseFn(s parsing.State) parsing.State {
	f := new(Fn)
	if s = f.Parse(s); !s.IsError() {
		p.defs = append(p.defs, f)
	}
	return s
}

func (f *Fn) Parse(s parsing.State) parsing.State {
	// TODO: Function body.
	return parsers.Seq(parsers.Word("def"), f.parseName, f.parseParams, parsers.Word(":"))(s)
}

func (f *Fn) parseName(s parsing.State) parsing.State {
	if s = parsers.Lowercase(s); !s.IsError() {
		f.Name = &Var{s.Text(s.Span())}
	}
	return s
}

func (f *Fn) parseParams(s parsing.State) parsing.State {
	return parsers.Choice(
		parsers.Seq(parsers.Word("("), parsers.Word(")")),
		parsers.Seq(
			parsers.Word("("),
			f.parseParam,
			parsers.Many(parsers.Seq(parsers.Word(","), f.parseParam)),
			parsers.Word(")"),
		),
	)(s)
}

func (f *Fn) parseParam(s parsing.State) parsing.State {
	p := new(Param)
	if s = p.Parse(s); !s.IsError() {
		f.Params = append(f.Params, p)
	}
	return s
}

func (p *Param) Parse(s parsing.State) parsing.State {
	return parsers.Seq(
		p.parseName,
		parsers.Word(":"),
		parsers.Choice(
			p.parseIntType,
			p.parseBoolType,
		),
	)(s)
}

func (p *Param) parseName(s parsing.State) parsing.State {
	if s = parsers.Lowercase(s); !s.IsError() {
		p.name = &Var{s.Text(s.Span())}
	}
	return s
}

func (p *Param) parseIntType(s parsing.State) parsing.State {
	if s = parsers.Word("int")(s); !s.IsError() {
		p.typ = new(theories.IntType)
	}
	return s
}

func (p *Param) parseBoolType(s parsing.State) parsing.State {
	if s = parsers.Word("bool")(s); !s.IsError() {
		p.typ = new(theories.BoolType)
	}
	return s
}
