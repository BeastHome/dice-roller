# AGENTS.md

Orientation for AI coding agents working in this repository. Tool-agnostic —
applies whether you're Claude, Cursor, Codex, Aider, or anything else.

## Project

`dice-roller` is a Go terminal application for evaluating tabletop-RPG dice
expressions. It ships one executable with two modes: a CLI for one-shot or
scripted evaluation and a tcell-based TUI for interactive rolling and session
history. The user is David Harris; the longer description, dice notation
reference, and platform-specific paths live in [README.md](README.md).

## Workspace layout

Single Go module (`github.com/showr/dice-roller`, Go 1.24). Layout:

- `main.go` — entrypoint: handles `--help` / `--version`, dispatches to TUI
  (no args) or CLI (with args).
- `cmd/cli/` — CLI runner (`RunCLI`, `PrintHelp`).
- `tui/` — tcell-based three-pane TUI (input / output / history), layout,
  event handling, rendering, history pane wiring.
- `internal/dice/` — dice engine. Two parse/eval paths exist:
  - **AST path**: recursive-descent parser (`parser_rd.go`) produces an
    AST (`ast_nodes.go`); `eval_ast.go` evaluates it directly, including
    arithmetic, grouping, and unary ops.
  - **Legacy path**: regex-based parser (`parser.go` + `parser_modifiers.go`)
    produces a flat `Expression`; `roller.go::EvaluateSingle` evaluates a
    single dice term with its modifiers. `bridge.go` lets the AST parser
    feed this evaluator for the dice-only case.
  Other contents: lexer/tokens, modifiers (keep/drop, explode, reroll,
  success), multi-roll, verbose breakdowns, platform-specific history
  paths, version.
- `internal/parse/` — input normalization shared by CLI and TUI (lowercase,
  spacing, grouping, flag extraction).
- `internal/history/` — persistent session history store (line-delimited
  JSON, file-backed).
- `internal/presentation/` — color scheme and CLI formatter.

Authoritative architectural intent: [CHARTER.md](CHARTER.md) and
[SEMANTIC_DECISIONS.md](SEMANTIC_DECISIONS.md). The roadmap is in
[IMPLEMENTATION_ROADMAP.md](IMPLEMENTATION_ROADMAP.md).

## Build, test, lint

| Action | Command |
|---|---|
| Build | `go build ./...` |
| Run all tests | `go test ./...` |
| Lint (must be clean) | `go vet ./...` |

Go 1.24+ (module declares `go 1.24.2`). Plain stable toolchain; no build
tags or generation steps. Single-binary build with `go build -o dice-roller .`
when producing a runnable artifact.

## Documentation

- [README.md](README.md) — user-facing usage: modes, flags, full dice
  notation reference, history paths.
- [CHARTER.md](CHARTER.md) — mission, product shape, scope, non-goals.
- [SEMANTIC_DECISIONS.md](SEMANTIC_DECISIONS.md) — numbered behavioral
  decisions (SD-001…); consult before changing parser/engine semantics,
  multi-roll behavior, verbose output, history format, or color rules.
- [IMPLEMENTATION_ROADMAP.md](IMPLEMENTATION_ROADMAP.md) — phase status
  and planned work.
- [CHANGELOG.md](CHANGELOG.md) — released and unreleased changes.

When a slice changes cross-cutting behaviour (dice syntax, output mode,
history schema, platform path selection), update the relevant doc(s) and
add a CHANGELOG entry in the same commit. New behavioral decisions get a
new `SD-NNN` entry in `SEMANTIC_DECISIONS.md`.

## Commit conventions

Match the existing `git log` style: short, sentence-style subjects, often
prefixed `docs:` / `fix:` when scoped, otherwise a plain imperative
summary. Real examples from the tree:

```
docs: add charter/semantic decisions/roadmap/changelog and README doc links
Added the build directory to the ignore list
Update drift check path to dev-tools/project-check.py
Corrected CLI coloring issue on the rolls line.
```

Subject under ~72 chars. Body only when the *why* isn't obvious from the
subject; one short paragraph is plenty.

**Do not add a `Co-Authored-By: <AI>` trailer** on commits in this repo.

## Slice discipline

Each commit should:

- leave the workspace compiling and `go test ./...` passing
- add a user-visible capability or semantic guarantee (or document
  reorganisation that has no behaviour change — call that out plainly in
  the commit message)
- update relevant docs when behaviour changes (README, SEMANTIC_DECISIONS,
  CHANGELOG)
- avoid mixing unrelated refactors with feature work

`go vet ./...` is treated as a gate — don't land a commit that regresses
it. The dice engine has both an AST path and a legacy regex/flat path
(SD-005); when touching either, run `go test ./internal/dice/...` to
exercise both. The pre-commit hook runs a separate drift check —
see [[project-drift-check-hook]] for the dependency.

## Things to avoid

- **PowerShell text rewrites on UTF-8 files.** Windows default codepage
  isn't UTF-8, so `Get-Content | … | Set-Content` corrupts smart quotes,
  em-dashes, CJK, accented characters. Use the `Edit` tool semantics (or
  `sed` via bash) for text edits, and Unicode escape syntax for non-ASCII
  string literals in Go source (`\u00XX` / `\U000XXXXX`).
- **Test names that mask intent.** Existing tests in `internal/dice/` and
  `tui/` favor descriptive names tied to the behavior under test
  (`Test_<Subject>_<Condition>` shape). Match the local style.
- **Touching the engine without consulting `SEMANTIC_DECISIONS.md`.**
  Parser fallback, normalization rules, verbose-vs-mechanics separation,
  history format, and color degradation are all documented decisions —
  changing them is a semantic change, not an implementation tweak.
- **Diverging CLI and TUI behavior.** Both paths share `internal/parse`
  normalization and the `dice.Engine` evaluator (SD-001, SD-002). Add
  shared logic to those packages rather than duplicating in `cmd/cli` or
  `tui`.
- **Sharing a `dice.Engine` across goroutines.** The engine's `*rand.Rand`
  is not safe for concurrent use; the dice package takes no locks (SD-010).
  Give each goroutine its own Engine, or wrap access in a mutex.
