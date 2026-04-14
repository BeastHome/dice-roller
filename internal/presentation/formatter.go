package presentation

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/showr/dice-roller/internal/dice"
)

// Formatter provides methods for formatting dice roll results for display.
type Formatter interface {
	// FormatSingle returns a single-line summary of a single roll result.
	FormatSingleSummary(r dice.Result) string

	// FormatMultiSummary returns a single-line summary of a multi-roll result with stats.
	FormatMultiSummary(mr dice.MultiRollResult) string

	// FormatVerbose returns a detailed breakdown of a single roll.
	FormatVerbose(r dice.Result) string

	// FormatMultiRollLine returns a single-line display of a roll result (for compact output).
	FormatMultiRollLine(r dice.Result, rollNum int) string
}

// DefaultFormatter provides standard text formatting for results.
type DefaultFormatter struct{}

// NewDefaultFormatter creates a new DefaultFormatter.
func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{}
}

// FormatSingleSummary returns a one-line summary: "expr | total=T | rolls=... kept=..."
func (f *DefaultFormatter) FormatSingleSummary(r dice.Result) string {
	return fmt.Sprintf(
		"%s | total=%d | rolls=%v kept=%v",
		r.Expression,
		r.Total,
		r.Rolls,
		r.Kept,
	)
}

// FormatMultiSummary returns a one-line summary with stats: "expr rolls=N | avg=X.XX | min=M | max=X"
func (f *DefaultFormatter) FormatMultiSummary(mr dice.MultiRollResult) string {
	fullExpr := dice.FormatMultiExpression(mr.Expression, len(mr.Rolls))
	if len(mr.Rolls) == 0 {
		return fullExpr
	}

	// Compute stats
	var sum int64
	min := mr.Rolls[0].Total
	max := mr.Rolls[0].Total

	for _, r := range mr.Rolls {
		sum += int64(r.Total)
		if r.Total < min {
			min = r.Total
		}
		if r.Total > max {
			max = r.Total
		}
	}

	avg := float64(sum) / float64(len(mr.Rolls))
	return fmt.Sprintf("%s | avg=%.2f | min=%d | max=%d", fullExpr, avg, min, max)
}

// FormatVerbose returns a detailed breakdown with all roll details.
func (f *DefaultFormatter) FormatVerbose(r dice.Result) string {
	if r.Verbose != "" {
		return r.Verbose
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Expression: %s\n", r.Expression))
	sb.WriteString(fmt.Sprintf("Rolls: %v\n", r.Rolls))

	if len(r.Rerolls) > 0 {
		sb.WriteString(fmt.Sprintf("Rerolls: %v\n", r.Rerolls))
	}

	if len(r.Exploded) > 0 {
		sb.WriteString(fmt.Sprintf("Exploded: %v\n", r.Exploded))
	}

	sb.WriteString(fmt.Sprintf("Kept: %v\n", r.Kept))

	if len(r.Dropped) > 0 {
		sb.WriteString(fmt.Sprintf("Dropped: %v\n", r.Dropped))
	}

	sb.WriteString(fmt.Sprintf("Total: %d\n", r.Total))

	if r.Successes > 0 {
		sb.WriteString(fmt.Sprintf("Successes: %d\n", r.Successes))
	}

	return sb.String()
}

// FormatMultiRollLine returns a compact line for a single roll in a multi-roll context.
func (f *DefaultFormatter) FormatMultiRollLine(r dice.Result, rollNum int) string {
	return fmt.Sprintf("Roll %d: %d", rollNum, r.Total)
}

// SimpleFormat returns just "expression -> total"
func SimpleFormat(r dice.Result) string {
	return fmt.Sprintf("%s -> %d", r.Expression, r.Total)
}

// ColoredFormatter provides color-aware formatting for CLI output
type ColoredFormatter struct {
	colors ColorScheme
}

// NewColoredFormatter creates a formatter with the given color scheme
func NewColoredFormatter(colors ColorScheme) *ColoredFormatter {
	return &ColoredFormatter{colors: colors}
}

// FormatCompactSingle returns: [expr] → total (optionally with effects)
func (f *ColoredFormatter) FormatCompactSingle(r dice.Result) string {
	var effects []string
	if len(r.Rerolls) > 0 {
		effects = append(effects, f.colors.Col(f.colors.Reroll, fmt.Sprintf("rerolled: %v", r.Rerolls)))
	}
	if len(r.Exploded) > 0 {
		effects = append(effects, f.colors.Col(f.colors.Exploded, fmt.Sprintf("exploded: %v", r.Exploded)))
	}

	effectsStr := ""
	if len(effects) > 0 {
		effectsStr = "  (" + strings.Join(effects, ", ") + ")"
	}

	return fmt.Sprintf("[%s] %s %s%s",
		r.Expression,
		f.colors.Col(f.colors.Bold, "→"),
		f.colors.Col(f.colors.Total, fmt.Sprintf("%d", r.Total)),
		effectsStr,
	)
}

// FormatCompactMulti returns a multi-roll summary with stats and totals strip
func (f *ColoredFormatter) FormatCompactMulti(mr dice.MultiRollResult) string {
	if len(mr.Rolls) == 0 {
		return fmt.Sprintf("[%s]", mr.Expression)
	}

	// Compute stats
	stats := computeMultiRollStatsForFormatter(mr.Rolls)

	header := fmt.Sprintf("[%s %s]  %s %.2f  %s %.2f  %s %d  %s %d  %s %.2f",
		mr.Expression,
		f.colors.Col(f.colors.Dim, fmt.Sprintf("rolls=%d", len(mr.Rolls))),
		f.colors.Col(f.colors.Stats, "AVG"), stats.average,
		f.colors.Col(f.colors.Stats, "MED"), stats.median,
		f.colors.Col(f.colors.Stats, "MIN"), stats.min,
		f.colors.Col(f.colors.Stats, "MAX"), stats.max,
		f.colors.Col(f.colors.Stats, "SD"), stats.stdDev,
	)

	// Build totals line with highest highlighted
	maxTotal := mr.Rolls[0].Total
	maxIdx := 0
	for i, r := range mr.Rolls {
		if r.Total > maxTotal {
			maxTotal = r.Total
			maxIdx = i
		}
	}

	var totalsSb strings.Builder
	totalsSb.WriteString(f.colors.Col(f.colors.Stats, "Totals:"))
	totalsSb.WriteString(f.colors.Col(f.colors.Stats, " "))
	for i, r := range mr.Rolls {
		if i > 0 {
			totalsSb.WriteString(f.colors.Col(f.colors.Stats, " "))
		}
		if i == maxIdx {
			totalsSb.WriteString(f.colors.Col(f.colors.Stats, f.colors.Col(f.colors.Reverse, fmt.Sprintf("%d", r.Total))))
		} else {
			totalsSb.WriteString(f.colors.Col(f.colors.Stats, fmt.Sprintf("%d", r.Total)))
		}
	}

	return fmt.Sprintf("%s\n%s", header, totalsSb.String())
}

// FormatVerboseSingle returns a detailed breakdown with semantic colors
func (f *ColoredFormatter) FormatVerboseSingle(r dice.Result) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s]\n", r.Expression))

	if len(r.Rolls) > 0 {
		sb.WriteString(f.formatValueLine("Raw rolls:", r.Rolls, f.colors.Stats))
	}

	if len(r.Rerolls) > 0 {
		sb.WriteString(f.formatValueLine("Rerolls:", r.Rerolls, f.colors.Reroll))
	}

	if len(r.Exploded) > 0 {
		sb.WriteString(f.formatValueLine("Exploded:", r.Exploded, f.colors.Exploded))
	}

	if len(r.Kept) > 0 {
		sb.WriteString(f.formatValueLine("Kept:", r.Kept, f.colors.Kept))
	}

	if len(r.Dropped) > 0 {
		sb.WriteString(f.formatValueLine("Dropped:", r.Dropped, f.colors.Dropped))
	}

	sb.WriteString(fmt.Sprintf("%-11s%s\n",
		"Total:",
		f.colors.Col(f.colors.Total, fmt.Sprintf("%d", r.Total)),
	))

	if r.Successes > 0 {
		sb.WriteString(fmt.Sprintf("%-11s%s\n",
			"Successes:",
			f.colors.Col(f.colors.Success, fmt.Sprintf("%d", r.Successes)),
		))
	}

	return sb.String()
}

// FormatVerboseMulti returns detailed breakdowns for each roll in a multi-roll result
func (f *ColoredFormatter) FormatVerboseMulti(mr dice.MultiRollResult) string {
	if len(mr.Rolls) == 0 {
		return fmt.Sprintf("[%s]", mr.Expression)
	}

	stats := computeMultiRollStatsForFormatter(mr.Rolls)
	header := fmt.Sprintf("[%s %s]  %s %.2f  %s %.2f  %s %d  %s %d  %s %.2f\n",
		mr.Expression,
		f.colors.Col(f.colors.Dim, fmt.Sprintf("rolls=%d", len(mr.Rolls))),
		f.colors.Col(f.colors.Stats, "AVG"), stats.average,
		f.colors.Col(f.colors.Stats, "MED"), stats.median,
		f.colors.Col(f.colors.Stats, "MIN"), stats.min,
		f.colors.Col(f.colors.Stats, "MAX"), stats.max,
		f.colors.Col(f.colors.Stats, "SD"), stats.stdDev,
	)

	// Build totals line
	maxTotal := mr.Rolls[0].Total
	maxIdx := 0
	for i, r := range mr.Rolls {
		if r.Total > maxTotal {
			maxTotal = r.Total
			maxIdx = i
		}
	}
	var totalsSb strings.Builder
	totalsSb.WriteString(f.colors.Col(f.colors.Stats, "Totals:"))
	totalsSb.WriteString(f.colors.Col(f.colors.Stats, " "))
	for i, r := range mr.Rolls {
		if i > 0 {
			totalsSb.WriteString(f.colors.Col(f.colors.Stats, " "))
		}
		if i == maxIdx {
			totalsSb.WriteString(f.colors.Col(f.colors.Stats, f.colors.Col(f.colors.Reverse, fmt.Sprintf("%d", r.Total))))
		} else {
			totalsSb.WriteString(f.colors.Col(f.colors.Stats, fmt.Sprintf("%d", r.Total)))
		}
	}

	var sb strings.Builder
	sb.WriteString(header)
	sb.WriteString(totalsSb.String())
	sb.WriteString("\n\n")

	// Detail for each roll
	for i, r := range mr.Rolls {
		sb.WriteString(fmt.Sprintf("%sRoll %d/%d:%s\n",
			f.colors.Col(f.colors.Bold, ">>> "),
			i+1,
			len(mr.Rolls),
			f.colors.Col(f.colors.Total, fmt.Sprintf(" %d", r.Total)),
		))

		if len(r.Rolls) > 0 {
			sb.WriteString(f.formatValueLine("  Raw rolls:", r.Rolls, f.colors.Stats))
		}

		if len(r.Rerolls) > 0 {
			sb.WriteString(f.formatValueLine("  Rerolls:", r.Rerolls, f.colors.Reroll))
		}

		if len(r.Exploded) > 0 {
			sb.WriteString(f.formatValueLine("  Exploded:", r.Exploded, f.colors.Exploded))
		}

		if len(r.Kept) > 0 {
			sb.WriteString(f.formatValueLine("  Kept:", r.Kept, f.colors.Kept))
		}

		if len(r.Dropped) > 0 {
			sb.WriteString(f.formatValueLine("  Dropped:", r.Dropped, f.colors.Dropped))
		}

		if r.Successes > 0 {
			sb.WriteString(fmt.Sprintf("  %-9s%s\n",
				"Successes:",
				f.colors.Col(f.colors.Success, fmt.Sprintf("%d", r.Successes)),
			))
		}

		if i < len(mr.Rolls)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// formatValueLine formats a label and values with color
func (f *ColoredFormatter) formatValueLine(label string, values []int, color string) string {
	if len(values) == 0 {
		return fmt.Sprintf("%-11s-\n", label)
	}
	valueStr := formatIntSlice(values)
	return fmt.Sprintf("%-11s%s\n", label, f.colors.Col(color, valueStr))
}

// formatIntSlice formats an int slice as space-separated string
func formatIntSlice(values []int) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ")
}

// multiRollStatsForFormatter holds computed stats
type multiRollStatsForFormatter struct {
	average float64
	median  float64
	stdDev  float64
	min     int
	max     int
}

// computeMultiRollStatsForFormatter computes stats for multi-roll display
func computeMultiRollStatsForFormatter(rolls []dice.Result) multiRollStatsForFormatter {
	if len(rolls) == 0 {
		return multiRollStatsForFormatter{}
	}

	totals := make([]int, len(rolls))
	minTotal := rolls[0].Total
	maxTotal := rolls[0].Total
	var sum float64

	for i, roll := range rolls {
		totals[i] = roll.Total
		sum += float64(roll.Total)
		if roll.Total < minTotal {
			minTotal = roll.Total
		}
		if roll.Total > maxTotal {
			maxTotal = roll.Total
		}
	}

	average := sum / float64(len(totals))
	median := computeMedianForFormatter(totals)
	stdDev := computeStdDevForFormatter(totals, average)

	return multiRollStatsForFormatter{
		average: average,
		median:  median,
		stdDev:  stdDev,
		min:     minTotal,
		max:     maxTotal,
	}
}

func computeMedianForFormatter(values []int) float64 {
	copyValues := append([]int(nil), values...)
	sort.Ints(copyValues)
	middle := len(copyValues) / 2
	if len(copyValues)%2 == 1 {
		return float64(copyValues[middle])
	}
	return float64(copyValues[middle-1]+copyValues[middle]) / 2
}

func computeStdDevForFormatter(values []int, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var variance float64
	for _, value := range values {
		delta := float64(value) - mean
		variance += delta * delta
	}
	variance /= float64(len(values))
	return math.Sqrt(variance)
}
