package conc

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

type Var struct{ string }

func (v *Var) String() string { return v.string }

type Expr interface{ isExpr() }

type (
	Int struct{ Text string }
	Str struct{ Text string }
)

func (*Int) isExpr() {}
func (*Str) isExpr() {}

type (
	UnitType struct{}
	IntType  struct{}
	BoolType struct{}
)

func (*UnitType) isExpr() {}
func (*IntType) isExpr()  {}
func (*BoolType) isExpr() {}
