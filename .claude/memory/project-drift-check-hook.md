---
name: project-drift-check-hook
description: Pre-commit hook depends on an external Python drift-check script outside the repo
metadata:
  type: project
---

The repo's `.githooks/pre-commit` invokes a Python drift checker located
at `S:/Dev/dev-tools/project-check.py` (run inside the venv at
`S:/DevTools/venv-drift/`). The script lives **outside this repo** —
it's a shared dev-tool across David's Go projects.

**Why:** keeping the drift checker outside the repo lets it evolve once
and apply to every project that installs the hook, instead of forking
the script per repo.

**How to apply:**
- A commit can fail with "Drift detected — commit blocked." or with a
  Python error like `No such file or directory` if the path is missing.
- Don't try to fix drift by editing the hook or bypassing it
  (`--no-verify` is off-limits per [[feedback-commit-style]]); fix the
  source-of-truth that the script is comparing against.
- On a fresh clone or a different machine, both `S:/Dev/dev-tools/` and
  the venv at `S:/DevTools/venv-drift/` need to exist for the hook to
  succeed. If they don't, surface that to David rather than disabling
  the hook silently.
- Last known successful run: 2026-05-31, commit `6a5fc3d` (the AGENTS/
  CLAUDE/memory bootstrap). Hook printed `No drift detected` then
  `No drift — commit allowed.`
