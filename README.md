# dice-roller

A terminal-based dice roller for tabletop RPGs. Supports a full dice notation syntax with modifiers, a text UI (TUI), and a CLI mode for scripted use.

**Version:** 2.0.0

---

## Modes

| Mode | How to invoke |
|------|--------------|
| TUI  | Run with no arguments: `dice-roller` |
| CLI  | Pass one or more expressions: `dice-roller 4d6k3` |

---

## CLI Usage

```
dice-roller <expression> [--verbose] [--multi N] [--no-color]
dice-roller "<expression> rolls=N" [--verbose] [--no-color]
```

### Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Show full roll breakdown (rerolls, explosions, kept/dropped) |
| `--multi N` | Repeat the expression N times |
| `--no-color` | Disable ANSI color output |
| `--help` | Show help text |
| `--version` | Show version |

Colors are auto-disabled when output is piped or redirected.

Multiple expressions can be passed and are evaluated in sequence:

```
dice-roller 2d20kh1 3d8! 5d10>=8 --no-color
```

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
| `NdX!` | Explode on max value |
| `NdX!T` | Explode on >= T |
| `NdX!!` | Compound explode on max value |

### Rerolls

| Syntax | Meaning |
|--------|---------|
| `NdXrT` | Reroll values <= T (replace) |
| `NdXroT` | Reroll once |
| `NdXraT` | Reroll and add (accumulate) |

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
