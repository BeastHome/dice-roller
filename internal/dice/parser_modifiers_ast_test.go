package dice

import "testing"

func TestParseTreeExpression_KeepModifier(t *testing.T) {
	tree, err := ParseTreeExpression("4d6k3")
	if err != nil {
		t.Fatalf("ParseTreeExpression returned error: %v", err)
	}

	diceNode, ok := tree.Expr.(*DiceNode)
	if !ok {
		t.Fatalf("expected *DiceNode, got %T", tree.Expr)
	}
	if len(diceNode.Modifiers) != 1 {
		t.Fatalf("expected 1 modifier, got %d", len(diceNode.Modifiers))
	}

	mod, ok := diceNode.Modifiers[0].(*DiceModifierAST)
	if !ok {
		t.Fatalf("expected *DiceModifierAST, got %T", diceNode.Modifiers[0])
	}
	if mod.Kind != ModKeepHigh || mod.Count != 3 {
		t.Fatalf("expected keep-high count 3, got kind=%v count=%d", mod.Kind, mod.Count)
	}
}

func TestParseTreeExpression_RerollAndSuccessModifiers(t *testing.T) {
	tree, err := ParseTreeExpression("5d10ro1>=8")
	if err != nil {
		t.Fatalf("ParseTreeExpression returned error: %v", err)
	}

	diceNode, ok := tree.Expr.(*DiceNode)
	if !ok {
		t.Fatalf("expected *DiceNode, got %T", tree.Expr)
	}
	if len(diceNode.Modifiers) != 2 {
		t.Fatalf("expected 2 modifiers, got %d", len(diceNode.Modifiers))
	}

	reroll := diceNode.Modifiers[0].(*DiceModifierAST)
	if reroll.Kind != ModRerollOnce || reroll.Threshold != 1 {
		t.Fatalf("expected reroll-once threshold 1, got kind=%v threshold=%d", reroll.Kind, reroll.Threshold)
	}

	success := diceNode.Modifiers[1].(*DiceModifierAST)
	if success.Kind != ModSuccessThreshold || success.Op != ">=" || success.Value != 8 {
		t.Fatalf("expected success >= 8, got kind=%v op=%q value=%d", success.Kind, success.Op, success.Value)
	}
}

func TestParseTreeExpression_ExplodeModifiers(t *testing.T) {
	tree, err := ParseTreeExpression("1d6!!")
	if err != nil {
		t.Fatalf("ParseTreeExpression returned error: %v", err)
	}

	diceNode, ok := tree.Expr.(*DiceNode)
	if !ok {
		t.Fatalf("expected *DiceNode, got %T", tree.Expr)
	}
	if len(diceNode.Modifiers) != 1 {
		t.Fatalf("expected 1 modifier, got %d", len(diceNode.Modifiers))
	}

	mod := diceNode.Modifiers[0].(*DiceModifierAST)
	if mod.Kind != ModExplodeCompound {
		t.Fatalf("expected compound explode modifier, got kind=%v", mod.Kind)
	}
}
