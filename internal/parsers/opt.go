package parsers

import "unicode"

type opts struct{}

func (*opts) IsSpace(r rune) bool { return unicode.IsSpace(r) }

func (*opts) IsNewline(r rune) bool { return r == '\n' }
