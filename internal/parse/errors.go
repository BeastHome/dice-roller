package parse

import "fmt"

type ParseError struct {
	Token   string
	Message string
	Hint    string
}

func (e *ParseError) Error() string {
	if e.Hint == "" {
		return fmt.Sprintf("parse error near %q: %s", e.Token, e.Message)
	}
	return fmt.Sprintf("parse error near %q: %s\nHint: %s", e.Token, e.Message, e.Hint)
}

func newParseError(token, msg, hint string) error {
	return &ParseError{
		Token:   token,
		Message: msg,
		Hint:    hint,
	}
}
