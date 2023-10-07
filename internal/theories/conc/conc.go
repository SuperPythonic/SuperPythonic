package conc

import "github.com/SuperPythonic/SuperPythonic/pkg/theories"

type (
	prog struct {
		defs *defs
	}

	defs struct {
		defs []theories.Def
	}
)

func (p *prog) At(idx int) theories.Def { return p.defs.defs[idx] }

func Prog() theories.Prog { return new(prog) }

type fn struct {
	name   theories.Var
	params theories.Params
}

func (p *fn) Name() theories.Var { return p.name }

func (p *fn) Params() theories.Params { return p.params }

func Fn() theories.Def { return new(fn) }

type lcVar struct{ name string }

func LowercaseVar() theories.Var { return new(lcVar) }

type (
	params struct {
		list *paramList
	}

	paramList struct {
		names []theories.Var
	}
)

func (p *params) At(idx int) theories.Var { return p.list.names[idx] }

func Params() theories.Params { return new(params) }
