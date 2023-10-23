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

func Newline(s parsing.State) parsing.State {
	return Word(string(s.Options().Newline()))(s)
}

func Entry(s parsing.State) parsing.State {
	parsers := []parsing.ParserFunc{Newline}
	for i := 0; i < s.Depth()+1; i++ {
		parsers = append(parsers, Word(s.IndentWord()))
	}
	if s = parsing.Atom(Seq(parsers...))(s); !s.IsError() {
		return s.WithEntry()
	}
	return s
}

func Indent(s parsing.State) parsing.State {
	var parsers []parsing.ParserFunc
	for i := 0; i < s.Depth(); i++ {
		parsers = append(parsers, Word(s.IndentWord()))
	}
	return parsing.Atom(Seq(parsers...))(s)
}

func Exit(s parsing.State) parsing.State { return Option(Newline)(s).WithExit() }

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
			s.Eat(r)
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

			s.Eat(r)
		}

		if start == s.Pos() {
			return s.WithError(start)
		}
		return s.WithSpan(start)
	}
}

var (
	Lowercase = id(func(r rune) bool { return unicode.IsLower(r) || r == '_' })
	Uppercase = id(func(r rune) bool { return unicode.IsUpper(r) || r == '_' })
)

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

var (
	dec             = decDigits
	decDigits       = Seq(decDigit, Many(Seq(Option(Word("_")), decDigit)))
	decDigit        = Range('0', '9')
	decNonZeroDigit = Range('1', '9')
	decPart         = Choice(
		Word("0"),
		Seq(Option(Word("0")), decNonZeroDigit, Option(Seq(Option(Word("_")), decDigits))),
	)
	expPart = Seq(
		Option(Choice(Word("e"), Word("E"))),
		Option(Choice(Word("-"), Word("+"))),
		decDigits,
	)

	bin      = Seq(Choice(Word("0b"), Word("0B")), binDigit, Many(Seq(Option(Word("_")), binDigit)))
	binDigit = Range('0', '1')

	oct      = Seq(Choice(Word("0o"), Word("0O")), octDigit, Many(Seq(Option(Word("_")), octDigit)))
	octDigit = Range('0', '7')

	hex      = Seq(Choice(Word("0x"), Word("0X")), hexDigit, Many(Seq(Option(Word("_")), hexDigit)))
	hexDigit = Choice(decDigit, Range('a', 'f'), Range('A', 'F'))
)

func Int(s parsing.State) parsing.State {
	start := s.Pos()
	if s = parsing.Atom(Choice(dec, bin, oct, hex))(s); !s.IsError() {
		s.WithSpan(start)
	}
	return s
}

func Float(s parsing.State) parsing.State {
	start := s.Pos()
	if s = parsing.Atom(Choice(
		Seq(decPart, expPart),
		Seq(Word("."), decDigits, Option(expPart)),
		Seq(decPart, Word("."), Option(decDigits), Option(expPart)),
	))(s); !s.IsError() {
		return s.WithSpan(start)
	}
	return s
}

func Str(s parsing.State) parsing.State {
	start := s.Pos()
	if s = parsing.Atom(Seq(Word(`"`), Many(Choice(unescaped, escaped)), Word(`"`)))(s); !s.IsError() {
		s.WithSpan(start)
	}
	return s
}

var (
	unescaped = More(Not(Choice(Word(`"`), Word(`\`))))
	escaped   = Seq(
		Word(`\`),
		Choice(
			Not(Choice(Word("x"), Word("u"), octDigit)),
			Occur(octDigit, 1, 3),
			Seq(Word("x"), Times(hexDigit, 2)),
			Seq(Word("u"), Times(hexDigit, 4)),
			Seq(Word("u{"), More(hexDigit), Word("}")),
		),
	)
)

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
			s.Next()
			return s
		}
		return s.WithError(pos)
	}
}
