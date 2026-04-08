package ascii

import (
	"strconv"
	"strings"
	"unicode"
)

func IsValid(s string) bool {
	for i := range len(s) {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

var ErrRange = strconv.ErrRange

type String struct {
	s string
}

func (s String) String() string {
	return s.s
}

func (s String) Len() int {
	return len(s.s)
}

func (s String) Get(i int) byte {
	return s.s[i]
}

func (s String) Slice(i, j int) String {
	return String{s: s.s[i:j]}
}

func NewString(b []byte) (String, error) {
	for _, c := range b {
		if c > unicode.MaxASCII {
			return String{}, ErrRange
		}
	}
	return String{s: string(b)}, nil
}

func (s String) TrimPrefix(prefix string) String {
	return String{s: strings.TrimPrefix(s.s, prefix)}
}

func (s String) TrimLeft(cutset string) String {
	return String{s: strings.TrimLeft(s.s, cutset)}
}
