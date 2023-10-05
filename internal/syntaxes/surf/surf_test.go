package surf

import "testing"

func TestParse(t *testing.T) {
	const (
		text = "  \n  def  main   "
	)
	s := Seq(
		Keyword(42, "def"),
		Keyword(69, "main"),
	).
		Run(NewState(text))
	println(s)
}
