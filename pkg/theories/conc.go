package theories

import "github.com/SuperPythonic/SuperPythonic/internal/theories/conc"

type (
	Prog = conc.Prog

	Def = conc.Def
	Fn  = conc.Fn

	Param = conc.Param

	Var = conc.Var

	Expr = conc.Expr

	Ref   = conc.Ref
	Unit  = conc.Unit
	Int   = conc.Int
	Float = conc.Float
	Str   = conc.Str
	Let   = conc.Let

	UnitType  = conc.UnitType
	BoolType  = conc.BoolType
	IntType   = conc.IntType
	FloatType = conc.FloatType
	StrType   = conc.StrType
)
