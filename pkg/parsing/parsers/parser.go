package parsers

import "github.com/SuperPythonic/SuperPythonic/internal/parsers"

type Word = parsers.Word

var (
	NewState = parsers.NewState
	Parse    = parsers.Parse

	Start = parsers.Start
	End   = parsers.End

	Range = parsers.Range

	Lowercase = parsers.Lowercase
	Uppercase = parsers.Uppercase
	CamelBack = parsers.CamelBack
	CamelCase = parsers.CamelCase

	OctDigit = parsers.OctDigit
	HexDigit = parsers.HexDigit
	Int      = parsers.Int
	Str      = parsers.Str
	Bool     = parsers.Bool
	Unit     = parsers.Unit

	Seq    = parsers.Seq
	Choice = parsers.Choice
	Many   = parsers.Many
	More   = parsers.More
	Option = parsers.Option
	Times  = parsers.Times
	Occur  = parsers.Occur
	Not    = parsers.Not
)
