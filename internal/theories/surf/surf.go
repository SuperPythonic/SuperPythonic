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
	return func(s parsing.State) parsing.State {
		return parsers.OnText(parsers.Lowercase, func(text string) { *dst = &conc.Var{Text: text} })(s)
	}
}

func Type(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			FnType(dst),
			PrimaryType(dst),
			parsers.Seq(parsers.Word("("), Type(dst), parsers.Word(")")),
		)(s)
	}
}

func FnType(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		f := new(conc.FnType)
		return parsers.On(parsers.Seq(
			FnTypeParams(&f.Params),
			parsers.Word("->"),
			Type(&f.Body),
		), func() { *dst = f })(s)
	}
}

func FnTypeParams(dst *[]*conc.Param) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			parsers.Word("()"),
			FnTypeParamOne(dst),
			parsers.Seq(
				parsers.Word("("),
				FnTypeParamMore(dst),
				parsers.Many(parsers.Seq(parsers.Word(","), FnTypeParamMore(dst))),
				parsers.Word(")"),
			),
		)(s)
	}
}

func FnTypeParamOne(dst *[]*conc.Param) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		p := &conc.Param{Name: conc.Unbound()}
		return parsers.On(PrimaryType(&p.Type), func() { *dst = append(*dst, p) })(s)
	}
}

func FnTypeParamMore(dst *[]*conc.Param) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		p := &conc.Param{Name: conc.Unbound()}
		return parsers.On(Type(&p.Type), func() { *dst = append(*dst, p) })(s)
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
			IdRef(dst),
		)(s)
	}
}

func Value(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			If(dst),
			InlineIf(dst),
			Lam(dst),
			PrimaryValue(dst),
		)(s)
	}
}

func InlineIf(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		i := new(conc.If)
		return parsers.On(
			parsers.Seq(
				PrimaryValue(&i.Then),
				parsers.Word("if"),
				Value(&i.Test),
				parsers.Word("else"),
				Value(&i.Else),
			),
			func() { *dst = i },
		)(s)
	}
}

func If(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		i := new(conc.If)
		return parsers.On(
			parsers.Seq(
				parsers.Word("if"),
				Value(&i.Test),
				parsers.Word(":"),
				parsers.Entry,
				Value(&i.Then),
				parsers.Exit,
				parsers.Indent,
				parsers.Word("else"),
				parsers.Word(":"),
				parsers.Entry,
				Value(&i.Else),
				parsers.Exit,
			),
			func() { *dst = i },
		)(s)
	}
}

func Lam(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		l := &conc.Lam{Name: conc.Unbound()}
		return parsers.On(
			parsers.Seq(
				parsers.Word("lambda"),
				parsers.Option(Lowercase(&l.Name)),
				parsers.Word(":"),
				Value(&l.Body),
			),
			func() { *dst = l },
		)(s)
	}
}

func PrimaryValue(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Choice(
			parsers.OnWord("()", func() { *dst = new(conc.Unit) }),
			parsers.OnWord("False", func() { *dst = conc.Bool(false) }),
			parsers.OnWord("True", func() { *dst = conc.Bool(true) }),
			parsers.OnText(parsers.Float, func(text string) { *dst = &conc.Float{Text: text} }),
			parsers.OnText(parsers.Int, func(text string) { *dst = &conc.Int{Text: text} }),
			parsers.OnText(parsers.Str, func(text string) { *dst = conc.Str(text) }),
			IdRef(dst),
			parsers.Seq(parsers.Word("("), Value(dst), parsers.Word(")")),
		)(s)
	}
}

func IdRef(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.OnText(
			parsers.Choice(parsers.Lowercase, parsers.CamelCase),
			func(text string) { *dst = conc.NewRef(text) },
		)(s)
	}
}

func Prog(dst *[]conc.Def) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		return parsers.Seq(parsers.Start, parsers.Many(parsers.Choice(parsers.Newline, FnDef(dst))), parsers.End)(s)
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
			FnBodyUnitLet(dst),
			parsers.Seq(parsers.Word("return"), parsers.Option(Value(dst))),
		)(s)
	}
}

func FnBodyLet(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		l := new(conc.Let)
		return parsers.On(
			parsers.Seq(
				Lowercase(&l.Name),
				parsers.Option(parsers.Seq(parsers.Word(":"), Type(&l.Type))),
				parsers.Word("="),
				Value(&l.Value),
				parsers.Newline,
				parsers.Indent,
				FnBody(&l.Body),
			),
			func() { *dst = l },
		)(s)
	}
}

func FnBodyUnitLet(dst *conc.Expr) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		l := new(conc.UnitLet)
		return parsers.On(
			parsers.Seq(Value(&l.Value), parsers.Newline, parsers.Indent, FnBody(&l.Body)),
			func() { *dst = l },
		)(s)
	}
}
