// Package dice is the dice-expression engine: lexer, parsers (AST and
// legacy regex), evaluator, modifiers, multi-roll orchestration, and
// result formatting hooks.
//
// Thread-safety: an Engine holds a *rand.Rand which is NOT safe for
// concurrent use. Each goroutine that calls RollOnce / RollN must have
// its own Engine, or callers must serialize access with their own
// mutex. The dice engine itself takes no locks. See SD-010 in
// SEMANTIC_DECISIONS.md for the rationale.
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

// RollOnce parses and evaluates a single dice expression string.
//
// Evaluation strategy: try the AST path (parser_rd + eval_ast) first;
// only fall back to the legacy regex parser if the AST parser cannot
// recognize the expression. Errors from a successful AST parse (e.g.,
// division by zero) are returned directly — they are NOT masked by
// retrying the legacy path.
func (e *Engine) RollOnce(expr string) (Result, error) {
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

// RollN evaluates the same expression `count` times and returns the
// aggregate result. For count <= 1, it returns a MultiRollResult with
// a single entry — callers that want a bare Result for the single case
// should call RollOnce directly.
func (e *Engine) RollN(expr string, count int) (MultiRollResult, error) {
	if count < 1 {
		count = 1
	}
	return EvaluateMulti(e, expr, count)
}
