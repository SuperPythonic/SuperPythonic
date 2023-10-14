package parsers

import "github.com/SuperPythonic/SuperPythonic/internal/parsers"

var (
	NewState = parsers.NewState
	Parse    = parsers.Parse

	On     = parsers.On
	OnText = parsers.OnText
	OnWord = parsers.OnWord

	Start = parsers.Start
	End   = parsers.End

	Range = parsers.Range
	Word  = parsers.Word

	Lowercase = parsers.Lowercase
	Uppercase = parsers.Uppercase
	CamelBack = parsers.CamelBack
	CamelCase = parsers.CamelCase

	Int = parsers.Int
	Str = parsers.Str

	Seq    = parsers.Seq
	Choice = parsers.Choice
	Many   = parsers.Many
	More   = parsers.More
	Option = parsers.Option
	Times  = parsers.Times
	Occur  = parsers.Occur
	Not    = parsers.Not
)
