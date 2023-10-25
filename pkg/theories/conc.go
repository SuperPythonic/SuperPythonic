package theories

import "github.com/SuperPythonic/SuperPythonic/internal/theories/conc"

type (
	Prog = conc.Prog

	Def = conc.Def
	Fn  = conc.Fn

	Param = conc.Param

	Var = conc.Var

	Expr = conc.Expr

	Unit  = conc.Unit
	Bool  = conc.Bool
	Int   = conc.Int
	Float = conc.Float
	Str   = conc.Str
	Lam   = conc.Lam
	Ref   = conc.Ref

	UnitType  = conc.UnitType
	BoolType  = conc.BoolType
	IntType   = conc.IntType
	FloatType = conc.FloatType
	StrType   = conc.StrType
	FnType    = conc.FnType

	App     = conc.App
	If      = conc.If
	Let     = conc.Let
	UnitLet = conc.UnitLet
)
