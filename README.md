# dice-roller

A terminal-based dice roller for tabletop RPGs. Supports a full dice notation syntax with modifiers, a text UI (TUI), and a CLI mode for scripted use.

**Version:** 2.1.0

---

## Modes

| Mode | How to invoke |
|------|--------------|
| TUI  | Run with no arguments: `dice-roller` |
| CLI  | Pass one or more expressions: `dice-roller 4d6k3` |

---

## CLI Usage

```
dice-roller <expression> [--verbose] [--multi N] [--no-color] [--seed N]
dice-roller "<expression> rolls=N" [--verbose] [--no-color] [--seed N]
```

### Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Show full roll breakdown (rerolls, explosions, kept/dropped) |
| `--multi N` | Repeat the expression N times |
| `--seed N` | Seed the RNG for deterministic output across the invocation |
| `--no-color` | Disable ANSI color output |
| `--help` | Show help text |
| `--version` | Show version |

Colors are auto-disabled when output is piped or redirected.

Multiple expressions can be passed and are evaluated in sequence:

```
dice-roller 2d20kh1 3d8! 5d10>=8 --no-color
```

If a single expression in a batch fails to parse or evaluate, the
error goes to stderr and the remaining expressions still run. The
process exit code is non-zero whenever any expression failed, so
scripts can branch on it.

### Multi-roll precedence

`--multi N` and the inline `rolls=N` suffix express the same idea
(repeat the expression N times). When both are present, the `--multi`
flag wins — `rolls=N` is only applied if `--multi` is at its default
of 1. Pick whichever fits your invocation; don't mix them.

### Deterministic output

With `--seed N`, the entire invocation — including all expressions and
all repetitions under `--multi` — is reproducible:

```
dice-roller 4d6k3 --multi 5 --seed 42 --no-color
```

…produces the same five totals every time you run it. The seed is a
signed integer; any non-zero value works.

---

## Dice Notation

### Core

| Syntax | Meaning |
|--------|---------|
| `NdX` | Roll N dice of size X |
| `dX` | Roll 1 die of size X |
| `NdX+Y` | Add or subtract a constant |
| `NdX+MdY` | Combine multiple roll terms |
| `(expr)*N` | Group arithmetic with parentheses |

### Keep / Drop

| Syntax | Meaning |
|--------|---------|
| `NdXkY` | Keep highest Y |
| `NdXklY` | Keep lowest Y |
| `NdXdhY` | Drop highest Y |
| `NdXdlY` | Drop lowest Y |

### Exploding Dice

| Syntax | Meaning |
|--------|---------|
| `NdX!` | Explode on max value (single pass — exploded dice don't re-explode) |
| `NdX!T` | Explode on >= T |
| `NdX!>T` | Same as `NdX!T`; alternate form accepted by the parser |
| `NdX!!` | Compound explode on max value (recursive — exploded dice can chain) |

Compound explode is bounded at 1000 expansions per roll to guard
against degenerate RNGs; an error is returned if the cap is hit.

### Rerolls

| Syntax | Meaning |
|--------|---------|
| `NdXrT` | Reroll values <= T once; the new value replaces the original |
| `NdXroT` | Reroll once; same single-pass semantics as `r` (the surfaced reroll values are tracked differently for display purposes) |
| `NdXraT` | Reroll values <= T once and *add* — the new die is appended rather than replacing the original |

All reroll variants are single-pass — a qualifying die is rerolled
exactly once, even if the replacement value would also qualify.

### Success Counting

| Syntax | Meaning |
|--------|---------|
| `NdX>=T` | Count successes >= T |
| `NdX<=T` | Count successes <= T |
| `NdX>T` | Count successes > T |
| `NdX<T` | Count successes < T |

### Multi-Roll

| Syntax | Meaning |
|--------|---------|
| `rolls=N` | Repeat the entire expression N times |
| `--multi N` | Same as `rolls=N` |

---

## Examples

```bash
# Standard D&D ability score roll (4d6 drop lowest)
dice-roller 4d6k3

# Grouped arithmetic
dice-roller "(2d6 + 1d4) * 2"

# Success pool with reroll-once, 10 repetitions
dice-roller "5d10ro1>=8 rolls=10"

# Advantage roll with verbose breakdown
dice-roller 4d6k3 --multi 10 --verbose

# Multiple expressions in one call
dice-roller 2d20kh1 3d8! 5d10>=8 --no-color
```

---

## TUI

Launch with no arguments. Three panes:

- **Input** — type a dice expression and press Enter
- **Output** — result of the current roll
- **History** — all rolls from past sessions, loaded on startup

Session history files are stored at:

- **Windows:** `%USERPROFILE%\Documents\dice-roller\history\`
- **Linux/macOS:** `~/.local/share/dice-roller/`

---

## Building

```bash
go build -o dice-roller .
```

Requires Go 1.24+.

For the test and lint gates contributors run before committing, see
[AGENTS.md § Build, test, lint](AGENTS.md#build-test-lint).

---

## Project Docs

- [CHARTER.md](CHARTER.md)
- [SEMANTIC_DECISIONS.md](SEMANTIC_DECISIONS.md)
- [IMPLEMENTATION_ROADMAP.md](IMPLEMENTATION_ROADMAP.md)
- [CHANGELOG.md](CHANGELOG.md)
- [EMBEDDING.md](EMBEDDING.md) — for external Go projects importing the dice engine.
