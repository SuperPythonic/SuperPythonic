package surf

import "unicode"

type TokenKind int

type Token struct {
	Kind                  TokenKind
	Start, End, Line, Col int
}

const (
	Error TokenKind = iota
	SOI
	EOI
	WHITESPACE
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
	return &state{
		opt:  opt,
		text: []rune(text),
		ln:   1,
		col:  1,
		cur:  &Token{Kind: SOI},
	}
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

func (p *state) Cur() *Token { return p.cur }

func (p *state) SkipSpaces() {
	for {
		if p.pos >= len(p.text) {
			return
		}

		r := p.text[p.pos]
		if !p.opt.IsSpace(r) {
			return
		}
		p.Rune(r)
	}
}

func SetError(s State, start int) State { return s.SetToken(Error, start) }
func IsSOI(s State) bool                { return s.Cur().Kind == SOI }
func IsEOI(s State) bool                { return s.Cur().Kind == EOI }
func IsError(s State) bool              { return s.Cur().Kind == Error }

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
		p.SetToken(EOI, p.pos)
	}

	return c == r
}

type Parser interface {
	Run(s State) State
}

type startParser struct{}

func Start() Parser { return new(startParser) }

func (p *startParser) Run(s State) State {
	if !IsSOI(s) {
		return SetError(s, s.Pos())
	}
	return s
}

type endParser struct{}

func End() Parser { return new(endParser) }

func (e *endParser) Run(s State) State {
	if !IsEOI(s) {
		return SetError(s, s.Pos())
	}
	return s
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
	for i, parser := range p.parsers {
		s = parser.Run(s)
		if IsError(s) {
			return s
		}
		if i+1 < len(p.parsers) {
			s.SkipSpaces()
		}
	}
	return s
}
