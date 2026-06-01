package dice

import "fmt"

// ModifierKind enumerates all supported modifier types.
type ModifierKind int

const (
	// Reroll family
	ModReroll ModifierKind = iota
	ModRerollOnce
	ModRerollAdd

	// Explode family
	ModExplode
	ModExplodeThreshold
	ModExplodeCompound

	// Keep/Drop family
	ModKeepHigh
	ModKeepLow
	ModDropHigh
	ModDropLow

	// Additive + Success
	ModAdditive
	ModSuccessThreshold
)

// Modifier represents a parsed modifier in the expression.
type Modifier struct {
	Kind      ModifierKind
	Value     int    // for +N, -N, success threshold value
	Threshold int    // for reroll/explode thresholds
	Op        string // for success thresholds (>=, <=, >, <)
	Count     int    // for keep/drop
}

//
// ------------------------------------------------------------
// SUCCESS HELPERS
// ------------------------------------------------------------
//

// countSuccesses counts how many values satisfy a threshold operator.
func countSuccesses(values []int, op string, target int) int {
	count := 0
	for _, v := range values {
		switch op {
		case ">":
			if v > target {
				count++
			}
		case ">=":
			if v >= target {
				count++
			}
		case "<":
			if v < target {
				count++
			}
		case "<=":
			if v <= target {
				count++
			}
		}
	}
	return count
}

//
// ------------------------------------------------------------
// NORMALIZATION + VALIDATION
// (moved here from parser.go — correct architectural home)
// ------------------------------------------------------------
//

// normalizeModifiers enforces last‑wins semantics for keep/drop and
// ensures only one reroll family member and one success threshold are
// applied. ModAdditive entries are silently dropped — arithmetic is
// the AST evaluator's responsibility (see eval_ast.go), so any
// additive that survives normalization here would be a no-op anyway.
func normalizeModifiers(mods []Modifier) []Modifier {
	out := []Modifier{}

	var (
		haveReroll     bool
		haveRerollOnce bool
		haveRerollAdd  bool

		keepDropIndex = -1
		haveSuccess   bool
	)

	for _, m := range mods {
		switch m.Kind {

		// Reroll family: only one allowed
		case ModReroll:
			if !haveReroll && !haveRerollOnce && !haveRerollAdd {
				out = append(out, m)
				haveReroll = true
			}

		case ModRerollOnce:
			if !haveReroll && !haveRerollOnce && !haveRerollAdd {
				out = append(out, m)
				haveRerollOnce = true
			}

		case ModRerollAdd:
			if !haveReroll && !haveRerollOnce && !haveRerollAdd {
				out = append(out, m)
				haveRerollAdd = true
			}

		// Keep/Drop family: last one wins
		case ModKeepHigh, ModKeepLow, ModDropHigh, ModDropLow:
			if keepDropIndex == -1 {
				keepDropIndex = len(out)
				out = append(out, m)
			} else {
				out[keepDropIndex] = m
			}

		// Success threshold: only one allowed
		case ModSuccessThreshold:
			if !haveSuccess {
				out = append(out, m)
				haveSuccess = true
			}

		// Additive: silently dropped (handled by AST evaluator).
		case ModAdditive:

		default:
			out = append(out, m)
		}
	}

	return out
}

// validateExpression ensures the parsed expression is semantically valid.
func validateExpression(expr Expression) error {
	if expr.Count <= 0 {
		return fmt.Errorf("dice count must be positive")
	}
	if expr.Sides <= 0 {
		return fmt.Errorf("dice sides must be positive")
	}

	for _, m := range expr.Modifiers {
		switch m.Kind {

		case ModKeepHigh, ModKeepLow, ModDropHigh, ModDropLow:
			if m.Count <= 0 {
				return fmt.Errorf("keep/drop count must be positive")
			}
			if m.Count > expr.Count {
				return fmt.Errorf("keep/drop count (%d) cannot exceed dice count (%d)", m.Count, expr.Count)
			}

		case ModReroll, ModRerollOnce, ModRerollAdd:
			if m.Threshold <= 0 {
				return fmt.Errorf("reroll threshold must be positive")
			}

		case ModExplodeThreshold:
			if m.Threshold <= 0 {
				return fmt.Errorf("explode threshold must be positive")
			}

		case ModSuccessThreshold:
			if m.Value <= 0 {
				return fmt.Errorf("success threshold must be positive")
			}
		}
	}

	return nil
}
