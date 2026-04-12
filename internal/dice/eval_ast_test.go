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
