package dice

import (
	"reflect"
	"testing"
)

func TestLex_BasicExpression(t *testing.T) {
	tokens, err := Lex("2d6 + (1d4 - 3)")
	if err != nil {
		t.Fatalf("Lex returned error: %v", err)
	}

	got := make([]TokenType, len(tokens))
	for i, tok := range tokens {
		got[i] = tok.Type
	}

	want := []TokenType{
		TokenNumber,
		TokenDice,
		TokenNumber,
		TokenPlus,
		TokenLParen,
		TokenNumber,
		TokenDice,
		TokenNumber,
		TokenMinus,
		TokenNumber,
		TokenRParen,
		TokenEOF,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected token types\nwant: %v\n got: %v", want, got)
	}
}

func TestLex_ImplicitDie(t *testing.T) {
	tokens, err := Lex("d20")
	if err != nil {
		t.Fatalf("Lex returned error: %v", err)
	}

	got := []TokenType{tokens[0].Type, tokens[1].Type, tokens[2].Type}
	want := []TokenType{TokenDice, TokenNumber, TokenEOF}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected token types\nwant: %v\n got: %v", want, got)
	}
}

func TestLex_KeepDropRerollModifierTokens(t *testing.T) {
	cases := []struct {
		input string
		want  []TokenType
	}{
		{"4d6k3", []TokenType{TokenNumber, TokenDice, TokenNumber, TokenKeep, TokenNumber, TokenEOF}},
		{"4d6kl2", []TokenType{TokenNumber, TokenDice, TokenNumber, TokenKeepLow, TokenNumber, TokenEOF}},
		{"4d6dh1", []TokenType{TokenNumber, TokenDice, TokenNumber, TokenDropHigh, TokenNumber, TokenEOF}},
		{"4d6dl1", []TokenType{TokenNumber, TokenDice, TokenNumber, TokenDropLow, TokenNumber, TokenEOF}},
		{"4d6r1", []TokenType{TokenNumber, TokenDice, TokenNumber, TokenReroll, TokenNumber, TokenEOF}},
		{"4d6ro1", []TokenType{TokenNumber, TokenDice, TokenNumber, TokenRerollOnce, TokenNumber, TokenEOF}},
		{"4d6ra1", []TokenType{TokenNumber, TokenDice, TokenNumber, TokenRerollAdd, TokenNumber, TokenEOF}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			tokens, err := Lex(tc.input)
			if err != nil {
				t.Fatalf("Lex(%q) returned error: %v", tc.input, err)
			}
			got := make([]TokenType, len(tokens))
			for i, tok := range tokens {
				got[i] = tok.Type
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Lex(%q) tokens\nwant: %v\n got: %v", tc.input, tc.want, got)
			}
		})
	}
}

func TestLex_ExplodeAndCompareOperators(t *testing.T) {
	cases := []struct {
		input string
		want  []TokenType
	}{
		{"d6!", []TokenType{TokenDice, TokenNumber, TokenBang, TokenEOF}},
		{"d6!!", []TokenType{TokenDice, TokenNumber, TokenDoubleBang, TokenEOF}},
		{"d6>=5", []TokenType{TokenDice, TokenNumber, TokenGreaterEqual, TokenNumber, TokenEOF}},
		{"d6<=2", []TokenType{TokenDice, TokenNumber, TokenLessEqual, TokenNumber, TokenEOF}},
		{"d6>4", []TokenType{TokenDice, TokenNumber, TokenGreater, TokenNumber, TokenEOF}},
		{"d6<3", []TokenType{TokenDice, TokenNumber, TokenLess, TokenNumber, TokenEOF}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			tokens, err := Lex(tc.input)
			if err != nil {
				t.Fatalf("Lex(%q) returned error: %v", tc.input, err)
			}
			got := make([]TokenType, len(tokens))
			for i, tok := range tokens {
				got[i] = tok.Type
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Lex(%q) tokens\nwant: %v\n got: %v", tc.input, tc.want, got)
			}
		})
	}
}

func TestLex_FullArithmeticAndGrouping(t *testing.T) {
	tokens, err := Lex("(2d6 * 3) / 2 + 1")
	if err != nil {
		t.Fatalf("Lex returned error: %v", err)
	}
	want := []TokenType{
		TokenLParen,
		TokenNumber, TokenDice, TokenNumber, TokenStar, TokenNumber,
		TokenRParen,
		TokenSlash, TokenNumber,
		TokenPlus, TokenNumber,
		TokenEOF,
	}
	got := make([]TokenType, len(tokens))
	for i, tok := range tokens {
		got[i] = tok.Type
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected tokens\nwant: %v\n got: %v", want, got)
	}
}

func TestLex_RejectsUnknownCharacter(t *testing.T) {
	if _, err := Lex("2d6 @ 3"); err == nil {
		t.Fatalf("expected error for unknown character '@', got nil")
	}
}
