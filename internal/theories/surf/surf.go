package surf

import (
	"github.com/SuperPythonic/SuperPythonic/internal/theories/conc"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
)

func Parse(text string) (*conc.Prog, parsing.State) {
	p := new(conc.Prog)
	return p, parsers.Parse(Prog(&p.Defs), text)
}

func Lowercase(dst **conc.Var) parsing.ParserFunc {
	return parsers.OnText(parsers.Lowercase, func(text string) { *dst = &conc.Var{Text: text} })
}

func TypeExpr(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			parsers.OnWord("unit", func() { *dst = new(conc.UnitType) }),
			parsers.OnWord("bool", func() { *dst = new(conc.BoolType) }),
			parsers.OnWord("float", func() { *dst = new(conc.FloatType) }),
			parsers.OnWord("int", func() { *dst = new(conc.IntType) }),
			parsers.OnWord("str", func() { *dst = new(conc.StrType) }),
		)(s)
	}
}

func ValueExpr(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			Lam(dst),
			parsers.OnWord("()", func() { *dst = new(conc.Unit) }),
			parsers.OnWord("False", func() { *dst = conc.Bool(false) }),
			parsers.OnWord("True", func() { *dst = conc.Bool(true) }),
			parsers.OnText(parsers.Float, func(text string) { *dst = &conc.Float{Text: text} }),
			parsers.OnText(parsers.Int, func(text string) { *dst = &conc.Int{Text: text} }),
			parsers.OnText(parsers.Str, func(text string) { *dst = conc.Str(text) }),

			parsers.OnText(IdRef, func(text string) { *dst = &conc.Ref{Name: &conc.Var{Text: text}} }),
			parsers.Seq(parsers.Word("("), ValueExpr(dst), parsers.Word(")")),
		)(s)
	}
}

func Lam(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		l := new(conc.Lam)
		if s = parsers.Seq(
			parsers.Word("lambda"),
			Lowercase(&l.Name),
			parsers.Word(":"),
			ValueExpr(&l.Body),
		)(s); !s.IsError() {
			*dst = l
		}
		return s
	}
}

var IdRef = parsers.Lowercase

func Prog(dst *[]conc.Def) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Seq(
			parsers.Start,
			parsers.Many(parsers.Newline),
			parsers.Many(parsers.Choice(FnDef(dst))),
			parsers.Many(parsers.Newline),
			parsers.End,
		)(s)
	}
}

func FnDef(dst *[]conc.Def) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		f := new(conc.Fn)
		return parsers.On(Fn(f), func() { *dst = append(*dst, f) })(s)
	}
}

func Fn(dst *conc.Fn) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		dst.R = new(conc.UnitType)
		return parsers.Seq(
			parsers.Word("def"),
			Lowercase(&dst.N),
			parsers.Choice(
				parsers.Seq(parsers.Word("("), parsers.Word(")")),
				parsers.Seq(
					parsers.Word("("),
					Params(&dst.Ps),
					parsers.Many(parsers.Seq(parsers.Word(","), Params(&dst.Ps))),
					parsers.Word(")"),
				),
			),
			parsers.Option(parsers.Seq(parsers.Word("->"), TypeExpr(&dst.R))),
			parsers.Word(":"),
			parsers.Entry,
			FnBody(&dst.Body),
			parsers.Exit,
		)(s)
	}
}

func Params(dst *[]*conc.Param) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		p := new(conc.Param)
		return parsers.On(
			parsers.Seq(Lowercase(&p.Name), parsers.Word(":"), TypeExpr(&p.Type)),
			func() { *dst = append(*dst, p) },
		)(s)
	}
}

func FnBody(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			FnBodyLet(dst),
			parsers.Seq(parsers.Word("return"), parsers.Option(ValueExpr(dst))),
		)(s)
	}
}

func FnBodyLet(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		l := new(conc.Let)
		if s = parsers.Seq(
			Lowercase(&l.Name),
			parsers.Option(parsers.Seq(parsers.Word(":"), TypeExpr(&l.Type))),
			parsers.Word("="),
			ValueExpr(&l.Value),
			parsers.Newline,
			parsers.Indent,
			FnBody(&l.Body),
		)(s); !s.IsError() {
			*dst = l
		}
		return s
	}
}
