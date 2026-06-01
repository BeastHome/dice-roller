---
name: reference-memory-layout
description: Memory canonical location and the AGENTS.md / CLAUDE.md layered convention
metadata:
  type: reference
---

Memory canonical location: `.claude/memory/` in the **project repo**
(not in `~/.claude/projects/<project>/memory/`). Git-tracked, portable
across machines, reviewable in PRs.

Layered orientation files at the project root:

- **`AGENTS.md`** — tool-agnostic project orientation (build/test/lint
  commands, crate map, doc map, commit conventions, things to avoid).
  Any AI coding agent reads this. Patterned after the emerging
  `AGENTS.md` convention across AI coding tools.
- **`CLAUDE.md`** — Claude-specific entry. References AGENTS.md for the
  shared context and points here (`.claude/memory/MEMORY.md`) for the
  session memory. Claude Code auto-loads `CLAUDE.md` at session start.

Session start protocol:

1. Read `CLAUDE.md` (auto-loaded by harness on Claude sessions).
2. Read `AGENTS.md` for project context.
3. Read `.claude/memory/MEMORY.md` and skim file titles; load the ones
   relevant to the task.

The legacy userdir location `~/.claude/projects/<project>/memory/` is
no longer canonical. Any new memory writes go to the project repo.

**Why this layout:** memory files belong with the code they describe.
Git-tracking the memory means it survives userdir corruption, clones
identically to any machine, and is reviewable alongside the slice that
produced it. The trade-off is losing the harness's automatic
relevance-based memory surfacing — Claude reads the project memory by
following the `CLAUDE.md` pointer at session start, rather than by
passive recall during work.

Related: [[user-david]].
