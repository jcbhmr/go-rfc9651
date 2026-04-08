package rfc9651

import "fmt"

type UnexpectedCharError struct {
	Found    byte
	Expected string
}

func (e *UnexpectedCharError) Error() string {
	return fmt.Sprintf("expected %s, found %q", e.Expected, e.Found)
}

type UnexpectedEOIError struct {
	Expected string
}

func (e *UnexpectedEOIError) Error() string {
	return fmt.Sprintf("expected %s, found end of input", e.Expected)
}
