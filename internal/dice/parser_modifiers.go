package dice

import (
	"fmt"
	"strconv"
	"strings"
)

func parseNextModifier(rest string) (Modifier, int, error) {
	if mod, consumed, ok, err := parseRerollModifier(rest); ok || err != nil {
		return mod, consumed, err
	}

	if mod, consumed, ok, err := parseExplodeModifier(rest); ok || err != nil {
		return mod, consumed, err
	}

	if mod, consumed, ok, err := parseKeepDropModifier(rest); ok || err != nil {
		return mod, consumed, err
	}

	if rest[0] == '+' || rest[0] == '-' {
		n, consumed, err := parseSignedNumber(rest)
		if err != nil {
			return Modifier{}, 0, err
		}
		return Modifier{Kind: ModAdditive, Value: n}, consumed, nil
	}

	if hasComparisonPrefix(rest) {
		op, n, consumed, err := parseComparison(rest)
		if err != nil {
			return Modifier{}, 0, err
		}
		return Modifier{Kind: ModSuccessThreshold, Op: op, Value: n}, consumed, nil
	}

	return Modifier{}, 0, fmt.Errorf("unexpected modifier sequence: %q", rest)
}

func parseRerollModifier(rest string) (Modifier, int, bool, error) {
	switch {
	case strings.HasPrefix(rest, "ro"):
		mod, consumed, err := parseThresholdModifier(rest, "ro", ModRerollOnce)
		return mod, consumed, true, err
	case strings.HasPrefix(rest, "ra"):
		mod, consumed, err := parseThresholdModifier(rest, "ra", ModRerollAdd)
		return mod, consumed, true, err
	case strings.HasPrefix(rest, "r"):
		mod, consumed, err := parseThresholdModifier(rest, "r", ModReroll)
		return mod, consumed, true, err
	default:
		return Modifier{}, 0, false, nil
	}
}

func parseExplodeModifier(rest string) (Modifier, int, bool, error) {
	switch {
	case strings.HasPrefix(rest, "!!"):
		return Modifier{Kind: ModExplodeCompound}, 2, true, nil
	case strings.HasPrefix(rest, "!>"):
		mod, consumed, err := parseThresholdModifier(rest, "!>", ModExplodeThreshold)
		return mod, consumed, true, err
	case strings.HasPrefix(rest, "!") && len(rest) > 1 && isDigit(rest[1]):
		mod, consumed, err := parseThresholdModifier(rest, "!", ModExplodeThreshold)
		return mod, consumed, true, err
	case strings.HasPrefix(rest, "!"):
		return Modifier{Kind: ModExplode}, 1, true, nil
	default:
		return Modifier{}, 0, false, nil
	}
}

func parseKeepDropModifier(rest string) (Modifier, int, bool, error) {
	switch {
	case strings.HasPrefix(rest, "kl"):
		mod, consumed, err := parseCountModifier(rest, "kl", ModKeepLow)
		return mod, consumed, true, err
	case strings.HasPrefix(rest, "k"):
		mod, consumed, err := parseCountModifier(rest, "k", ModKeepHigh)
		return mod, consumed, true, err
	case strings.HasPrefix(rest, "dh"):
		mod, consumed, err := parseCountModifier(rest, "dh", ModDropHigh)
		return mod, consumed, true, err
	case strings.HasPrefix(rest, "dl"):
		mod, consumed, err := parseCountModifier(rest, "dl", ModDropLow)
		return mod, consumed, true, err
	default:
		return Modifier{}, 0, false, nil
	}
}

func parseThresholdModifier(rest, prefix string, kind ModifierKind) (Modifier, int, error) {
	n, consumed, err := parseNumberAfter(rest, prefix)
	if err != nil {
		return Modifier{}, 0, err
	}
	return Modifier{Kind: kind, Threshold: n}, consumed, nil
}

func parseCountModifier(rest, prefix string, kind ModifierKind) (Modifier, int, error) {
	n, consumed, err := parseNumberAfter(rest, prefix)
	if err != nil {
		return Modifier{}, 0, err
	}
	return Modifier{Kind: kind, Count: n}, consumed, nil
}

func hasComparisonPrefix(s string) bool {
	return strings.HasPrefix(s, ">=") ||
		strings.HasPrefix(s, "<=") ||
		strings.HasPrefix(s, ">") ||
		strings.HasPrefix(s, "<")
}

func parseNumberAfter(s, prefix string) (int, int, error) {
	rest := s[len(prefix):]
	i := 0
	for i < len(rest) && isDigit(rest[i]) {
		i++
	}
	if i == 0 {
		return 0, 0, fmt.Errorf("expected number after %q", prefix)
	}
	n, _ := strconv.Atoi(rest[:i])
	return n, len(prefix) + i, nil
}

func parseSignedNumber(s string) (int, int, error) {
	sign := 1
	i := 0
	if s[0] == '-' {
		sign = -1
		i = 1
	} else if s[0] == '+' {
		i = 1
	}
	start := i
	for i < len(s) && isDigit(s[i]) {
		i++
	}
	if i == start {
		return 0, 0, fmt.Errorf("expected signed number")
	}
	n, _ := strconv.Atoi(s[start:i])
	return sign * n, i, nil
}

func parseComparison(s string) (string, int, int, error) {
	var op string
	consumed := 0

	switch {
	case strings.HasPrefix(s, ">="):
		op = ">="
		consumed = 2
	case strings.HasPrefix(s, "<="):
		op = "<="
		consumed = 2
	case strings.HasPrefix(s, ">"):
		op = ">"
		consumed = 1
	case strings.HasPrefix(s, "<"):
		op = "<"
		consumed = 1
	default:
		return "", 0, 0, fmt.Errorf("invalid comparison operator")
	}

	rest := s[consumed:]
	i := 0
	for i < len(rest) && isDigit(rest[i]) {
		i++
	}
	if i == 0 {
		return "", 0, 0, fmt.Errorf("expected number after comparison operator")
	}

	n, _ := strconv.Atoi(rest[:i])
	return op, n, consumed + i, nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
