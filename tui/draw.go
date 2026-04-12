package tui

import "github.com/gdamore/tcell/v2"

// ------------------------------------------------------------
// Draw Box
// ------------------------------------------------------------
func (a *app) drawBox(x1, y1, x2, y2 int, title string) {
	if x2 <= x1 || y2 <= y1 {
		return
	}

	horiz := '─'
	vert := '│'
	tl := '┌'
	tr := '┐'
	bl := '└'
	br := '┘'

	a.screen.SetContent(x1, y1, tl, nil, tcell.StyleDefault)
	a.screen.SetContent(x2, y1, tr, nil, tcell.StyleDefault)
	a.screen.SetContent(x1, y2, bl, nil, tcell.StyleDefault)
	a.screen.SetContent(x2, y2, br, nil, tcell.StyleDefault)

	for x := x1 + 1; x < x2; x++ {
		a.screen.SetContent(x, y1, horiz, nil, tcell.StyleDefault)
		a.screen.SetContent(x, y2, horiz, nil, tcell.StyleDefault)
	}

	for y := y1 + 1; y < y2; y++ {
		a.screen.SetContent(x1, y, vert, nil, tcell.StyleDefault)
		a.screen.SetContent(x2, y, vert, nil, tcell.StyleDefault)
	}

	if title != "" {
		runes := []rune(title)
		tx := x1 + 2
		for i, r := range runes {
			if tx+i >= x2 {
				break
			}
			a.screen.SetContent(tx+i, y1, r, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
		}
	}
}

// ------------------------------------------------------------
// Quit Dialog
// ------------------------------------------------------------
func (a *app) drawQuitDialog() {
	w, h := a.screen.Size()
	msg := "Quit? (y/n)"
	boxWidth := len(msg) + 4
	boxHeight := 3

	x1 := (w - boxWidth) / 2
	y1 := (h - boxHeight) / 2
	x2 := x1 + boxWidth - 1
	y2 := y1 + boxHeight - 1

	a.drawBox(x1, y1, x2, y2, "")

	style := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	for i, r := range msg {
		a.screen.SetContent(x1+2+i, y1+1, r, nil, style)
	}
}

// ------------------------------------------------------------
// Centered Text Helper
// ------------------------------------------------------------
func (a *app) drawCenteredText(msg string, w, h int) {
	x := (w - len(msg)) / 2
	y := h / 2
	for i, r := range msg {
		a.screen.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}
}

func drawStyledText(s tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
