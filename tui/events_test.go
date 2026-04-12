package tui

import (
	"testing"

	"github.com/showr/dice-roller/internal/dice"
)

func TestMoveSelectedOutputRoll_ChangesWithinBounds(t *testing.T) {
	a := &app{
		historyOffset: 0,
		historyResults: []interface{}{
			dice.MultiRollResult{Rolls: []dice.Result{{Total: 10}, {Total: 14}, {Total: 12}}},
		},
		selectedOutputRoll: 1,
	}

	if !a.moveSelectedOutputRoll(1) {
		t.Fatalf("expected movement to succeed")
	}
	if a.selectedOutputRoll != 2 {
		t.Fatalf("expected selectedOutputRoll=2, got %d", a.selectedOutputRoll)
	}
	if a.moveSelectedOutputRoll(1) {
		t.Fatalf("expected out-of-range movement to fail")
	}
}
