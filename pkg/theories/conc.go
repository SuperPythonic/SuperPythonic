package theories

import "fmt"

type Prog interface {
	Defs() []Def
}

type (
	Def interface {
		Global() Var
		Tele() []Param

		isDef()
	}

	Decl struct {
		Name   Var
		Params []Param
	}
)

func (d *Decl) Global() Var   { return d.Name }
func (d *Decl) Tele() []Param { return d.Params }
func (*Decl) isDef()          {}

type Param interface {
	Local() Var
	Type() Expr
}

type Var interface{ fmt.Stringer }

type (
	Expr interface{ isExpr() }

	IntType  struct{}
	BoolType struct{}
)

func (*IntType) isExpr()  {}
func (*BoolType) isExpr() {}
