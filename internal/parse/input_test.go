package parse

import (
	"reflect"
	"testing"
)

func TestParseLine_ArithmeticWithSpacesStaysSingleExpression(t *testing.T) {
	parsed, err := ParseLine("(d20 + 2) * 3")
	if err != nil {
		t.Fatalf("ParseLine returned error: %v", err)
	}

	want := []string{"(d20 + 2) * 3"}
	if !reflect.DeepEqual(parsed.Expressions, want) {
		t.Fatalf("unexpected expressions\nwant: %#v\n got: %#v", want, parsed.Expressions)
	}
}

func TestParseArgs_MultipleExpressionsStillSplit(t *testing.T) {
	parsed, err := ParseArgs([]string{"2d20kh1", "3d8!", "5d10>=8"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	want := []string{"2d20kh1", "3d8!", "5d10>=8"}
	if !reflect.DeepEqual(parsed.Expressions, want) {
		t.Fatalf("unexpected expressions\nwant: %#v\n got: %#v", want, parsed.Expressions)
	}
}

func TestParseArgs_ImplicitDiceExpressionsStillSplit(t *testing.T) {
	parsed, err := ParseArgs([]string{"d20", "d6"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	want := []string{"d20", "d6"}
	if !reflect.DeepEqual(parsed.Expressions, want) {
		t.Fatalf("unexpected expressions\nwant: %#v\n got: %#v", want, parsed.Expressions)
	}
}
