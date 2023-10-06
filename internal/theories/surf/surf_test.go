package surf

import (
	"testing"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
)

func TestParse(t *testing.T) {
	const (
		text = "  \n  def  main   "
	)
	if s := parsers.Seq(
		parsers.Start(),
		parsers.Keyword(42, "def"),
		parsers.Keyword(69, "main"),
		parsers.End(),
	).
		Run(parsers.NewState(text)); s.Cur().Kind != parsing.EOI {
		t.Fatal(s)
	}
}
