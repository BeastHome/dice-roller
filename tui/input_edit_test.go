package tui

import (
	"testing"

	"github.com/showr/dice-roller/internal/dice"
)

func TestHandleRecallFromHistory_UsesHighlightedEntryWithoutDuplicatingRolls(t *testing.T) {
	a := &app{
		activePane:    0,
		historyOffset: 0,
		historyResults: []interface{}{
			dice.MultiRollResult{
				Expression: "1d1 rolls=2",
				Rolls: []dice.Result{
					{Expression: "1d1", Total: 1},
					{Expression: "1d1", Total: 1},
				},
			},
		},
	}

	a.handleRecallFromHistory()

	if a.input != "1d1 rolls=2" {
		t.Fatalf("expected recalled input %q, got %q", "1d1 rolls=2", a.input)
	}
	if a.activePane != 0 {
		t.Fatalf("expected input pane to become active, got %d", a.activePane)
	}
}
