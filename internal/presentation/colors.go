package presentation

import (
	"fmt"
	"os"
)

// ANSI color codes for terminal output
const (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorDim     = "\033[2m"
	colorReverse = "\033[7m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorMagenta = "\033[35m"
	colorWhite   = "\033[37m"
	colorGray    = "\033[90m"
)

// ColorScheme holds color codes for different roll elements
type ColorScheme struct {
	Reset    string
	Bold     string
	Dim      string
	Reverse  string
	Stats    string // Cyan
	Kept     string // Green
	Dropped  string // Gray
	Reroll   string // Yellow
	Exploded string // Magenta + Bold
	Total    string // White + Bold
	Success  string // Green
}

// GetColorScheme returns the appropriate color scheme based on TTY detection and colorDisabled flag
func GetColorScheme(colorDisabled bool) ColorScheme {
	if colorDisabled || !isTTY() {
		return ColorScheme{
			Reset:    "",
			Bold:     "",
			Dim:      "",
			Reverse:  "",
			Stats:    "",
			Kept:     "",
			Dropped:  "",
			Reroll:   "",
			Exploded: "",
			Total:    "",
			Success:  "",
		}
	}

	return ColorScheme{
		Reset:    colorReset,
		Bold:     colorBold,
		Dim:      colorDim,
		Reverse:  colorReverse,
		Stats:    colorCyan,
		Kept:     colorGreen,
		Dropped:  colorGray,
		Reroll:   colorYellow,
		Exploded: colorMagenta + colorBold,
		Total:    colorWhite + colorBold,
		Success:  colorGreen,
	}
}

// isTTY checks if stdout is connected to a terminal
func isTTY() bool {
	fi, _ := os.Stdout.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// ColF formats a string with color
func (cs ColorScheme) ColF(color string, format string, args ...interface{}) string {
	if color == "" {
		return fmt.Sprintf(format, args...)
	}
	return fmt.Sprintf("%s%s%s", color, fmt.Sprintf(format, args...), cs.Reset)
}

// Col applies color to a string
func (cs ColorScheme) Col(color string, text string) string {
	if color == "" {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, cs.Reset)
}
