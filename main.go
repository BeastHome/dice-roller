package main

import (
	"fmt"
	"os"

	"github.com/showr/dice-roller/cmd/cli"
	"github.com/showr/dice-roller/internal/dice"
	"github.com/showr/dice-roller/tui"
)

func main() {
	// Handle CLI help / version before doing anything else so they work
	// even if dice-engine construction would fail.
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "--help" || arg == "-h" {
			cli.PrintHelp()
			return
		}
		if arg == "--version" {
			fmt.Println("dice-roller version", dice.Version)
			return
		}
	}

	// No arguments → launch TUI with a time-seeded engine (TUI does
	// not currently honor --seed; it's a CLI-only flag).
	if len(os.Args) == 1 {
		_ = tui.RunTUI(dice.NewEngine())
		return
	}

	// Arguments → CLI mode constructs its own engine from parsed args,
	// so --seed takes effect.
	cli.RunCLI(os.Args[1:])
}
