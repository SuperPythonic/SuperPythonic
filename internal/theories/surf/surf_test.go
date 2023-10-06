package surf

import (
	"testing"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
)

func TestParse(t *testing.T) {
	const text = "  \n  class  "
	if s := Prog().Parse(parsers.NewState(text)); s.Cur().Kind != parsing.EOI {
		t.Fatal(s)
	}
}
