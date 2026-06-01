# Embedding the dice engine

This document is for external consumers who want to import the dice
engine — `github.com/showr/dice-roller/dice` — into another Go project
(games, simulators, analytics, scripts, …). It describes the stable
API surface, the contracts it makes (determinism, thread-safety, error
handling), and the surface area that's exported but not part of the
stable contract.

If you're working *inside* this repo (CLI, TUI, future tooling),
[AGENTS.md](AGENTS.md) is the right starting point; this document is
about what we promise to people who don't have that context.

---

## Quick start

```go
package main

import (
    "fmt"

    "github.com/showr/dice-roller/dice"
)

func main() {
    engine := dice.NewEngine() // time-seeded

    res, err := engine.RollOnce("4d6k3")
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    fmt.Printf("expression: %s\n", res.Expression)
    fmt.Printf("rolls:      %v\n", res.Rolls)   // post-reroll, pre-explode
    fmt.Printf("kept:       %v\n", res.Kept)   // final dice contributing to Total
    fmt.Printf("total:      %d\n", res.Total)
}
```

For reproducible output (replay, tests, deterministic simulations):

```go
engine := dice.NewEngineWithSeed(42)
```

For N independent evaluations of the same expression:

```go
mr, err := engine.RollN("4d6k3", 100)
// mr.Rolls is a []Result of length 100
stats := dice.StatsFromResults(mr.Rolls)
fmt.Printf("avg=%.2f  median=%.2f  stddev=%.2f\n",
    stats.Average, stats.Median, stats.StdDev)
```

---

## Public API surface (stable)

### Engine construction

| Function | Purpose |
|---|---|
| `NewEngine() *Engine` | Time-seeded RNG. Use for interactive / non-reproducible work. |
| `NewEngineWithSeed(seed int64) *Engine` | Explicit deterministic seed. Preferred for tests, replay, simulations. |
| `NewEngineWithOptions(opt EngineOptions) *Engine` | Options bag. Currently only `Seed int64` — note `Seed: 0` means "use time-based default"; pass a non-zero seed or use `NewEngineWithSeed(0)` to seed with literal zero. |

### Evaluation

| Method | Returns | Purpose |
|---|---|---|
| `(e *Engine) RollOnce(expr string)` | `(Result, error)` | Parse + evaluate a single expression. Arithmetic, grouping, and all dice modifiers supported. |
| `(e *Engine) RollN(expr string, count int)` | `(MultiRollResult, error)` | N independent evaluations of the same expression. `count <= 1` is normalized to 1. |

### Result types

```go
type Result struct {
    Expression string // the original expression text
    Rolls      []int  // post-reroll values (PRE-explode). See "Field semantics" below.
    Rerolls    []int  // values recorded by the reroll pass
    Exploded   []int  // additional values produced by explosions
    Kept       []int  // post-explode + post-keep/drop slice — the final dice
    Dropped    []int  // values removed by keep/drop
    Total      int    // sum of Kept
    Successes  int    // count of dice satisfying >=T, >T, <=T, <T (zero if no success modifier)
    Verbose    string // human-readable breakdown, auto-attached after RollOnce
}

type MultiRollResult struct {
    Expression string   `json:"expression"`           // includes the "rolls=N" suffix
    Rolls      []Result `json:"rolls"`
    Summary    string   `json:"summary,omitempty"`    // optional caller-set summary line
}
```

**Field semantics** (important — these are NOT what casual reading
might suggest):

- `Rolls` is the **post-reroll, pre-explode** state. For most modifiers
  its length equals the original dice count. For *reroll-add* it also
  includes the appended dice. Explosions never appear in `Rolls` —
  they live in `Exploded` and reach `Total` via `Kept`.
- `Kept` is what `Total` is summed from. It's the slice you want for
  most "what did the player actually roll?" questions.
- `Total` is the sum of `Kept`. Surrounding arithmetic (`2d6 + 3`,
  `(d20+5)*2`) is handled by the AST evaluator and folded into the
  outermost `Result.Total` — sub-expression `Result`s contain only
  their own subtree's contribution.

### Stats helper

```go
type Stats struct {
    Count   int
    Sum     int
    Min     int
    Max     int
    Average float64
    Median  float64
    StdDev  float64 // POPULATION stddev (divide by N), not sample (N-1)
}

func ComputeStats(totals []int) Stats
func StatsFromResults(rolls []Result) Stats
```

`StdDev` uses the population formula because a dice batch *is* the full
population being summarized, not a sample. If you need sample variance,
recompute from `Stats.Sum` and the raw totals.

### Other exported helpers

| Symbol | Purpose |
|---|---|
| `FormatMultiExpression(expr string, count int) string` | Ensures `rolls=N` suffix appears exactly once. Returns `"rolls=N"` alone when `expr` is empty. |
| `AttachVerbose(r *Result)` | Fills `Result.Verbose` with a human-readable breakdown. `RollOnce` calls this automatically; you typically don't need to. |
| `HistoryDir() string` | OS-aware default path for session history (Documents/dice-roller/history on Windows; XDG-style on Unix). Useful if you're persisting history alongside our session format. |
| `HelpText []string`, `HelpLines() []string` | The CLI/TUI help content — list of pre-formatted lines. Useful if you're building a help screen and want to defer to the engine's syntax reference. |
| `Version` constant | The engine's version string. |

---

## Expression grammar

The full syntax (NdX, keep/drop, exploding, rerolls, success counting,
arithmetic, grouping, multi-roll) is documented in
[README.md § Dice Notation](README.md#dice-notation). The engine
accepts the same input as the CLI.

A few embedded-consumer notes that aren't obvious from the user-facing
table:

- Whitespace around binary operators is permitted; the lexer trims it.
- Input is normalized to lowercase before parsing, so `4D6K3` and
  `4d6k3` are equivalent.
- Arithmetic outside dice terms (`+ - * /` plus parentheses) is
  fully supported and handled by the AST evaluator.
- The `rolls=N` inline suffix is recognized by the CLI's argument
  parser, not by `RollOnce` itself — embedded consumers should pass
  N to `RollN(expr, N)` directly rather than embedding it in the
  expression string.

---

## Determinism and seeding

- `NewEngineWithSeed(seed)` produces a deterministic sequence: the
  same seed + the same sequence of `RollOnce` / `RollN` calls
  produces byte-identical `Result.Rolls`, `Kept`, `Total`, etc.
- The RNG is `math/rand`'s default source. The Go team treats its
  output as stable across minor releases; if you need cross-language
  reproducibility or cryptographic guarantees, wrap your own RNG and
  evaluate via `EvaluateSingle(rng *rand.Rand, expr Expression)`
  directly (see "Lower-level API" below).
- Two engines with the same seed are independent — interleaving calls
  on engine A vs engine B does **not** consume the same RNG sequence
  twice.

---

## Thread-safety

**The Engine is NOT safe for concurrent use.** It holds a
`*math/rand.Rand` which itself isn't thread-safe, and the dice
package takes no locks (see [SD-010](SEMANTIC_DECISIONS.md#sd-010-engine-is-single-goroutine-callers-serialize-concurrency)).

For concurrent callers, pick whichever fits your access pattern:

```go
// One engine per goroutine — no contention, each has its own seed.
go func() {
    e := dice.NewEngineWithSeed(seedA)
    // ... use e
}()

// Or share with explicit synchronization:
var mu sync.Mutex
engine := dice.NewEngineWithSeed(42)
roll := func(expr string) (dice.Result, error) {
    mu.Lock()
    defer mu.Unlock()
    return engine.RollOnce(expr)
}
```

The second pattern is fine for low-frequency calls (turn-based games,
periodic checks). For high-throughput parallel rolling, one engine
per goroutine is faster.

---

## Error handling

`RollOnce` and `RollN` return errors for:

- **Parse failures** — malformed expressions, unknown characters.
  Wrapped as `parse error: <detail>`.
- **Validation failures** — expressions that parse but are semantically
  invalid (e.g., `4d6k5` keeps more dice than exist; reroll/explode
  thresholds <= 0).
- **Evaluation failures** — division by zero, compound explode chains
  exceeding `1000` expansions (the safety ceiling against adversarial
  RNG; see [SD-011](SEMANTIC_DECISIONS.md)).

A returned error means the `Result` is zero-valued — don't read its
fields. Errors are plain `error` values; you can `errors.Is` /
`errors.As` against the wrapped causes when relevant.

---

## Stability and versioning

What we treat as **stable** (won't change without a major version bump
and a CHANGELOG entry):

- `Engine`, `NewEngine`, `NewEngineWithSeed`, `NewEngineWithOptions`,
  `EngineOptions`.
- `Engine.RollOnce`, `Engine.RollN`.
- `Result`, `MultiRollResult` field set and types.
- `Stats`, `ComputeStats`, `StatsFromResults`.
- `FormatMultiExpression`, `HistoryDir`, `Version`.
- The expression grammar accepted by the parser.
- The determinism guarantee (same seed → same rolls).
- The thread-safety contract (single-goroutine, [SD-010](SEMANTIC_DECISIONS.md#sd-010-engine-is-single-goroutine-callers-serialize-concurrency)).

What is **exported but not stable** — feel free to use it, but expect
it may change without ceremony:

- `Lex`, `Token`, `TokenType`, the AST node types (`DiceNode`,
  `BinaryNode`, etc.), `ParseExpression`, `ParseTreeExpression`,
  `EvaluateSingle`, `EvaluateMulti`, `EvaluateParseTree`,
  `BuildExpressionFromTree`, `Expression`, `Modifier`, `ModifierKind`,
  the `Mod*` constants.

These are exported because internal cross-file callers need them;
they're not part of the contract we're making with consumers.

---

## Out of scope

Things the dice engine does **not** do (and isn't going to grow into
in any near term):

- **Networked / shared sessions.** No client-server protocol.
- **Character sheets, campaign management, or rule-system specifics.**
  The engine knows about dice; it doesn't know that 18 + Str-mod
  passes a DC 15 check.
- **Statistical analysis dashboards.** `Stats` gives you the
  per-batch summary; anything richer (distribution charts, Monte
  Carlo dashboards) is your application's job.
- **GUI rendering.** The TUI in this repo uses tcell, but the engine
  itself has no display concerns. Embedded consumers ship their own
  rendering.

---

## Lower-level API

If you need control beyond `Engine.RollOnce` / `RollN` — for example,
to plug in your own `*rand.Rand`, or to walk the parse tree before
evaluation — the building blocks are exported but **not stable**:

```go
// Bring your own RNG.
rng := rand.New(rand.NewSource(seed))
expr, err := dice.ParseExpression("4d6k3")
if err != nil { /* ... */ }
res, err := dice.EvaluateSingle(rng, expr)

// Or work with the AST directly.
tree, err := dice.ParseTreeExpression("(2d6 + 3) * 2")
if err != nil { /* ... */ }
res, err := dice.EvaluateParseTree(rng, tree)
```

If you find yourself reaching for these regularly, file an issue —
that's a signal we should promote some of this surface to the stable
tier.

---

## Worked example: an ability check

A generic pattern most consumers will recognize: roll a check, compare
against a difficulty number, surface success/failure plus the raw
result for display.

```go
type CheckOutcome struct {
    Roll     int
    Modifier int
    Total    int
    DC       int
    Success  bool
    Detail   string // human-readable, e.g. for a log pane
}

func AbilityCheck(engine *dice.Engine, modifier, dc int) (CheckOutcome, error) {
    res, err := engine.RollOnce("1d20")
    if err != nil {
        return CheckOutcome{}, err
    }
    total := res.Total + modifier
    return CheckOutcome{
        Roll:     res.Total,
        Modifier: modifier,
        Total:    total,
        DC:       dc,
        Success:  total >= dc,
        Detail:   res.Verbose,
    }, nil
}
```

Reproducibility for tests:

```go
func TestAbilityCheck_Deterministic(t *testing.T) {
    engine := dice.NewEngineWithSeed(7)
    out, err := AbilityCheck(engine, 3, 15)
    if err != nil { t.Fatal(err) }
    // Same seed → same Roll value every time.
    if out.Roll < 1 || out.Roll > 20 {
        t.Fatalf("d20 out of range: %d", out.Roll)
    }
}
```

For "advantage" / "disadvantage" patterns, `RollOnce("2d20kh1")` /
`RollOnce("2d20kl1")` express them directly without leaving the
engine API.

---

## Feedback and contact

If something in the stable surface bites you, the engine doesn't do
what its docs say, or you need a contract guarantee that isn't here
yet — open an issue. The engine's design intentionally favors a small
stable core; growing it is fine, but we'd rather grow it deliberately.
