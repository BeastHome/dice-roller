package dice

import (
	"fmt"
	"regexp"
	"strconv"
)

var reBaseExpression = regexp.MustCompile(`^(\d+)[dD](\d+)`)

// ParseExpression parses a single normalized dice expression.
// All user-facing normalization, grouping, and flag handling
// is done in internal/parse. This function assumes input is clean.
func ParseExpression(input string) (Expression, error) {
	if tree, err := ParseTreeExpression(input); err == nil {
		if expr, bridgeErr := BuildExpressionFromTree(tree); bridgeErr == nil {
			return expr, nil
		}
	}

	expr, rest, err := parseBaseExpression(input)
	if err != nil {
		return Expression{}, err
	}

	for len(rest) > 0 {
		mod, consumed, err := parseNextModifier(rest)
		if err != nil {
			return Expression{}, err
		}

		expr.Modifiers = append(expr.Modifiers, mod)
		rest = rest[consumed:]
	}

	expr.Modifiers = normalizeModifiers(expr.Modifiers)
	if err := validateExpression(expr); err != nil {
		return Expression{}, err
	}

	return expr, nil
}

func parseBaseExpression(raw string) (Expression, string, error) {
	m := reBaseExpression.FindStringSubmatch(raw)
	if m == nil {
		return Expression{}, "", fmt.Errorf("invalid dice expression")
	}

	count, _ := strconv.Atoi(m[1])
	sides, _ := strconv.Atoi(m[2])

	expr := Expression{
		Raw:   raw,
		Count: count,
		Sides: sides,
	}

	return expr, raw[len(m[0]):], nil
}
