package parsers

import "github.com/SuperPythonic/SuperPythonic/pkg/parsing"

type State struct {
	opt  parsing.Options
	text []rune

	pos, ln, col int

	cur       *parsing.Token
	committed []*parsing.Token
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
func (p *State) Reset(pos, ln, col int) {
	p.pos = pos
	p.ln = ln
	p.col = col
	p.cur = nil
}

func (p *State) Rune(r rune) bool {
	if c, ok := p.Next(); ok {
		return c == r
	}
	return false
}

func (p *State) Next() (rune, bool) {
	if p.pos >= len(p.text) {
		return 0, false
	}

	c := p.text[p.pos]
	p.pos++
	p.col++

	if p.opt.IsNewline(c) {
		p.ln++
		p.col = 1
	}

	if p.pos == len(p.text) {
		p.Set(parsing.EOI, p.pos)
	}

	return c, true
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

func (p *State) Commit() parsing.State {
	if p.cur == nil {
		panic("current token is nil")
	}
	if p.cur.Kind == parsing.Error {
		panic("cannot commit error token")
	}
	p.committed = append(p.committed, p.cur)
	return p
}

func (p *State) Cur() *parsing.Token { return p.cur }

func (p *State) All() []*parsing.Token { return p.committed }

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
		p.Rune(r)
	}
}

func SetError(s parsing.State, start int) parsing.State { return s.Set(parsing.Error, start) }
func IsSOI(s parsing.State) bool                        { return s.Cur().Kind == parsing.SOI }
func IsEOI(s parsing.State) bool                        { return s.Cur().Kind == parsing.EOI }
func IsError(s parsing.State) bool                      { return s.Cur().Kind == parsing.Error }
