package ascii

import (
	"strings"
	"unicode"
)

type StringBuilder struct {
	builder strings.Builder
}

func (sb *StringBuilder) Len() int {
	return sb.builder.Len()
}

func (sb *StringBuilder) Cap() int {
	return sb.builder.Cap()
}

func (sb *StringBuilder) Write(b []byte) (int, error) {
	for i, c := range b {
		if c > unicode.MaxASCII {
			return i, ErrRange
		}
	}
	return sb.builder.Write(b)
}

func (sb *StringBuilder) WriteByte(c byte) error {
	if c > unicode.MaxASCII {
		return ErrRange
	}
	return sb.builder.WriteByte(c)
}

func (sb *StringBuilder) WriteString(s string) (int, error) {
	for i := range len(s) {
		if s[i] > unicode.MaxASCII {
			return i, ErrRange
		}
	}
	return sb.builder.WriteString(s)
}

func (sb *StringBuilder) WriteASCIIString(s String) (int, error) {
	return sb.builder.WriteString(s.s)
}

func (sb *StringBuilder) String() string {
	return sb.builder.String()
}

func (sb *StringBuilder) ASCIIString() String {
	return String{s: sb.builder.String()}
}
