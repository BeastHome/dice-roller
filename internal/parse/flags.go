package parse

import (
	"fmt"
	"strconv"
	"strings"
)

type parseState struct {
	// flags
	multi   int
	verbose bool
	noColor bool
	seed    *int

	// raw expression-like tokens (before grouping/normalization)
	rawExprTokens []string
}

func newParseState() parseState {
	return parseState{
		multi:         1,
		verbose:       false,
		noColor:       false,
		seed:          nil,
		rawExprTokens: []string{},
	}
}

// parseFlagToken inspects args[i] and, if it is a flag, updates state.
// It returns how many tokens were consumed (1 or 2) or 0 if not a flag.
func parseFlagToken(state *parseState, args []string, i int) (int, error) {
	raw := strings.TrimSpace(args[i])

	// Help/version are intentionally NOT handled here; CLI/TUI can
	// short-circuit before calling ParseArgs/ParseLine if desired.

	// --verbose
	if raw == "--verbose" {
		state.verbose = true
		return 1, nil
	}

	// --no-color
	if raw == "--no-color" {
		state.noColor = true
		return 1, nil
	}

	// --multi or --multi=N
	if strings.HasPrefix(raw, "--multi") {
		// --multi=N
		if strings.HasPrefix(raw, "--multi=") {
			val := strings.TrimPrefix(raw, "--multi=")
			return 1, setMulti(state, val)
		}

		// --multi N
		if i+1 >= len(args) {
			return 0, newParseError(raw, "missing value for --multi", "")
		}
		val := strings.TrimSpace(args[i+1])
		return 2, setMulti(state, val)
	}

	// --seed or --seed=N
	if strings.HasPrefix(raw, "--seed") {
		// --seed=N
		if strings.HasPrefix(raw, "--seed=") {
			val := strings.TrimPrefix(raw, "--seed=")
			return 1, setSeed(state, raw, val)
		}

		// --seed N
		if i+1 >= len(args) {
			return 0, newParseError(raw, "missing value for --seed", "")
		}
		val := strings.TrimSpace(args[i+1])
		return 2, setSeed(state, raw, val)
	}

	// Not a flag
	return 0, nil
}

func setMulti(state *parseState, val string) error {
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return newParseError("multi", fmt.Sprintf("invalid multi value %q", val), "multi must be a positive integer")
	}
	state.multi = n
	return nil
}

func setSeed(state *parseState, raw, val string) error {
	n, err := strconv.Atoi(val)
	if err != nil {
		return newParseError(raw, fmt.Sprintf("invalid seed value %q", val), "seed must be an integer")
	}
	state.seed = &n
	return nil
}
