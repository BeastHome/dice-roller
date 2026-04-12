package dice

import (
	"fmt"
	"strconv"
)

// ParseTreeExpression parses an input string into the new AST-based parse tree.
// This Step 1 parser currently handles numbers, dice terms, arithmetic, and grouping.
func ParseTreeExpression(input string) (*ParseTree, error) {
	tokens, err := Lex(input)
	if err != nil {
		return nil, err
	}

	p := rdParser{source: input, tokens: tokens}
	expr, err := p.parseExpression(PrecLowest)
	if err != nil {
		return nil, err
	}

	if tok := p.current(); tok.Type != TokenEOF {
		return nil, fmt.Errorf("unexpected token %q at position %d", tok.Lexeme, tok.Pos)
	}

	return &ParseTree{Source: input, Expr: expr}, nil
}

type rdParser struct {
	source string
	tokens []Token
	pos    int
}

func (p *rdParser) parseExpression(minPrec Precedence) (ExprNode, error) {
	left, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.current()
		if !tok.Type.IsBinaryOperator() {
			break
		}

		prec := tok.Type.Precedence()
		if prec < minPrec {
			break
		}

		op := p.advance()
		right, err := p.parseExpression(prec + 1)
		if err != nil {
			return nil, err
		}

		left = &BinaryNode{
			Left:  left,
			Op:    op,
			Right: right,
			Range: combineSpans(left.Span(), right.Span()),
		}
	}

	return left, nil
}

func (p *rdParser) parsePrefix() (ExprNode, error) {
	tok := p.current()

	switch tok.Type {
	case TokenNumber:
		return p.parseNumberOrDice()
	case TokenDice:
		return p.parseImplicitDice()
	case TokenPlus, TokenMinus:
		op := p.advance()
		right, err := p.parseExpression(PrecPrefix)
		if err != nil {
			return nil, err
		}
		return &UnaryNode{
			Op:    op,
			Right: right,
			Range: Span{Start: op.Pos, End: right.Span().End},
		}, nil
	case TokenLParen:
		return p.parseGroupedExpression()
	case TokenEOF:
		return nil, fmt.Errorf("unexpected end of input")
	default:
		return nil, fmt.Errorf("unexpected token %q at position %d", tok.Lexeme, tok.Pos)
	}
}

func (p *rdParser) parseNumberOrDice() (ExprNode, error) {
	numberTok := p.advance()
	value, err := strconv.Atoi(numberTok.Lexeme)
	if err != nil {
		return nil, fmt.Errorf("invalid number %q at position %d", numberTok.Lexeme, numberTok.Pos)
	}

	number := &NumberNode{Value: value, Range: spanFromToken(numberTok)}
	if p.current().Type == TokenDice {
		return p.finishDiceNode(number, p.advance())
	}

	return number, nil
}

func (p *rdParser) parseImplicitDice() (ExprNode, error) {
	diceTok := p.advance()
	count := &NumberNode{
		Value: 1,
		Range: Span{Start: diceTok.Pos, End: diceTok.Pos},
	}
	return p.finishDiceNode(count, diceTok)
}

func (p *rdParser) finishDiceNode(count ExprNode, diceTok Token) (ExprNode, error) {
	sidesTok := p.current()
	if sidesTok.Type != TokenNumber {
		return nil, fmt.Errorf("expected dice sides after %q at position %d", diceTok.Lexeme, diceTok.Pos)
	}
	p.advance()

	sidesValue, err := strconv.Atoi(sidesTok.Lexeme)
	if err != nil {
		return nil, fmt.Errorf("invalid dice sides %q at position %d", sidesTok.Lexeme, sidesTok.Pos)
	}

	start := count.Span().Start
	if count.Span().Len() == 0 {
		start = diceTok.Pos
	}

	node := &DiceNode{
		Count:     count,
		Sides:     &NumberNode{Value: sidesValue, Range: spanFromToken(sidesTok)},
		Modifiers: nil,
		Range:     Span{Start: start, End: sidesTok.End()},
	}

	for p.current().Type.IsModifierToken() {
		mod, err := p.parseDiceModifier()
		if err != nil {
			return nil, err
		}
		node.Modifiers = append(node.Modifiers, mod)
		node.Range.End = mod.Span().End
	}

	return node, nil
}

func (p *rdParser) parseGroupedExpression() (ExprNode, error) {
	open := p.advance()
	inner, err := p.parseExpression(PrecLowest)
	if err != nil {
		return nil, err
	}

	closeTok := p.current()
	if closeTok.Type != TokenRParen {
		return nil, fmt.Errorf("expected ')' at position %d", closeTok.Pos)
	}
	p.advance()

	return &GroupNode{
		Inner: inner,
		Range: Span{Start: open.Pos, End: closeTok.End()},
	}, nil
}

func (p *rdParser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF, Pos: len(p.source)}
	}
	return p.tokens[p.pos]
}

func (p *rdParser) advance() Token {
	tok := p.current()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

func (p *rdParser) parseDiceModifier() (ModifierNode, error) {
	switch p.current().Type {
	case TokenKeep:
		return p.parseCountModifier(ModKeepHigh)
	case TokenKeepLow:
		return p.parseCountModifier(ModKeepLow)
	case TokenDropHigh:
		return p.parseCountModifier(ModDropHigh)
	case TokenDropLow:
		return p.parseCountModifier(ModDropLow)
	case TokenReroll:
		return p.parseThresholdModifier(ModReroll)
	case TokenRerollOnce:
		return p.parseThresholdModifier(ModRerollOnce)
	case TokenRerollAdd:
		return p.parseThresholdModifier(ModRerollAdd)
	case TokenBang:
		return p.parseExplodeModifier(false)
	case TokenDoubleBang:
		return p.parseExplodeModifier(true)
	case TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual:
		return p.parseSuccessModifier()
	default:
		tok := p.current()
		return nil, fmt.Errorf("unexpected modifier token %q at position %d", tok.Lexeme, tok.Pos)
	}
}

func (p *rdParser) parseCountModifier(kind ModifierKind) (ModifierNode, error) {
	tok := p.advance()
	numTok, value, err := p.consumeNumberToken(tok.Lexeme)
	if err != nil {
		return nil, err
	}

	return &DiceModifierAST{
		Kind:  kind,
		Token: tok,
		Count: value,
		Range: Span{Start: tok.Pos, End: numTok.End()},
	}, nil
}

func (p *rdParser) parseThresholdModifier(kind ModifierKind) (ModifierNode, error) {
	tok := p.advance()
	numTok, value, err := p.consumeNumberToken(tok.Lexeme)
	if err != nil {
		return nil, err
	}

	return &DiceModifierAST{
		Kind:      kind,
		Token:     tok,
		Threshold: value,
		Range:     Span{Start: tok.Pos, End: numTok.End()},
	}, nil
}

func (p *rdParser) parseExplodeModifier(compound bool) (ModifierNode, error) {
	tok := p.advance()
	kind := ModExplode
	if compound {
		kind = ModExplodeCompound
	}

	mod := &DiceModifierAST{
		Kind:  kind,
		Token: tok,
		Range: spanFromToken(tok),
	}

	switch p.current().Type {
	case TokenNumber:
		numTok, value, err := p.consumeNumberToken(tok.Lexeme)
		if err != nil {
			return nil, err
		}
		mod.Threshold = value
		mod.Range.End = numTok.End()
		if !compound {
			mod.Kind = ModExplodeThreshold
		}

	case TokenGreater, TokenGreaterEqual:
		opTok := p.advance()
		numTok, value, err := p.consumeNumberToken(opTok.Lexeme)
		if err != nil {
			return nil, err
		}
		mod.Op = opTok.Lexeme
		mod.Threshold = value
		mod.Range.End = numTok.End()
		if !compound {
			mod.Kind = ModExplodeThreshold
		}
	}

	return mod, nil
}

func (p *rdParser) parseSuccessModifier() (ModifierNode, error) {
	tok := p.advance()
	numTok, value, err := p.consumeNumberToken(tok.Lexeme)
	if err != nil {
		return nil, err
	}

	return &DiceModifierAST{
		Kind:  ModSuccessThreshold,
		Token: tok,
		Value: value,
		Op:    tok.Lexeme,
		Range: Span{Start: tok.Pos, End: numTok.End()},
	}, nil
}

func (p *rdParser) consumeNumberToken(after string) (Token, int, error) {
	tok := p.current()
	if tok.Type != TokenNumber {
		return Token{}, 0, fmt.Errorf("expected number after %q at position %d", after, tok.Pos)
	}
	p.advance()

	value, err := strconv.Atoi(tok.Lexeme)
	if err != nil {
		return Token{}, 0, fmt.Errorf("invalid number %q at position %d", tok.Lexeme, tok.Pos)
	}
	return tok, value, nil
}

func spanFromToken(tok Token) Span {
	return Span{Start: tok.Pos, End: tok.End()}
}

func combineSpans(a, b Span) Span {
	return Span{Start: a.Start, End: b.End}
}
