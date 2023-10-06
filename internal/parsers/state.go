package parsers

import "github.com/SuperPythonic/SuperPythonic/pkg/parsing"

type State struct {
	opt          parsing.Options
	text         []rune
	pos, ln, col int
	cur          *parsing.Token
}

func NewState(text string) *State {
	return NewStateWith(text, new(opt))
}

func NewStateWith(text string, opt parsing.Options) *State {
	return &State{
		opt:  opt,
		text: []rune(text),
		ln:   1,
		col:  1,
		cur:  &parsing.Token{Kind: parsing.SOI},
	}
}

func (p *State) Pos() int { return p.pos }

func (p *State) SetToken(kind parsing.TokenKind, start int) parsing.State {
	p.cur = &parsing.Token{
		Kind:  kind,
		Start: start,
		End:   p.pos,
		Line:  p.ln,
		Col:   p.col - (p.pos - start),
	}
	return p
}

func (p *State) Cur() *parsing.Token { return p.cur }

func (p *State) SkipSpaces() {
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

func (p *State) Rune(r rune) bool {
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
		p.SetToken(parsing.EOI, p.pos)
	}

	return c == r
}

func SetError(s parsing.State, start int) parsing.State { return s.SetToken(parsing.Error, start) }
func IsSOI(s parsing.State) bool                        { return s.Cur().Kind == parsing.SOI }
func IsEOI(s parsing.State) bool                        { return s.Cur().Kind == parsing.EOI }
func IsError(s parsing.State) bool                      { return s.Cur().Kind == parsing.Error }
