package parsers

import "unicode"

type opt struct{}

func (*opt) IsSpace(r rune) bool { return unicode.IsSpace(r) }

func (*opt) IsNewline(r rune) bool { return r == '\n' }
