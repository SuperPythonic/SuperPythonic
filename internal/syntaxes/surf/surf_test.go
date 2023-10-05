package surf

import "testing"

func TestParse(t *testing.T) {
	const text = "def"
	s := Keyword(114514, "def").Run(NewState(text))
	println(s)
}
