package surf

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
	Match(r rune) bool

	SetToken(kind TokenKind, start int) State
}

type state struct {
	text         []rune
	pos, ln, col int
	cur          *Token
}

func NewState(text string) State {
	return &state{text: []rune(text), ln: 1, col: 1}
}

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

func (p *state) Match(r rune) bool {
	if p.pos >= len(p.text) {
		return false
	}

	ok := p.text[p.pos] == r

	p.pos++
	if p.pos == len(p.text) {
		p.SetToken(EOF, p.pos)
	}

	return ok
}

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
		if !p.Match(c) {
			return SetError(p, start)
		}
	}
	return p.SetToken(h.kind, start)
}

func Keyword(kind TokenKind, word string) Parser { return &keywordParser{kind, word} }
