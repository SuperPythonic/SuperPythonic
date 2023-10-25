package parsers

import (
	"strings"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

type State struct {
	opts   parsing.Options
	text   []rune
	point  parsing.Point
	isErr  bool
	atomic bool
	depth  int
	indent string
}

func NewState(text string) *State { return NewStateWith(text, new(opts)) }

func NewStateWith(text string, opt parsing.Options) *State {
	return &State{
		opts:   opt,
		text:   []rune(text),
		point:  parsing.Point{Line: 1, Col: 1},
		indent: strings.Repeat(string(opt.Indent()), opt.IndentWordN()),
	}
}

func Parse(parser parsing.ParserFunc, text string) parsing.State {
	return parser(NewState(text))
}

func OnState(parser parsing.ParserFunc, f func(s parsing.State)) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		if s = parser(s); !s.IsError() {
			f(s)
		}
		return s
	}
}

func On(parser parsing.ParserFunc, f func()) parsing.ParserFunc {
	return OnState(parser, func(parsing.State) { f() })
}

func OnWord(w string, f func()) parsing.ParserFunc { return On(Word(w), f) }

func (s *State) Options() parsing.Options { return s.opts }

func (s *State) Point() parsing.Point { return s.point }
func (s *State) Restore(p parsing.Point) {
	s.point = p
	s.isErr = false
}

func (s *State) Peek() (rune, bool) {
	if s.point.Pos >= len(s.text) {
		return 0, false
	}
	return s.text[s.point.Pos], true
}

func (s *State) Next() (rune, bool) {
	r, ok := s.Peek()
	if !ok {
		return 0, false
	}

	s.point.Pos++
	s.point.Col++

	if r == s.opts.Newline() {
		s.point.Line++
		s.point.Col = 1
	}

	return r, true
}

func (s *State) Eat(e rune) bool {
	if r, ok := s.Next(); ok {
		return r == e
	}
	return false
}

func (s *State) IsError() bool { return s.isErr }
func (s *State) IsSOI() bool   { return s.point.Pos == 0 }
func (s *State) IsEOI() bool   { return s.point.Pos == len(s.text) }

func (s *State) WithOK() parsing.State {
	s.isErr = false
	return s
}

func (s *State) WithError() parsing.State {
	s.isErr = true
	return s
}

func (s *State) Span(start parsing.Point) parsing.Span {
	return parsing.Span{Start: start, End: s.Point()}
}

func (s *State) Text(sp parsing.Span) string { return string(s.text[sp.Start.Pos:sp.End.Pos]) }

func (s *State) IsAtomic() bool        { return s.atomic }
func (s *State) SetAtomic(atomic bool) { s.atomic = atomic }

func (s *State) SkipSpaces() {
	if s.atomic {
		return
	}
	for {
		if r, ok := s.Peek(); ok && s.opts.IsSpace(r) {
			s.Eat(r)
			continue
		}
		break
	}
}

func (s *State) IndentWord() string       { return s.indent }
func (s *State) Depth() int               { return s.depth }
func (s *State) WithEntry() parsing.State { s.depth++; return s }
func (s *State) WithExit() parsing.State {
	if s.depth == 0 {
		return s.WithError()
	}
	s.depth--
	return s
}
