package tui

import (
	"strings"
	"testing"

	"github.com/showr/dice-roller/internal/dice"
)

func TestBuildMultiRollOutput_IncludesStatsTotalsAndDetail(t *testing.T) {
	mr := dice.MultiRollResult{
		Expression: "4d6k3 rolls=3",
		Rolls: []dice.Result{
			{Expression: "4d6k3", Rolls: []int{6, 5, 4, 2}, Kept: []int{6, 5, 4}, Dropped: []int{2}, Total: 15},
			{Expression: "4d6k3", Rolls: []int{6, 4, 3, 1}, Kept: []int{6, 4, 3}, Dropped: []int{1}, Total: 13},
			{Expression: "4d6k3", Rolls: []int{5, 5, 2, 1}, Kept: []int{5, 5, 2}, Dropped: []int{1}, Total: 12},
		},
	}

	lines, hitZones, selected := buildMultiRollOutput(mr, -1)
	if selected != 0 {
		t.Fatalf("expected highest total roll to be selected, got %d", selected)
	}
	if len(hitZones) != 3 {
		t.Fatalf("expected 3 hit zones, got %d", len(hitZones))
	}
	if !strings.Contains(lines[0].text, "AVG") || !strings.Contains(lines[0].text, "MED") || !strings.Contains(lines[0].text, "SD") {
		t.Fatalf("expected stats in header, got %q", lines[0].text)
	}
	if strings.Contains(lines[1].text, "Roll 1") {
		t.Fatalf("did not expect numbered roll labels in totals line: %q", lines[1].text)
	}
	if !strings.HasPrefix(lines[1].text, "Totals: ") {
		t.Fatalf("expected totals line, got %q", lines[1].text)
	}
	if !strings.HasPrefix(lines[2].text, "Trend: ") {
		t.Fatalf("expected trend line, got %q", lines[2].text)
	}
	if !strings.HasPrefix(lines[3].text, "Freq:  ") {
		t.Fatalf("expected frequency line, got %q", lines[3].text)
	}
	foundDetail := false
	for _, line := range lines {
		if strings.HasPrefix(line.text, "Selected roll detail") {
			foundDetail = true
			break
		}
	}
	if !foundDetail {
		t.Fatalf("expected selected roll detail block")
	}
}

func TestClampRollSelection_PicksHighestWhenOutOfRange(t *testing.T) {
	rolls := []dice.Result{{Total: 4}, {Total: 8}, {Total: 6}}
	if got := clampRollSelection(99, rolls); got != 1 {
		t.Fatalf("expected highest total index 1, got %d", got)
	}
}
