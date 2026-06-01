# Memory index for dice-roller

- [User: David](user-david.md) — developer; prefers small tested slices.
- [Feedback: commit from main project](feedback-commit-from-main-project.md) — land commits on `main`, not the worktree branch.
- [Feedback: commit style](feedback-commit-style.md) — short conventional-commit subjects; no `Co-Authored-By` trailer.
- [Feedback: no PowerShell UTF-8 rewrites](feedback-no-powershell-utf8-rewrites.md) — `Get-Content`/`Set-Content` corrupts non-ASCII; use `Edit` or `\u{}` escapes.
- [Reference: memory layout](reference-memory-layout.md) — canonical memory location is this directory; AGENTS.md / CLAUDE.md layered pattern.
- [Project: drift-check hook](project-drift-check-hook.md) — pre-commit hook depends on `S:/Dev/dev-tools/project-check.py` + a shared venv.
- [Feedback: engine API naming](feedback-engine-api-naming.md) — engine methods are `RollOnce` / `RollN`; never bare `Roll`.
