package theories

import "fmt"

type Prog interface {
	Defs() []Def
}

type (
	Def interface {
		Name() Var
		Params() []Param
		Ret() Expr

		isDef()
	}

	Decl struct {
		N Var
		P []Param
		R Expr
	}
)

func (d *Decl) Name() Var       { return d.N }
func (d *Decl) Params() []Param { return d.P }
func (d *Decl) Ret() Expr       { return d.R }
func (*Decl) isDef()            {}

type Param interface {
	Name() Var
	Type() Expr
}

type Var interface{ fmt.Stringer }

type Expr interface{ isExpr() }

type (
	UnitType struct{}
	IntType  struct{}
	BoolType struct{}
)

func (*UnitType) isExpr() {}
func (*IntType) isExpr()  {}
func (*BoolType) isExpr() {}
