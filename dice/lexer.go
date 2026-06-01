package dice

import (
	"fmt"
	"strings"
)

// Lex tokenizes a dice expression into a stream of parser tokens.
func Lex(input string) ([]Token, error) {
	l := lexer{input: input}
	tokens := make([]Token, 0, len(input)+1)

	for {
		tok, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}

	return tokens, nil
}

type lexer struct {
	input string
	pos   int
}

func (l *lexer) nextToken() (Token, error) {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF, Pos: l.pos}, nil
	}

	start := l.pos
	ch := l.input[l.pos]

	switch ch {
	case '+':
		l.pos++
		return Token{Type: TokenPlus, Lexeme: "+", Pos: start}, nil
	case '-':
		l.pos++
		return Token{Type: TokenMinus, Lexeme: "-", Pos: start}, nil
	case '*':
		l.pos++
		return Token{Type: TokenStar, Lexeme: "*", Pos: start}, nil
	case '/':
		l.pos++
		return Token{Type: TokenSlash, Lexeme: "/", Pos: start}, nil
	case '(':
		l.pos++
		return Token{Type: TokenLParen, Lexeme: "(", Pos: start}, nil
	case ')':
		l.pos++
		return Token{Type: TokenRParen, Lexeme: ")", Pos: start}, nil
	case ',':
		l.pos++
		return Token{Type: TokenComma, Lexeme: ",", Pos: start}, nil
	case '=':
		l.pos++
		return Token{Type: TokenAssign, Lexeme: "=", Pos: start}, nil
	case '!':
		if l.match("!!") {
			l.pos += 2
			return Token{Type: TokenDoubleBang, Lexeme: "!!", Pos: start}, nil
		}
		l.pos++
		return Token{Type: TokenBang, Lexeme: "!", Pos: start}, nil
	case '>':
		if l.match(">=") {
			l.pos += 2
			return Token{Type: TokenGreaterEqual, Lexeme: ">=", Pos: start}, nil
		}
		l.pos++
		return Token{Type: TokenGreater, Lexeme: ">", Pos: start}, nil
	case '<':
		if l.match("<=") {
			l.pos += 2
			return Token{Type: TokenLessEqual, Lexeme: "<=", Pos: start}, nil
		}
		l.pos++
		return Token{Type: TokenLess, Lexeme: "<", Pos: start}, nil
	}

	if isDigitByte(ch) {
		return l.readNumber()
	}
	if isLetterByte(ch) {
		return l.readWordToken()
	}

	return Token{}, fmt.Errorf("unexpected character %q at position %d", ch, start)
}

func (l *lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\n', '\r':
			l.pos++
		default:
			return
		}
	}
}

func (l *lexer) readNumber() (Token, error) {
	start := l.pos
	for l.pos < len(l.input) && isDigitByte(l.input[l.pos]) {
		l.pos++
	}
	return Token{Type: TokenNumber, Lexeme: l.input[start:l.pos], Pos: start}, nil
}

func (l *lexer) readWordToken() (Token, error) {
	start := l.pos
	remaining := strings.ToLower(l.input[l.pos:])

	switch {
	case strings.HasPrefix(remaining, "kl"):
		l.pos += 2
		return Token{Type: TokenKeepLow, Lexeme: l.input[start:l.pos], Pos: start}, nil
	case strings.HasPrefix(remaining, "dh"):
		l.pos += 2
		return Token{Type: TokenDropHigh, Lexeme: l.input[start:l.pos], Pos: start}, nil
	case strings.HasPrefix(remaining, "dl"):
		l.pos += 2
		return Token{Type: TokenDropLow, Lexeme: l.input[start:l.pos], Pos: start}, nil
	case strings.HasPrefix(remaining, "ro"):
		l.pos += 2
		return Token{Type: TokenRerollOnce, Lexeme: l.input[start:l.pos], Pos: start}, nil
	case strings.HasPrefix(remaining, "ra"):
		l.pos += 2
		return Token{Type: TokenRerollAdd, Lexeme: l.input[start:l.pos], Pos: start}, nil
	case remaining[0] == 'k':
		l.pos++
		return Token{Type: TokenKeep, Lexeme: l.input[start:l.pos], Pos: start}, nil
	case remaining[0] == 'r':
		l.pos++
		return Token{Type: TokenReroll, Lexeme: l.input[start:l.pos], Pos: start}, nil
	case remaining[0] == 'd':
		l.pos++
		return Token{Type: TokenDice, Lexeme: l.input[start:l.pos], Pos: start}, nil
	default:
		for l.pos < len(l.input) && (isLetterByte(l.input[l.pos]) || isDigitByte(l.input[l.pos]) || l.input[l.pos] == '_') {
			l.pos++
		}
		return Token{Type: TokenIdentifier, Lexeme: l.input[start:l.pos], Pos: start}, nil
	}
}

func (l *lexer) match(prefix string) bool {
	return strings.HasPrefix(l.input[l.pos:], prefix)
}

func isDigitByte(b byte) bool {
	return b >= '0' && b <= '9'
}

func isLetterByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
