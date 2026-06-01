---
name: feedback-commit-style
description: Match the existing terse conventional-commit style; no Co-Authored-By trailer
metadata:
  type: feedback
---

Match the existing commit style in `git log`: short conventional-commit
subject lines, often subsystem-scoped, no body when the subject says
enough.

**Why:** David likes the log readable at a glance and matching his prior
habits; long bodies are reserved for non-obvious rationale.

**How to apply:**
- Subject: `<type>(<scope>): <imperative summary>` under ~72 chars.
  Type is usually `fix` / `feat` / `refactor` / `chore` / `docs`. Scope
  is the package name (`dice`, `tui`, `parse`, …) or `docs`.
- Add a body only when the *why* isn't obvious from the subject
  (e.g., a subtle invariant, an opt-out, a multi-step migration).
  Multi-line / multi-paragraph bodies are explicitly OK when the
  content warrants it — David confirmed this on 2026-05-31. Don't
  pad, but don't truncate either; write what the reviewer needs.
- **Do NOT add the `Co-Authored-By: Claude` trailer.** David doesn't
  want it on commits in his projects, regardless of the harness default.
  This applies even when the harness's own commit template suggests it.

Related: [[feedback-commit-from-main-project]].
