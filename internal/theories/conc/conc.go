package conc

import (
	"fmt"

	"github.com/SuperPythonic/SuperPythonic/internal/theories"
)

type (
	Prog interface {
		fmt.Stringer

		Defs() []Def
	}

	Def interface {
		fmt.Stringer

		Name() theories.Var
		Params() []Param
		Ret() Expr
		Body() Expr
	}

	Param interface {
		Name() theories.Var
		Type() Expr
	}

	Expr interface {
		fmt.Stringer
	}
)
