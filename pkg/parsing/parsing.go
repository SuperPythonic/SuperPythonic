package parsing

type TokenKind int

const (
	Error TokenKind = iota
	SOI
	EOI
	Keyword
)

type Token struct {
	Kind                  TokenKind
	Start, End, Line, Col int
}

type State interface {
	Pos() int
	Loc() (pos, ln, col int)
	SetLoc(pos, ln, col int)
	Rune(r rune) bool
	SetToken(kind TokenKind, start int) State // TODO: append to token stream
	Cur() *Token
	SkipSpaces()
}

type Options interface {
	IsSpace(r rune) bool
	IsNewline(r rune) bool
}

type Parser interface{ Parse(s State) State }
