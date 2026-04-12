package dice

// Expression is the parsed form of a dice expression.
// It currently serves as both the AST root and the
// evaluation configuration for a single roll.
type Expression struct {
	Raw       string
	Count     int
	Sides     int
	Modifiers []Modifier
}
