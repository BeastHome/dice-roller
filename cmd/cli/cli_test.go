package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunCLI_SameSeedProducesByteIdenticalOutput(t *testing.T) {
	// The contract --seed and any embedded scripting depend on: identical
	// args → identical stdout. Run twice and compare.
	var outA, outB, errBuf bytes.Buffer
	args := []string{"4d6k3", "2d20", "5d10>=7", "--seed", "42", "--no-color"}

	codeA := RunCLI(args, &outA, &errBuf)
	codeB := RunCLI(args, &outB, &errBuf)

	if codeA != 0 || codeB != 0 {
		t.Fatalf("expected exit 0 from both runs; got A=%d B=%d (stderr=%q)", codeA, codeB, errBuf.String())
	}
	if outA.String() != outB.String() {
		t.Fatalf("same-seed CLI output diverged\nA: %q\nB: %q", outA.String(), outB.String())
	}
}

func TestRunCLI_InvalidExpressionExitsNonZero(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := RunCLI([]string{"not-a-dice-expression"}, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("expected non-zero exit for invalid expression; stdout=%q stderr=%q",
			stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "Error evaluating") {
		t.Fatalf("expected error message on stderr, got: %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected stdout to be empty on failure, got: %q", stdout.String())
	}
}

func TestRunCLI_ContinuesPastBadExpressionInBatch(t *testing.T) {
	// One bad expression must not abort the entire batch. Valid
	// expressions still produce output; exit code is non-zero so
	// scripts can detect the partial failure.
	var stdout, stderr bytes.Buffer
	args := []string{"not-dice", "1d1", "--seed", "1", "--no-color"}
	code := RunCLI(args, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("expected non-zero exit (one expression failed); got 0")
	}
	if !strings.Contains(stdout.String(), "[1d1]") {
		t.Fatalf("expected valid expression to still produce output; stdout=%q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "not-dice") {
		t.Fatalf("expected stderr to reference the failed expression; got: %q", stderr.String())
	}
}

func TestRunCLI_MalformedFlagExitsNonZero(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := RunCLI([]string{"--multi", "abc", "1d6"}, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("expected non-zero exit for malformed --multi value")
	}
	if !strings.Contains(stderr.String(), "Error:") {
		t.Fatalf("expected 'Error:' prefix on stderr, got: %q", stderr.String())
	}
}

func TestRunCLI_NoArgsReturnsZeroWithNoOutput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := RunCLI([]string{}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit 0 for no args, got %d", code)
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("expected no output for empty args; stdout=%q stderr=%q",
			stdout.String(), stderr.String())
	}
}

func TestRunCLI_HelpFlagWritesToStdoutAndExitsZero(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := RunCLI([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit 0 for --help, got %d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected --help to write nothing to stderr, got: %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Dice Roller CLI") {
		t.Fatalf("expected help header on stdout, got: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "--seed") {
		t.Fatalf("expected --seed flag documented in help")
	}
}

func TestRunCLI_VersionFlagPrintsVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := RunCLI([]string{"--version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit 0 for --version, got %d", code)
	}
	if !strings.Contains(stdout.String(), "dice-roller version") {
		t.Fatalf("expected version line on stdout, got: %q", stdout.String())
	}
}

func TestRunCLI_MultiModeProducesMultiSummary(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"1d6", "--multi", "5", "--seed", "1", "--no-color"}
	code := RunCLI(args, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d (stderr=%q)", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "rolls=5") {
		t.Fatalf("expected multi summary with rolls=5; got: %q", out)
	}
	if !strings.Contains(out, "AVG") || !strings.Contains(out, "MIN") || !strings.Contains(out, "MAX") {
		t.Fatalf("expected multi stats (AVG/MIN/MAX) in output; got: %q", out)
	}
}

func TestRunCLI_HelpTextExampleExpressionsAllEvaluate(t *testing.T) {
	// Every expression shown in PrintHelp's Examples section must actually
	// evaluate — otherwise the help is advertising broken syntax. This is
	// the regression guard for the 2d20kh1 keep-high parser bug.
	examples := []string{
		"4d6k3",
		"2d20kh1",
		"3d8!",
		"5d10>=8",
		"5d10ro1>=8",
	}
	for _, expr := range examples {
		t.Run(expr, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := RunCLI([]string{expr, "--seed", "1", "--no-color"}, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("help example %q failed: code=%d stderr=%q", expr, code, stderr.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("help example %q wrote to stderr: %q", expr, stderr.String())
			}
		})
	}
}

func TestRunCLI_VerboseModeIncludesDetailLines(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"4d6k3", "--verbose", "--seed", "1", "--no-color"}
	code := RunCLI(args, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := stdout.String()
	for _, want := range []string{"Raw rolls", "Kept", "Total"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected verbose output to include %q; got: %q", want, out)
		}
	}
}

func TestRunCLI_ErrorsGoToStderrNotStdout(t *testing.T) {
	// Scriptability contract: stdout is clean of error messages so
	// scripts can pipe results without filtering.
	var stdout, stderr bytes.Buffer
	args := []string{"garbage", "1d1", "--seed", "1", "--no-color"}
	_ = RunCLI(args, &stdout, &stderr)

	if strings.Contains(stdout.String(), "Error") {
		t.Fatalf("error text leaked into stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Error") {
		t.Fatalf("expected error on stderr, got: %q", stderr.String())
	}
}
