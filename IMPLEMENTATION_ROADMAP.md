# Implementation Roadmap

This roadmap is retroactive and lightweight. It captures the state implied by the current v2.x codebase and likely near-term work.

## Phase 0 - Foundation (completed)
- [x] Project scaffold and Go module setup.
- [x] Core dice parser/evaluator packages.
- [x] Basic CLI entrypoint and expression execution.

## Phase 1 - Dice Language Coverage (completed)
- [x] Standard roll terms (`NdX`, arithmetic, grouping).
- [x] Keep/drop variants.
- [x] Exploding dice variants.
- [x] Reroll variants.
- [x] Success-threshold counting.
- [x] Multi-roll support (`rolls=N`, `--multi`).

## Phase 2 - UX and Presentation (completed)
- [x] Compact and verbose formatting modes.
- [x] ANSI color scheme with no-color behavior.
- [x] Help and version output in CLI.

## Phase 3 - TUI and Persistence (completed)
- [x] Three-pane TUI (input/output/history).
- [x] History persistence and reload on startup.
- [x] History browsing synced with output detail.

## Phase 4 - Test and Stability Hardening (completed in v2.1.0)
- [x] Parser and evaluator tests for major paths.
- [x] TUI interaction-focused tests for event and formatting behavior.
- [x] Characterization tests for normalize / validate / EvaluateSingle / lexer (v2.1.0).
- [x] History and presentation packages brought from 0% to >80% coverage (v2.1.0).
- [x] CLI scriptability tests against an `io.Writer` buffer (v2.1.0).
- [ ] Expand parser fuzz/negative-case coverage.
- [ ] Add regression suite for cross-mode output parity.

## Phase 5 - Next Improvements
- [x] Deterministic seeding surfaced as `--seed N` and `dice.NewEngineWithSeed` (v2.1.0).
- [x] Compatibility notes for embedded use — `dice` is a public package; see [EMBEDDING.md](EMBEDDING.md) (v2.1.0).
- [ ] Structured machine-readable output mode for CLI automation.
- [ ] Optional probability/statistics helpers for repeated rolls beyond `dice.Stats`.

## Phase 6 - Known cleanups (planned)
- [ ] Reroll semantics review: `r`, `ro`, and `ra` are documented as
      distinct but the engine currently treats `r` and `ro` as single-
      pass with only tracking differences. Decide whether `r` should
      become "reroll until > T" or whether the docs should continue to
      describe the single-pass reality.
- [ ] `FormatMultiExpression("", 0)` returns `"rolls=0"` rather than an
      empty or bare-expression form. Decide on a more useful behavior
      for the empty case.
