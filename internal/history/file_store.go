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

	"github.com/showr/dice-roller/dice"
)

// FileStore implements Store using JSON files on disk.
//
// Each Append opens, writes one line, and closes the file — appropriate
// for interactive or low-frequency persistence (the TUI does one Append
// per user gesture). High-throughput embedded consumers should implement
// their own Store with a long-lived open handle and explicit flush.
type FileStore struct {
	// baseDir overrides the default dice.HistoryDir() when non-empty.
	// Set via NewFileStoreInDir; tests and embedded consumers with
	// their own history root use this.
	baseDir     string
	currentPath string
}

// NewFileStore creates a file-based history store rooted at the OS
// default location returned by dice.HistoryDir().
func NewFileStore() *FileStore {
	return &FileStore{}
}

// NewFileStoreInDir creates a file-based history store rooted at the
// given directory. Used by tests and by embedded consumers that want
// their own history root.
func NewFileStoreInDir(dir string) *FileStore {
	return &FileStore{baseDir: dir}
}

// dir returns the base directory for session files.
func (fs *FileStore) dir() string {
	if fs.baseDir != "" {
		return fs.baseDir
	}
	return dice.HistoryDir()
}

// NewSession creates a new session file and returns its path and handle.
func (fs *FileStore) NewSession(_ string) (string, *os.File, error) {
	dir := fs.dir()

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

// AppendSingle writes a single-roll result to the current session file.
func (fs *FileStore) AppendSingle(r dice.Result) error {
	return fs.appendJSON(r)
}

// AppendMulti writes a multi-roll result to the current session file.
// The expression is normalized via dice.FormatMultiExpression to ensure
// the rolls=N suffix is present exactly once.
func (fs *FileStore) AppendMulti(mr dice.MultiRollResult) error {
	wrapper := struct {
		Expression string        `json:"expression"`
		Rolls      []dice.Result `json:"rolls"`
		Summary    string        `json:"summary"`
	}{
		Expression: dice.FormatMultiExpression(mr.Expression, len(mr.Rolls)),
		Rolls:      mr.Rolls,
		Summary:    mr.Summary,
	}
	return fs.appendJSON(wrapper)
}

// appendJSON serializes v as a single JSON line appended to the current
// session file.
func (fs *FileStore) appendJSON(v interface{}) error {
	if fs.currentPath == "" {
		return fmt.Errorf("no active session")
	}
	f, err := os.OpenFile(fs.currentPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = f.Write(append(data, '\n'))
	return err
}

// Load reads all results from a session file. Lines that fail to
// parse become string entries describing the failure, so the caller
// can surface them rather than silently dropping data.
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

		// Try multi-roll wrapper format first (discriminator: non-empty Rolls).
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

		// Single-roll format (discriminator: non-empty Expression).
		var single dice.Result
		if err := json.Unmarshal([]byte(line), &single); err == nil && single.Expression != "" {
			out = append(out, single)
			continue
		}

		out = append(out, fmt.Sprintf("Invalid history entry in %s: %s", filepath.Base(path), line))
	}

	return out, nil
}
