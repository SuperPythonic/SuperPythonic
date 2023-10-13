package parsers

import "testing"

func TestInt(t *testing.T) {
	if s := Int().Parse(NewState("0b1_01_0")); s.IsError() {
		t.Fatal(s)
	}
	if s := Int().Parse(NewState("0o6_7_6_7")); s.IsError() {
		t.Fatal(s)
	}
	if s := Int().Parse(NewState("0xdead_beef")); s.IsError() {
		t.Fatal(s)
	}
}
