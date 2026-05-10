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

## Phase 4 - Test and Stability Hardening (in progress)
- [x] Parser and evaluator tests for major paths.
- [x] TUI interaction-focused tests for event and formatting behavior.
- [ ] Expand parser fuzz/negative-case coverage.
- [ ] Add regression suite for cross-mode output parity.

## Phase 5 - Next Improvements (planned)
- [ ] Optional deterministic seeding surfaced consistently in user-facing docs/flags.
- [ ] Structured machine-readable output mode for CLI automation.
- [ ] Optional probability/statistics helpers for repeated rolls.
- [ ] More explicit compatibility notes for imported/embedded use in other projects.
