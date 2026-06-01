package cli

import (
	"fmt"
	"io"

	"github.com/showr/dice-roller/internal/dice"
	"github.com/showr/dice-roller/internal/parse"
	"github.com/showr/dice-roller/internal/presentation"
)

// PrintHelp writes the CLI help text to w.
func PrintHelp(w io.Writer) {
	fmt.Fprintln(w, "Dice Roller CLI")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Version:", dice.Version)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  dice-roller <expression> [--verbose] [--multi N] [--no-color] [--seed N]")
	fmt.Fprintln(w, "  dice-roller \"<expression> rolls=N\" [--verbose] [--no-color] [--seed N]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  dice-roller 4d6k3")
	fmt.Fprintln(w, "  dice-roller \"(2d6 + 1d4) * 2\"")
	fmt.Fprintln(w, "  dice-roller \"5d10ro1>=8 rolls=10\"")
	fmt.Fprintln(w, "  dice-roller 4d6k3 --multi 10 --verbose")
	fmt.Fprintln(w, "  dice-roller 2d20kh1 3d8! 5d10>=8 --no-color")
	fmt.Fprintln(w, "  dice-roller 4d6k3 --seed 42")
	fmt.Fprintln(w)

	for _, line := range dice.HelpLines() {
		fmt.Fprintln(w, line)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Notes:")
	fmt.Fprintln(w, "  - Arithmetic and grouping are supported, e.g. \"(2d6 + 1d4) * 2\".")
	fmt.Fprintln(w, "  - In a shell, quote expressions that contain spaces or parentheses.")
	fmt.Fprintln(w, "  - Multiple expressions may be passed and are evaluated in order.")
	fmt.Fprintln(w, "  - Verbose mode prints rerolls, explosions, kept/dropped dice, and totals.")
	fmt.Fprintln(w, "  - Colors are auto-disabled if output is piped or redirected.")
	fmt.Fprintln(w, "  - --seed produces deterministic output for the entire invocation.")
}

// RunCLI parses args, evaluates each expression, and writes results to
// stdout / errors to stderr. Returns a process exit code: 0 on success,
// 1 if any expression failed or args couldn't be parsed. One bad
// expression doesn't abort the batch — the remaining expressions still
// evaluate, but the exit code reflects the partial failure so callers
// can detect it.
func RunCLI(args []string, stdout, stderr io.Writer) int {
	// Help / version short-circuit before parsing (so they work even
	// if other args are malformed).
	for _, a := range args {
		switch a {
		case "--help", "-h":
			PrintHelp(stdout)
			return 0
		case "--version":
			fmt.Fprintln(stdout, "dice-roller version", dice.Version)
			return 0
		}
	}

	parsed, err := parse.ParseArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}

	if len(parsed.Expressions) == 0 {
		return 0
	}

	engine := engineFromParsed(parsed)
	colors := presentation.GetColorScheme(parsed.NoColor)
	coloredFormatter := presentation.NewColoredFormatter(colors)

	exitCode := 0
	for _, expr := range parsed.Expressions {
		if parsed.Multi <= 1 {
			r, err := engine.RollOnce(expr)
			if err != nil {
				fmt.Fprintf(stderr, "Error evaluating %q: %v\n", expr, err)
				exitCode = 1
				continue
			}
			if parsed.Verbose {
				fmt.Fprint(stdout, coloredFormatter.FormatVerboseSingle(r))
			} else {
				fmt.Fprintln(stdout, coloredFormatter.FormatCompactSingle(r))
			}
			continue
		}

		mr, err := engine.RollN(expr, parsed.Multi)
		if err != nil {
			fmt.Fprintf(stderr, "Error evaluating %q: %v\n", expr, err)
			exitCode = 1
			continue
		}
		if parsed.Verbose {
			fmt.Fprint(stdout, coloredFormatter.FormatVerboseMulti(mr))
		} else {
			fmt.Fprintln(stdout, coloredFormatter.FormatCompactMulti(mr))
		}
	}
	return exitCode
}

// engineFromParsed honors --seed when present; otherwise time-based default.
func engineFromParsed(p parse.ParsedInput) *dice.Engine {
	if p.Seed != nil {
		return dice.NewEngineWithSeed(int64(*p.Seed))
	}
	return dice.NewEngine()
}
