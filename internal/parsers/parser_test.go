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

func TestStr(t *testing.T) {
	ss := []string{
		`""`,
		`"hello"`,
		`"Hello, \1 \a \123 \xfF \ubeeF \u{deadbeeF}!"`,
	}
	for _, text := range ss {
		if s := Str().Parse(NewState(text)); s.IsError() {
			t.Fatal(s)
		}
	}
}
