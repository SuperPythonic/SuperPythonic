package theories

import "github.com/SuperPythonic/SuperPythonic/internal/theories/conc"

type (
	Prog = conc.Prog

	Def = conc.Def
	Fn  = conc.Fn

	Param = conc.Param

	Var = conc.Var

	Expr = conc.Expr

	UnitType struct{}
	IntType  struct{}
	BoolType struct{}
)

func (*UnitType) isExpr() {}
func (*IntType) isExpr()  {}
func (*BoolType) isExpr() {}
