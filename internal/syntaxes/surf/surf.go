package surf

import "unicode"

type TokenKind int

type Token struct {
	Kind                  TokenKind
	Start, End, Line, Col int
}

const (
	Error TokenKind = iota
	EOF
	WHITESPACE
	NEWLINE
)

type State interface {
	Pos() int
	Rune(r rune) bool
	Token(kind TokenKind) bool

	SetToken(kind TokenKind, start int) State
}

type Options interface {
	IsSpace(r rune) bool
	IsNewline(r rune) bool
}

type state struct {
	opt          Options
	text         []rune
	pos, ln, col int
	cur          *Token
}

func NewState(text string) State {
	return NewStateWith(text, new(opt))
}

func NewStateWith(text string, opt Options) State {
	return &state{opt: opt, text: []rune(text), ln: 1, col: 1}
}

type opt struct{}

func (*opt) IsSpace(r rune) bool { return unicode.IsSpace(r) }

func (*opt) IsNewline(r rune) bool { return r == '\n' }

func (p *state) Pos() int { return p.pos }

func (p *state) SetToken(kind TokenKind, start int) State {
	p.cur = &Token{
		Kind:  kind,
		Start: start,
		End:   p.pos,
		Line:  p.ln,
		Col:   p.col,
	}
	return p
}

func SetError(p State, start int) State { return p.SetToken(Error, start) }

func (p *state) Rune(r rune) bool {
	if p.pos >= len(p.text) {
		return false
	}

	c := p.text[p.pos]
	p.pos++
	p.col++

	if p.opt.IsNewline(c) {
		p.ln++
		p.col = 1
	}

	if p.pos == len(p.text) {
		p.SetToken(EOF, p.pos)
	}

	return c == r
}

func (p *state) Token(kind TokenKind) bool { return p.cur.Kind == kind }

type Parser interface {
	Run(p State) State
}

type keywordParser struct {
	kind TokenKind
	word string
}

func (h *keywordParser) Run(p State) State {
	start := p.Pos()
	for _, c := range h.word {
		if !p.Rune(c) {
			return SetError(p, start)
		}
	}
	return p.SetToken(h.kind, start)
}

func Keyword(kind TokenKind, word string) Parser { return &keywordParser{kind, word} }
