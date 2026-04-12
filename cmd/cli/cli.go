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
	fmt.Println("  dice-roller <expression> [--verbose] [--multi N] [--no-color]")
	fmt.Println("  dice-roller \"<expression> rolls=N\" [--verbose] [--no-color]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  dice-roller 4d6k3")
	fmt.Println("  dice-roller \"(2d6 + 1d4) * 2\"")
	fmt.Println("  dice-roller \"5d10ro1>=8 rolls=10\"")
	fmt.Println("  dice-roller 4d6k3 --multi 10 --verbose")
	fmt.Println("  dice-roller 2d20kh1 3d8! 5d10>=8 --no-color")
	fmt.Println()

	for _, line := range dice.HelpText {
		fmt.Println(line)
	}

	fmt.Println()
	fmt.Println("Notes:")
	fmt.Println("  - Arithmetic and grouping are supported, e.g. \"(2d6 + 1d4) * 2\".")
	fmt.Println("  - In a shell, quote expressions that contain spaces or parentheses.")
	fmt.Println("  - Multiple expressions may be passed and are evaluated in order.")
	fmt.Println("  - Verbose mode prints rerolls, explosions, kept/dropped dice, and totals.")
	fmt.Println("  - Colors are auto-disabled if output is piped or redirected.")
}

func RunCLI(engine *dice.Engine, args []string) {

	// Handle help/version BEFORE parsing
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

	// Shared parser handles flags, grouping, normalization
	parsed, err := parse.ParseArgs(args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// No expressions? Nothing to do.
	if len(parsed.Expressions) == 0 {
		return
	}

	// Create colored formatter with appropriate color scheme
	colors := presentation.GetColorScheme(parsed.NoColor)
	coloredFormatter := presentation.NewColoredFormatter(colors)

	// Evaluate each normalized expression
	for _, expr := range parsed.Expressions {

		result, err := engine.Evaluate(expr, parsed.Multi)
		if err != nil {
			fmt.Printf("Error evaluating %q: %v\n", expr, err)
			return
		}

		switch v := result.(type) {

		case dice.Result:
			if parsed.Verbose {
				fmt.Print(coloredFormatter.FormatVerboseSingle(v))
			} else {
				fmt.Println(coloredFormatter.FormatCompactSingle(v))
			}

		case dice.MultiRollResult:
			if parsed.Verbose {
				fmt.Print(coloredFormatter.FormatVerboseMulti(v))
			} else {
				fmt.Println(coloredFormatter.FormatCompactMulti(v))
			}
		}
	}
}
