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
	Text() string

	IsAtomic() bool
	SetAtomic(atomic bool)
	SkipSpaces()
}

type Options interface {
	IsSpace(r rune) bool
	IsNewline(r rune) bool
	IndentWord() string
	OnNewline(state State)
}

type ParserFunc func(s State) State

func Atom(parser ParserFunc) ParserFunc {
	return func(s State) State {
		prev := s.IsAtomic()
		s.SetAtomic(true)
		defer s.SetAtomic(prev)
		return parser(s)
	}
}
