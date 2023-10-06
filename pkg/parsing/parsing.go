package parsing

type TokenKind int

type Token struct {
	Kind                  TokenKind
	Start, End, Line, Col int
}

const (
	Error TokenKind = iota
	SOI
	EOI
	NEWLINE
)

type State interface {
	Pos() int
	Rune(r rune) bool
	SetToken(kind TokenKind, start int) State
	Cur() *Token
	SkipSpaces()
}

type Options interface {
	IsSpace(r rune) bool
	IsNewline(r rune) bool
}

type Parser interface {
	Run(s State) State
}
