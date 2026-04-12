package tui

import (
	"fmt"

	"github.com/showr/dice-roller/internal/dice"
)

// ------------------------------------------------------------
// Sync Output Pane With Selected History Entry
// ------------------------------------------------------------
func (a *app) syncOutputFromHistory() {
	if a.historyOffset < 0 || a.historyOffset >= len(a.historyResults) {
		return
	}

	entry := a.historyResults[a.historyOffset]

	switch v := entry.(type) {

	case dice.Result:
		a.selectedOutputRoll = -1
		// Single-roll: show verbose if available
		if v.Verbose != "" {
			a.outputLines = splitLines(v.Verbose)
		} else {
			a.outputLines = []string{
				fmt.Sprintf("%s -> %d", v.Expression, v.Total),
			}
		}

	case dice.MultiRollResult:
		a.selectedOutputRoll = defaultSelectedRollIndex(v)
		// Multi-roll: show the last roll's verbose output
		if len(v.Rolls) == 0 {
			a.outputLines = []string{"(no rolls)"}
			break
		}

		last := v.Rolls[len(v.Rolls)-1]
		if last.Verbose != "" {
			a.outputLines = splitLines(last.Verbose)
		} else {
			a.outputLines = []string{
				fmt.Sprintf("%s -> %d", dice.FormatMultiExpression(v.Expression, len(v.Rolls)), last.Total),
			}
		}

	default:
		a.outputLines = []string{"(unknown history entry)"}
	}

	a.outputOffset = 0
}

// ------------------------------------------------------------
// History Pane Scrolling (Mouse Wheel)
// ------------------------------------------------------------
func (a *app) handleMouseScroll(delta int) {
	switch a.activePane {

	case 0: // Input pane scroll (help area)
		if delta < 0 && a.helpOffset > 0 {
			a.helpOffset--
			a.redraw()
		}

	case 1: // Output pane scroll
		if delta < 0 && a.outputOffset > 0 {
			a.outputOffset--
		}
		if delta > 0 && a.outputOffset+1 < len(a.outputLines) {
			a.outputOffset++
		}

	case 2: // History pane scroll
		if delta < 0 && a.historyOffset > 0 {
			a.historyOffset--
			a.syncOutputFromHistory()
		}
		if delta > 0 && a.historyOffset+1 < len(a.historyLines) {
			a.historyOffset++
			a.syncOutputFromHistory()
		}
	}
}

// ------------------------------------------------------------
// Utility Helpers (local to history UI)
// ------------------------------------------------------------
func splitLines(s string) []string {
	var out []string
	start := 0
	for i, r := range s {
		if r == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
