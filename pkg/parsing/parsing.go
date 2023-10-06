package parsing

type TokenKind int

const (
	Error TokenKind = iota
	SOI
	EOI
	Keyword
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
	Reset(pos, ln, col int)

	Rune(r rune) bool
	Set(kind TokenKind, start int) State
	Commit() State
	Cur() *Token

	SkipSpaces()
}

type Options interface {
	IsSpace(r rune) bool
	IsNewline(r rune) bool
}

type Parser interface{ Parse(s State) State }
