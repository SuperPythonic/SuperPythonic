package conc

import (
	"fmt"
	"testing"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

func TestParse(t *testing.T) {
	const text = `
def fn_0 ( ) :

   def fn_1 ( aaa )    :

def fn_2 ( bbb,  ccc ):
    `
	prog, s := Parse(text)
	if s.Cur().Kind == parsing.Error {
		t.Fatal(s)
	}
	fmt.Println(prog)
}
