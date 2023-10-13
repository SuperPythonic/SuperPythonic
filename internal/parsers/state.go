package parsers

import "github.com/SuperPythonic/SuperPythonic/pkg/parsing"

const noErrAt = -1

type State struct {
	opts         parsing.Options
	text         []rune
	pos, ln, col int
	cur          *parsing.Span
	errAt        int
}

func NewState(text string) *State { return NewStateWith(text, new(opts)) }

func NewStateWith(text string, opt parsing.Options) *State {
	return &State{opts: opt, text: []rune(text), ln: 1, col: 1, errAt: noErrAt}
}

func Parse(parser parsing.Parser, text string) parsing.State {
	return parser.Parse(NewState(text))
}

func (s *State) Pos() int { return s.pos }

func (s *State) Dump() (pos, ln, col int, span *parsing.Span) { return s.pos, s.ln, s.col, s.cur }
func (s *State) Restore(pos, ln, col int, span *parsing.Span) {
	s.pos = pos
	s.ln = ln
	s.col = col
	s.cur = span
	s.errAt = noErrAt
}

func (s *State) Peek() (rune, bool) {
	if s.pos >= len(s.text) {
		return 0, false
	}
	return s.text[s.pos], true
}

func (s *State) Next() (rune, bool) {
	r, ok := s.Peek()
	if !ok {
		return 0, false
	}

	s.pos++
	s.col++

	if s.opts.IsNewline(r) {
		s.ln++
		s.col = 1
	}

	return r, true
}

func (s *State) Eat(e rune) bool {
	if r, ok := s.Next(); ok {
		return r == e
	}
	return false
}

func (s *State) IsError() bool { return s.errAt != noErrAt }

func (s *State) IsSOI() bool { return s.pos == 0 }

func (s *State) IsEOI() bool { return s.pos == len(s.text) }

func (s *State) WithSpan(start int) parsing.State {
	s.cur = &parsing.Span{
		Start: start,
		End:   s.pos,
		Line:  s.ln,
		Col:   s.col - (s.pos - start),
	}
	s.errAt = noErrAt
	return s
}

func (s *State) WithError(start int) parsing.State {
	s.errAt = start
	return s
}

func (s *State) Span() *parsing.Span { return s.cur }

func (s *State) Text(span *parsing.Span) string { return string(s.text[span.Start:span.End]) }

func (s *State) SkipSpaces() {
	for {
		if s.pos >= len(s.text) {
			return
		}

		r := s.text[s.pos]
		if !s.opts.IsSpace(r) {
			return
		}
		s.Eat(r)
	}
}
