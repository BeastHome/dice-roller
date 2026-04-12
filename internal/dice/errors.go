package dice

import "fmt"

// ParseError represents a failure during expression parsing.
type ParseError struct {
	Msg string
}

func (e ParseError) Error() string { return fmt.Sprintf("parse error: %s", e.Msg) }

// EvalError represents a failure during evaluation.
type EvalError struct {
	Msg string
}

func (e EvalError) Error() string { return fmt.Sprintf("evaluation error: %s", e.Msg) }
