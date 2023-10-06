package surf

import (
	"fmt"
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
	for _, token := range s.Committed() {
		fmt.Println(s.Text(token))
	}
	if s.Cur().Kind == parsing.Error {
		t.Fatal(s)
	}
}
