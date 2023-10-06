package parsing

type TokenKind int

const (
	Error TokenKind = iota
	SOI
	EOI
	Keyword
	Ident
)

func (k TokenKind) String() string {
	switch k {
	case Error:
		return "Error"
	case SOI:
		return "SOI"
	case EOI:
		return "EOI"
	case Keyword:
		return "Keyword"
	case Ident:
		return "Ident"
	}
	panic(k)
}

type Token struct {
	Kind                  TokenKind
	Start, End, Line, Col int
}

type State interface {
	Pos() int
	Loc() (pos, ln, col int)
	Reset(pos, ln, col int, t *Token)

	Peek() (rune, bool)
	Next() (rune, bool)
	Eat(e rune) bool

	Set(kind TokenKind, start int) State
	Cur() *Token
	Text(t *Token) string

	SkipSpaces()
}

func SetError(s State, start int) State { return s.Set(Error, start) }
func IsSOI(s State) bool                { return s.Cur().Kind == SOI }
func IsEOI(s State) bool                { return s.Cur().Kind == EOI }
func IsError(s State) bool              { return s.Cur().Kind == Error }

type Options interface {
	IsSpace(r rune) bool
	IsNewline(r rune) bool
}

type Parser interface{ Parse(s State) State }
