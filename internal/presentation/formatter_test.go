package presentation

import (
	"strings"
	"testing"

	"github.com/showr/dice-roller/internal/dice"
)

// plainFormatter returns a ColoredFormatter wired with the
// color-disabled scheme so output strings are stable and assertable.
func plainFormatter() *ColoredFormatter {
	return NewColoredFormatter(GetColorScheme(true))
}

func TestColoredFormatter_CompactSingle_NoEffects(t *testing.T) {
	f := plainFormatter()
	r := dice.Result{Expression: "4d6k3", Total: 15}
	got := f.FormatCompactSingle(r)
	want := "[4d6k3] → 15"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestColoredFormatter_CompactSingle_WithRerollsAndExploded(t *testing.T) {
	f := plainFormatter()
	r := dice.Result{
		Expression: "4d6r1!",
		Rerolls:    []int{1},
		Exploded:   []int{6},
		Total:      18,
	}
	got := f.FormatCompactSingle(r)
	if !strings.Contains(got, "[4d6r1!] → 18") {
		t.Fatalf("expected base format with arrow and total, got %q", got)
	}
	if !strings.Contains(got, "rerolled: [1]") {
		t.Fatalf("expected rerolled effect in output, got %q", got)
	}
	if !strings.Contains(got, "exploded: [6]") {
		t.Fatalf("expected exploded effect in output, got %q", got)
	}
}

func TestColoredFormatter_CompactMulti_IncludesAllStats(t *testing.T) {
	f := plainFormatter()
	mr := dice.MultiRollResult{
		Expression: "1d6 rolls=3",
		Rolls: []dice.Result{
			{Total: 4},
			{Total: 2},
			{Total: 6},
		},
	}
	got := f.FormatCompactMulti(mr)
	for _, want := range []string{"AVG 4.00", "MED 4.00", "MIN 2", "MAX 6", "SD"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
	if !strings.Contains(got, "Totals:") {
		t.Fatalf("expected Totals: line, got %q", got)
	}
}

func TestColoredFormatter_CompactMulti_EmptyRollsReturnsBareExpression(t *testing.T) {
	f := plainFormatter()
	mr := dice.MultiRollResult{Expression: "4d6"}
	got := f.FormatCompactMulti(mr)
	want := "[4d6]"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestColoredFormatter_VerboseSingle_LabelsAndValues(t *testing.T) {
	f := plainFormatter()
	r := dice.Result{
		Expression: "4d6k3",
		Rolls:      []int{6, 5, 4, 2},
		Kept:       []int{6, 5, 4},
		Dropped:    []int{2},
		Total:      15,
	}
	got := f.FormatVerboseSingle(r)
	for _, want := range []string{
		"[4d6k3]",
		"Raw rolls:",
		"6 5 4 2",
		"Kept:",
		"6 5 4",
		"Dropped:",
		"Total:",
		"15",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in verbose output, got:\n%s", want, got)
		}
	}
}

func TestColoredFormatter_VerboseSingle_SuccessesLineOnlyWhenPositive(t *testing.T) {
	f := plainFormatter()
	withZero := f.FormatVerboseSingle(dice.Result{Expression: "1d6", Total: 4, Kept: []int{4}})
	if strings.Contains(withZero, "Successes:") {
		t.Fatalf("Successes: line should be omitted when count=0, got:\n%s", withZero)
	}

	withSome := f.FormatVerboseSingle(dice.Result{Expression: "5d6>=4", Total: 3, Kept: []int{4, 5, 6}, Successes: 3})
	if !strings.Contains(withSome, "Successes:") {
		t.Fatalf("Successes: line should be present when count>0, got:\n%s", withSome)
	}
}

func TestColoredFormatter_VerboseMulti_PerRollDetailAndSeparators(t *testing.T) {
	f := plainFormatter()
	mr := dice.MultiRollResult{
		Expression: "1d6 rolls=2",
		Rolls: []dice.Result{
			{Expression: "1d6", Total: 3, Rolls: []int{3}, Kept: []int{3}},
			{Expression: "1d6", Total: 5, Rolls: []int{5}, Kept: []int{5}},
		},
	}
	got := f.FormatVerboseMulti(mr)
	if !strings.Contains(got, "Roll 1/2") || !strings.Contains(got, "Roll 2/2") {
		t.Fatalf("expected per-roll Roll N/M markers, got:\n%s", got)
	}
	if !strings.Contains(got, "AVG 4.00") {
		t.Fatalf("expected aggregate stats header, got:\n%s", got)
	}
}

func TestDefaultFormatter_SingleSummary(t *testing.T) {
	f := NewDefaultFormatter()
	r := dice.Result{
		Expression: "4d6k3",
		Rolls:      []int{6, 5, 4, 2},
		Kept:       []int{6, 5, 4},
		Total:      15,
	}
	got := f.FormatSingleSummary(r)
	want := "4d6k3 | total=15 | rolls=[6 5 4 2] kept=[6 5 4]"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDefaultFormatter_MultiSummary_AppendsRollsSuffixAndStats(t *testing.T) {
	f := NewDefaultFormatter()
	mr := dice.MultiRollResult{
		Expression: "1d6",
		Rolls: []dice.Result{
			{Total: 3}, {Total: 5}, {Total: 4},
		},
	}
	got := f.FormatMultiSummary(mr)
	want := "1d6 rolls=3 | avg=4.00 | min=3 | max=5"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDefaultFormatter_MultiSummary_EmptyRollsReturnsExpressionWithRollsZero(t *testing.T) {
	// FormatMultiExpression unconditionally appends "rolls=N" when the
	// expression doesn't already contain it, including N=0. Characterize
	// that quirk so it isn't accidentally "fixed" in a refactor — and so
	// the next person knows the empty-rolls case looks like "4d6 rolls=0",
	// not bare "4d6".
	f := NewDefaultFormatter()
	got := f.FormatMultiSummary(dice.MultiRollResult{Expression: "4d6", Rolls: nil})
	if got != "4d6 rolls=0" {
		t.Fatalf("expected %q, got %q", "4d6 rolls=0", got)
	}
}

func TestDefaultFormatter_VerboseUsesAttachedVerboseIfPresent(t *testing.T) {
	f := NewDefaultFormatter()
	r := dice.Result{Expression: "1d6", Total: 4, Verbose: "PRE-RENDERED VERBOSE\n"}
	got := f.FormatVerbose(r)
	if got != "PRE-RENDERED VERBOSE\n" {
		t.Fatalf("expected pre-rendered verbose to be returned as-is, got %q", got)
	}
}

func TestDefaultFormatter_VerboseRendersStructFieldsWhenAttachedEmpty(t *testing.T) {
	f := NewDefaultFormatter()
	r := dice.Result{
		Expression: "4d6k3",
		Rolls:      []int{6, 5, 4, 2},
		Kept:       []int{6, 5, 4},
		Dropped:    []int{2},
		Total:      15,
	}
	got := f.FormatVerbose(r)
	for _, want := range []string{"Expression: 4d6k3", "Rolls: [6 5 4 2]", "Kept: [6 5 4]", "Dropped: [2]", "Total: 15"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in verbose output, got:\n%s", want, got)
		}
	}
}

func TestDefaultFormatter_MultiRollLine(t *testing.T) {
	f := NewDefaultFormatter()
	got := f.FormatMultiRollLine(dice.Result{Total: 17}, 4)
	want := "Roll 4: 17"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSimpleFormat(t *testing.T) {
	got := SimpleFormat(dice.Result{Expression: "2d20kh1", Total: 18})
	want := "2d20kh1 -> 18"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestGetColorScheme_DisabledReturnsEmptyCodes(t *testing.T) {
	cs := GetColorScheme(true)
	for _, code := range []string{cs.Reset, cs.Bold, cs.Dim, cs.Stats, cs.Kept, cs.Dropped, cs.Reroll, cs.Exploded, cs.Total, cs.Success} {
		if code != "" {
			t.Fatalf("expected empty color code when disabled, got %q", code)
		}
	}
}

func TestColorScheme_ColAndColFNoOpWithEmptyColor(t *testing.T) {
	cs := GetColorScheme(true)
	if got := cs.Col("", "hello"); got != "hello" {
		t.Fatalf("Col with empty color should return text unchanged, got %q", got)
	}
	if got := cs.ColF("", "value=%d", 42); got != "value=42" {
		t.Fatalf("ColF with empty color should return formatted text only, got %q", got)
	}
}
