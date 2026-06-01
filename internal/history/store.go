package history

import (
	"os"

	"github.com/showr/dice-roller/internal/dice"
)

// Store defines the interface for persisting and retrieving dice roll
// results. Implementations can use files, databases, or other backends.
//
// Append is split into typed methods (AppendSingle, AppendMulti) rather
// than a single Append(interface{}) so the compiler verifies the call
// site sends the right type. Load still returns interface{} because a
// stored session file is genuinely heterogeneous (it can contain a mix
// of single and multi-roll entries from one user's session).
type Store interface {
	// AppendSingle records a single-roll result to the current session.
	AppendSingle(r dice.Result) error

	// AppendMulti records a multi-roll result to the current session.
	AppendMulti(mr dice.MultiRollResult) error

	// Load retrieves all results from a stored session file or path.
	// The returned slice may contain dice.Result, dice.MultiRollResult,
	// or a string entry describing an invalid line that couldn't be
	// parsed back.
	Load(path string) ([]interface{}, error)

	// NewSession creates a new session and returns its path and a file
	// handle. The session path is used for subsequent Append operations.
	NewSession(expr string) (string, *os.File, error)

	// SetSession sets the current session path for Append operations.
	SetSession(path string)

	// CurrentSession returns the current session path.
	CurrentSession() string
}
