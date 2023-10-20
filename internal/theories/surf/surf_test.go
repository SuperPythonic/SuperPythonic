package surf

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	texts := []string{
		`
def f0() -> () -> int:
	return lambda: x

def f1() -> int -> int:
	return lambda x: x

def f2() -> (int) -> int:
	return lambda x: x

def f3() -> (int, int) -> int:
	return lambda x: x

def f4() -> (int, int, str) -> int:
	return lambda x: x

def f5() -> (int -> str) -> int:
	return lambda x: x

def f6() -> (int -> str) -> (int -> str):
	return lambda x: x
`,
		`
def f() -> str:
	return ("hi")
`,
		`
def f() -> str:
	a: str = "hi"
	b: int = 42
	c: float = 42.0
	return c
`,
	}
	for _, text := range texts {
		prog, s := Parse(text)
		if s.IsError() {
			t.Fatal(s)
		}
		fmt.Println(prog)
	}
}
