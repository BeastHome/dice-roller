package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/showr/dice-roller/internal/dice"
)

func fillRect(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			s.SetContent(x, y, ' ', nil, style)
		}
	}
}

// ------------------------------------------------------------
// Draw Input Pane
// ------------------------------------------------------------
func (a *app) drawInput(x1, y1, x2, y2 int) {
	s := a.screen
	fillRect(s, x1, y1, x2, y2-1, styleHelpLabel)

	height := y2 - y1 + 1
	helpHeight := height - 2
	hlines := buildHelpLines()

	visible := hlines
	if len(hlines) > helpHeight {
		end := min(len(hlines), a.helpOffset+helpHeight)
		visible = hlines[a.helpOffset:end]
	}
	for i, hl := range visible {
		drawStyledText(s, x1, y1+i, hl.text, hl.style)
	}

	sepY := y1 + helpHeight
	for x := x1; x <= x2; x++ {
		s.SetContent(x, sepY, '─', nil, styleHelpHeading)
	}

	if a.helpOffset > 0 {
		s.SetContent(x2, y1, '↑', nil, styleScrollIndicator)
	}
	if a.helpOffset+helpHeight < len(hlines) {
		s.SetContent(x2, sepY-1, '↓', nil, styleScrollIndicator)
	}

	rollY := sepY + 1
	fillRect(s, x1, rollY, x2, rollY, styleInputText)
	drawStyledText(s, x1, rollY, "Roll:", styleRollLabel)
	drawStyledText(s, x1+6, rollY, a.input, styleInputText)
	s.ShowCursor(x1+6+a.cursorPos, rollY)
}

func buildHelpLines() []styledLine {
	lines := make([]styledLine, 0, len(dice.HelpText))

	for _, line := range dice.HelpText {
		switch {
		case line == "":
			lines = append(lines, styledLine{text: "", style: styleHelpLabel})
		case strings.HasSuffix(line, ":"):
			lines = append(lines, styledLine{text: line, style: styleHelpHeading})
		default:
			lines = append(lines, styledLine{text: line, style: styleHelpExample})
		}
	}

	return lines
}

// ------------------------------------------------------------
// Draw Output Pane
// ------------------------------------------------------------
func (a *app) drawOutputWithScroll(x1, y1, x2, y2 int) {
	s := a.screen
	fillRect(s, x1, y1, x2, y2, styleOutputText)
	a.outputHitZones = nil

	width := x2 - x1 + 1
	height := y2 - y1 + 1
	if width <= 0 || height <= 0 {
		return
	}

	styledLines := a.currentOutputStyledLines()

	start := a.outputOffset
	end := min(len(styledLines), start+height)
	row := 0
	for i := start; i < end; i++ {
		line := styledLines[i]
		runes := []rune(line.text)
		if len(runes) > width {
			runes = runes[:width]
		}

		for col := 0; col < width; col++ {
			ch := ' '
			style := styleOutputText
			for _, span := range line.spans {
				if col >= span.start && col < span.end {
					style = span.style
					break
				}
			}
			if col < len(runes) {
				ch = runes[col]
			}
			s.SetContent(x1+col, y1+row, ch, nil, style)
		}
		row++
	}

	if a.outputOffset > 0 {
		s.SetContent(x2, y1, '↑', nil, styleScrollIndicator)
	}
	if a.outputOffset+height < len(styledLines) {
		s.SetContent(x2, y2, '↓', nil, styleScrollIndicator)
	}
}

func (a *app) currentOutputStyledLines() []outputStyledLine {
	if a.historyOffset < 0 || a.historyOffset >= len(a.historyResults) {
		styled := make([]outputStyledLine, len(a.outputLines))
		for i, line := range a.outputLines {
			styled[i] = outputStyledLine{text: line}
		}
		return styled
	}

	styledLines, hitZones, selectedRoll := buildOutputView(a.historyResults[a.historyOffset], a.selectedOutputRoll)
	a.outputLines = a.outputLines[:0]
	for _, line := range styledLines {
		a.outputLines = append(a.outputLines, line.text)
	}
	a.outputHitZones = hitZones
	a.selectedOutputRoll = selectedRoll
	return styledLines
}

// ------------------------------------------------------------
// Draw History Pane
// ------------------------------------------------------------
func (a *app) drawHistoryWithHighlight(x1, y1, x2, y2 int) {
	fillRect(a.screen, x1, y1, x2, y2, styleHistoryText)

	width := x2 - x1 + 1
	height := y2 - y1 + 1
	if width <= 0 || height <= 0 {
		return
	}

	start := a.historyOffset
	if start < 0 {
		start = 0
	}
	end := min(len(a.historyLines), start+height)

	row := 0
	for i := start; i < end; i++ {
		runes := []rune(a.historyLines[i])
		if len(runes) > width {
			runes = runes[:width]
		}

		style := styleHistoryText
		if i == a.historyOffset {
			style = historyHighlightStyle
		}

		for col := 0; col < width; col++ {
			ch := ' '
			if col < len(runes) {
				ch = runes[col]
			}
			a.screen.SetContent(x1+col, y1+row, ch, nil, style)
		}
		row++
	}

	if a.historyOffset > 0 {
		a.screen.SetContent(x2, y1, '↑', nil, styleScrollIndicator)
	}
	if a.historyOffset+height < len(a.historyLines) {
		a.screen.SetContent(x2, y2, '↓', nil, styleScrollIndicator)
	}
}
