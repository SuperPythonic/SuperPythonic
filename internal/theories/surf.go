package theories

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/theories"
)

func Parse(text string) parsing.State {
	return parsers.Parse(Prog(), text)
}

func Prog() parsing.Parser {
	return parsers.Seq(parsers.Start(), parsers.Many(parsers.Choice(Fn())), parsers.End())
}

type fn struct {
	name   *lcVar
	params *params
}

func (p *fn) Parse(s parsing.State) parsing.State {
	s = parsers.
		Seq(
			parsers.Keyword("def"),
			p.name,
			p.params,
			parsers.Keyword(":"),
			// TODO: Function body.
		).
		Parse(s)
	return s
}

func (p *fn) Name() theories.Var { return p.name }

func (p *fn) Params() theories.Params { return p.params }

func Fn() parsing.Parser { return &fn{name: new(lcVar), params: new(params)} }

type lcVar struct{ name string }

func (p *lcVar) Parse(s parsing.State) parsing.State {
	if s = parsers.Lowercase().Parse(s); !parsing.IsError(s) {
		p.name = s.Text(s.Cur())
	}
	return s
}

type (
	params struct {
		list *paramList
	}

	paramList struct {
		names []*lcVar
	}
)

func (p *params) Parse(s parsing.State) parsing.State {
	p.list = new(paramList)
	return parsers.Choice(
		parsers.Seq(parsers.Keyword("("), parsers.Keyword(")")),
		parsers.Seq(
			parsers.Keyword("("),
			p.list,
			parsers.Many(parsers.Seq(parsers.Keyword(","), p.list)),
			parsers.Keyword(")"),
		),
	).Parse(s)
}

func (p *params) At(idx int) theories.Var { return p.list.names[idx] }

func (p *paramList) Parse(s parsing.State) parsing.State {
	name := new(lcVar)
	if s = name.Parse(s); !parsing.IsError(s) {
		p.names = append(p.names, name)
	}
	return s
}
