package parse

import (
	"strings"
	"unicode"
)

// Hybrid normalization:
// - canonicalize syntax (case, spacing inside expressions, operators)
// - preserve harmless formatting (spacing around groups, etc.) where reasonable.
func normalizeExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return ""
	}

	// Lowercase everything for consistency.
	expr = strings.ToLower(expr)

	// Normalize spaces inside parentheses and around delimiters.
	expr = normalizeSpaces(expr)

	// Normalize comparison operators.
	expr = strings.ReplaceAll(expr, "=>", ">=")
	expr = strings.ReplaceAll(expr, "=<", "<=")

	// Normalize grouping delimiters inside parentheses: comma -> semicolon.
	expr = normalizeGroupDelimiters(expr)

	// Remove redundant outer parentheses.
	expr = trimRedundantParens(expr)

	return expr
}

func normalizeSpaces(s string) string {
	var b strings.Builder
	var prev rune
	for _, r := range s {
		if unicode.IsSpace(r) {
			// collapse multiple spaces to a single space
			if prev != ' ' {
				b.WriteRune(' ')
				prev = ' '
			}
			continue
		}
		b.WriteRune(r)
		prev = r
	}
	return strings.TrimSpace(b.String())
}

func normalizeGroupDelimiters(s string) string {
	var b strings.Builder
	depth := 0

	for _, r := range s {
		switch r {
		case '(':
			depth++
			b.WriteRune(r)
		case ')':
			if depth > 0 {
				depth--
			}
			b.WriteRune(r)
		case ',':
			if depth > 0 {
				// inside group, normalize to semicolon
				b.WriteRune(';')
			} else {
				// outside group, leave as-is (already split by splitHybrid)
				b.WriteRune(r)
			}
		default:
			b.WriteRune(r)
		}
	}

	return b.String()
}

func trimRedundantParens(s string) string {
	for {
		s = strings.TrimSpace(s)
		if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
			return s
		}

		// Check if outer parens actually wrap the whole expression.
		depth := 0
		valid := true
		for i, ch := range s {
			if ch == '(' {
				depth++
			} else if ch == ')' {
				depth--
				if depth == 0 && i != len(s)-1 {
					valid = false
					break
				}
			}
		}
		if !valid || depth != 0 {
			return s
		}

		// Strip one layer and try again.
		s = s[1 : len(s)-1]
	}
}
