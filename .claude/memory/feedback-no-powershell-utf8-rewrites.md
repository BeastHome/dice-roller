---
name: feedback-no-powershell-utf8-rewrites
description: Don't use PowerShell Get-Content/Set-Content to rewrite files containing non-ASCII text
metadata:
  type: feedback
---

On Windows, `Get-Content -Raw` reads in the system codepage (often
CP-1252), not UTF-8. Piping that through a regex replace and back with
`Set-Content -Encoding utf8` corrupts every non-ASCII byte in the file
(em-dash, CJK, smart quotes all become mojibake).

**Why:** the rewrite looks successful — the file has a UTF-8 BOM and
the right encoding tag — but the bytes were already mis-decoded before
being re-encoded. The mojibake then panics non-ASCII tests at runtime.

**How to apply:**
- For multi-file regex rewrites use the `Edit` tool (UTF-8 safe) with
  `replace_all: true`, or use ripgrep + `sed` via WSL bash. Don't use
  PowerShell `Get-Content`/`Set-Content` for text replacement on source
  files.
- When tests with non-ASCII literals are needed, prefer the language's
  Unicode escape syntax (`\u00XX` / `\U000XXXXX` in Go) over pasting
  the literal character, so a future encoding mishap doesn't corrupt
  them silently.
- If PowerShell really is the only option, set
  `[Console]::InputEncoding = [Console]::OutputEncoding =
  [Text.UTF8Encoding]::new()` and use
  `[IO.File]::ReadAllText(path, [Text.UTF8Encoding]::new())` /
  `WriteAllText`. Still error-prone — just use `Edit`.
