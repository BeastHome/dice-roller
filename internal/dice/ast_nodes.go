package dice

// Span describes a half-open byte range in the original input: [Start, End).
type Span struct {
	Start int
	End   int
}

func (s Span) Len() int {
	return s.End - s.Start
}

// Node is the common interface for all future parser AST nodes.
type Node interface {
	Span() Span
}

// ExprNode marks AST nodes that evaluate to a numeric expression value.
type ExprNode interface {
	Node
	exprNode()
}

// ModifierNode marks AST nodes that attach to a DiceNode.
type ModifierNode interface {
	Node
	modifierNode()
}

// ParseTree is the root object returned by the future recursive parser.
// The current evaluator still uses Expression for compatibility.
type ParseTree struct {
	Source string
	Expr   ExprNode
}

// NumberNode represents a plain integer literal.
type NumberNode struct {
	Value int
	Range Span
}

func (n *NumberNode) Span() Span { return n.Range }
func (n *NumberNode) exprNode()  {}

// UnaryNode represents prefix expressions like -1 or +2.
type UnaryNode struct {
	Op    Token
	Right ExprNode
	Range Span
}

func (n *UnaryNode) Span() Span { return n.Range }
func (n *UnaryNode) exprNode()  {}

// BinaryNode represents infix arithmetic like 2d6 + 3.
type BinaryNode struct {
	Left  ExprNode
	Op    Token
	Right ExprNode
	Range Span
}

func (n *BinaryNode) Span() Span { return n.Range }
func (n *BinaryNode) exprNode()  {}

// GroupNode preserves explicit parenthesized grouping.
type GroupNode struct {
	Inner ExprNode
	Range Span
}

func (n *GroupNode) Span() Span { return n.Range }
func (n *GroupNode) exprNode()  {}

// DiceNode represents a dice term such as d20, 4d6, or 2d20kh1.
type DiceNode struct {
	Count     ExprNode
	Sides     ExprNode
	Modifiers []ModifierNode
	Range     Span
}

func (n *DiceNode) Span() Span { return n.Range }
func (n *DiceNode) exprNode()  {}

// DiceModifierAST represents a modifier attached to a dice term.
type DiceModifierAST struct {
	Kind      ModifierKind
	Token     Token
	Value     int
	Threshold int
	Count     int
	Op        string
	Range     Span
}

func (n *DiceModifierAST) Span() Span    { return n.Range }
func (n *DiceModifierAST) modifierNode() {}
