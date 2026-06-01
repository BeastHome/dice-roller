package dice

import "testing"

func TestEvaluateMulti_RunsExpressionNTimes(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})
	mr, err := EvaluateMulti(engine, "1d6", 5)
	if err != nil {
		t.Fatalf("EvaluateMulti returned error: %v", err)
	}
	if len(mr.Rolls) != 5 {
		t.Fatalf("expected 5 rolls, got %d", len(mr.Rolls))
	}
	if mr.Expression != "1d6" {
		t.Fatalf("expected expression %q, got %q", "1d6", mr.Expression)
	}
}

func TestEvaluateMulti_PropagatesParseError(t *testing.T) {
	engine := NewEngineWithOptions(EngineOptions{Seed: 1})
	if _, err := EvaluateMulti(engine, "not-a-dice-expression", 3); err == nil {
		t.Fatalf("expected parse error to propagate")
	}
}

func TestFormatMultiExpression_EmptyInputProducesBareSuffix(t *testing.T) {
	got := FormatMultiExpression("", 3)
	if got != "rolls=3" {
		t.Fatalf("expected %q, got %q", "rolls=3", got)
	}
}

func TestFormatMultiExpression_AppendsRollsSuffix(t *testing.T) {
	got := FormatMultiExpression("4d6k3", 10)
	if got != "4d6k3 rolls=10" {
		t.Fatalf("expected %q, got %q", "4d6k3 rolls=10", got)
	}
}

func TestFormatMultiExpression_PreservesExistingRollsSuffix(t *testing.T) {
	// If the expression already mentions rolls=, don't double-append.
	got := FormatMultiExpression("4d6 rolls=10", 5)
	if got != "4d6 rolls=10" {
		t.Fatalf("expected %q (no duplicate suffix), got %q", "4d6 rolls=10", got)
	}
}

func TestFormatMultiExpression_TrimsWhitespace(t *testing.T) {
	got := FormatMultiExpression("  4d6k3  ", 3)
	if got != "4d6k3 rolls=3" {
		t.Fatalf("expected %q, got %q", "4d6k3 rolls=3", got)
	}
}
