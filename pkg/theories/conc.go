package theories

import "github.com/SuperPythonic/SuperPythonic/pkg/parsing"

type Prog interface {
	parsing.Parser

	At(idx int) Def
}

type Def interface {
	parsing.Parser

	Name() Var
	Params() Params
}

type Params interface {
	parsing.Parser

	At(idx int) Var
}

type Var interface {
	parsing.Parser
}
