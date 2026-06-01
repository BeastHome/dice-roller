package dice

import "testing"

func TestParseTreeExpression_BinaryDiceSum(t *testing.T) {
	tree, err := ParseTreeExpression("2d6 + 3")
	if err != nil {
		t.Fatalf("ParseTreeExpression returned error: %v", err)
	}

	root, ok := tree.Expr.(*BinaryNode)
	if !ok {
		t.Fatalf("expected *BinaryNode, got %T", tree.Expr)
	}
	if root.Op.Type != TokenPlus {
		t.Fatalf("expected plus operator, got %v", root.Op.Type)
	}
	if _, ok := root.Left.(*DiceNode); !ok {
		t.Fatalf("expected left side to be *DiceNode, got %T", root.Left)
	}
	if right, ok := root.Right.(*NumberNode); !ok || right.Value != 3 {
		t.Fatalf("expected right side to be NumberNode(3), got %T %#v", root.Right, root.Right)
	}
}

func TestParseTreeExpression_GroupingAndImplicitDie(t *testing.T) {
	tree, err := ParseTreeExpression("(d20 + 2) * 3")
	if err != nil {
		t.Fatalf("ParseTreeExpression returned error: %v", err)
	}

	root, ok := tree.Expr.(*BinaryNode)
	if !ok {
		t.Fatalf("expected *BinaryNode, got %T", tree.Expr)
	}
	if root.Op.Type != TokenStar {
		t.Fatalf("expected multiply operator, got %v", root.Op.Type)
	}

	group, ok := root.Left.(*GroupNode)
	if !ok {
		t.Fatalf("expected grouped left side, got %T", root.Left)
	}

	inner, ok := group.Inner.(*BinaryNode)
	if !ok {
		t.Fatalf("expected inner grouped expression to be *BinaryNode, got %T", group.Inner)
	}

	diceNode, ok := inner.Left.(*DiceNode)
	if !ok {
		t.Fatalf("expected grouped left term to be *DiceNode, got %T", inner.Left)
	}

	count, ok := diceNode.Count.(*NumberNode)
	if !ok || count.Value != 1 {
		t.Fatalf("expected implicit die count to be 1, got %T %#v", diceNode.Count, diceNode.Count)
	}
}
