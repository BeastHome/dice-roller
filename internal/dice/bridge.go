package dice

import (
	"fmt"
	"strings"
)

// BuildExpressionFromTree converts a simple dice-only AST into the legacy
// Expression type used by the current evaluator. Rich arithmetic trees are
// intentionally rejected until AST evaluation is introduced.
func BuildExpressionFromTree(tree *ParseTree) (Expression, error) {
	if tree == nil || tree.Expr == nil {
		return Expression{}, fmt.Errorf("empty parse tree")
	}
	return buildExpressionFromNode(tree.Source, tree.Expr)
}

func buildExpressionFromNode(source string, node ExprNode) (Expression, error) {
	switch n := node.(type) {
	case *DiceNode:
		count, err := extractLiteralInt(n.Count, "dice count")
		if err != nil {
			return Expression{}, err
		}
		sides, err := extractLiteralInt(n.Sides, "dice sides")
		if err != nil {
			return Expression{}, err
		}

		expr := Expression{
			Raw:   sliceSourceBySpan(source, n.Range),
			Count: count,
			Sides: sides,
		}

		for _, modNode := range n.Modifiers {
			mod, err := buildModifierFromNode(modNode)
			if err != nil {
				return Expression{}, err
			}
			expr.Modifiers = append(expr.Modifiers, mod)
		}

		expr.Modifiers = normalizeModifiers(expr.Modifiers)
		if err := validateExpression(expr); err != nil {
			return Expression{}, err
		}
		return expr, nil

	case *GroupNode:
		return buildExpressionFromNode(source, n.Inner)

	default:
		return Expression{}, fmt.Errorf("expression %T is not yet supported by the legacy evaluator bridge", node)
	}
}

func buildModifierFromNode(node ModifierNode) (Modifier, error) {
	mod, ok := node.(*DiceModifierAST)
	if !ok {
		return Modifier{}, fmt.Errorf("unsupported modifier node %T", node)
	}

	return Modifier{
		Kind:      mod.Kind,
		Value:     mod.Value,
		Threshold: mod.Threshold,
		Op:        mod.Op,
		Count:     mod.Count,
	}, nil
}

func extractLiteralInt(node ExprNode, label string) (int, error) {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value, nil
	case *GroupNode:
		return extractLiteralInt(n.Inner, label)
	default:
		return 0, fmt.Errorf("%s must be a literal number, got %T", label, node)
	}
}

func sliceSourceBySpan(source string, span Span) string {
	if source == "" {
		return ""
	}
	if span.Start < 0 {
		span.Start = 0
	}
	if span.End > len(source) {
		span.End = len(source)
	}
	if span.Start >= span.End {
		return ""
	}
	return strings.TrimSpace(source[span.Start:span.End])
}
