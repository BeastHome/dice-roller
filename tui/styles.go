package tui

import "github.com/gdamore/tcell/v2"

type styledLine struct {
	text  string
	style tcell.Style
}

var historyHighlightStyle = tcell.StyleDefault.
	Foreground(tcell.ColorBlack).
	Background(tcell.ColorYellow)

var (
	styleHelpHeading        = tcell.StyleDefault.Foreground(tcell.ColorLightCyan)
	styleHelpLabel          = tcell.StyleDefault.Foreground(tcell.ColorLightBlue)
	styleHelpExample        = tcell.StyleDefault.Foreground(tcell.ColorLightGreen)
	styleRollLabel          = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	styleInputText          = tcell.StyleDefault.Foreground(tcell.ColorWhite)
	styleScrollIndicator    = tcell.StyleDefault.Foreground(tcell.ColorLightPink).Bold(true)
	styleOutputText         = tcell.StyleDefault.Foreground(tcell.ColorLightGreen)
	styleOutputHeading      = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	styleOutputSelectedRoll = tcell.StyleDefault.Foreground(tcell.ColorWhite).Reverse(true)
	styleHistoryText        = tcell.StyleDefault.Foreground(tcell.ColorFuchsia)

	// Semantic output styles (dark mode, black background)
	styleOutputStats    = tcell.StyleDefault.Foreground(tcell.ColorLightCyan)
	styleOutputKept     = tcell.StyleDefault.Foreground(tcell.ColorLightGreen)
	styleOutputDropped  = tcell.StyleDefault.Foreground(tcell.ColorGray)
	styleOutputReroll   = tcell.StyleDefault.Foreground(tcell.ColorYellow)
	styleOutputExploded = tcell.StyleDefault.Foreground(tcell.ColorLightPink).Bold(true)
	styleOutputTotal    = tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	styleOutputSuccess  = tcell.StyleDefault.Foreground(tcell.ColorLightGreen)
)
