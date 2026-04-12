package dice

import (
	"math/rand"
	"sort"
)

// EvaluateSingle runs the full evaluation pipeline on a parsed Expression.
// Multi-roll orchestration lives in multiroll.go.
func EvaluateSingle(rng *rand.Rand, expr Expression) Result {
	// ------------------------------------------------------------
	// 1. INITIAL ROLLS
	// ------------------------------------------------------------
	rolls := make([]int, expr.Count)
	for i := 0; i < expr.Count; i++ {
		rolls[i] = rng.Intn(expr.Sides) + 1
	}

	var (
		rerolls  []int
		exploded []int
	)

	// ------------------------------------------------------------
	// 2. REROLL PASS
	// ------------------------------------------------------------
	rerollReplaceApplied := false
	rerollOnceApplied := false
	rerollAddApplied := false

	for _, mod := range expr.Modifiers {
		switch mod.Kind {

		case ModReroll:
			if !rerollReplaceApplied {
				rolls, rerolls = doRerollReplace(rng, rolls, expr.Sides, mod.Threshold, rerolls)
				rerollReplaceApplied = true
			}

		case ModRerollOnce:
			if !rerollOnceApplied {
				rolls, rerolls = doRerollOnce(rng, rolls, expr.Sides, mod.Threshold, rerolls)
				rerollOnceApplied = true
			}

		case ModRerollAdd:
			if !rerollAddApplied {
				rolls, rerolls = doRerollAdd(rng, rolls, expr.Sides, mod.Threshold, rerolls)
				rerollAddApplied = true
			}
		}
	}

	// Freeze post-reroll state for result reporting
	postReroll := append([]int(nil), rolls...)

	// ------------------------------------------------------------
	// 3. EXPLODING PASS
	// ------------------------------------------------------------
	for _, mod := range expr.Modifiers {
		switch mod.Kind {
		case ModExplode:
			rolls, exploded = doExplodeSimple(rng, rolls, expr.Sides, exploded)
		case ModExplodeThreshold:
			rolls, exploded = doExplodeThreshold(rng, rolls, expr.Sides, mod.Threshold, exploded)
		case ModExplodeCompound:
			rolls, exploded = doExplodeCompound(rng, rolls, expr.Sides, exploded)
		}
	}

	// ------------------------------------------------------------
	// 4. KEEP/DROP PASS
	// ------------------------------------------------------------
	kept := append([]int(nil), rolls...)
	var dropped []int

	for _, mod := range expr.Modifiers {
		switch mod.Kind {
		case ModKeepHigh:
			kept, dropped = keepHighest(kept, mod.Count)
		case ModKeepLow:
			kept, dropped = keepLowest(kept, mod.Count)
		case ModDropHigh:
			kept, dropped = dropHighest(kept, mod.Count)
		case ModDropLow:
			kept, dropped = dropLowest(kept, mod.Count)
		}
	}

	// ------------------------------------------------------------
	// 5. ADDITIVE PASS
	// ------------------------------------------------------------
	total := 0
	for _, v := range kept {
		total += v
		for _, mod := range expr.Modifiers {
			if mod.Kind == ModAdditive {
				total = applyAdditive(total, mod.Value)
			}
		}
	}
	// ------------------------------------------------------------
	// 6. SUCCESS THRESHOLDS
	// ------------------------------------------------------------
	successes := 0
	for _, mod := range expr.Modifiers {
		if mod.Kind == ModSuccessThreshold {
			successes = countSuccesses(kept, mod.Op, mod.Value)
			break
		}
	}

	// ------------------------------------------------------------
	// 7. BUILD RESULT
	// ------------------------------------------------------------
	return Result{
		Expression: expr.Raw,
		Rolls:      postReroll,
		Rerolls:    rerolls,
		Exploded:   exploded,
		Kept:       kept,
		Dropped:    dropped,
		Total:      total,
		Successes:  successes,
	}
}

//
// REROLL HELPERS
//

func doRerollReplace(rng *rand.Rand, rolls []int, sides, threshold int, existing []int) ([]int, []int) {
	rerolls := existing
	for i, v := range rolls {
		if v <= threshold {
			rerolls = append(rerolls, v)
			rolls[i] = rng.Intn(sides) + 1
		}
	}
	return rolls, rerolls
}

func doRerollOnce(rng *rand.Rand, rolls []int, sides, threshold int, existing []int) ([]int, []int) {
	rerolls := existing
	for i, v := range rolls {
		if v <= threshold {
			newVal := rng.Intn(sides) + 1
			rerolls = append(rerolls, newVal)
			rolls[i] = newVal
		}
	}
	return rolls, rerolls
}

func doRerollAdd(rng *rand.Rand, rolls []int, sides, threshold int, existing []int) ([]int, []int) {
	rerolls := existing
	out := rolls
	for _, v := range rolls {
		if v <= threshold {
			newVal := rng.Intn(sides) + 1
			rerolls = append(rerolls, newVal)
			out = append(out, newVal)
		}
	}
	return out, rerolls
}

//
// EXPLODE HELPERS
//

func doExplodeSimple(rng *rand.Rand, rolls []int, sides int, acc []int) ([]int, []int) {
	out := append([]int(nil), rolls...)
	for _, v := range rolls {
		if v == sides {
			newv := rng.Intn(sides) + 1
			acc = append(acc, newv)
			out = append(out, newv)
		}
	}
	return out, acc
}

func doExplodeThreshold(rng *rand.Rand, rolls []int, sides int, threshold int, acc []int) ([]int, []int) {
	out := append([]int(nil), rolls...)
	for _, v := range rolls {
		if v >= threshold {
			newv := rng.Intn(sides) + 1
			acc = append(acc, newv)
			out = append(out, newv)
		}
	}
	return out, acc
}

func doExplodeCompound(rng *rand.Rand, rolls []int, sides int, acc []int) ([]int, []int) {
	out := append([]int(nil), rolls...)
	queue := append([]int(nil), rolls...)

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		if v == sides {
			newv := rng.Intn(sides) + 1
			acc = append(acc, newv)
			out = append(out, newv)
			queue = append(queue, newv)
		}
	}

	return out, acc
}

//
// KEEP/DROP HELPERS
//

func keepHighest(values []int, n int) ([]int, []int) {
	if n >= len(values) {
		return append([]int(nil), values...), nil
	}
	sorted := append([]int(nil), values...)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	return sorted[:n], sorted[n:]
}

func keepLowest(values []int, n int) ([]int, []int) {
	if n >= len(values) {
		return append([]int(nil), values...), nil
	}
	sorted := append([]int(nil), values...)
	sort.Ints(sorted)
	return sorted[:n], sorted[n:]
}

func dropHighest(values []int, n int) ([]int, []int) {
	if n >= len(values) {
		return nil, append([]int(nil), values...)
	}
	sorted := append([]int(nil), values...)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	return sorted[n:], sorted[:n]
}

func dropLowest(values []int, n int) ([]int, []int) {
	if n >= len(values) {
		return nil, append([]int(nil), values...)
	}
	sorted := append([]int(nil), values...)
	sort.Ints(sorted)
	return sorted[n:], sorted[:n]
}
