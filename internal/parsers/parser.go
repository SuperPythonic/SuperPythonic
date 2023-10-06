package parsers

import "github.com/SuperPythonic/SuperPythonic/pkg/parsing"

type start struct{}

func Start() parsing.Parser { return new(start) }

func (p *start) Run(s parsing.State) parsing.State {
	if !IsSOI(s) {
		return SetError(s, s.Pos())
	}
	return s
}

type end struct{}

func End() parsing.Parser { return new(end) }

func (e *end) Run(s parsing.State) parsing.State {
	if !IsEOI(s) {
		return SetError(s, s.Pos())
	}
	return s
}

type kw struct {
	kind parsing.TokenKind
	word string
}

func Keyword(kind parsing.TokenKind, word string) parsing.Parser { return &kw{kind, word} }

func (p *kw) Run(s parsing.State) parsing.State {
	start := s.Pos()
	for _, c := range p.word {
		if !s.Rune(c) {
			return SetError(s, start)
		}
	}
	return s.SetToken(p.kind, start)
}

type seq struct{ parsers []parsing.Parser }

func Seq(parsers ...parsing.Parser) parsing.Parser { return &seq{parsers} }

func (p *seq) Run(s parsing.State) parsing.State {
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
