package surf

import (
	"testing"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

func TestParse(t *testing.T) {
	const text = `
def fn_0 ( ) :

   def fn_1 ( aaa )    :

def fn_2 ( aaa, bbb ):
    `
	s := Parse(text)
	if s.Cur().Kind == parsing.Error {
		t.Fatal(s)
	}
}
