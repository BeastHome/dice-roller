# Semantic Decisions

This document records behavioral decisions for dice-roller so future changes stay deliberate.

## SD-001: Shared engine across CLI and TUI
- Status: accepted
- Decision: Both modes call the same `dice.Engine` evaluation API.
- Rationale: Keeps roll semantics consistent regardless of interface.

## SD-002: One parser path for command-line and typed input
- Status: accepted
- Decision: CLI args and TUI line input both flow through `internal/parse` normalization.
- Rationale: Avoids feature drift and duplicated edge-case handling.

## SD-003: Expression normalization is case-insensitive
- Status: accepted
- Decision: Input is normalized to lowercase and spacing/operator forms are canonicalized before evaluation.
- Rationale: Reduces user friction and parser ambiguity from stylistic input differences.

## SD-004: Inline and flag-based multi-roll are equivalent
- Status: accepted
- Decision: `rolls=N` and `--multi N` represent the same behavior. Explicit flag parsing remains first-class.
- Rationale: Supports both shell workflows and familiar inline RPG notation.

## SD-005: Parser/evaluator fallback path is allowed
- Status: accepted
- Decision: Engine may attempt parse-tree evaluation first, then fall back to AST parsing/evaluation.
- Rationale: Preserves compatibility while parser internals evolve.

## SD-006: Verbose output is presentation-only
- Status: accepted
- Decision: Verbose mode changes rendering detail, not roll mechanics or totals.
- Rationale: Users should trust that display mode does not alter outcomes.

## SD-007: Session history is append-oriented JSON
- Status: accepted
- Decision: TUI history storage appends one JSON record per line to per-session files.
- Rationale: Simple recovery, easy inspection, and resilient partial writes.

## SD-008: Platform-aware history path selection
- Status: accepted
- Decision: History storage path is OS-specific (`Documents` on Windows, data dir path on Unix-like systems).
- Rationale: Respect platform conventions and avoid hardcoded single-OS assumptions.

## SD-009: Color output should degrade safely
- Status: accepted
- Decision: CLI color is disabled when output is redirected unless explicitly forced by mode behavior.
- Rationale: Keep pipelines and logs clean.

## SD-010: Engine is single-goroutine; callers serialize concurrency
- Status: accepted
- Decision: `dice.Engine` holds a `*rand.Rand` and takes no internal locks. Each goroutine that calls `Roll` / `RollMany` must have its own Engine, or callers must serialize access with their own mutex.
- Rationale: Most consumers (CLI, TUI, turn-based games) are single-goroutine and shouldn't pay the cost of lock contention. Concurrent callers can wrap an Engine with a Mutex or construct one Engine per goroutine — both are cheap. If/when a concurrent-by-default constructor is needed it can be added as `NewSyncEngine` without disturbing the existing API.

## SD-011: Arithmetic is the AST evaluator's responsibility
- Status: accepted
- Decision: Operators `+ - * /` and grouping are evaluated by `eval_ast.go::combineBinaryResults`, never by `roller.go::EvaluateSingle`. The latter ignores any `ModAdditive` entries it receives and returns `Total = sum(Kept)` for a single dice term in isolation.
- Rationale: Splits concerns cleanly between the AST and legacy paths. EvaluateSingle previously had a buggy nested-loop additive application that was masked only because the AST path always handled `+N` first; removing the duplicate eliminates the bug class entirely and makes the legacy path's job — one dice term + its modifiers — unambiguous.
