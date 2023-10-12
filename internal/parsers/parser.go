package parsers

import (
	"fmt"
	"unicode"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

type start struct{}

func Start() parsing.Parser { return new(start) }

func (p *start) Parse(s parsing.State) parsing.State {
	if !s.IsSOI() {
		return s.WithError(s.Pos())
	}
	return s
}

type end struct{}

func End() parsing.Parser { return new(end) }

func (e *end) Parse(s parsing.State) parsing.State {
	if !s.IsEOI() {
		return s.WithError(s.Pos())
	}
	return s
}

type ran struct{ from, to rune }

func Range(from, to rune) parsing.Parser {
	if from >= to {
		panic(fmt.Sprintf("invalid range [%c, %c]", from, to))
	}
	return &ran{from, to}
}

func (p *ran) Parse(s parsing.State) parsing.State {
	r, ok := s.Peek()
	if !ok {
		return s
	}
	if p.from <= r && r >= p.to {
		_, _ = s.Next()
		return s
	}
	return s.WithError(s.Pos())
}

type kw string

func Keyword(word string) parsing.Parser { return kw(word) }

func (p kw) Parse(s parsing.State) parsing.State {
	start := s.Pos()
	for _, c := range p {
		if !s.Eat(c) {
			return s.WithError(start)
		}
	}
	return s.WithSpan(start)
}

type RuneFilter interface{ Valid(r rune) bool }

type id struct{ f RuneFilter }

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

		if !unicode.IsDigit(r) && !p.f.Valid(r) {
			break
		}

		_, _ = s.Next()
	}

	if start == s.Pos() {
		return s.WithError(start)
	}
	return s.WithSpan(start)
}

type lc struct{}

func Lowercase() parsing.Parser { return &id{new(lc)} }
func (*lc) Valid(r rune) bool   { return unicode.IsLower(r) || r == '_' }

type up struct{}

func Uppercase() parsing.Parser { return &id{new(up)} }
func (*up) Valid(r rune) bool   { return unicode.IsUpper(r) || r == '_' }

type cb struct{ first bool }

func CamelBack() parsing.Parser { return &id{&cb{true}} }
func (c *cb) Valid(r rune) bool {
	if c.first {
		c.first = false
		return unicode.IsLower(r)
	}
	return unicode.IsLetter(r)
}

type cc struct{ first bool }

func CamelCase() parsing.Parser { return &id{&cc{true}} }
func (c *cc) Valid(r rune) bool {
	if c.first {
		c.first = false
		return unicode.IsUpper(r)
	}
	return unicode.IsLetter(r)
}

type long struct{}

func Long() parsing.Parser { return new(long) }

func (*long) Parse(s parsing.State) parsing.State {
	start := s.Pos()
	if s = Choice(new(bin), new(oct), new(hex)).Parse(s); !s.IsError() {
		s.WithSpan(start)
	}
	return s
}

type bin struct{}

func (*bin) Parse(s parsing.State) parsing.State {
	return Seq(
		Choice(Keyword("0b"), Keyword("0B")),
		Range('0', '1'),
		Many(
			Seq(
				Option(Keyword("_")),
				Range('0', '1'),
			),
		),
	).Parse(s)
}

type oct struct{}

func (*oct) Parse(s parsing.State) parsing.State {
	//TODO implement me
	panic("implement me")
}

type hex struct{}

func (*hex) Parse(s parsing.State) parsing.State {
	//TODO implement me
	panic("implement me")
}

type seq struct{ parsers []parsing.Parser }

func Seq(parsers ...parsing.Parser) parsing.Parser {
	if len(parsers) == 0 {
		panic("at least 1 parser expected")
	}
	return &seq{parsers}
}

func (p *seq) Parse(s parsing.State) parsing.State {
	for i, parser := range p.parsers {
		if s = parser.Parse(s); s.IsError() {
			return s
		}
		if i+1 < len(p.parsers) {
			s.SkipSpaces()
		}
	}
	return s
}

type choice struct{ parsers []parsing.Parser }

func Choice(parsers ...parsing.Parser) parsing.Parser {
	if len(parsers) == 0 {
		panic("at least 1 choice expected")
	}
	return &choice{parsers}
}

func (p *choice) Parse(s parsing.State) parsing.State {
	pos, ln, col, prev := s.Dump()
	// TODO: Collect all failed states for error messages.
	for _, parser := range p.parsers {
		if s = parser.Parse(s); !s.IsError() {
			return s
		}
		s.Restore(pos, ln, col, prev)
	}
	return s
}

type many struct{ parser parsing.Parser }

func Many(parser parsing.Parser) parsing.Parser { return &many{parser} }

func (p *many) Parse(s parsing.State) parsing.State {
	var (
		pos, ln, col int
		prev         *parsing.Span
	)
	for {
		pos, ln, col, prev = s.Dump()
		if s = p.parser.Parse(s); s.IsError() {
			s.Restore(pos, ln, col, prev)
			break
		}
		if s.IsEOI() {
			break
		}
		s.SkipSpaces()
	}
	return s
}

type opt struct{ parser parsing.Parser }

func Option(parser parsing.Parser) parsing.Parser { return &opt{parser} }

func (p *opt) Parse(s parsing.State) parsing.State {
	pos, ln, col, prev := s.Dump()
	if s = p.parser.Parse(s); s.IsError() {
		s.Restore(pos, ln, col, prev)
	}
	return s
}
