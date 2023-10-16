package conc

import (
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
	return "hi"
`,
	}
	for _, text := range texts {
		if _, s := Parse(text); s.IsError() {
			t.Fatal(s)
		}
	}
}
