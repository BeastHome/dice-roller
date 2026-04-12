package history

import "os"

// Store defines the interface for persisting and retrieving dice roll results.
// Implementations can use files, databases, or other backends.
type Store interface {
	// Append adds a single or multi-roll result to the current session.
	Append(result interface{}) error

	// Load retrieves all results from a stored session file or path.
	Load(path string) ([]interface{}, error)

	// NewSession creates a new session and returns its path and a file handle.
	// The session path is used for subsequent Append operations.
	NewSession(expr string) (string, *os.File, error)

	// SetSession sets the current session path for Append operations.
	SetSession(path string)

	// CurrentSession returns the current session path.
	CurrentSession() string
}
