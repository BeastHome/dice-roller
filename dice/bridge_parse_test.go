package dice

import "testing"

func TestParseExpression_UsesBridgeForModifierCases(t *testing.T) {
	expr, err := ParseExpression("5d10ro1>=8")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}

	if expr.Count != 5 || expr.Sides != 10 {
		t.Fatalf("expected count=5 sides=10, got count=%d sides=%d", expr.Count, expr.Sides)
	}
	if len(expr.Modifiers) != 2 {
		t.Fatalf("expected 2 modifiers, got %d", len(expr.Modifiers))
	}
}
