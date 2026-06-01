package dice

import "testing"

func TestEngineRoll_EvaluatesArithmeticDiceExpression(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})

	res, err := engine.Roll("1d1+2d1")
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

	res, err := engine.Roll("(d1+2)*3")
	if err != nil {
		t.Fatalf("Roll returned error: %v", err)
	}

	if res.Total != 9 {
		t.Fatalf("expected total 9, got %d", res.Total)
	}
}

func TestEngineRoll_DivisionByZeroReturnsError(t *testing.T) {
	// Slice B will also tighten this: today the AST eval error is masked
	// by a fallback to the legacy parser, which then fails with a parse
	// error. Either way the user sees an error — assert that without
	// pinning the exact message.
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})
	if _, err := engine.Roll("1d1 / (1d1 - 1)"); err == nil {
		t.Fatalf("expected error from division by zero, got nil")
	}
}

func TestEngineRoll_UnaryMinusNegatesTotal(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})
	res, err := engine.Roll("-1d1")
	if err != nil {
		t.Fatalf("Roll returned error: %v", err)
	}
	if res.Total != -1 {
		t.Fatalf("expected total -1, got %d", res.Total)
	}
}

func TestEngineRoll_GroupedLiteralEvaluatesInner(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})
	res, err := engine.Roll("(5)")
	if err != nil {
		t.Fatalf("Roll returned error: %v", err)
	}
	if res.Total != 5 {
		t.Fatalf("expected total 5, got %d", res.Total)
	}
}
