package history

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/showr/dice-roller/internal/dice"
)

func newSessionInTempDir(t *testing.T) (*FileStore, string) {
	t.Helper()
	fs := NewFileStoreInDir(t.TempDir())
	path, f, err := fs.NewSession("test")
	if err != nil {
		t.Fatalf("NewSession: %v", err)
	}
	f.Close()
	return fs, path
}

func TestFileStore_NewSessionCreatesFileInBaseDir(t *testing.T) {
	dir := t.TempDir()
	fs := NewFileStoreInDir(dir)
	path, f, err := fs.NewSession("4d6")
	if err != nil {
		t.Fatalf("NewSession returned error: %v", err)
	}
	defer f.Close()

	if filepath.Dir(path) != dir {
		t.Fatalf("expected session under %q, got %q", dir, path)
	}
	if fs.CurrentSession() != path {
		t.Fatalf("CurrentSession()=%q, NewSession returned %q", fs.CurrentSession(), path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("session file should exist: %v", err)
	}
}

func TestFileStore_AppendWithoutSessionReturnsError(t *testing.T) {
	fs := NewFileStoreInDir(t.TempDir())
	err := fs.AppendSingle(dice.Result{Expression: "1d6", Total: 4})
	if err == nil || !strings.Contains(err.Error(), "no active session") {
		t.Fatalf("expected 'no active session' error, got %v", err)
	}
}

func TestFileStore_RoundTripSingle(t *testing.T) {
	fs, path := newSessionInTempDir(t)
	r := dice.Result{
		Expression: "4d6k3",
		Rolls:      []int{6, 5, 4, 2},
		Kept:       []int{6, 5, 4},
		Dropped:    []int{2},
		Total:      15,
	}
	if err := fs.AppendSingle(r); err != nil {
		t.Fatalf("AppendSingle: %v", err)
	}

	entries, err := fs.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d: %v", len(entries), entries)
	}
	got, ok := entries[0].(dice.Result)
	if !ok {
		t.Fatalf("expected dice.Result, got %T", entries[0])
	}
	if got.Expression != r.Expression || got.Total != r.Total {
		t.Fatalf("round-trip mismatch: want %#v, got %#v", r, got)
	}
}

func TestFileStore_RoundTripMulti(t *testing.T) {
	fs, path := newSessionInTempDir(t)
	mr := dice.MultiRollResult{
		Expression: "1d6 rolls=3",
		Rolls: []dice.Result{
			{Expression: "1d6", Total: 3, Rolls: []int{3}, Kept: []int{3}},
			{Expression: "1d6", Total: 5, Rolls: []int{5}, Kept: []int{5}},
			{Expression: "1d6", Total: 1, Rolls: []int{1}, Kept: []int{1}},
		},
		Summary: "1d6 rolls=3 | avg=3.00 | min=1 | max=5",
	}
	if err := fs.AppendMulti(mr); err != nil {
		t.Fatalf("AppendMulti: %v", err)
	}

	entries, err := fs.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	got, ok := entries[0].(dice.MultiRollResult)
	if !ok {
		t.Fatalf("expected dice.MultiRollResult, got %T", entries[0])
	}
	if len(got.Rolls) != 3 {
		t.Fatalf("expected 3 rolls, got %d", len(got.Rolls))
	}
	if got.Summary != mr.Summary {
		t.Fatalf("summary mismatch: want %q, got %q", mr.Summary, got.Summary)
	}
}

func TestFileStore_MixedFilePreservesOrderAndTypes(t *testing.T) {
	fs, path := newSessionInTempDir(t)

	if err := fs.AppendSingle(dice.Result{Expression: "d20", Total: 17}); err != nil {
		t.Fatalf("AppendSingle: %v", err)
	}
	if err := fs.AppendMulti(dice.MultiRollResult{
		Expression: "1d1 rolls=2",
		Rolls:      []dice.Result{{Expression: "1d1", Total: 1}, {Expression: "1d1", Total: 1}},
	}); err != nil {
		t.Fatalf("AppendMulti: %v", err)
	}
	if err := fs.AppendSingle(dice.Result{Expression: "d6", Total: 4}); err != nil {
		t.Fatalf("AppendSingle: %v", err)
	}

	entries, err := fs.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if _, ok := entries[0].(dice.Result); !ok {
		t.Fatalf("entry 0: expected dice.Result, got %T", entries[0])
	}
	if _, ok := entries[1].(dice.MultiRollResult); !ok {
		t.Fatalf("entry 1: expected dice.MultiRollResult, got %T", entries[1])
	}
	if _, ok := entries[2].(dice.Result); !ok {
		t.Fatalf("entry 2: expected dice.Result, got %T", entries[2])
	}
}

func TestFileStore_LoadInvalidLineInjectsStringEntry(t *testing.T) {
	fs, path := newSessionInTempDir(t)

	// Write a mix: valid entry, then a garbage line, then another valid entry.
	if err := fs.AppendSingle(dice.Result{Expression: "d20", Total: 10}); err != nil {
		t.Fatalf("AppendSingle: %v", err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("open for garbage write: %v", err)
	}
	if _, err := f.WriteString("this-is-not-json\n"); err != nil {
		t.Fatalf("write garbage: %v", err)
	}
	f.Close()
	if err := fs.AppendSingle(dice.Result{Expression: "d6", Total: 5}); err != nil {
		t.Fatalf("AppendSingle: %v", err)
	}

	entries, err := fs.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries (2 valid + 1 synthetic), got %d", len(entries))
	}
	synthetic, ok := entries[1].(string)
	if !ok {
		t.Fatalf("entry 1: expected string synthetic entry, got %T", entries[1])
	}
	if !strings.Contains(synthetic, "Invalid history entry") {
		t.Fatalf("expected synthetic 'Invalid history entry' message, got %q", synthetic)
	}
}

func TestFileStore_LoadEmptyFileReturnsNoEntries(t *testing.T) {
	_, path := newSessionInTempDir(t)
	// File was created empty by NewSession.
	entries, err := NewFileStoreInDir(filepath.Dir(path)).Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries from empty file, got %d", len(entries))
	}
}

func TestFileStore_SetSessionThenAppendUsesGivenPath(t *testing.T) {
	dir := t.TempDir()
	fs := NewFileStoreInDir(dir)

	// Create a file by hand.
	path := filepath.Join(dir, "manual-session.json")
	if err := os.WriteFile(path, nil, 0644); err != nil {
		t.Fatalf("create manual file: %v", err)
	}
	fs.SetSession(path)
	if fs.CurrentSession() != path {
		t.Fatalf("CurrentSession=%q, want %q", fs.CurrentSession(), path)
	}

	if err := fs.AppendSingle(dice.Result{Expression: "d4", Total: 3}); err != nil {
		t.Fatalf("AppendSingle: %v", err)
	}
	entries, err := fs.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}
