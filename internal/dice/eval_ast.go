package dice

import (
	"fmt"
	"math/rand"
	"strings"
)

// EvaluateParseTree evaluates an AST-based parse tree and produces the legacy
// Result shape expected by the rest of the application.
func EvaluateParseTree(rng *rand.Rand, tree *ParseTree) (Result, error) {
	if tree == nil || tree.Expr == nil {
		return Result{}, fmt.Errorf("empty parse tree")
	}

	res, err := evaluateExprNode(rng, tree.Source, tree.Expr)
	if err != nil {
		return Result{}, err
	}

	res.Expression = strings.TrimSpace(tree.Source)
	return res, nil
}

func evaluateExprNode(rng *rand.Rand, source string, node ExprNode) (Result, error) {
	switch n := node.(type) {
	case *NumberNode:
		return Result{
			Expression: sliceSourceBySpan(source, n.Range),
			Total:      n.Value,
		}, nil

	case *DiceNode:
		expr, err := buildExpressionFromNode(source, n)
		if err != nil {
			return Result{}, err
		}
		res := EvaluateSingle(rng, expr)
		res.Expression = expr.Raw
		return res, nil

	case *UnaryNode:
		res, err := evaluateExprNode(rng, source, n.Right)
		if err != nil {
			return Result{}, err
		}
		switch n.Op.Type {
		case TokenPlus:
			res.Expression = sliceSourceBySpan(source, n.Range)
			return res, nil
		case TokenMinus:
			res.Total = -res.Total
			res.Expression = sliceSourceBySpan(source, n.Range)
			return res, nil
		default:
			return Result{}, fmt.Errorf("unsupported unary operator %q", n.Op.Lexeme)
		}

	case *GroupNode:
		res, err := evaluateExprNode(rng, source, n.Inner)
		if err != nil {
			return Result{}, err
		}
		res.Expression = sliceSourceBySpan(source, n.Range)
		return res, nil

	case *BinaryNode:
		left, err := evaluateExprNode(rng, source, n.Left)
		if err != nil {
			return Result{}, err
		}
		right, err := evaluateExprNode(rng, source, n.Right)
		if err != nil {
			return Result{}, err
		}
		return combineBinaryResults(source, n, left, right)

	default:
		return Result{}, fmt.Errorf("unsupported AST node %T", node)
	}
}

func combineBinaryResults(source string, node *BinaryNode, left, right Result) (Result, error) {
	res := Result{
		Expression: sliceSourceBySpan(source, node.Range),
		Rolls:      append(copyInts(left.Rolls), right.Rolls...),
		Rerolls:    append(copyInts(left.Rerolls), right.Rerolls...),
		Exploded:   append(copyInts(left.Exploded), right.Exploded...),
		Kept:       append(copyInts(left.Kept), right.Kept...),
		Dropped:    append(copyInts(left.Dropped), right.Dropped...),
		Successes:  left.Successes + right.Successes,
	}

	switch node.Op.Type {
	case TokenPlus:
		res.Total = left.Total + right.Total
	case TokenMinus:
		res.Total = left.Total - right.Total
	case TokenStar:
		res.Total = left.Total * right.Total
	case TokenSlash:
		if right.Total == 0 {
			return Result{}, fmt.Errorf("division by zero")
		}
		res.Total = left.Total / right.Total
	default:
		return Result{}, fmt.Errorf("unsupported binary operator %q", node.Op.Lexeme)
	}

	return res, nil
}

func copyInts(src []int) []int {
	if len(src) == 0 {
		return nil
	}
	out := make([]int, len(src))
	copy(out, src)
	return out
}
