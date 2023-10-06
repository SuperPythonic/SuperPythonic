package parsers

import "github.com/SuperPythonic/SuperPythonic/pkg/parsing"

type start struct{}

func Start() parsing.Parser { return new(start) }

func (p *start) Parse(s parsing.State) parsing.State {
	if !IsSOI(s) {
		return SetError(s, s.Pos())
	}
	return s
}

type end struct{}

func End() parsing.Parser { return new(end) }

func (e *end) Parse(s parsing.State) parsing.State {
	if !IsEOI(s) {
		return SetError(s, s.Pos())
	}
	return s
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
	return s.SetToken(parsing.Keyword, start)
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
	// TODO: Collect all failed states for error messages.
	for _, parser := range p.parsers {
		if s = parser.Parse(s); !IsError(s) {
			return s
		}
		s.SetLoc(pos, ln, col)
	}
	return s
}
