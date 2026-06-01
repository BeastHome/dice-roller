---
name: feedback-commit-from-main-project
description: Land commits in the main project checkout, not inside the worktree branch
metadata:
  type: feedback
---

When the harness spawns me into a worktree under `.claude/worktrees/<slug>`,
David wants the actual edits and commits to land in the **main project
checkout** on `main`, not on the worktree's `claude/<slug>` branch.

**Why:** he doesn't want a sprawl of throwaway worktree branches; the
main branch is the one he reviews and pushes from.

**How to apply:**
- Use the worktree for read-only scanning / exploration if needed, but
  make edits with `Edit`/`Write` against absolute paths under the main
  project root.
- Run the project's build tool (`cargo`, `go`, `npm`, …) and `git`
  commands against the main project path.
- Do NOT switch branches in the worktree to try to commit there.

Related: [[feedback-commit-style]].
