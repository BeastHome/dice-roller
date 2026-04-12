package parse

import (
	"strconv"
	"strings"
)

// processExpressions applies grouping, inline rolls, and normalization
// to the raw expression-like tokens collected in parseState.
func processExpressions(state *parseState) ([]string, error) {
	if len(state.rawExprTokens) == 0 {
		return nil, nil
	}

	// First, join tokens with spaces to get a single logical input string.
	joined := strings.Join(state.rawExprTokens, " ")

	// Extract inline rolls=N (applies to whole input unless overridden by --multi).
	joined = extractInlineRolls(state, joined)

	// Now perform hybrid grouping and splitting into individual expressions.
	exprs := splitHybrid(joined)

	// Normalize each expression.
	var out []string
	for _, e := range exprs {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		norm := normalizeExpression(e)
		if norm != "" {
			out = append(out, norm)
		}
	}

	return out, nil
}

// extractInlineRolls finds "rolls=N" patterns and updates state.multi
// if --multi was not already set explicitly.
func extractInlineRolls(state *parseState, input string) string {
	// Simple approach: look for "rolls=" tokens separated by whitespace.
	parts := strings.Fields(input)
	var kept []string

	for _, p := range parts {
		if strings.HasPrefix(p, "rolls=") {
			val := strings.TrimPrefix(p, "rolls=")
			// Only apply if multi is still default (1).
			if state.multi == 1 {
				if n, err := strconv.Atoi(val); err == nil && n > 0 {
					state.multi = n
				}
			}
			continue
		}
		kept = append(kept, p)
	}

	return strings.Join(kept, " ")
}

// splitHybrid implements the hybrid grouping rules:
// - parentheses = explicit grouping
// - semicolon/comma inside parentheses split within the group
// - whitespace outside parentheses separates expressions
// - semicolon/comma outside parentheses also separate expressions
func splitHybrid(input string) []string {
	var result []string
	var current strings.Builder

	depth := 0

	flush := func() {
		s := strings.TrimSpace(current.String())
		if s != "" {
			result = append(result, s)
		}
		current.Reset()
	}

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case '(':
			depth++
			current.WriteByte(ch)

		case ')':
			if depth > 0 {
				depth--
			}
			current.WriteByte(ch)

		case ';', ',':
			if depth == 0 {
				flush()
			} else {
				current.WriteByte(ch)
			}

		case ' ', '\t', '\n', '\r':
			if depth > 0 {
				current.WriteByte(ch)
				continue
			}

			j := i
			for j < len(input) && isSpaceByte(input[j]) {
				j++
			}

			prev, hasPrev := lastNonSpaceByte(current.String())
			next, hasNext := nextNonSpaceByte(input, j)
			if hasPrev && hasNext && shouldPreserveSpace(prev, next) {
				if current.Len() > 0 {
					current.WriteByte(' ')
				}
			} else {
				flush()
			}
			i = j - 1

		default:
			current.WriteByte(ch)
		}
	}

	flush()
	return result
}

func isSpaceByte(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func nextNonSpaceByte(input string, start int) (byte, bool) {
	for i := start; i < len(input); i++ {
		if !isSpaceByte(input[i]) {
			return input[i], true
		}
	}
	return 0, false
}

func lastNonSpaceByte(s string) (byte, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if !isSpaceByte(s[i]) {
			return s[i], true
		}
	}
	return 0, false
}

func shouldPreserveSpace(prev, next byte) bool {
	if isOperatorOrGroupingByte(prev) || isOperatorOrGroupingByte(next) {
		return true
	}

	if isModifierLeadByte(next) {
		return true
	}

	if isNumericModifierLeadByte(prev) && next >= '0' && next <= '9' {
		return true
	}

	return false
}

func isOperatorOrGroupingByte(b byte) bool {
	switch b {
	case '+', '-', '*', '/', '(', ')', '<', '>', '=':
		return true
	default:
		return false
	}
}

func isModifierLeadByte(b byte) bool {
	switch b {
	case '!', 'k', 'K', 'r', 'R', 'l', 'L', 'h', 'H', 'a', 'A', 'o', 'O':
		return true
	default:
		return false
	}
}

func isNumericModifierLeadByte(b byte) bool {
	switch b {
	case 'k', 'K', 'r', 'R', 'l', 'L', 'h', 'H', 'a', 'A', 'o', 'O':
		return true
	default:
		return false
	}
}
