package tui

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/showr/dice-roller/internal/dice"
)

type outputSpan struct {
	start int
	end   int
	style tcell.Style
}

type outputStyledLine struct {
	text  string
	spans []outputSpan
}

type outputHitZone struct {
	row       int
	startCol  int
	endCol    int
	rollIndex int
}

func buildOutputView(entry interface{}, selectedRoll int) ([]outputStyledLine, []outputHitZone, int) {
	switch v := entry.(type) {
	case dice.Result:
		return buildSingleRollOutput(v), nil, -1
	case dice.MultiRollResult:
		return buildMultiRollOutput(v, selectedRoll)
	default:
		return []outputStyledLine{{text: "(unknown history entry)"}}, nil, -1
	}
}

func buildSingleRollOutput(r dice.Result) []outputStyledLine {
	lines := []outputStyledLine{
		{text: fmt.Sprintf("[%s]  TOTAL %d", r.Expression, r.Total)},
		{text: formatValueLine("Rolls:", r.Rolls)},
		{text: formatValueLine("Kept:", r.Kept)},
	}
	if len(r.Dropped) > 0 {
		lines = append(lines, outputStyledLine{text: formatValueLine("Dropped:", r.Dropped)})
	}
	if len(r.Rerolls) > 0 || len(r.Exploded) > 0 {
		lines = append(lines, outputStyledLine{text: fmt.Sprintf("Effects:  rerolls %s  exploded %s", formatShortValueBlock(r.Rerolls), formatShortValueBlock(r.Exploded))})
	}
	if r.Successes > 0 {
		lines = append(lines, outputStyledLine{text: fmt.Sprintf("Successes: %d", r.Successes)})
	}
	return lines
}

func buildMultiRollOutput(mr dice.MultiRollResult, selectedRoll int) ([]outputStyledLine, []outputHitZone, int) {
	if len(mr.Rolls) == 0 {
		return []outputStyledLine{{text: "(no rolls)"}}, nil, -1
	}

	selectedRoll = clampRollSelection(selectedRoll, mr.Rolls)
	stats := computeMultiRollStats(mr.Rolls)
	headerText := fmt.Sprintf("[%s]  AVG %.2f  MED %.2f  MIN %d  MAX %d  SD %.2f", dice.FormatMultiExpression(mr.Expression, len(mr.Rolls)), stats.average, stats.median, stats.min, stats.max, stats.stdDev)
	header := outputStyledLine{text: headerText, spans: []outputSpan{{start: 0, end: len(headerText), style: styleOutputHeading}}}

	statsLine, statsZones := buildTotalsLine(mr.Rolls, selectedRoll)
	trendText := fmt.Sprintf("Trend: %s", stats.sparkline)
	trend := outputStyledLine{text: trendText, spans: []outputSpan{{start: 7, end: len(trendText), style: styleOutputStats}}}
	freqText := fmt.Sprintf("Freq:  %s", strings.Join(stats.frequency, "  "))
	freq := outputStyledLine{text: freqText, spans: []outputSpan{{start: 7, end: len(freqText), style: styleOutputStats}}}
	detail := buildSelectedRollDetail(mr.Rolls[selectedRoll], selectedRoll, len(mr.Rolls))

	lines := []outputStyledLine{header, statsLine, trend, freq, {text: ""}}
	lines = append(lines, detail...)
	return lines, statsZones, selectedRoll
}

func buildTotalsLine(rolls []dice.Result, selectedRoll int) (outputStyledLine, []outputHitZone) {
	var builder strings.Builder
	builder.WriteString("Totals: ")
	line := outputStyledLine{}
	hitZones := make([]outputHitZone, 0, len(rolls))

	for i, roll := range rolls {
		if i > 0 {
			builder.WriteByte(' ')
		}
		start := builder.Len()
		text := fmt.Sprintf("%d", roll.Total)
		builder.WriteString(text)
		end := builder.Len()
		if i == selectedRoll {
			line.spans = append(line.spans, outputSpan{start: start, end: end, style: styleOutputSelectedRoll})
		} else {
			line.spans = append(line.spans, outputSpan{start: start, end: end, style: tcell.StyleDefault.Foreground(tcell.ColorWhite)})
		}
		hitZones = append(hitZones, outputHitZone{row: 1, startCol: start, endCol: end, rollIndex: i})
	}

	line.text = builder.String()
	return line, hitZones
}

func buildSelectedRollDetail(r dice.Result, index, total int) []outputStyledLine {
	headerText := fmt.Sprintf("Selected roll detail (%d/%d):", index+1, total)
	lines := []outputStyledLine{
		{text: headerText, spans: []outputSpan{{start: 0, end: len(headerText), style: styleOutputHeading}}},
		buildValueLine("Raw rolls:", r.Rolls, tcell.StyleDefault.Foreground(tcell.ColorWhite)),
		buildValueLine("Rerolls:", r.Rerolls, styleOutputReroll),
		buildValueLine("Exploded:", r.Exploded, styleOutputExploded),
		buildValueLine("Kept:", r.Kept, styleOutputKept),
		buildValueLine("Dropped:", r.Dropped, styleOutputDropped),
	}

	totalText := fmt.Sprintf("Total:     %d", r.Total)
	lines = append(lines, outputStyledLine{
		text:  totalText,
		spans: []outputSpan{{start: 11, end: len(totalText), style: styleOutputTotal}},
	})

	if r.Successes > 0 {
		successText := fmt.Sprintf("Successes: %d", r.Successes)
		lines = append(lines, outputStyledLine{
			text:  successText,
			spans: []outputSpan{{start: 11, end: len(successText), style: styleOutputSuccess}},
		})
	}
	return lines
}

func buildValueLine(label string, values []int, style tcell.Style) outputStyledLine {
	text := fmt.Sprintf("%-10s %s", label, formatShortValueBlock(values))
	line := outputStyledLine{text: text}
	if len(values) > 0 {
		// Values start after 10-char label + 1 space
		line.spans = append(line.spans, outputSpan{start: 11, end: len(text), style: style})
	}
	return line
}

func formatValueLine(label string, values []int) string {
	return fmt.Sprintf("%-10s %s", label, formatShortValueBlock(values))
}

func formatShortValueBlock(values []int) string {
	if len(values) == 0 {
		return "-"
	}
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprintf("%d", value)
	}
	return strings.Join(parts, " ")
}

type multiRollStats struct {
	average   float64
	median    float64
	stdDev    float64
	min       int
	max       int
	sparkline string
	frequency []string
}

func computeMultiRollStats(rolls []dice.Result) multiRollStats {
	totals := make([]int, len(rolls))
	freqMap := map[int]int{}
	minTotal := rolls[0].Total
	maxTotal := rolls[0].Total
	var sum float64

	for i, roll := range rolls {
		totals[i] = roll.Total
		freqMap[roll.Total]++
		sum += float64(roll.Total)
		if roll.Total < minTotal {
			minTotal = roll.Total
		}
		if roll.Total > maxTotal {
			maxTotal = roll.Total
		}
	}

	average := sum / float64(len(totals))
	median := computeMedian(totals)
	stdDev := computeStdDev(totals, average)
	frequency := make([]string, 0, len(freqMap))
	keys := make([]int, 0, len(freqMap))
	for total := range freqMap {
		keys = append(keys, total)
	}
	sort.Ints(keys)
	for _, total := range keys {
		frequency = append(frequency, fmt.Sprintf("%d:%d", total, freqMap[total]))
	}

	return multiRollStats{
		average:   average,
		median:    median,
		stdDev:    stdDev,
		min:       minTotal,
		max:       maxTotal,
		sparkline: buildSparkline(totals),
		frequency: frequency,
	}
}

func computeMedian(values []int) float64 {
	copyValues := append([]int(nil), values...)
	sort.Ints(copyValues)
	middle := len(copyValues) / 2
	if len(copyValues)%2 == 1 {
		return float64(copyValues[middle])
	}
	return float64(copyValues[middle-1]+copyValues[middle]) / 2
}

func computeStdDev(values []int, mean float64) float64 {
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

func buildSparkline(values []int) string {
	const bars = "▁▂▃▄▅▆▇█"
	if len(values) == 0 {
		return ""
	}
	minValue := values[0]
	maxValue := values[0]
	for _, value := range values[1:] {
		if value < minValue {
			minValue = value
		}
		if value > maxValue {
			maxValue = value
		}
	}
	if minValue == maxValue {
		return strings.Repeat(string([]rune(bars)[4]), len(values))
	}

	var builder strings.Builder
	runes := []rune(bars)
	for _, value := range values {
		ratio := float64(value-minValue) / float64(maxValue-minValue)
		index := int(math.Round(ratio * float64(len(runes)-1)))
		if index < 0 {
			index = 0
		}
		if index >= len(runes) {
			index = len(runes) - 1
		}
		builder.WriteRune(runes[index])
	}
	return builder.String()
}

func defaultSelectedRollIndex(mr dice.MultiRollResult) int {
	if len(mr.Rolls) == 0 {
		return -1
	}
	best := 0
	for i := 1; i < len(mr.Rolls); i++ {
		if mr.Rolls[i].Total > mr.Rolls[best].Total {
			best = i
		}
	}
	return best
}

func clampRollSelection(selected int, rolls []dice.Result) int {
	if len(rolls) == 0 {
		return -1
	}
	if selected < 0 || selected >= len(rolls) {
		return defaultSelectedRollIndex(dice.MultiRollResult{Rolls: rolls})
	}
	return selected
}
