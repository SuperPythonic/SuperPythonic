package theories

import "fmt"

type Prog interface {
	Defs() []Def
}

type Def interface {
	Name() Var
	Params() []Var
}

type Var interface {
	fmt.Stringer
}
