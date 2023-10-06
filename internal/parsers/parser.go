package parsers

import "github.com/SuperPythonic/SuperPythonic/pkg/parsing"

type startParser struct{}

func Start() parsing.Parser { return new(startParser) }

func (p *startParser) Run(s parsing.State) parsing.State {
	if !IsSOI(s) {
		return SetError(s, s.Pos())
	}
	return s
}

type endParser struct{}

func End() parsing.Parser { return new(endParser) }

func (e *endParser) Run(s parsing.State) parsing.State {
	if !IsEOI(s) {
		return SetError(s, s.Pos())
	}
	return s
}

type keywordParser struct {
	kind parsing.TokenKind
	word string
}

func Keyword(kind parsing.TokenKind, word string) parsing.Parser { return &keywordParser{kind, word} }

func (p *keywordParser) Run(s parsing.State) parsing.State {
	start := s.Pos()
	for _, c := range p.word {
		if !s.Rune(c) {
			return SetError(s, start)
		}
	}
	return s.SetToken(p.kind, start)
}

type seqParser struct {
	parsers []parsing.Parser
}

func Seq(parsers ...parsing.Parser) parsing.Parser {
	return &seqParser{parsers}
}

func (p *seqParser) Run(s parsing.State) parsing.State {
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
