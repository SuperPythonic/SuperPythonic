package surf

import "testing"

func TestParse(t *testing.T) {
	const (
		text = "  \n  def  main   "
	)
	if s := Seq(
		Start(),
		Keyword(42, "def"),
		Keyword(69, "main"),
		End(),
	).
		Run(NewState(text)); s.Cur().Kind != EOI {
		t.Fatal(s)
	}
}
