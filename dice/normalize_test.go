package dice

import (
	"reflect"
	"testing"
)

func TestNormalizeModifiers_LastKeepDropWins(t *testing.T) {
	mods := []Modifier{
		{Kind: ModKeepHigh, Count: 3},
		{Kind: ModKeepLow, Count: 2},
		{Kind: ModDropHigh, Count: 1},
	}
	out := normalizeModifiers(mods)
	if len(out) != 1 {
		t.Fatalf("expected single keep/drop modifier, got %d", len(out))
	}
	if out[0].Kind != ModDropHigh || out[0].Count != 1 {
		t.Fatalf("expected last keep/drop to win, got %#v", out[0])
	}
}

func TestNormalizeModifiers_AdditivesAreDropped(t *testing.T) {
	// Slice B removed additive handling from EvaluateSingle (arithmetic
	// is the AST evaluator's job). normalizeModifiers now drops all
	// ModAdditive entries — they're a no-op in the legacy path.
	mods := []Modifier{
		{Kind: ModAdditive, Value: 3},
		{Kind: ModAdditive, Value: -5},
		{Kind: ModAdditive, Value: 7},
	}
	out := normalizeModifiers(mods)
	if len(out) != 0 {
		t.Fatalf("expected all additives to be dropped, got %d: %#v", len(out), out)
	}
}

func TestNormalizeModifiers_FirstRerollFamilyMemberWins(t *testing.T) {
	mods := []Modifier{
		{Kind: ModRerollOnce, Threshold: 1},
		{Kind: ModReroll, Threshold: 3},
		{Kind: ModRerollAdd, Threshold: 5},
	}
	out := normalizeModifiers(mods)
	if len(out) != 1 {
		t.Fatalf("expected one reroll modifier, got %d", len(out))
	}
	if out[0].Kind != ModRerollOnce || out[0].Threshold != 1 {
		t.Fatalf("expected first-encountered reroll-once to win, got %#v", out[0])
	}
}

func TestNormalizeModifiers_FirstSuccessThresholdWins(t *testing.T) {
	mods := []Modifier{
		{Kind: ModSuccessThreshold, Op: ">=", Value: 8},
		{Kind: ModSuccessThreshold, Op: ">", Value: 9},
	}
	out := normalizeModifiers(mods)
	if len(out) != 1 {
		t.Fatalf("expected one success threshold, got %d", len(out))
	}
	if out[0].Op != ">=" || out[0].Value != 8 {
		t.Fatalf("expected first success >=8 to win, got %#v", out[0])
	}
}

func TestNormalizeModifiers_PreservesOrderOfUnnormalizedKinds(t *testing.T) {
	mods := []Modifier{
		{Kind: ModExplode},
		{Kind: ModKeepHigh, Count: 3},
		{Kind: ModExplodeThreshold, Threshold: 5},
	}
	out := normalizeModifiers(mods)
	kinds := make([]ModifierKind, len(out))
	for i, m := range out {
		kinds[i] = m.Kind
	}
	want := []ModifierKind{ModExplode, ModKeepHigh, ModExplodeThreshold}
	if !reflect.DeepEqual(kinds, want) {
		t.Fatalf("unexpected modifier order\nwant: %v\n got: %v", want, kinds)
	}
}
