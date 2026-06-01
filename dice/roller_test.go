package dice

import (
	"math/rand"
	"strings"
	"testing"
)

// newSeededRNG returns a deterministic RNG for tests.
func newSeededRNG(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func sumInts(v []int) int {
	total := 0
	for _, x := range v {
		total += x
	}
	return total
}

// mustEvaluate wraps EvaluateSingle for tests that expect success.
func mustEvaluate(t *testing.T, rng *rand.Rand, expr Expression) Result {
	t.Helper()
	res, err := EvaluateSingle(rng, expr)
	if err != nil {
		t.Fatalf("EvaluateSingle returned error: %v", err)
	}
	return res
}

func TestEvaluateSingle_BasicNdXKeepsAllRollsInRange(t *testing.T) {
	rng := newSeededRNG(1)
	expr := Expression{Raw: "4d6", Count: 4, Sides: 6}
	res := mustEvaluate(t, rng, expr)

	if len(res.Rolls) != 4 {
		t.Fatalf("expected 4 rolls, got %d", len(res.Rolls))
	}
	for _, v := range res.Rolls {
		if v < 1 || v > 6 {
			t.Fatalf("roll out of range [1,6]: %d", v)
		}
	}
	if len(res.Kept) != 4 {
		t.Fatalf("expected all kept when no keep/drop modifier, got %d kept", len(res.Kept))
	}
	if res.Total != sumInts(res.Kept) {
		t.Fatalf("total %d != sum of kept %v", res.Total, res.Kept)
	}
}

func TestEvaluateSingle_KeepHighest(t *testing.T) {
	rng := newSeededRNG(2)
	expr := Expression{
		Raw: "4d6k3", Count: 4, Sides: 6,
		Modifiers: []Modifier{{Kind: ModKeepHigh, Count: 3}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Kept) != 3 || len(res.Dropped) != 1 {
		t.Fatalf("expected 3 kept and 1 dropped, got %d/%d", len(res.Kept), len(res.Dropped))
	}
	for _, k := range res.Kept {
		for _, d := range res.Dropped {
			if k < d {
				t.Fatalf("kept %d less than dropped %d (kept=%v dropped=%v)", k, d, res.Kept, res.Dropped)
			}
		}
	}
	if res.Total != sumInts(res.Kept) {
		t.Fatalf("total %d != sum of kept %v", res.Total, res.Kept)
	}
}

func TestEvaluateSingle_KeepLowest(t *testing.T) {
	rng := newSeededRNG(3)
	expr := Expression{
		Raw: "4d6kl2", Count: 4, Sides: 6,
		Modifiers: []Modifier{{Kind: ModKeepLow, Count: 2}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Kept) != 2 || len(res.Dropped) != 2 {
		t.Fatalf("expected 2 kept and 2 dropped, got %d/%d", len(res.Kept), len(res.Dropped))
	}
	for _, k := range res.Kept {
		for _, d := range res.Dropped {
			if k > d {
				t.Fatalf("kept %d greater than dropped %d", k, d)
			}
		}
	}
}

func TestEvaluateSingle_DropHighest(t *testing.T) {
	rng := newSeededRNG(4)
	expr := Expression{
		Raw: "4d6dh1", Count: 4, Sides: 6,
		Modifiers: []Modifier{{Kind: ModDropHigh, Count: 1}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Kept) != 3 || len(res.Dropped) != 1 {
		t.Fatalf("expected 3 kept and 1 dropped, got %d/%d", len(res.Kept), len(res.Dropped))
	}
	for _, k := range res.Kept {
		if k > res.Dropped[0] {
			t.Fatalf("kept %d greater than dropped %d", k, res.Dropped[0])
		}
	}
}

func TestEvaluateSingle_DropLowest(t *testing.T) {
	rng := newSeededRNG(5)
	expr := Expression{
		Raw: "4d6dl1", Count: 4, Sides: 6,
		Modifiers: []Modifier{{Kind: ModDropLow, Count: 1}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Kept) != 3 || len(res.Dropped) != 1 {
		t.Fatalf("expected 3 kept and 1 dropped, got %d/%d", len(res.Kept), len(res.Dropped))
	}
	for _, k := range res.Kept {
		if k < res.Dropped[0] {
			t.Fatalf("kept %d less than dropped %d", k, res.Dropped[0])
		}
	}
}

func TestEvaluateSingle_RerollReplace_TracksOldValue(t *testing.T) {
	// threshold=6 on d6 → every die rerolls exactly once. doRerollReplace
	// records the OLD (discarded) value in res.Rerolls.
	rng := newSeededRNG(6)
	expr := Expression{
		Raw: "5d6r6", Count: 5, Sides: 6,
		Modifiers: []Modifier{{Kind: ModReroll, Threshold: 6}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Rerolls) != 5 {
		t.Fatalf("expected 5 reroll entries (all triggered), got %d", len(res.Rerolls))
	}
	if len(res.Rolls) != 5 {
		t.Fatalf("expected 5 final rolls, got %d", len(res.Rolls))
	}
	for _, v := range res.Rerolls {
		if v < 1 || v > 6 {
			t.Fatalf("reroll value out of range: %d", v)
		}
	}
}

func TestEvaluateSingle_RerollOnce_TracksNewValue(t *testing.T) {
	// threshold=6 on d6 → all triggered. doRerollOnce records the NEW
	// (placed) value in res.Rerolls, so Rolls[i] == Rerolls[i].
	rng := newSeededRNG(7)
	expr := Expression{
		Raw: "5d6ro6", Count: 5, Sides: 6,
		Modifiers: []Modifier{{Kind: ModRerollOnce, Threshold: 6}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Rerolls) != 5 {
		t.Fatalf("expected 5 reroll entries, got %d", len(res.Rerolls))
	}
	for i, v := range res.Rolls {
		if v != res.Rerolls[i] {
			t.Fatalf("reroll-once should record NEW value; rolls=%v rerolls=%v", res.Rolls, res.Rerolls)
		}
	}
}

func TestEvaluateSingle_RerollAdd_AppendsExtraDice(t *testing.T) {
	// threshold=6 on d6 → every original die triggers, so doRerollAdd
	// appends one new die per trigger. Final rolls = 3 + 3 = 6.
	rng := newSeededRNG(8)
	expr := Expression{
		Raw: "3d6ra6", Count: 3, Sides: 6,
		Modifiers: []Modifier{{Kind: ModRerollAdd, Threshold: 6}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Rolls) != 6 {
		t.Fatalf("expected 6 final rolls (3 original + 3 added), got %d: %v", len(res.Rolls), res.Rolls)
	}
	if len(res.Rerolls) != 3 {
		t.Fatalf("expected 3 reroll entries, got %d", len(res.Rerolls))
	}
}

func TestEvaluateSingle_ExplodeSimple_SinglePassOnMax(t *testing.T) {
	// d1 always rolls 1 (max). Result field semantics:
	//   Rolls    = post-reroll state (pre-explode) → len 1
	//   Exploded = the explosion values            → len 1
	//   Kept     = post-explode final state        → len 2
	// ModExplode does ONE pass; the new die from the explosion is not
	// itself examined for re-explosion (that's ModExplodeCompound).
	rng := newSeededRNG(9)
	expr := Expression{
		Raw: "1d1!", Count: 1, Sides: 1,
		Modifiers: []Modifier{{Kind: ModExplode}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Rolls) != 1 {
		t.Fatalf("expected Rolls to hold 1 post-reroll value, got %d: %v", len(res.Rolls), res.Rolls)
	}
	if len(res.Exploded) != 1 {
		t.Fatalf("expected 1 explosion entry, got %d", len(res.Exploded))
	}
	if len(res.Kept) != 2 {
		t.Fatalf("expected Kept to hold original + explosion (len 2), got %d: %v", len(res.Kept), res.Kept)
	}
}

func TestEvaluateSingle_ExplodeThreshold_AddsOnePerQualifyingInitialDie(t *testing.T) {
	rng := newSeededRNG(10)
	expr := Expression{
		Raw: "5d6!5", Count: 5, Sides: 6,
		Modifiers: []Modifier{{Kind: ModExplodeThreshold, Threshold: 5}},
	}
	res := mustEvaluate(t, rng, expr)

	qualifying := 0
	for _, v := range res.Rolls {
		if v >= 5 {
			qualifying++
		}
	}
	if len(res.Rolls) != 5 {
		t.Fatalf("expected Rolls to hold 5 post-reroll values (pre-explode), got %d: %v", len(res.Rolls), res.Rolls)
	}
	if len(res.Kept) != 5+qualifying {
		t.Fatalf("expected Kept len = 5 + %d explosions, got %d (rolls=%v kept=%v)",
			qualifying, len(res.Kept), res.Rolls, res.Kept)
	}
	if len(res.Exploded) != qualifying {
		t.Fatalf("expected %d explosion entries, got %d", qualifying, len(res.Exploded))
	}
}

func TestEvaluateSingle_ExplodeCompound_ExercisesBfsPathSafely(t *testing.T) {
	// Use d100 so consecutive max-rolls are astronomically unlikely with
	// a stock RNG. This characterizes that the BFS code path runs and
	// terminates under normal conditions. (A MaxDepth ceiling for
	// adversarial RNG is added in Slice B.)
	rng := newSeededRNG(1)
	expr := Expression{
		Raw: "5d100!!", Count: 5, Sides: 100,
		Modifiers: []Modifier{{Kind: ModExplodeCompound}},
	}
	res := mustEvaluate(t, rng, expr)

	if len(res.Rolls) < 5 {
		t.Fatalf("expected at least 5 rolls (originals), got %d", len(res.Rolls))
	}
}

func TestEvaluateSingle_SuccessThreshold_CountsKeptValues(t *testing.T) {
	rng := newSeededRNG(11)
	expr := Expression{
		Raw: "10d10>=7", Count: 10, Sides: 10,
		Modifiers: []Modifier{{Kind: ModSuccessThreshold, Op: ">=", Value: 7}},
	}
	res := mustEvaluate(t, rng, expr)

	expected := 0
	for _, v := range res.Kept {
		if v >= 7 {
			expected++
		}
	}
	if res.Successes != expected {
		t.Fatalf("expected %d successes (count of kept >= 7), got %d (kept=%v)", expected, res.Successes, res.Kept)
	}
}

func TestEvaluateSingle_SuccessThreshold_AllFourOperators(t *testing.T) {
	cases := []struct {
		op        string
		target    int
		predicate func(int) bool
	}{
		{">", 5, func(v int) bool { return v > 5 }},
		{">=", 5, func(v int) bool { return v >= 5 }},
		{"<", 3, func(v int) bool { return v < 3 }},
		{"<=", 3, func(v int) bool { return v <= 3 }},
	}
	for _, tc := range cases {
		t.Run(tc.op, func(t *testing.T) {
			rng := newSeededRNG(12)
			expr := Expression{
				Raw: "8d6", Count: 8, Sides: 6,
				Modifiers: []Modifier{{Kind: ModSuccessThreshold, Op: tc.op, Value: tc.target}},
			}
			res := mustEvaluate(t, rng, expr)
			expected := 0
			for _, v := range res.Kept {
				if tc.predicate(v) {
					expected++
				}
			}
			if res.Successes != expected {
				t.Fatalf("op %s target %d: expected %d successes, got %d (kept=%v)",
					tc.op, tc.target, expected, res.Successes, res.Kept)
			}
		})
	}
}

// Slice B fixes — companion tests for the changes that landed in
// the same slice as roller.go's refactor.

func TestEvaluateSingle_AdditiveModifierIsNoOp(t *testing.T) {
	// Slice B removed additive handling from EvaluateSingle — arithmetic
	// is the AST evaluator's job. A direct caller that constructs an
	// Expression with ModAdditive should see the additive ignored
	// (total = sum of kept, no +N applied here).
	rng := newSeededRNG(1)
	expr := Expression{
		Raw: "1d1", Count: 1, Sides: 1,
		Modifiers: []Modifier{{Kind: ModAdditive, Value: 99}},
	}
	res := mustEvaluate(t, rng, expr)
	if res.Total != 1 {
		t.Fatalf("expected total=1 (kept only, additive ignored); got %d", res.Total)
	}
}

func TestEvaluateSingle_CompoundExplodeAdversarialRngHitsDepthCap(t *testing.T) {
	// d1 always rolls 1 (max). Without maxExplodeDepth this would
	// loop forever; with the cap it returns an error after the
	// configured number of explosions.
	rng := newSeededRNG(1)
	expr := Expression{
		Raw: "1d1!!", Count: 1, Sides: 1,
		Modifiers: []Modifier{{Kind: ModExplodeCompound}},
	}
	_, err := EvaluateSingle(rng, expr)
	if err == nil {
		t.Fatalf("expected max-depth error for adversarial d1!!, got nil")
	}
	if !strings.Contains(err.Error(), "exceeded max depth") {
		t.Fatalf("expected max-depth error message, got: %v", err)
	}
}

func TestEngineRoll_CompoundExplodePropagatesDepthError(t *testing.T) {
	// End-to-end: 1d1!! through Engine.Roll should surface the
	// max-depth error rather than hanging or masking it.
	engine := NewEngineWithSeed(1)
	_, err := engine.RollOnce("1d1!!")
	if err == nil {
		t.Fatalf("expected max-depth error through Engine.Roll, got nil")
	}
}
