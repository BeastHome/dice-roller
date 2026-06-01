package dice

import (
	"fmt"
	"math/rand"
)

// Engine is the public entry point for evaluating dice expressions.
type Engine struct {
	rng *rand.Rand
}

// NewEngine constructs a new Engine with a time-based RNG seed. For
// reproducible output (replay, tests, --seed), use NewEngineWithSeed
// or NewEngineWithOptions.
func NewEngine() *Engine {
	return &Engine{
		rng: newRNG(defaultSeed()),
	}
}

// Roll parses and evaluates a single dice expression string.
//
// Evaluation strategy: try the AST path (parser_rd + eval_ast) first;
// only fall back to the legacy regex parser if the AST parser cannot
// recognize the expression. Errors from a successful AST parse (e.g.,
// division by zero) are returned directly — they are NOT masked by
// retrying the legacy path.
func (e *Engine) Roll(expr string) (Result, error) {
	if tree, parseErr := ParseTreeExpression(expr); parseErr == nil {
		res, evalErr := EvaluateParseTree(e.rng, tree)
		if evalErr != nil {
			return Result{}, evalErr
		}
		AttachVerbose(&res)
		return res, nil
	}

	ast, err := ParseExpression(expr)
	if err != nil {
		return Result{}, fmt.Errorf("parse error: %w", err)
	}

	res, err := EvaluateSingle(e.rng, ast)
	if err != nil {
		return Result{}, err
	}
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
