# dice-roller Charter

## Mission
Build a dependable, scriptable dice expression evaluator with both terminal-interactive and CLI workflows, suitable for tabletop play and automation.

## Product Shape
- One executable, two user modes:
  - TUI mode for interactive rolling and session review.
  - CLI mode for one-shot and multi-expression evaluation.
- Expression language supports arithmetic plus common TTRPG modifiers.
- Output supports compact summaries and verbose roll breakdowns.

## Core Principles
- Deterministic parsing and readable errors over clever but fragile syntax shortcuts.
- Shared behavior between CLI and TUI where practical (same parser and engine).
- Keep storage and terminal behavior platform-aware from the start.
- Preserve a clean engine boundary so other projects can embed the dice evaluator.

## Current Scope (v2.x)
- Dice notation: NdX, dX, arithmetic, grouping.
- Modifiers: keep/drop, exploding, rerolls, success counting.
- Multi-roll: inline `rolls=N` and CLI `--multi N`.
- Deterministic seeding via `--seed N` and `dice.NewEngineWithSeed` (v2.1.0+).
- TUI panes: input, output, history.
- Session history persisted as line-delimited JSON entries.
- Engine importable by external Go projects as `github.com/showr/dice-roller/dice` (v2.1.0+); embedded-use contract is documented in [EMBEDDING.md](EMBEDDING.md).

## Non-Goals (for now)
- Networked/shared sessions.
- GUI beyond terminal UI.
- Rule-system specific character sheets or campaign management.
- Probabilistic analytics tooling (distribution charts, Monte Carlo dashboards).

## Platform and Tooling
- Language: Go.
- Terminal UI: tcell v2.
- Target environments: Windows, Linux, and other Go-supported platforms where terminal features are compatible.

## Evolution Policy
The parser and evaluator are considered project infrastructure: improvements should favor backward compatibility for common expressions, and breaking syntax changes should be called out explicitly in the changelog.
