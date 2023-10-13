package parsers

import (
	"github.com/SuperPythonic/SuperPythonic/internal/parsers"
)

type Keyword = parsers.Keyword

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

	Int = parsers.Int

	Seq    = parsers.Seq
	Choice = parsers.Choice
	Many   = parsers.Many
	Option = parsers.Option
	Times  = parsers.Times
	Until  = parsers.Until
)
