package cli

import (
	"fmt"

	"github.com/showr/dice-roller/internal/dice"
	"github.com/showr/dice-roller/internal/parse"
	"github.com/showr/dice-roller/internal/presentation"
)

func PrintHelp() {
	fmt.Println("Dice Roller CLI")
	fmt.Println()
	fmt.Println("Version:", dice.Version)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  dice-roller <expression> [--verbose] [--multi N] [--no-color] [--seed N]")
	fmt.Println("  dice-roller \"<expression> rolls=N\" [--verbose] [--no-color] [--seed N]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  dice-roller 4d6k3")
	fmt.Println("  dice-roller \"(2d6 + 1d4) * 2\"")
	fmt.Println("  dice-roller \"5d10ro1>=8 rolls=10\"")
	fmt.Println("  dice-roller 4d6k3 --multi 10 --verbose")
	fmt.Println("  dice-roller 2d20kh1 3d8! 5d10>=8 --no-color")
	fmt.Println("  dice-roller 4d6k3 --seed 42")
	fmt.Println()

	for _, line := range dice.HelpLines() {
		fmt.Println(line)
	}

	fmt.Println()
	fmt.Println("Notes:")
	fmt.Println("  - Arithmetic and grouping are supported, e.g. \"(2d6 + 1d4) * 2\".")
	fmt.Println("  - In a shell, quote expressions that contain spaces or parentheses.")
	fmt.Println("  - Multiple expressions may be passed and are evaluated in order.")
	fmt.Println("  - Verbose mode prints rerolls, explosions, kept/dropped dice, and totals.")
	fmt.Println("  - Colors are auto-disabled if output is piped or redirected.")
	fmt.Println("  - --seed produces deterministic output for the entire invocation.")
}

// RunCLI parses args, constructs an Engine (seeded if --seed was given),
// and evaluates each expression. The engine is built here — not in main —
// so that --seed flows from parsed input into NewEngineWithSeed.
func RunCLI(args []string) {
	// Handle help/version BEFORE parsing.
	for _, a := range args {
		switch a {
		case "--help", "-h":
			PrintHelp()
			return
		case "--version":
			fmt.Println("dice-roller version", dice.Version)
			return
		}
	}

	parsed, err := parse.ParseArgs(args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(parsed.Expressions) == 0 {
		return
	}

	engine := engineFromParsed(parsed)
	colors := presentation.GetColorScheme(parsed.NoColor)
	coloredFormatter := presentation.NewColoredFormatter(colors)

	for _, expr := range parsed.Expressions {
		if parsed.Multi <= 1 {
			r, err := engine.Roll(expr)
			if err != nil {
				fmt.Printf("Error evaluating %q: %v\n", expr, err)
				continue
			}
			if parsed.Verbose {
				fmt.Print(coloredFormatter.FormatVerboseSingle(r))
			} else {
				fmt.Println(coloredFormatter.FormatCompactSingle(r))
			}
			continue
		}

		mr, err := engine.RollMany(expr, parsed.Multi)
		if err != nil {
			fmt.Printf("Error evaluating %q: %v\n", expr, err)
			continue
		}
		if parsed.Verbose {
			fmt.Print(coloredFormatter.FormatVerboseMulti(mr))
		} else {
			fmt.Println(coloredFormatter.FormatCompactMulti(mr))
		}
	}
}

// engineFromParsed honors --seed when present; otherwise uses a
// time-based default.
func engineFromParsed(p parse.ParsedInput) *dice.Engine {
	if p.Seed != nil {
		return dice.NewEngineWithSeed(int64(*p.Seed))
	}
	return dice.NewEngine()
}
