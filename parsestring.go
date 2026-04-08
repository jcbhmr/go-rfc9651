package rfc9651

import (
	"go.jcbhmr.com/rfc9651/internal/ascii"
)

// https://httpwg.org/specs/rfc9651.html#parse-string
func ParseString(inputString *ascii.String) (unquotedString ascii.String, err error) {
	var outputString ascii.StringBuilder

	if inputString.Get(0) != '"' {
		return ascii.String{}, &UnexpectedCharError{Found: inputString.Get(0), Expected: `'"'`}
	}

	*inputString = inputString.Slice(1, inputString.Len())

	for inputString.Len() > 0 {
		char := inputString.Get(0)
		*inputString = inputString.Slice(1, inputString.Len())

		if char == '\\' {
			if inputString.Len() == 0 {
				return ascii.String{}, &UnexpectedEOIError{Expected: "character after '\\'"}
			}

			nextChar := inputString.Get(0)
			*inputString = inputString.Slice(1, inputString.Len())

			if !(nextChar == '"' || nextChar == '\\') {
				return ascii.String{}, &UnexpectedCharError{Found: nextChar, Expected: `'"' or '\\'`}
			}

			outputString.WriteByte(nextChar)
		} else if char == '"' {
			return outputString.ASCIIString(), nil
		} else if ('\x00' <= char && char <= '\x1F') || ('\x7F' <= char && char <= '\xFF') {
			return ascii.String{}, &UnexpectedCharError{Found: char, Expected: "ASCII character"}
		} else {
			outputString.WriteByte(char)
		}
	}

	return ascii.String{}, &UnexpectedEOIError{Expected: `'"'`}
}
