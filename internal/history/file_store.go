package history

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/showr/dice-roller/internal/dice"
)

// FileStore implements Store using JSON files on disk.
type FileStore struct {
	currentPath string
}

// NewFileStore creates a new file-based history store.
func NewFileStore() *FileStore {
	return &FileStore{}
}

// NewSession creates a new session file and returns its path and handle.
func (fs *FileStore) NewSession(_ string) (string, *os.File, error) {
	dir := dice.HistoryDir()

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", nil, err
	}

	ts := time.Now().Format("2006-01-02_15-04-05")
	prefix := fmt.Sprintf("%s_%s", ts, generateShortID())

	for i := 0; i < 9999; i++ {
		name := prefix
		if i > 0 {
			name = fmt.Sprintf("%s_%d", prefix, i)
		}
		name += ".json"

		path := filepath.Join(dir, name)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			f, err := os.Create(path)
			if err != nil {
				return "", nil, err
			}
			fs.currentPath = path
			return path, f, nil
		}
	}

	return "", nil, fmt.Errorf("unable to allocate session file")
}

func generateShortID() string {
	b := make([]byte, 3) // 6 hex chars
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%06x", time.Now().UnixNano()&0xFFFFFF)
	}
	return hex.EncodeToString(b)
}

// SetSession sets the current session path.
func (fs *FileStore) SetSession(path string) {
	fs.currentPath = path
}

// CurrentSession returns the current session path.
func (fs *FileStore) CurrentSession() string {
	return fs.currentPath
}

// Append writes a result to the current session file.
func (fs *FileStore) Append(result interface{}) error {
	if fs.currentPath == "" {
		return fmt.Errorf("no active session")
	}

	f, err := os.OpenFile(fs.currentPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	switch r := result.(type) {
	case dice.Result:
		data, err := json.Marshal(r)
		if err != nil {
			return err
		}
		_, err = f.Write(append(data, '\n'))
		return err

	case dice.MultiRollResult:
		wrapper := struct {
			Expression string        `json:"expression"`
			Rolls      []dice.Result `json:"rolls"`
			Summary    string        `json:"summary"`
		}{
			Expression: dice.FormatMultiExpression(r.Expression, len(r.Rolls)),
			Rolls:      r.Rolls,
			Summary:    r.Summary,
		}

		data, err := json.Marshal(wrapper)
		if err != nil {
			return err
		}
		_, err = f.Write(append(data, '\n'))
		return err

	default:
		return fmt.Errorf("Append: unsupported result type %T", result)
	}
}

// Load reads all results from a session file.
func (fs *FileStore) Load(path string) ([]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	out := make([]interface{}, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try multi-roll wrapper format
		var wrapper struct {
			Expression string        `json:"expression"`
			Rolls      []dice.Result `json:"rolls"`
			Summary    string        `json:"summary"`
		}

		if err := json.Unmarshal([]byte(line), &wrapper); err == nil && len(wrapper.Rolls) > 0 {
			out = append(out, dice.MultiRollResult{
				Expression: wrapper.Expression,
				Rolls:      wrapper.Rolls,
				Summary:    wrapper.Summary,
			})
			continue
		}

		// Try single-roll format
		var single dice.Result
		if err := json.Unmarshal([]byte(line), &single); err == nil && single.Expression != "" {
			out = append(out, single)
			continue
		}

		// Invalid line — inject synthetic entry
		out = append(out, fmt.Sprintf("Invalid history entry in %s: %s", filepath.Base(path), line))
	}

	return out, nil
}
