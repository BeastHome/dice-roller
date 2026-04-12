package tui

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/showr/dice-roller/internal/dice"
)

func (a *app) redrawIfNeeded() {
	if a.screen != nil {
		a.redraw()
	}
}

func (a *app) finishScreen() {
	if a.screen != nil {
		a.screen.Fini()
	}
}

// ------------------------------------------------------------
// Keyboard Event Handling
// ------------------------------------------------------------
func (a *app) handleKeyEvent(e *tcell.EventKey) bool {
	if a.quitting {
		return a.handleQuitDialogKey(e)
	}

	switch e.Key() {
	case tcell.KeyCtrlC:
		a.finishScreen()
		return true

	case tcell.KeyEscape:
		a.quitting = true
		a.redrawIfNeeded()

	case tcell.KeyEnter:
		a.activePane = 0
		a.cursorPos = utf8.RuneCountInString(a.input)
		a.handleEnter()
		a.redrawIfNeeded()

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		a.deleteRuneLeft()
		a.redrawIfNeeded()

	case tcell.KeyLeft:
		if a.activePane == 1 {
			a.moveSelectedOutputRoll(-1)
		} else if a.cursorPos > 0 {
			a.cursorPos--
		}
		a.redrawIfNeeded()

	case tcell.KeyRight:
		if a.activePane == 1 {
			a.moveSelectedOutputRoll(1)
		} else if a.cursorPos < utf8.RuneCountInString(a.input) {
			a.cursorPos++
		}
		a.redrawIfNeeded()

	case tcell.KeyUp:
		a.handleUpArrow()

	case tcell.KeyDown:
		a.handleDownArrow()

	case tcell.KeyPgUp:
		if a.activePane == 0 && a.helpOffset > 0 {
			a.helpOffset--
			a.redrawIfNeeded()
			return false
		}
		if a.activePane == 1 && a.outputOffset > 0 {
			a.outputOffset--
			a.redrawIfNeeded()
		}

	case tcell.KeyPgDn:
		if a.activePane == 0 {
			hlines := buildHelpLines()
			if a.helpOffset+1 < len(hlines) {
				a.helpOffset++
				a.redrawIfNeeded()
				return false
			}
		}
		if a.activePane == 1 && a.outputOffset+1 < len(a.outputLines) {
			a.outputOffset++
			a.redrawIfNeeded()
		}

	case tcell.KeyCtrlR:
		a.handleRecallFromHistory()

	case tcell.KeyRune:
		r := e.Rune()
		if r == 'q' || r == 'Q' {
			a.quitting = true
			a.redrawIfNeeded()
			return false
		}
		a.insertRune(r)
		a.redrawIfNeeded()
	}

	return false
}

func (a *app) handleQuitDialogKey(e *tcell.EventKey) bool {
	switch e.Key() {
	case tcell.KeyRune:
		switch e.Rune() {
		case 'y', 'Y':
			a.finishScreen()
			return true
		case 'n', 'N':
			a.quitting = false
			a.redrawIfNeeded()
		}

	case tcell.KeyEscape:
		a.quitting = false
		a.redrawIfNeeded()
	}

	return false
}

// ------------------------------------------------------------
// Mouse Event Handling
// ------------------------------------------------------------
func (a *app) handleMouseEvent(e *tcell.EventMouse) {
	if a.screen == nil {
		return
	}

	mx, my := e.Position()
	w, h := a.screen.Size()
	_, outputPane, _ := computePaneLayout(w, h)
	a.setActivePaneFromPoint(mx, my, w, h)
	if a.activePane == 1 && e.Buttons()&tcell.Button1 != 0 {
		a.selectOutputRollAt(mx-(outputPane.x1+1), my-(outputPane.y1+1))
	}

	switch e.Buttons() {
	case tcell.WheelUp:
		if a.activePane == 0 && a.helpOffset > 0 {
			a.helpOffset--
			a.redrawIfNeeded()
			return
		}
		a.handleMouseScroll(-1)

	case tcell.WheelDown:
		if a.activePane == 0 {
			hlines := buildHelpLines()
			if a.helpOffset+1 < len(hlines) {
				a.helpOffset++
				a.redrawIfNeeded()
				return
			}
		}
		a.handleMouseScroll(1)
	}

	a.redrawIfNeeded()
}

func (a *app) setActivePaneFromPoint(mx, my, width, height int) {
	inputPane, outputPane, historyPane := computePaneLayout(width, height)

	switch {
	case mx >= inputPane.x1 && mx <= inputPane.x2 && my >= inputPane.y1 && my <= inputPane.y2:
		if a.activePane != 0 {
			a.outputOffset = 0
		}
		a.activePane = 0

	case mx >= outputPane.x1 && mx <= outputPane.x2 && my >= outputPane.y1 && my <= outputPane.y2:
		if a.activePane != 1 {
			a.outputOffset = 0
		}
		a.activePane = 1

	case mx >= historyPane.x1 && mx <= historyPane.x2 && my >= historyPane.y1 && my <= historyPane.y2:
		if a.activePane != 2 {
			a.outputOffset = 0
		}
		a.activePane = 2
	}
}

func (a *app) moveSelectedOutputRoll(delta int) bool {
	mr, ok := a.currentMultiRollResult()
	if !ok {
		return false
	}
	next := clampRollSelection(a.selectedOutputRoll, mr.Rolls) + delta
	if next < 0 || next >= len(mr.Rolls) {
		return false
	}
	a.selectedOutputRoll = next
	return true
}

func (a *app) selectOutputRollAt(col, row int) bool {
	for _, zone := range a.outputHitZones {
		if zone.row == row && col >= zone.startCol && col < zone.endCol {
			a.selectedOutputRoll = zone.rollIndex
			return true
		}
	}
	return false
}

func (a *app) currentMultiRollResult() (dice.MultiRollResult, bool) {
	if a.historyOffset < 0 || a.historyOffset >= len(a.historyResults) {
		return dice.MultiRollResult{}, false
	}
	mr, ok := a.historyResults[a.historyOffset].(dice.MultiRollResult)
	return mr, ok
}
