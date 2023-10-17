package parsers

import "unicode"

type opts struct{}

func (o *opts) IsSpace(r rune) bool { return r != o.Newline() && r != o.Indent() && unicode.IsSpace(r) }
func (*opts) Newline() rune         { return '\n' }
func (*opts) Indent() rune          { return '\t' }
func (*opts) IndentWordN() int      { return 1 }
