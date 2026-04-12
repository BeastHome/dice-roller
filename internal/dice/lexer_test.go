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
