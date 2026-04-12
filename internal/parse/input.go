package parse

import (
	"strings"
)

// ParseArgs is the entry point for CLI-style argument parsing.
// It handles flags, inline rolls, grouping, and normalization,
// and returns a ParsedInput with normalized expressions.
func ParseArgs(args []string) (ParsedInput, error) {
	state := newParseState()

	for i := 0; i < len(args); i++ {
		raw := strings.TrimSpace(args[i])
		if raw == "" {
			continue
		}

		// Flags first
		consumed, err := parseFlagToken(&state, args, i)
		if err != nil {
			return ParsedInput{}, err
		}
		if consumed > 0 {
			i += consumed - 1
			continue
		}

		// Non-flag token: treat as expression-ish input
		state.rawExprTokens = append(state.rawExprTokens, raw)
	}

	// Now process expressions: grouping + inline rolls + normalization
	exprs, err := processExpressions(&state)
	if err != nil {
		return ParsedInput{}, err
	}

	out := ParsedInput{
		Expressions: exprs,
		Multi:       state.multi,
		Verbose:     state.verbose,
		NoColor:     state.noColor,
		Seed:        state.seed,
	}

	if out.Multi <= 0 {
		out.Multi = 1
	}

	return out, nil
}

// ParseLine is the entry point for TUI-style single-line input.
// It tokenizes the line and then delegates to the same logic as ParseArgs,
// so flags like --multi/--seed work the same in both TUI and CLI modes.
func ParseLine(line string) (ParsedInput, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return ParsedInput{
			Expressions: nil,
			Multi:       1,
			Verbose:     false,
			NoColor:     false,
			Seed:        nil,
		}, nil
	}

	args := strings.Fields(line)
	return ParseArgs(args)
}
