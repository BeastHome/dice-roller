package dice

import (
	"strings"
	"testing"
)

func TestValidateExpression_RejectsNonPositiveCountAndSides(t *testing.T) {
	cases := []struct {
		name string
		expr Expression
		want string
	}{
		{"zero count", Expression{Count: 0, Sides: 6}, "count must be positive"},
		{"negative count", Expression{Count: -1, Sides: 6}, "count must be positive"},
		{"zero sides", Expression{Count: 1, Sides: 0}, "sides must be positive"},
		{"negative sides", Expression{Count: 1, Sides: -6}, "sides must be positive"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateExpression(tc.expr)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func TestValidateExpression_KeepDropCountBounded(t *testing.T) {
	t.Run("zero count", func(t *testing.T) {
		err := validateExpression(Expression{
			Count: 4, Sides: 6,
			Modifiers: []Modifier{{Kind: ModKeepHigh, Count: 0}},
		})
		if err == nil || !strings.Contains(err.Error(), "keep/drop count must be positive") {
			t.Fatalf("expected positive-count error, got %v", err)
		}
	})
	t.Run("exceeds dice count", func(t *testing.T) {
		err := validateExpression(Expression{
			Count: 3, Sides: 6,
			Modifiers: []Modifier{{Kind: ModKeepHigh, Count: 4}},
		})
		if err == nil || !strings.Contains(err.Error(), "cannot exceed dice count") {
			t.Fatalf("expected exceeds-dice-count error, got %v", err)
		}
	})
}

func TestValidateExpression_RerollAndExplodeAndSuccessThresholdsBounded(t *testing.T) {
	cases := []struct {
		name string
		mod  Modifier
		want string
	}{
		{"reroll zero", Modifier{Kind: ModReroll, Threshold: 0}, "reroll threshold"},
		{"reroll-once negative", Modifier{Kind: ModRerollOnce, Threshold: -1}, "reroll threshold"},
		{"reroll-add zero", Modifier{Kind: ModRerollAdd, Threshold: 0}, "reroll threshold"},
		{"explode threshold zero", Modifier{Kind: ModExplodeThreshold, Threshold: 0}, "explode threshold"},
		{"success threshold zero", Modifier{Kind: ModSuccessThreshold, Op: ">=", Value: 0}, "success threshold"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateExpression(Expression{
				Count: 4, Sides: 6,
				Modifiers: []Modifier{tc.mod},
			})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func TestValidateExpression_AcceptsValidCombination(t *testing.T) {
	err := validateExpression(Expression{
		Count: 4, Sides: 6,
		Modifiers: []Modifier{
			{Kind: ModKeepHigh, Count: 3},
			{Kind: ModReroll, Threshold: 1},
			{Kind: ModSuccessThreshold, Op: ">=", Value: 4},
		},
	})
	if err != nil {
		t.Fatalf("expected no error for valid expression, got %v", err)
	}
}
