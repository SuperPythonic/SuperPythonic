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
	Unit struct{}
	Int  struct{ Text string }
	Str  struct{ Text string }
)

func (*Unit) isExpr() {}
func (*Int) isExpr()  {}
func (*Str) isExpr()  {}

type (
	UnitType struct{}
	BoolType struct{}
	IntType  struct{}
	StrType  struct{}
)

func (*UnitType) isExpr() {}
func (*BoolType) isExpr() {}
func (*IntType) isExpr()  {}
func (*StrType) isExpr()  {}
