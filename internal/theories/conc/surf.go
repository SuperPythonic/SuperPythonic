package conc

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/theories"
)

func Lowercase(dst *theories.Var) parsing.ParserFunc {
	return parsers.On(parsers.Lowercase, func(s parsing.State) { *dst = &Var{s.Text()} })
}

func Type(dst *theories.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			intType(dst),
			boolType(dst),
		)(s)
	}
}

func intType(dst *theories.Expr) parsing.ParserFunc {
	return parsers.On(parsers.Word("int"), func(parsing.State) { *dst = new(theories.IntType) })
}

func boolType(dst *theories.Expr) parsing.ParserFunc {
	return parsers.On(parsers.Word("bool"), func(parsing.State) { *dst = new(theories.BoolType) })
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
	return parsers.On(f.Parse, func(parsing.State) { p.defs = append(p.defs, f) })(s)
}

func (f *Fn) Parse(s parsing.State) parsing.State {
	// TODO: Function body.
	return parsers.Seq(parsers.Word("def"), Lowercase(&f.Name), f.parseParams, parsers.Word(":"))(s)
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
	return parsers.On(p.Parse, func(parsing.State) { f.Params = append(f.Params, p) })(s)
}

func (p *Param) Parse(s parsing.State) parsing.State {
	return parsers.Seq(Lowercase(&p.name), parsers.Word(":"), Type(&p.typ))(s)
}
