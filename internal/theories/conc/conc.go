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

func Unbound() *Var           { return &Var{"_"} }
func (v *Var) String() string { return v.Text }

type Expr interface{ isExpr() }

type (
	Unit  struct{}
	Bool  bool
	Int   struct{ Text string }
	Float struct{ Text string }
	Str   string
	Let   struct {
		Name              *Var
		Type, Value, Body Expr
	}
	Lam struct {
		Name *Var
		Body Expr
	}
	Ref struct{ Name *Var }
)

func NewRef(text string) *Ref { return &Ref{&Var{text}} }

func (*Unit) isExpr()  {}
func (Bool) isExpr()   {}
func (*Int) isExpr()   {}
func (*Float) isExpr() {}
func (Str) isExpr()    {}
func (*Let) isExpr()   {}
func (*Lam) isExpr()   {}
func (*Ref) isExpr()   {}

type (
	UnitType  struct{}
	BoolType  struct{}
	FloatType struct{}
	IntType   struct{}
	StrType   struct{}
	FnType    struct {
		Params []*Param
		Body   Expr
	}
)

func (*UnitType) isExpr()  {}
func (*BoolType) isExpr()  {}
func (*FloatType) isExpr() {}
func (*IntType) isExpr()   {}
func (*StrType) isExpr()   {}
func (*FnType) isExpr()    {}
