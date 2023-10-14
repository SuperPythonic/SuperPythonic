package conc

import "github.com/SuperPythonic/SuperPythonic/pkg/theories"

type Prog struct{ defs []theories.Def }

func (p *Prog) Defs() []theories.Def { return p.defs }

type Fn struct {
	theories.Decl
}

type Param struct {
	name theories.Var
	typ  theories.Expr
}

func (p *Param) Local() theories.Var { return p.name }
func (p *Param) Type() theories.Expr { return p.typ }

type Var struct{ string }

func (v *Var) String() string { return v.string }
