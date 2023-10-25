package parsers

import "github.com/SuperPythonic/SuperPythonic/internal/parsers"

var (
	NewState = parsers.NewState
	Parse    = parsers.Parse

	OnState = parsers.OnState
	On      = parsers.On
	OnWord  = parsers.OnWord

	Start = parsers.Start
	End   = parsers.End

	Newline = parsers.Newline
	Entry   = parsers.Entry
	Indent  = parsers.Indent
	Exit    = parsers.Exit

	Range = parsers.Range
	Word  = parsers.Word

	Lowercase = parsers.Lowercase
	Uppercase = parsers.Uppercase
	CamelBack = parsers.CamelBack
	CamelCase = parsers.CamelCase

	Int   = parsers.Int
	Float = parsers.Float
	Str   = parsers.Str

	Seq    = parsers.Seq
	Choice = parsers.Choice
	Many   = parsers.Many
	More   = parsers.More
	Option = parsers.Option
	Times  = parsers.Times
	Occur  = parsers.Occur
	Not    = parsers.Not
)
