package conc

import (
	"fmt"

	"github.com/SuperPythonic/SuperPythonic/internal/syntaxes"
)

type (
	Prog interface {
		fmt.Stringer

		Defs() []Def
	}

	Def interface {
		fmt.Stringer

		Name() syntaxes.Var
		Params() []Param
		Ret() Term
		Body() Term
	}

	Param interface {
		Name() syntaxes.Var
		Type() Term
	}

	Term interface {
		fmt.Stringer
	}
)
