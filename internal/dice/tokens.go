package dice

// TokenType identifies the kind of lexical token produced by the
// future dice expression lexer.
type TokenType int

const (
	TokenIllegal TokenType = iota
	TokenEOF

	// Literals / generic identifiers.
	TokenNumber
	TokenIdentifier

	// Core dice grammar.
	TokenDice // d / D

	// Arithmetic.
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash

	// Grouping / separators.
	TokenLParen
	TokenRParen
	TokenComma
	TokenAssign

	// Comparison / explode operators.
	TokenBang
	TokenDoubleBang
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual

	// Dice modifier keywords.
	TokenKeep
	TokenKeepLow
	TokenDropHigh
	TokenDropLow
	TokenReroll
	TokenRerollOnce
	TokenRerollAdd
)

// Precedence describes operator binding strength for the future
// recursive descent / Pratt-style expression parser.
type Precedence int

const (
	PrecLowest Precedence = iota
	PrecSum
	PrecProduct
	PrecPrefix
)

// Token represents a single lexed token from the source input.
type Token struct {
	Type   TokenType
	Lexeme string
	Pos    int // byte offset in the original input
}

func (t Token) End() int {
	return t.Pos + len(t.Lexeme)
}

func (tt TokenType) String() string {
	switch tt {
	case TokenIllegal:
		return "ILLEGAL"
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "NUMBER"
	case TokenIdentifier:
		return "IDENT"
	case TokenDice:
		return "DICE"
	case TokenPlus:
		return "PLUS"
	case TokenMinus:
		return "MINUS"
	case TokenStar:
		return "STAR"
	case TokenSlash:
		return "SLASH"
	case TokenLParen:
		return "LPAREN"
	case TokenRParen:
		return "RPAREN"
	case TokenComma:
		return "COMMA"
	case TokenAssign:
		return "ASSIGN"
	case TokenBang:
		return "BANG"
	case TokenDoubleBang:
		return "DOUBLE_BANG"
	case TokenGreater:
		return "GREATER"
	case TokenGreaterEqual:
		return "GREATER_EQUAL"
	case TokenLess:
		return "LESS"
	case TokenLessEqual:
		return "LESS_EQUAL"
	case TokenKeep:
		return "KEEP"
	case TokenKeepLow:
		return "KEEP_LOW"
	case TokenDropHigh:
		return "DROP_HIGH"
	case TokenDropLow:
		return "DROP_LOW"
	case TokenReroll:
		return "REROLL"
	case TokenRerollOnce:
		return "REROLL_ONCE"
	case TokenRerollAdd:
		return "REROLL_ADD"
	default:
		return "UNKNOWN"
	}
}

func (tt TokenType) Precedence() Precedence {
	switch tt {
	case TokenPlus, TokenMinus:
		return PrecSum
	case TokenStar, TokenSlash:
		return PrecProduct
	default:
		return PrecLowest
	}
}

func (tt TokenType) IsBinaryOperator() bool {
	switch tt {
	case TokenPlus, TokenMinus, TokenStar, TokenSlash:
		return true
	default:
		return false
	}
}

func (tt TokenType) IsComparisonOperator() bool {
	switch tt {
	case TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual:
		return true
	default:
		return false
	}
}

func (tt TokenType) IsModifierToken() bool {
	switch tt {
	case TokenBang, TokenDoubleBang,
		TokenKeep, TokenKeepLow,
		TokenDropHigh, TokenDropLow,
		TokenReroll, TokenRerollOnce, TokenRerollAdd,
		TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual:
		return true
	default:
		return false
	}
}
