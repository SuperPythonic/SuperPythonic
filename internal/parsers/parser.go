package parsers

import (
	"fmt"
	"unicode"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

type start struct{}

func Start() parsing.Parser { return (*start)(nil) }

func (*start) Parse(s parsing.State) parsing.State {
	if !s.IsSOI() {
		return s.WithError(s.Pos())
	}
	return s
}

type end struct{}

func End() parsing.Parser { return (*end)(nil) }

func (*end) Parse(s parsing.State) parsing.State {
	if !s.IsEOI() {
		return s.WithError(s.Pos())
	}
	return s
}

type ran struct{ from, to rune }

func Range(from, to rune) parsing.Parser {
	if from >= to {
		panic(fmt.Sprintf("invalid rune range [%c, %c]", from, to))
	}
	return &ran{from, to}
}

func (p *ran) Parse(s parsing.State) parsing.State {
	r, ok := s.Peek()
	if !ok {
		return s.WithError(s.Pos())
	}
	if p.from <= r && r <= p.to {
		_, _ = s.Next()
		return s
	}
	return s.WithError(s.Pos())
}

type Word string

func (p Word) Parse(s parsing.State) parsing.State {
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

func Lowercase() parsing.Parser { return &id{(*lc)(nil)} }
func (*lc) Valid(r rune) bool   { return unicode.IsLower(r) || r == '_' }

type up struct{}

func Uppercase() parsing.Parser { return &id{(*up)(nil)} }
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

type octDigit struct{}

func OctDigit() parsing.Parser                        { return (*octDigit)(nil) }
func (*octDigit) Parse(s parsing.State) parsing.State { return Range('0', '7').Parse(s) }

type hexDigit struct{}

func HexDigit() parsing.Parser { return (*hexDigit)(nil) }

func (*hexDigit) Parse(s parsing.State) parsing.State {
	return Choice(Range('0', '9'), Range('a', 'f'), Range('A', 'F')).Parse(s)
}

type i struct{}

func Int() parsing.Parser { return (*i)(nil) }

func (*i) Parse(s parsing.State) parsing.State {
	start := s.Pos()
	if s = Choice((*bin)(nil), (*oct)(nil), (*hex)(nil)).Parse(s); !s.IsError() {
		s.WithSpan(start)
	}
	return s
}

type bin struct{}

func (*bin) Parse(s parsing.State) parsing.State {
	return Seq(
		Choice(Word("0b"), Word("0B")),
		Range('0', '1'),
		Many(Seq(Option(Word("_")), Range('0', '1'))),
	).Parse(s)
}

type oct struct{}

func (*oct) Parse(s parsing.State) parsing.State {
	return Seq(
		Choice(Word("0o"), Word("0O")),
		OctDigit(),
		Many(Seq(Option(Word("_")), OctDigit())),
	).Parse(s)
}

type hex struct{}

func (p *hex) Parse(s parsing.State) parsing.State {
	return Seq(
		Choice(Word("0x"), Word("0X")),
		HexDigit(),
		Many(Seq(Option(Word("_")), HexDigit())),
	).Parse(s)
}

type (
	str              struct{}
	unescapedStrPart struct{}
	escapedStrPart   struct{}
)

func Str() parsing.Parser { return (*str)(nil) }

func (*str) Parse(s parsing.State) parsing.State {
	return Seq(
		Word(`"`),
		Many(Choice((*unescapedStrPart)(nil), (*escapedStrPart)(nil))),
		Word(`"`),
	).Parse(s)
}

func (*unescapedStrPart) Parse(s parsing.State) parsing.State {
	return More(Not(Choice(Word(`"`), Word(`\`)))).Parse(s)
}

func (*escapedStrPart) Parse(s parsing.State) parsing.State {
	return Seq(
		Word(`\`),
		Choice(
			Not(Choice(Word("x"), Word("u"), OctDigit())),
			Occur(OctDigit(), 1, 3),
			Seq(Word("x"), Times(HexDigit(), 2)),
			Seq(Word("u"), Times(HexDigit(), 4)),
			Seq(Word("u{"), More(HexDigit()), Word("}")),
		),
	).Parse(s)
}

type b struct{}

func Bool() parsing.Parser { return (*b)(nil) }

func (*b) Parse(s parsing.State) parsing.State { return Choice(Word("False"), Word("True")).Parse(s) }

type unit struct{}

func Unit() parsing.Parser { return (*unit)(nil) }

func (*unit) Parse(s parsing.State) parsing.State { return Word("()").Parse(s) }

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
	return s.WithError(pos)
}

type many struct {
	more   bool
	parser parsing.Parser
}

func Many(parser parsing.Parser) parsing.Parser { return &many{parser: parser} }
func More(parser parsing.Parser) parsing.Parser { return &many{true, parser} }

func (p *many) Parse(s parsing.State) parsing.State {
	var (
		pos, ln, col int
		prev         *parsing.Span
		hasOne       bool
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
		hasOne = true
	}
	if p.more && !hasOne {
		return s.WithError(pos)
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

type times struct {
	parser parsing.Parser
	n      int
}

func Times(parser parsing.Parser, n int) parsing.Parser {
	if n < 2 {
		panic("at least 2 times expected")
	}
	return &times{parser, n}
}

func (p *times) Parse(s parsing.State) parsing.State {
	for i := 0; i < p.n; i++ {
		if s = p.parser.Parse(s); s.IsError() {
			return s.WithError(s.Pos())
		}
	}
	return s
}

type occur struct {
	parser   parsing.Parser
	from, to int
}

func Occur(parser parsing.Parser, from, to int) parsing.Parser {
	if from < 0 || to < 0 || from >= to {
		panic(fmt.Sprintf("invalid occur range [%d, %d]", from, to))
	}
	return &occur{parser, from, to}
}

func (p *occur) Parse(s parsing.State) parsing.State {
	var (
		start        = s.Pos()
		pos, ln, col int
		prev         *parsing.Span
	)
	var n int
	for {
		pos, ln, col, prev = s.Dump()
		if s = p.parser.Parse(s); s.IsError() {
			s.Restore(pos, ln, col, prev)
			break
		}
		n++
	}
	if p.from <= n && n <= p.to {
		return s
	}
	return s.WithError(start)
}

type not struct{ parser parsing.Parser }

func Not(parser parsing.Parser) parsing.Parser { return &not{parser} }

func (p *not) Parse(s parsing.State) parsing.State {
	pos, ln, col, prev := s.Dump()
	s = p.parser.Parse(s)
	notMatched := s.IsError()
	s.Restore(pos, ln, col, prev)
	if notMatched {
		_, _ = s.Next()
		return s
	}
	return s.WithError(pos)
}
