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
	SetToken(kind TokenKind, start int) State
	IsError() bool
	SkipSpaces()
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
		Col:   p.col - (p.pos - start),
	}
	return p
}

func (p *state) IsError() bool { return p.cur.Kind == Error }

func (p *state) SkipSpaces() {
	if p.pos >= len(p.text) {
		return
	}

	for {
		r := p.text[p.pos]
		if !p.opt.IsSpace(r) {
			return
		}
		p.Rune(r)
	}
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

type Parser interface {
	Run(s State) State
}

type keywordParser struct {
	kind TokenKind
	word string
}

func Keyword(kind TokenKind, word string) Parser { return &keywordParser{kind, word} }

func (p *keywordParser) Run(s State) State {
	start := s.Pos()
	for _, c := range p.word {
		if !s.Rune(c) {
			return SetError(s, start)
		}
	}
	return s.SetToken(p.kind, start)
}

type seqParser struct {
	parsers []Parser
}

func Seq(parsers ...Parser) Parser {
	return &seqParser{parsers}
}

func (p *seqParser) Run(s State) State {
	for _, parser := range p.parsers {
		s.SkipSpaces()
		s = parser.Run(s)
		if s.IsError() {
			return s
		}
	}
	return s
}
