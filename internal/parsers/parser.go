package parsers

import (
	"fmt"
	"unicode"

	"github.com/SuperPythonic/SuperPythonic/pkg/parsing"
)

func Start(s parsing.State) parsing.State {
	if !s.IsSOI() {
		return s.WithError(s.Pos())
	}
	return s
}

func End(s parsing.State) parsing.State {
	if !s.IsEOI() {
		return s.WithError(s.Pos())
	}
	return s
}

func Range(from, to rune) parsing.ParserFunc {
	if from >= to {
		panic(fmt.Sprintf("invalid rune range [%c, %c]", from, to))
	}
	return func(s parsing.State) parsing.State {
		r, ok := s.Peek()
		if !ok {
			return s.WithError(s.Pos())
		}
		if from <= r && r <= to {
			_, _ = s.Next()
			return s
		}
		return s.WithError(s.Pos())
	}
}

func Word(w string) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		start := s.Pos()
		for _, c := range w {
			if !s.Eat(c) {
				return s.WithError(start)
			}
		}
		return s.WithSpan(start)
	}
}

func id(isValid func(r rune) bool) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
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

			if !unicode.IsDigit(r) && !isValid(r) {
				break
			}

			_, _ = s.Next()
		}

		if start == s.Pos() {
			return s.WithError(start)
		}
		return s.WithSpan(start)
	}
}

func Lowercase(s parsing.State) parsing.State {
	return id(func(r rune) bool { return unicode.IsLower(r) || r == '_' })(s)
}

func Uppercase(s parsing.State) parsing.State {
	return id(func(r rune) bool { return unicode.IsUpper(r) || r == '_' })(s)
}

func CamelBack(s parsing.State) parsing.State {
	first := true
	return id(func(r rune) bool {
		if first {
			first = false
			return unicode.IsLower(r)
		}
		return unicode.IsLetter(r)
	})(s)
}

func CamelCase(s parsing.State) parsing.State {
	first := true
	return id(func(r rune) bool {
		if first {
			first = false
			return unicode.IsUpper(r)
		}
		return unicode.IsLetter(r)
	})(s)
}

func octDigit(s parsing.State) parsing.State { return Range('0', '7')(s) }

func hexDigit(s parsing.State) parsing.State {
	return Choice(Range('0', '9'), Range('a', 'f'), Range('A', 'F'))(s)
}

func Int(s parsing.State) parsing.State {
	start := s.Pos()
	if s = parsing.Atom(Choice(bin, oct, hex))(s); !s.IsError() {
		s.WithSpan(start)
	}
	return s
}

func bin(s parsing.State) parsing.State {
	return Seq(
		Choice(Word("0b"), Word("0B")),
		Range('0', '1'),
		Many(Seq(Option(Word("_")), Range('0', '1'))),
	)(s)
}

func oct(s parsing.State) parsing.State {
	return Seq(
		Choice(Word("0o"), Word("0O")),
		octDigit,
		Many(Seq(Option(Word("_")), octDigit)),
	)(s)
}

func hex(s parsing.State) parsing.State {
	return Seq(
		Choice(Word("0x"), Word("0X")),
		hexDigit,
		Many(Seq(Option(Word("_")), hexDigit)),
	)(s)
}

func Str(s parsing.State) parsing.State {
	start := s.Pos()
	if s = parsing.Atom(Seq(
		Word(`"`),
		Many(Choice(unescapedStrPart, escapedStrPart)),
		Word(`"`),
	))(s); !s.IsError() {
		s.WithSpan(start)
	}
	return s
}

func unescapedStrPart(s parsing.State) parsing.State {
	return More(Not(Choice(Word(`"`), Word(`\`))))(s)
}

func escapedStrPart(s parsing.State) parsing.State {
	return Seq(
		Word(`\`),
		Choice(
			Not(Choice(Word("x"), Word("u"), octDigit)),
			Occur(octDigit, 1, 3),
			Seq(Word("x"), Times(hexDigit, 2)),
			Seq(Word("u"), Times(hexDigit, 4)),
			Seq(Word("u{"), More(hexDigit), Word("}")),
		),
	)(s)
}

func Seq(parsers ...parsing.ParserFunc) parsing.ParserFunc {
	if len(parsers) == 0 {
		panic("at least 1 parser expected")
	}
	return func(s parsing.State) parsing.State {
		for i, parser := range parsers {
			if s = parser(s); s.IsError() {
				return s
			}
			if i+1 < len(parsers) {
				s.SkipSpaces()
			}
		}
		return s
	}
}

func Choice(parsers ...parsing.ParserFunc) parsing.ParserFunc {
	if len(parsers) == 0 {
		panic("at least 1 choice expected")
	}
	return func(s parsing.State) parsing.State {
		pos, ln, col, prev := s.Dump()
		// TODO: Collect all failed states for error messages.
		for _, parser := range parsers {
			if s = parser(s); !s.IsError() {
				return s
			}
			s.Restore(pos, ln, col, prev)
		}
		return s.WithError(pos)
	}
}

func More(parser parsing.ParserFunc) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		var (
			pos, ln, col int
			prev         *parsing.Span
			hasOne       bool
		)
		for {
			pos, ln, col, prev = s.Dump()
			if s = parser(s); s.IsError() {
				s.Restore(pos, ln, col, prev)
				break
			}
			if s.IsEOI() {
				break
			}
			s.SkipSpaces()
			hasOne = true
		}
		if !hasOne {
			return s.WithError(pos)
		}
		return s
	}
}

func Many(parser parsing.ParserFunc) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		var (
			pos, ln, col int
			prev         *parsing.Span
		)
		for {
			pos, ln, col, prev = s.Dump()
			if s = parser(s); s.IsError() {
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
}

func Option(parser parsing.ParserFunc) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		pos, ln, col, prev := s.Dump()
		if s = parser(s); s.IsError() {
			s.Restore(pos, ln, col, prev)
		}
		return s
	}
}

func Times(parser parsing.ParserFunc, n int) parsing.ParserFunc {
	if n < 2 {
		panic("at least 2 times expected")
	}
	return func(s parsing.State) parsing.State {
		for i := 0; i < n; i++ {
			if s = parser(s); s.IsError() {
				return s.WithError(s.Pos())
			}
		}
		return s
	}
}

func Occur(parser parsing.ParserFunc, from, to int) parsing.ParserFunc {
	if from < 0 || to < 0 || from >= to {
		panic(fmt.Sprintf("invalid occur range [%d, %d]", from, to))
	}
	return func(s parsing.State) parsing.State {
		var (
			start        = s.Pos()
			pos, ln, col int
			prev         *parsing.Span
		)
		var n int
		for {
			pos, ln, col, prev = s.Dump()
			if s = parser(s); s.IsError() {
				s.Restore(pos, ln, col, prev)
				break
			}
			n++
		}
		if from <= n && n <= to {
			return s
		}
		return s.WithError(start)
	}
}

func Not(parser parsing.ParserFunc) parsing.ParserFunc {
	return func(s parsing.State) parsing.State {
		pos, ln, col, prev := s.Dump()
		s = parser(s)
		notMatched := s.IsError()
		s.Restore(pos, ln, col, prev)
		if notMatched {
			_, _ = s.Next()
			return s
		}
		return s.WithError(pos)
	}
}
