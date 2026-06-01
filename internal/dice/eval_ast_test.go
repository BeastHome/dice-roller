package dice

import (
	"strings"
	"testing"
)

func TestEngineRoll_EvaluatesArithmeticDiceExpression(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})

	res, err := engine.RollOnce("1d1+2d1")
	if err != nil {
		t.Fatalf("Roll returned error: %v", err)
	}

	if res.Total != 3 {
		t.Fatalf("expected total 3, got %d", res.Total)
	}
	if len(res.Rolls) != 3 {
		t.Fatalf("expected 3 underlying rolls, got %d", len(res.Rolls))
	}
}

func TestEngineRoll_EvaluatesGroupedArithmeticExpression(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})

	res, err := engine.RollOnce("(d1+2)*3")
	if err != nil {
		t.Fatalf("Roll returned error: %v", err)
	}

	if res.Total != 9 {
		t.Fatalf("expected total 9, got %d", res.Total)
	}
}

func TestEngineRoll_DivisionByZeroReturnsEvalError(t *testing.T) {
	// Slice B tightened error propagation: when the AST path parses
	// successfully but evaluation fails, the eval error is returned
	// directly rather than masked by a fallback to the legacy parser.
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})
	_, err := engine.RollOnce("1d1 / (1d1 - 1)")
	if err == nil {
		t.Fatalf("expected error from division by zero, got nil")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Fatalf("expected error mentioning 'division by zero', got: %v", err)
	}
}

func TestEngineRoll_UnaryMinusNegatesTotal(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})
	res, err := engine.RollOnce("-1d1")
	if err != nil {
		t.Fatalf("Roll returned error: %v", err)
	}
	if res.Total != -1 {
		t.Fatalf("expected total -1, got %d", res.Total)
	}
}

func TestEngineRoll_GroupedLiteralEvaluatesInner(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})
	res, err := engine.RollOnce("(5)")
	if err != nil {
		t.Fatalf("Roll returned error: %v", err)
	}
	if res.Total != 5 {
		t.Fatalf("expected total 5, got %d", res.Total)
	}
}
