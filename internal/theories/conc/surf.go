package conc

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/theories"
)

func Lowercase(dst *theories.Var) parsing.ParserFunc {
	return parsers.OnText(parsers.Lowercase, func(text string) { *dst = &Var{text} })
}

func Type(dst *theories.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			parsers.OnWord("int", func() { *dst = new(theories.IntType) }),
			parsers.OnWord("bool", func() { *dst = new(theories.BoolType) }),
			parsers.OnWord("unit", func() { *dst = new(theories.UnitType) }),
		)(s)
	}
}

func Parse(text string) (theories.Prog, parsing.State) {
	p := new(Prog)
	return p, parsers.Parse(p.Parse, text)
}

func (p *Prog) Parse(s parsing.State) parsing.State {
	return parsers.Seq(parsers.Start, parsers.Many(parsers.Choice(p.parseFn)), parsers.End)(s)
}

func (p *Prog) parseFn(s parsing.State) parsing.State {
	f := new(Fn)
	return parsers.On(f.Parse, func() { p.defs = append(p.defs, f) })(s)
}

func (f *Fn) Parse(s parsing.State) parsing.State {
	f.R = new(theories.UnitType)
	return parsers.Seq(
		parsers.Word("def"),
		Lowercase(&f.N),
		f.parseParams,
		parsers.Option(parsers.Seq(parsers.Word("->"), Type(&f.R))),
		parsers.Word(":"),
		// TODO: Function body.
	)(s)
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
	return parsers.On(p.Parse, func() { f.P = append(f.P, p) })(s)
}

func (p *Param) Parse(s parsing.State) parsing.State {
	return parsers.Seq(Lowercase(&p.name), parsers.Word(":"), Type(&p.typ))(s)
}
