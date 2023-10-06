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
		if !s.Rune(c) {
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
		r, ok := s.Next()
		if !ok {
			break
		}

		if first && unicode.IsDigit(r) {
			break
		}
		first = false

		if !p.f.IsRuneValid(r) {
			break
		}
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
	// TODO: Collect all failed states for error messages.
	for _, parser := range p.parsers {
		if s = parser.Parse(s); !IsError(s) {
			return s
		}
		s.Reset(pos, ln, col)
	}
	return s
}
