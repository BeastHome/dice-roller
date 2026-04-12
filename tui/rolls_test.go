package tui

import (
	"os"
	"testing"

	"github.com/showr/dice-roller/internal/dice"
	"github.com/showr/dice-roller/internal/presentation"
)

type stubStore struct {
	currentPath     string
	newSessionCalls int
	sessionExpr     string
}

func (s *stubStore) Append(result interface{}) error {
	return nil
}

func (s *stubStore) Load(path string) ([]interface{}, error) {
	return nil, nil
}

func (s *stubStore) NewSession(expr string) (string, *os.File, error) {
	s.newSessionCalls++
	s.sessionExpr = expr
	f, err := os.CreateTemp("", "dice-roller-history-*.json")
	if err != nil {
		return "", nil, err
	}
	s.currentPath = f.Name()
	return f.Name(), f, nil
}

func (s *stubStore) SetSession(path string) {
	s.currentPath = path
}

func (s *stubStore) CurrentSession() string {
	return s.currentPath
}

func TestHandleEnter_HelpDoesNotCreateSessionFile(t *testing.T) {
	store := &stubStore{}
	a := &app{
		input:        "--help",
		historyStore: store,
		outputLines:  []string{},
	}

	a.handleEnter()

	if store.newSessionCalls != 0 {
		t.Fatalf("expected no session file for --help, got %d creation call(s)", store.newSessionCalls)
	}
}

func TestHandleEnter_MultiRollUsesRawInputForSessionName(t *testing.T) {
	store := &stubStore{}
	a := &app{
		input:        "1d1 rolls=3",
		engine:       dice.NewEngineWithOptions(dice.EngineOptions{Seed: 1}),
		historyStore: store,
		formatter:    presentation.NewDefaultFormatter(),
		outputLines:  []string{},
	}

	a.handleEnter()

	if store.sessionExpr != "1d1 rolls=3" {
		t.Fatalf("expected session name input %q, got %q", "1d1 rolls=3", store.sessionExpr)
	}
}

func TestHandleConsolidatedMultiRoll_StoresExpressionWithRollCount(t *testing.T) {
	store := &stubStore{}
	a := &app{
		historyStore: store,
		formatter:    presentation.NewDefaultFormatter(),
		outputLines:  []string{},
	}

	mr := &dice.MultiRollResult{
		Expression: "1d1",
		Rolls: []dice.Result{
			{Expression: "1d1", Rolls: []int{1}, Kept: []int{1}, Total: 1},
			{Expression: "1d1", Rolls: []int{1}, Kept: []int{1}, Total: 1},
		},
	}

	a.handleConsolidatedMultiRoll(mr)

	got, ok := a.historyResults[len(a.historyResults)-1].(dice.MultiRollResult)
	if !ok {
		t.Fatalf("expected last history result to be a dice.MultiRollResult")
	}
	if got.Expression != "1d1 rolls=2" {
		t.Fatalf("expected stored expression %q, got %q", "1d1 rolls=2", got.Expression)
	}
}
