# CLAUDE.md

Claude-specific entry point for this repo. Anything that applies to AI
agents generally lives in [AGENTS.md](AGENTS.md); read that first.

## Session start

Read [AGENTS.md](AGENTS.md) for project context and shared conventions.
Then read [.claude/memory/MEMORY.md](.claude/memory/MEMORY.md) — that
index lists the current memory files. Follow the links to whichever
files look relevant to the task at hand; they record current focus,
agreed plans, lessons learned, and the user's working preferences.

## Memory protocol

Memory lives in `.claude/memory/` as one fact per `.md` file with YAML
frontmatter. `MEMORY.md` is the index — one line per memory, never the
fact content itself.

```markdown
---
name: <short-kebab-case-slug>
description: <one-line summary — used to decide relevance during recall>
metadata:
  type: user | feedback | project | reference
---

<the fact; for feedback/project, follow with **Why:** and
**How to apply:** lines. Link related memories with [[their-name]].>
```

After writing a new memory file:

1. Add a one-line pointer in `MEMORY.md` (`- [Title](file.md) — hook`).
2. Check for an existing file that already covers the same fact and
   update or delete rather than duplicating.

Memory categories:

- `user` — who the user is (role, expertise, preferences).
- `feedback` — guidance the user has given on how to work, both
  corrections and confirmed approaches; include the why.
- `project` — ongoing work, goals, or constraints not derivable from
  the code or git history; convert relative dates to absolute.
- `reference` — pointers to external resources or distilled findings
  (URLs, dashboards, file:line lists).

Don't save what the repo already records (code structure, past fixes,
git history). If asked to remember one of those, ask what was
non-obvious about it and save that instead.

## Worktree / commit handling

The harness sometimes spawns me into a worktree under
`.claude/worktrees/<slug>/`. **Land commits on the main project checkout
(on `main`), not on the worktree's `claude/<slug>` branch.** Treat the
worktree as scratch space for exploration; do edits with absolute paths
under the main project root and run build / git commands against that
path.

## Test list / progress tracking

Use the `TaskCreate` / `TaskUpdate` tools when work is multi-step. For
shorter work, narrate in the response; not every commit needs a task.

## Userdir memory

There's a vestigial `~/.claude/projects/<project>/memory/` directory in
the harness that may have been the original location of memory files.
**The project-repo `.claude/memory/` is canonical.** If the userdir
location exists and differs, prefer the project-repo files; if asked to
save a memory, save it in `.claude/memory/` and commit it with the
slice it relates to.
