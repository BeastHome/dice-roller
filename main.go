package main

import (
	"fmt"
	"os"

	"github.com/showr/dice-roller/cmd/cli"
	"github.com/showr/dice-roller/internal/dice"
	"github.com/showr/dice-roller/tui"
)

func main() {
	engine := dice.NewEngine()

	// Handle CLI help
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

	// No arguments → launch TUI
	if len(os.Args) == 1 {
		_ = tui.RunTUI(engine)
		return
	}

	// Arguments → run CLI mode
	cli.RunCLI(engine, os.Args[1:])
}
