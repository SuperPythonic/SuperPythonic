package parsing

type Span struct{ Start, End, Line, Col int }

type State interface {
	Pos() int
	Dump() (pos, ln, col int, span *Span)
	Restore(pos, ln, col int, span *Span)

	Peek() (rune, bool)
	Next() (rune, bool)
	Eat(e rune) bool

	IsError() bool
	IsSOI() bool
	IsEOI() bool

	WithSpan(start int) State
	WithError(start int) State
	Span() *Span
	Text(span *Span) string

	SkipSpaces()
}

type Options interface {
	IsSpace(r rune) bool
	IsNewline(r rune) bool
}

type ParserFunc func(s State) State
