package surf

import "testing"

func TestParse(t *testing.T) {
	const (
		text = "def"
		kind = 114514
	)
	ok := Keyword(kind, "def").Run(NewState(text)).Token(kind)
	println(ok)
}
