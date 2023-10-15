package conc

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	const text = `
def fn_0 ( ) -> str:
	return "hello"

   def fn_1 ( aaa :  unit )    :
return   ()

def fn_2 ( bbb  : int,  ccc: bool ) -> int:
			return  42

def fn_3() -> bool: return True

def fn_4(): return
    `
	prog, s := Parse(text)
	if s.IsError() {
		t.Fatal(s)
	}
	fmt.Println(prog)
}
