package conc

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	const text = `
def fn_0 ( ) :

   def fn_1 ( aaa :   int )    :

def fn_2 ( bbb  : int,  ccc: bool ):
    `
	prog, s := Parse(text)
	if s.IsError() {
		t.Fatal(s)
	}
	fmt.Println(prog)
}
