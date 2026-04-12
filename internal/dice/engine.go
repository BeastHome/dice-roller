package dice

import (
	"fmt"
	"math/rand"
)

// Engine is the public entry point for evaluating dice expressions.
type Engine struct {
	rng *rand.Rand
}

// NewEngine constructs a new Engine with a seeded RNG.
func NewEngine() *Engine {
	return &Engine{
		rng: newRNG(defaultSeed()),
	}
}

// Roll parses and evaluates a single dice expression string.
func (e *Engine) Roll(expr string) (Result, error) {
	if tree, err := ParseTreeExpression(expr); err == nil {
		if res, evalErr := EvaluateParseTree(e.rng, tree); evalErr == nil {
			AttachVerbose(&res)
			return res, nil
		}
	}

	ast, err := ParseExpression(expr)
	if err != nil {
		return Result{}, fmt.Errorf("parse error: %w", err)
	}

	res := EvaluateSingle(e.rng, ast)
	AttachVerbose(&res)
	return res, nil
}

// Evaluate is the unified API for single or multi‑roll evaluation.
func (e *Engine) Evaluate(expr string, count int) (interface{}, error) {
	if count <= 1 {
		return e.Roll(expr)
	}
	return EvaluateMulti(e, expr, count)
}
