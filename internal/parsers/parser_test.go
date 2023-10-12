package parsers

import (
	"testing"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

func TestLong(t *testing.T) {
	if s := Long().Parse(NewState("0b1_01_0")); parsing.IsError(s) {
		t.Fatal(s)
	}
}
