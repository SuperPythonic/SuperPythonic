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

func Type(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(Pi(dst), PrimaryType(dst))(s)
	}
}

func Pi(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		f := &conc.FuncType{Name: conc.Unbound()}
		if s = parsers.Seq(PrimaryType(&f.Type), parsers.Word("->"), Type(&f.Body))(s); !s.IsError() {
			*dst = f
		}
		return s
	}
}

func PrimaryType(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			parsers.OnWord("unit", func() { *dst = new(conc.UnitType) }),
			parsers.OnWord("bool", func() { *dst = new(conc.BoolType) }),
			parsers.OnWord("float", func() { *dst = new(conc.FloatType) }),
			parsers.OnWord("int", func() { *dst = new(conc.IntType) }),
			parsers.OnWord("str", func() { *dst = new(conc.StrType) }),

			parsers.OnText(IdRef, func(text string) { *dst = conc.NewRef(text) }),
		)(s)
	}
}

func Value(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			Lam(dst),
			parsers.OnWord("()", func() { *dst = new(conc.Unit) }),
			parsers.OnWord("False", func() { *dst = conc.Bool(false) }),
			parsers.OnWord("True", func() { *dst = conc.Bool(true) }),
			parsers.OnText(parsers.Float, func(text string) { *dst = &conc.Float{Text: text} }),
			parsers.OnText(parsers.Int, func(text string) { *dst = &conc.Int{Text: text} }),
			parsers.OnText(parsers.Str, func(text string) { *dst = conc.Str(text) }),

			parsers.OnText(IdRef, func(text string) { *dst = conc.NewRef(text) }),
			parsers.Seq(parsers.Word("("), Value(dst), parsers.Word(")")),
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
			Value(&l.Body),
		)(s); !s.IsError() {
			*dst = l
		}
		return s
	}
}

var IdRef = parsers.Choice(parsers.Lowercase, parsers.CamelCase)

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
			parsers.Option(parsers.Seq(parsers.Word("->"), Type(&dst.R))),
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
			parsers.Seq(Lowercase(&p.Name), parsers.Word(":"), Type(&p.Type)),
			func() { *dst = append(*dst, p) },
		)(s)
	}
}

func FnBody(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			FnBodyLet(dst),
			parsers.Seq(parsers.Word("return"), parsers.Option(Value(dst))),
		)(s)
	}
}

func FnBodyLet(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		l := new(conc.Let)
		if s = parsers.Seq(
			Lowercase(&l.Name),
			parsers.Option(parsers.Seq(parsers.Word(":"), Type(&l.Type))),
			parsers.Word("="),
			Value(&l.Value),
			parsers.Newline,
			parsers.Indent,
			FnBody(&l.Body),
		)(s); !s.IsError() {
			*dst = l
		}
		return s
	}
}
