package dice

import (
	"fmt"
	"strings"
)

// AttachVerbose builds a human-readable breakdown of the roll.
// This function is intentionally simple: it formats the Result
// without performing any evaluation logic.
func AttachVerbose(r *Result) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Expression: %s\n", r.Expression))
	sb.WriteString(fmt.Sprintf("Rolls: %v\n", r.Rolls))

	if len(r.Rerolls) > 0 {
		sb.WriteString(fmt.Sprintf("Rerolls: %v\n", r.Rerolls))
	}

	if len(r.Exploded) > 0 {
		sb.WriteString(fmt.Sprintf("Exploded: %v\n", r.Exploded))
	}

	sb.WriteString(fmt.Sprintf("Kept: %v\n", r.Kept))

	if len(r.Dropped) > 0 {
		sb.WriteString(fmt.Sprintf("Dropped: %v\n", r.Dropped))
	}

	sb.WriteString(fmt.Sprintf("Total: %d\n", r.Total))

	if r.Successes > 0 {
		sb.WriteString(fmt.Sprintf("Successes: %d\n", r.Successes))
	}

	r.Verbose = sb.String()
}
