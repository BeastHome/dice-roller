---
name: feedback-engine-api-naming
description: Use suffixed engine method names (RollOnce, RollN) — never bare Roll
metadata:
  type: feedback
---

The public dice-engine evaluation methods are `RollOnce(expr) (Result, error)`
and `RollN(expr, count) (MultiRollResult, error)`. Don't introduce or rename
back to a bare `Roll` form.

**Why:** The dice engine is being imported by David's roguelike (and
likely future projects). Embedded consumers asked for unambiguous naming
because:
- Bare `Roll(expr)` reads like "do a roll" but is overloaded with
  batch/multi connotations in common parlance.
- The suffixed pair `RollOnce` / `RollN` makes the intended cardinality
  explicit at every call site.
- The original audit doc proposed `RollOnce` / `RollN`; an earlier slice
  chose `Roll` / `RollMany`, and the embedding consumer flagged the
  drift on first import. SD-012 in `SEMANTIC_DECISIONS.md` records the
  reversal.

**How to apply:**
- When adding new evaluation entry points (e.g., a future
  `RollManyWithSeed` or `RollOnceContext`), follow the same suffix
  convention: lead with the verb, then the cardinality / variant tag.
- When proposing API changes that touch these methods, surface the
  rename impact for embedded consumers, not just the in-tree CLI/TUI.
- Don't quietly reintroduce `Roll` as an alias "for convenience" — it
  reopens the ambiguity SD-012 closed.

Related: [[reference-memory-layout]], SD-010 (thread-safety), SD-012.
