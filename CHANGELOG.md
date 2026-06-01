# Changelog

All notable changes to dice-roller should be documented in this file.

The format follows simple `Added`, `Changed`, and `Fixed` sections.

## 2.1.0

A post-audit hardening pass. The user-facing CLI behavior is largely
unchanged; the engine becomes a public package suitable for embedding,
and a handful of latent bugs are fixed.

### Added

- **`--seed N` CLI flag is now wired end-to-end.** The parser accepted
  it in 2.0 but it never reached the engine — output was always time-
  seeded. With this release, `dice-roller 4d6k3 --multi 5 --seed 42`
  produces the same five totals every run.
- **Public engine package** at `github.com/showr/dice-roller/dice`.
  The dice engine (`internal/dice` in 2.0) was promoted out of
  `internal/`, making it importable by external Go projects. See
  [EMBEDDING.md](EMBEDDING.md) for the stable API surface, thread-
  safety contract, determinism guarantees, and a worked example.
- **`Engine.RollOnce(expr) (Result, error)`** and **`Engine.RollN(expr, n) (MultiRollResult, error)`** —
  typed evaluation methods on the engine, replacing the older
  `Engine.Evaluate(expr, count) (interface{}, error)` that forced
  callers into a type switch.
- **`dice.NewEngineWithSeed(seed int64)`** — explicit deterministic
  constructor. Companion to `NewEngine()` (time-seeded) and
  `NewEngineWithOptions(EngineOptions{...})`.
- **`dice.Stats`, `dice.ComputeStats([]int)`, `dice.StatsFromResults([]Result)`** —
  shared min/max/avg/median/stddev helper that previously existed in
  three separate places in the codebase.
- **Typed history Store methods**: `AppendSingle(Result)` and
  `AppendMulti(MultiRollResult)` replace the untyped `Append(interface{})`.
- **`history.NewFileStoreInDir(dir)`** — base-directory override for
  tests and embedded consumers that want their own history root.
- **[AGENTS.md](AGENTS.md), [CLAUDE.md](CLAUDE.md), `.claude/memory/`** —
  layered orientation files for AI coding agents working in this
  repo, with project-tracked memory.
- **[EMBEDDING.md](EMBEDDING.md)** — target-agnostic embedded-use
  guide.
- **SD-010** (engine is single-goroutine), **SD-011** (arithmetic
  belongs to the AST evaluator), **SD-012** (`RollOnce`/`RollN`
  naming convention) added to `SEMANTIC_DECISIONS.md`.

### Changed

- **CLI errors now go to stderr.** In 2.0 they shared stdout. Scripts
  piping `dice-roller` output no longer have to filter error text.
- **CLI exit code is non-zero on parse or evaluation failure.** In 2.0
  every invocation exited 0 regardless. Useful for scripts and CI.
- **One bad expression no longer aborts a batch.** In 2.0,
  `dice-roller garbage 4d6 2d20` would stop at `garbage` and never
  roll the rest. Now the bad expression's error is reported on stderr
  and the remaining expressions still run; the exit code reflects
  the partial failure.
- **CLI help text** now documents `--seed` and notes that `--multi`
  takes precedence over the inline `rolls=N` suffix when both are
  present.
- **`Result.Rolls` field documentation** was corrected. The previous
  comment ("initial + post-reroll/explode values") was inaccurate.
  `Result.Rolls` holds the post-reroll, **pre-explode** state;
  explosions live in `Result.Exploded` and reach `Result.Total` via
  `Result.Kept`. The field itself is unchanged.
- **SD-005** reworded for accuracy. The two evaluation paths are
  *AST evaluator* and *legacy single-term flat evaluator*, not
  "parse-tree vs AST" as the original wording implied — both paths
  use ASTs.

### Fixed

- **`compound explode` (`NdX!!`) is bounded at 1000 expansions per
  roll.** Under adversarial or degenerate RNGs (e.g., `1d1!!`, which
  always rolls the max) the previous implementation would loop
  forever. The cap counts only the new dice produced by explosions —
  legitimate large rolls like `5000d100!!` are unaffected — and
  returns an error if exceeded.
- **AST evaluation errors are no longer masked.** When the AST parser
  succeeded but evaluation failed (e.g., division by zero), the
  engine previously discarded the eval error and retried with the
  legacy parser, which then reported `"parse error: invalid dice
  expression"`. Now `1d1 / (1d1 - 1)` correctly reports
  `division by zero`.
- **Additive nested-loop bug in the legacy evaluator** removed.
  `roller.go::EvaluateSingle` had been applying `+N` once per kept
  die rather than once per expression, yielding `sum(kept) + N*len(kept)`
  instead of `sum(kept) + N`. The bug was latent because
  `Engine.RollOnce` routes `+N` through the AST evaluator before
  the legacy path ever sees it, but a direct `EvaluateSingle` caller
  with `ModAdditive` would have hit it. The legacy path now ignores
  `ModAdditive` entirely — arithmetic is the AST evaluator's job.

### Internal

- **Test coverage** grew substantially: `internal/dice` 39.6% → 66.9%,
  `internal/history` 0% → 83.9%, `internal/presentation` 0% → 87.7%,
  `cmd/cli` 0% → 93.8%.
- **`internal/dice/util.go` deleted.** `min`/`max` shadowed Go 1.21+
  builtins; `sum`/`countIf` were unused.
- **`EngineOptions.Verbose` and `EngineOptions.MaxDepth`** removed
  (TODO fields with no callers).
- **`cmd/cli.RunCLI` refactored** to accept `io.Writer`s and return
  an `int` exit code, enabling unit tests against an `io.Writer`
  buffer.

## 2.0.0

### Added
- Unified Go dice engine with support for common tabletop dice notation and modifiers.
- CLI mode for one-shot and multi-expression evaluation.
- TUI mode with input, output, and history panes.
- Verbose output mode with detailed roll breakdowns.
- Session history persistence with startup reload.
- Platform-aware history paths and color presentation handling.

### Changed
- Parser and evaluation workflow consolidated so CLI and TUI share core behavior.
