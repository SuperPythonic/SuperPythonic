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
	).
		Run(NewState(text)); s.Cur().Kind != 69 {
		t.Fatal(s)
	}
}
