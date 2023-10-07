package conc

import (
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
	"github.com/SuperPythonic/SuperPythonic/pkg/theories"
)

func Parse(text string) (theories.Prog, parsing.State) {
	prog := Prog()
	return prog, parsers.Parse(prog, text)
}

func (p *prog) Parse(s parsing.State) parsing.State {
	p.defs = new(defs)
	return parsers.Seq(parsers.Start(), parsers.Many(parsers.Choice(p.defs)), parsers.End()).Parse(s)
}

func (p *defs) Parse(s parsing.State) parsing.State {
	def := Fn()
	if s = def.Parse(s); !parsing.IsError(s) {
		p.defs = append(p.defs, def)
	}
	return s
}

func (p *fn) Parse(s parsing.State) parsing.State {
	p.name = LowercaseVar()
	p.params = Params()
	// TODO: Function body.
	return parsers.Seq(parsers.Keyword("def"), p.name, p.params, parsers.Keyword(":")).Parse(s)
}

func (p *lcVar) Parse(s parsing.State) parsing.State {
	if s = parsers.Lowercase().Parse(s); !parsing.IsError(s) {
		p.name = s.Text(s.Cur())
	}
	return s
}

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

func (p *paramList) Parse(s parsing.State) parsing.State {
	name := LowercaseVar()
	if s = name.Parse(s); !parsing.IsError(s) {
		p.names = append(p.names, name)
	}
	return s
}
