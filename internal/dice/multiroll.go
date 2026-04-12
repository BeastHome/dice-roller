package dice

// EvaluateMulti performs N independent evaluations of the same expression.
// This function is called by Engine.Evaluate() when count > 1.
func EvaluateMulti(e *Engine, expr string, count int) (MultiRollResult, error) {
	results := make([]Result, 0, count)

	for i := 0; i < count; i++ {
		r, err := e.Roll(expr)
		if err != nil {
			return MultiRollResult{}, err
		}
		results = append(results, r)
	}

	return MultiRollResult{
		Expression: expr,
		Rolls:      results,
	}, nil
}
