package conc

import "github.com/SuperPythonic/SuperPythonic/pkg/parsing"

type Prog struct{ Defs []Def }

type (
	Def interface {
		Name() *Var
		Params() []*Param
		Ret() Expr

		isDef()
	}

	Decl struct {
		N  *Var
		Ps []*Param
		R  Expr
	}
)

func (d *Decl) Name() *Var       { return d.N }
func (d *Decl) Params() []*Param { return d.Ps }
func (d *Decl) Ret() Expr        { return d.R }
func (*Decl) isDef()             {}

type Fn struct {
	Decl
	Body Expr
}

type Param struct {
	Name *Var
	Type Expr
}

type Var struct{ Text string }

func Unbound() *Var           { return &Var{"_"} }
func (v *Var) String() string { return v.Text }

type Expr interface{ parsing.Spanner }

type (
	Unit struct{ parsing.Span }
	Bool struct {
		parsing.Span
		Value bool
	}
	Int   struct{ parsing.Span }
	Float struct{ parsing.Span }
	Str   struct{ parsing.Span }
	Lam   struct {
		parsing.Span
		Name *Var
		Body Expr
	}
	Ref struct {
		parsing.Span
		Name *Var
	}
)

func NewRef(span parsing.Span, text string) *Ref { return &Ref{Span: span, Name: &Var{text}} }

type (
	UnitType  struct{ parsing.Span }
	BoolType  struct{ parsing.Span }
	FloatType struct{ parsing.Span }
	IntType   struct{ parsing.Span }
	StrType   struct{ parsing.Span }
	FnType    struct {
		parsing.Span
		Params []*Param
		Body   Expr
	}
)

type (
	App struct {
		parsing.Span
		Fn   Expr
		Args []Expr
	}
	If struct {
		parsing.Span
		Test, Then, Else Expr
	}
	Let struct {
		parsing.Span
		Name              *Var
		Type, Value, Body Expr
	}
	UnitLet struct {
		parsing.Span
		Value, Body Expr
	}
)
