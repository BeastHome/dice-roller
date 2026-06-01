package dice

import "fmt"

// HelpText is the unified help content shared by CLI and TUI.
var HelpText = []string{
	"Core Syntax:",
	"  NdX          roll N dice of size X",
	"  dX           roll 1 die of size X",
	"  NdX+Y        add or subtract a constant",
	"  NdX+MdY      combine multiple roll terms",
	"  (expr)*N     group arithmetic with parentheses",
	"  Spaces around + - * / are allowed",
	"",
	"Keep / Drop:",
	"  NdXkY        keep highest Y",
	"  NdXklY       keep lowest Y",
	"  NdXdhY       drop highest Y",
	"  NdXdlY       drop lowest Y",
	"",
	"Exploding Dice:",
	"  NdX!         explode on max value",
	"  NdX!T        explode on >= T",
	"  NdX!>T       also accepted for explode threshold",
	"  NdX!!        compound explode on max value",
	"",
	"Rerolls:",
	"  NdXrT        reroll values <= T (replace)",
	"  NdXroT       reroll once",
	"  NdXraT       reroll and add (accumulate)",
	"",
	"Success Counting:",
	"  NdX>=T       count successes >= T",
	"  NdX<=T       count successes <= T",
	"  NdX>T        count successes > T",
	"  NdX<T        count successes < T",
	"",
	"Multi-Roll:",
	"  rolls=N      repeat the entire expression N times",
	"  --multi N    same as rolls=N",
	"",
	"Flags:",
	"  --verbose    show full roll breakdown",
	"  --multi N    repeat expression N times",
	"  --help       show this help message",
	"  --version    show version information",
}

// HelpLines returns shared help content with runtime-specific storage details.
func HelpLines() []string {
	lines := append([]string{}, HelpText...)
	lines = append(lines,
		"",
		"Storage:",
		fmt.Sprintf("  Session log/history files: %s", HistoryDir()),
	)
	return lines
}
