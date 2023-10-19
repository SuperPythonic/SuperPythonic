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

type Var struct{ Text string }

func (v *Var) String() string { return v.Text }

type Expr interface{ isExpr() }

type (
	Ref  struct{ Name *Var }
	Unit struct{}
	Bool bool
	Int  struct{ Text string }
	Str  string
	Let  struct {
		Name              *Var
		Type, Value, Body Expr
	}
)

func (*Ref) isExpr()  {}
func (*Unit) isExpr() {}
func (Bool) isExpr()  {}
func (*Int) isExpr()  {}
func (Str) isExpr()   {}
func (*Let) isExpr()  {}

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
