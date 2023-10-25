package parsing

type Point struct{ Pos, Line, Col int }

type (
	Spanner interface{ StartEnd() Span }
	Span    struct{ Start, End Point }
)

func (s Span) StartEnd() Span { return s }

type State interface {
	Options() Options

	Point() Point
	Restore(p Point)

	Peek() (rune, bool)
	Next() (rune, bool)
	Eat(e rune) bool

	IsError() bool
	IsSOI() bool
	IsEOI() bool

	WithOK() State
	WithError() State
	Span(start Point) Span
	Text(sp Span) string

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
