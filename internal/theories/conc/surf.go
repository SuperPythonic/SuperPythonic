package conc

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
)

func Lowercase(dst **Var) parsing.ParserFunc {
	return parsers.OnText(parsers.Lowercase, func(text string) { *dst = &Var{text} })
}

func TypeExpr(dst *Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			parsers.OnWord("unit", func() { *dst = new(UnitType) }),
			parsers.OnWord("bool", func() { *dst = new(BoolType) }),
			parsers.OnWord("int", func() { *dst = new(IntType) }),
			parsers.OnWord("str", func() { *dst = new(StrType) }),
		)(s)
	}
}

func ValueExpr(dst *Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			parsers.OnText(parsers.Lowercase, func(text string) { *dst = &Ref{&Var{text}} }),
			parsers.OnWord("()", func() { *dst = new(Unit) }),
			parsers.OnWord("False", func() { *dst = Bool(false) }),
			parsers.OnWord("True", func() { *dst = Bool(true) }),
			parsers.OnText(parsers.Int, func(text string) { *dst = &Int{text} }),
			parsers.OnText(parsers.Str, func(text string) { *dst = &Str{text} }),
		)(s)
	}
}

func Parse(text string) (*Prog, parsing.State) {
	p := new(Prog)
	return p, parsers.Parse(p.Parse, text)
}

func (p *Prog) Parse(s parsing.State) parsing.State {
	return parsers.Seq(
		parsers.Start,
		parsers.Many(parsers.Newline),
		parsers.Many(parsers.Choice(p.parseFn)),
		parsers.Many(parsers.Newline),
		parsers.End,
	)(s)
}

func (p *Prog) parseFn(s parsing.State) parsing.State {
	f := new(Fn)
	return parsers.On(f.Parse, func() { p.Defs = append(p.Defs, f) })(s)
}

func (f *Fn) Parse(s parsing.State) parsing.State {
	f.R = new(UnitType)
	return parsers.Seq(
		parsers.Word("def"),
		Lowercase(&f.N),
		f.parseParams,
		parsers.Option(parsers.Seq(parsers.Word("->"), TypeExpr(&f.R))),
		parsers.Word(":"),
		parsers.Entry,
		f.parseBody,
		parsers.Exit,
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
	return parsers.On(p.Parse, func() { f.Ps = append(f.Ps, p) })(s)
}

func (f *Fn) parseBody(s parsing.State) parsing.State {
	f.Body = new(Unit)
	return parsers.Choice(f.parseBodyLet, f.parseBodyRet)(s)
}

func (f *Fn) parseBodyLet(s parsing.State) parsing.State {
	l := new(Let)
	return parsers.Seq(
		Lowercase(&l.Name),
		parsers.Option(parsers.Seq(parsers.Word(":"), TypeExpr(&l.Type))),
		parsers.Word("="),
		ValueExpr(&l.Value),
		parsers.Newline,
		parsers.Indent,
		f.parseBody,
	)(s)
}

func (f *Fn) parseBodyRet(s parsing.State) parsing.State {
	return parsers.Seq(
		parsers.Word("return"),
		parsers.Option(ValueExpr(&f.Body)),
	)(s)
}

func (p *Param) Parse(s parsing.State) parsing.State {
	return parsers.Seq(Lowercase(&p.Name), parsers.Word(":"), TypeExpr(&p.Type))(s)
}
