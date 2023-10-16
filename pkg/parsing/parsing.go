package parsing

type Span struct{ Start, End, Line, Col int }

type State interface {
	Options() Options

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

	IndentWord() string
	Depth() int
	WithEntry() State
	WithExit() State
}

type Options interface {
	IsSpace(r rune) bool
	Newline() rune
	Indent() rune
	IndentWordN() int
	IsKeyword(text string) bool
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
