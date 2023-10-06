package surf

import (
	"testing"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
	"github.com/SuperPythonic/SuperPythonic/pkg/parsing/parsers"
)

func TestParse(t *testing.T) {
	const text = "  \n  def   hello_world  "
	s := Prog().Parse(parsers.NewState(text))
	if s.Cur().Kind != parsing.EOI {
		t.Fatal(s)
	}
}
