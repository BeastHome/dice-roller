package parse

// ParsedInput is the normalized, shared representation of user input
// used by both CLI and TUI.
type ParsedInput struct {
	Expressions []string
	Multi       int
	Verbose     bool
	NoColor     bool
	Seed        *int
}
