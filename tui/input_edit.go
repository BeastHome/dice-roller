package tui

import (
	"unicode/utf8"

	"github.com/showr/dice-roller/internal/dice"
)

func (a *app) clearInput() {
	a.input = ""
	a.cursorPos = 0
}

func (a *app) setInputText(text string) {
	a.input = text
	a.cursorPos = utf8.RuneCountInString(text)
}

// ------------------------------------------------------------
// Input Editing Helpers (Rune-safe)
// ------------------------------------------------------------
func (a *app) insertRune(r rune) {
	runes := []rune(a.input)

	if a.cursorPos < 0 {
		a.cursorPos = 0
	}
	if a.cursorPos > len(runes) {
		a.cursorPos = len(runes)
	}

	runes = append(runes[:a.cursorPos], append([]rune{r}, runes[a.cursorPos:]...)...)
	a.input = string(runes)
	a.cursorPos++
}

func (a *app) deleteRuneLeft() {
	if a.cursorPos == 0 {
		return
	}

	runes := []rune(a.input)
	if a.cursorPos > len(runes) {
		a.cursorPos = len(runes)
	}

	pos := a.cursorPos
	runes = append(runes[:pos-1], runes[pos:]...)
	a.input = string(runes)
	a.cursorPos--
}

// ------------------------------------------------------------
// Input History Navigation
// ------------------------------------------------------------
func (a *app) handleUpArrow() {
	if a.activePane == 0 && a.recallPreviousInput() {
		a.redrawIfNeeded()
		return
	}

	if a.activePane == 2 && a.moveHistorySelection(-1) {
		a.redrawIfNeeded()
	}
}

func (a *app) handleDownArrow() {
	if a.activePane == 0 && a.recallNextInput() {
		a.redrawIfNeeded()
		return
	}

	if a.activePane == 2 && a.moveHistorySelection(1) {
		a.redrawIfNeeded()
	}
}

func (a *app) recallPreviousInput() bool {
	if len(a.inputHistory) == 0 {
		return false
	}

	if a.inputHistoryIndex == -1 {
		a.inputHistoryIndex = len(a.inputHistory) - 1
	} else if a.inputHistoryIndex > 0 {
		a.inputHistoryIndex--
	}

	a.setInputText(a.inputHistory[a.inputHistoryIndex])
	return true
}

func (a *app) recallNextInput() bool {
	if a.inputHistoryIndex < 0 {
		return false
	}

	a.inputHistoryIndex++
	if a.inputHistoryIndex >= len(a.inputHistory) {
		a.inputHistoryIndex = -1
		a.clearInput()
		return true
	}

	a.setInputText(a.inputHistory[a.inputHistoryIndex])
	return true
}

func (a *app) moveHistorySelection(delta int) bool {
	next := a.historyOffset + delta
	if next < 0 || next >= len(a.historyLines) {
		return false
	}

	a.historyOffset = next
	a.syncOutputFromHistory()
	return true
}

// ------------------------------------------------------------
// Ctrl+R — Recall Expression From History
// ------------------------------------------------------------
func (a *app) handleRecallFromHistory() {
	if a.historyOffset < 0 || a.historyOffset >= len(a.historyResults) {
		return
	}

	var expr string

	switch v := a.historyResults[a.historyOffset].(type) {
	case dice.Result:
		expr = v.Expression
	case dice.MultiRollResult:
		expr = dice.FormatMultiExpression(v.Expression, len(v.Rolls))
	default:
		return
	}

	a.setInputText(expr)
	a.activePane = 0
	a.redrawIfNeeded()
}
