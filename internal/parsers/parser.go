package parsers

import (
	"unicode"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

type start struct{}

func Start() parsing.Parser { return new(start) }

func (p *start) Parse(s parsing.State) parsing.State {
	if !parsing.IsSOI(s) {
		return parsing.SetError(s, s.Pos())
	}
	return s
}

type end struct{}

func End() parsing.Parser { return new(end) }

func (e *end) Parse(s parsing.State) parsing.State {
	if !parsing.IsEOI(s) {
		return parsing.SetError(s, s.Pos())
	}
	return s
}

type kw string

func Keyword(word string) parsing.Parser { return kw(word) }

func (p kw) Parse(s parsing.State) parsing.State {
	start := s.Pos()
	for _, c := range p {
		if !s.Eat(c) {
			return parsing.SetError(s, start)
		}
	}
	return s.Set(parsing.Keyword, start)
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
		return parsing.SetError(s, start)
	}
	return s.Set(parsing.Ident, start)
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
	return Choice(new(bin), new(oct), new(hex)).Parse(s)
}

type bin struct{}

func (*bin) Parse(s parsing.State) parsing.State {
	//TODO implement me
	panic("implement me")
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

func Seq(parsers ...parsing.Parser) parsing.Parser { return &seq{parsers} }

func (p *seq) Parse(s parsing.State) parsing.State {
	if len(p.parsers) == 0 {
		panic("at least 1 parser expected")
	}
	for i, parser := range p.parsers {
		if s = parser.Parse(s); parsing.IsError(s) {
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
	if len(p.parsers) == 0 {
		panic("at least 1 choice expected")
	}
	pos, ln, col := s.Loc()
	prev := s.Cur()
	// TODO: Collect all failed states for error messages.
	for _, parser := range p.parsers {
		if s = parser.Parse(s); !parsing.IsError(s) {
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
		if s = p.parser.Parse(s); parsing.IsError(s) {
			s.Reset(pos, ln, col, prev)
			break
		}
		if parsing.IsEOI(s) {
			break
		}
		s.SkipSpaces()
	}
	return s
}

type opt struct{ parser parsing.Parser }

func Optional(parser parsing.Parser) parsing.Parser { return &opt{parser} }

func (p *opt) Parse(s parsing.State) parsing.State {
	pos, ln, col := s.Loc()
	prev := s.Cur()
	if s = p.parser.Parse(s); parsing.IsError(s) {
		s.Reset(pos, ln, col, prev)
	}
	return s
}
