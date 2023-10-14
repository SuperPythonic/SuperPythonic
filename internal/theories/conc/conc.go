package conc

import "github.com/SuperPythonic/SuperPythonic/pkg/theories"

type prog struct{ defs []theories.Def }

func (p *prog) Defs() []theories.Def { return p.defs }

type fn struct {
	name   theories.Var
	params []theories.Var
}

func (p *fn) Name() theories.Var     { return p.name }
func (p *fn) Params() []theories.Var { return p.params }

type Var struct{ string }

func (v *Var) String() string { return v.string }
