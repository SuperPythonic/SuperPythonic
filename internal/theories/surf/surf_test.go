package surf

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	texts := []string{
		`
def f0() -> str:
	return "hi"
`,
		`
def f1() -> str:
	a: str = "hi"
	b: int = 42
	return b
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
