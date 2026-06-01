package dice

import (
	"reflect"
	"testing"
)

func TestEngineRoll_DeterministicWithSameSeed(t *testing.T) {
	// Two engines with the same seed must produce byte-identical roll
	// sequences. This is the contract --seed and embedded use depend on.
	a := NewEngineWithOptions(EngineOptions{Seed: 42})
	b := NewEngineWithOptions(EngineOptions{Seed: 42})

	exprs := []string{"4d6", "2d20", "1d100", "5d10", "8d8"}
	var rollsA, rollsB []int
	for _, e := range exprs {
		ra, err := a.Roll(e)
		if err != nil {
			t.Fatalf("engine A roll %q: %v", e, err)
		}
		rb, err := b.Roll(e)
		if err != nil {
			t.Fatalf("engine B roll %q: %v", e, err)
		}
		rollsA = append(rollsA, ra.Rolls...)
		rollsB = append(rollsB, rb.Rolls...)
	}
	if !reflect.DeepEqual(rollsA, rollsB) {
		t.Fatalf("same-seed engines diverged\nA: %v\nB: %v", rollsA, rollsB)
	}
}

func TestEngineRoll_DifferentSeedsProduceDifferentSequences(t *testing.T) {
	// Sanity check that seed actually controls output. With two different
	// seeds and 50 1d100 rolls each, identical sequences would mean the
	// RNG isn't actually being seeded distinctly.
	a := NewEngineWithOptions(EngineOptions{Seed: 1})
	b := NewEngineWithOptions(EngineOptions{Seed: 2})

	const trials = 50
	ra := make([]int, trials)
	rb := make([]int, trials)
	for i := 0; i < trials; i++ {
		resA, _ := a.Roll("1d100")
		resB, _ := b.Roll("1d100")
		ra[i] = resA.Total
		rb[i] = resB.Total
	}
	if reflect.DeepEqual(ra, rb) {
		t.Fatalf("expected divergent sequences for different seeds, got identical: %v", ra)
	}
}
