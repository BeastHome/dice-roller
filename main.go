package main

import (
	"os"

	"github.com/showr/dice-roller/cmd/cli"
	"github.com/showr/dice-roller/internal/dice"
	"github.com/showr/dice-roller/tui"
)

func main() {
	// No arguments → launch TUI with a time-seeded engine (TUI does
	// not currently honor --seed; it's a CLI-only flag).
	if len(os.Args) == 1 {
		_ = tui.RunTUI(dice.NewEngine())
		return
	}

	// Arguments → CLI mode. RunCLI handles --help / --version,
	// parses args, builds an engine (honoring --seed), and returns
	// an exit code we propagate so scripts can branch on it.
	os.Exit(cli.RunCLI(os.Args[1:], os.Stdout, os.Stderr))
}
