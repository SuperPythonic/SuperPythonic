package parsers

import "testing"

func TestLong(t *testing.T) {
	if s := Long().Parse(NewState("0b1_01_0")); s.IsError() {
		t.Fatal(s)
	}
}
