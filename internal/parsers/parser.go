package parsers

import (
	"unicode"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

type start struct{}

func Start() parsing.Parser { return new(start) }

func (p *start) Parse(s parsing.State) parsing.State {
	if !IsSOI(s) {
		return SetError(s, s.Pos())
	}
	return s.Commit()
}

type end struct{}

func End() parsing.Parser { return new(end) }

func (e *end) Parse(s parsing.State) parsing.State {
	if !IsEOI(s) {
		return SetError(s, s.Pos())
	}
	return s.Commit()
}

type kw string

func Keyword(word string) parsing.Parser { return kw(word) }

func (p kw) Parse(s parsing.State) parsing.State {
	start := s.Pos()
	for _, c := range p {
		if !s.Eat(c) {
			return SetError(s, start)
		}
	}
	return s.Set(parsing.Keyword, start).Commit()
}

type Ident interface{ IsRuneValid(r rune) bool }

type id struct{ f Ident }

func (p *id) Parse(s parsing.State) parsing.State {
	start := s.Pos()
	first := true

	for {
		r, ok := s.Peek()
		if !ok {
			break
		}

		if first && unicode.IsDigit(r) {
			break
		}
		first = false

		if !unicode.IsDigit(r) && !p.f.IsRuneValid(r) {
			break
		}

		_, _ = s.Next()
	}

	if start == s.Pos() {
		return SetError(s, start)
	}
	return s.Set(parsing.Ident, start).Commit()
}

type lc struct{}

func Lowercase() parsing.Parser     { return &id{new(lc)} }
func (*lc) IsRuneValid(r rune) bool { return unicode.IsLower(r) || r == '_' }

type up struct{}

func Uppercase() parsing.Parser     { return &id{new(up)} }
func (*up) IsRuneValid(r rune) bool { return unicode.IsUpper(r) || r == '_' }

type cb struct{ first bool }

func CamelBack() parsing.Parser { return &id{&cb{true}} }
func (c *cb) IsRuneValid(r rune) bool {
	if c.first {
		c.first = false
		return unicode.IsLower(r)
	}
	return unicode.IsLetter(r)
}

type cc struct{ first bool }

func CamelCase() parsing.Parser { return &id{&cc{true}} }
func (c *cc) IsRuneValid(r rune) bool {
	if c.first {
		c.first = false
		return unicode.IsUpper(r)
	}
	return unicode.IsLetter(r)
}

type seq struct{ parsers []parsing.Parser }

func Seq(parsers ...parsing.Parser) parsing.Parser { return &seq{parsers} }

func (p *seq) Parse(s parsing.State) parsing.State {
	for i, parser := range p.parsers {
		if s = parser.Parse(s); IsError(s) {
			return s
		}
		if i+1 < len(p.parsers) {
			s.SkipSpaces()
		}
	}
	return s
}

type choice struct{ parsers []parsing.Parser }

func Choice(parsers ...parsing.Parser) parsing.Parser { return &choice{parsers} }

func (p *choice) Parse(s parsing.State) parsing.State {
	pos, ln, col := s.Loc()
	prev := s.Cur()
	// TODO: Collect all failed states for error messages.
	for _, parser := range p.parsers {
		if s = parser.Parse(s); !IsError(s) {
			return s
		}
		s.Reset(pos, ln, col, prev)
	}
	return s
}

type many struct{ parser parsing.Parser }

func Many(parser parsing.Parser) parsing.Parser { return &many{parser} }

func (p *many) Parse(s parsing.State) parsing.State {
	var (
		pos, ln, col int
		prev         *parsing.Token
	)
	for {
		pos, ln, col = s.Loc()
		prev = s.Cur()
		if s = p.parser.Parse(s); IsError(s) {
			s.Reset(pos, ln, col, prev)
			break
		}
		s.SkipSpaces()
	}
	return s
}
