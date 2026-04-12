package dice

import "testing"

func TestBuildExpressionFromTree_DiceWithModifiers(t *testing.T) {
	tree, err := ParseTreeExpression("4d6k3")
	if err != nil {
		t.Fatalf("ParseTreeExpression returned error: %v", err)
	}

	expr, err := BuildExpressionFromTree(tree)
	if err != nil {
		t.Fatalf("BuildExpressionFromTree returned error: %v", err)
	}

	if expr.Raw != "4d6k3" {
		t.Fatalf("expected raw expression %q, got %q", "4d6k3", expr.Raw)
	}
	if expr.Count != 4 || expr.Sides != 6 {
		t.Fatalf("expected count=4 and sides=6, got count=%d sides=%d", expr.Count, expr.Sides)
	}
	if len(expr.Modifiers) != 1 {
		t.Fatalf("expected 1 modifier, got %d", len(expr.Modifiers))
	}
	if expr.Modifiers[0].Kind != ModKeepHigh || expr.Modifiers[0].Count != 3 {
		t.Fatalf("expected keep-high 3, got %#v", expr.Modifiers[0])
	}
}

func TestBuildExpressionFromTree_RejectsArithmeticForLegacyBridge(t *testing.T) {
	tree, err := ParseTreeExpression("2d6 + 3")
	if err != nil {
		t.Fatalf("ParseTreeExpression returned error: %v", err)
	}

	if _, err := BuildExpressionFromTree(tree); err == nil {
		t.Fatalf("expected bridge to reject arithmetic tree")
	}
}
