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
		cur:  &parsing.Token{Kind: parsing.SOI, Line: 1, Col: 1},
	}
}

func (p *State) Pos() int { return p.pos }

func (p *State) Loc() (pos, ln, col int) { return p.pos, p.ln, p.col }
func (p *State) Reset(pos, ln, col int, t *parsing.Token) {
	p.pos = pos
	p.ln = ln
	p.col = col
	p.cur = t
}

func (p *State) Peek() (rune, bool) {
	if p.pos >= len(p.text) {
		return 0, false
	}
	return p.text[p.pos], true
}

func (p *State) Next() (rune, bool) {
	r, ok := p.Peek()
	if !ok {
		return 0, false
	}

	p.pos++
	p.col++

	if p.opt.IsNewline(r) {
		p.ln++
		p.col = 1
	}

	if p.pos == len(p.text) {
		p.Set(parsing.EOI, p.pos)
	}

	return r, true
}

func (p *State) Eat(e rune) bool {
	if r, ok := p.Next(); ok {
		return r == e
	}
	return false
}

func (p *State) Set(kind parsing.TokenKind, start int) parsing.State {
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

func (p *State) Text(t *parsing.Token) string { return string(p.text[t.Start:t.End]) }

func (p *State) SkipSpaces() {
	for {
		if p.pos >= len(p.text) {
			return
		}

		r := p.text[p.pos]
		if !p.opt.IsSpace(r) {
			return
		}
		p.Eat(r)
	}
}
