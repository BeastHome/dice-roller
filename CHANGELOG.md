# Changelog

All notable changes to dice-roller should be documented in this file.

The format follows simple `Added`, `Changed`, and `Fixed` sections.

## Unreleased

### Added
- Retroactive architecture documentation set:
  - `CHARTER.md`
  - `SEMANTIC_DECISIONS.md`
  - `IMPLEMENTATION_ROADMAP.md`

## 2.0.0

### Added
- Unified Go dice engine with support for common tabletop dice notation and modifiers.
- CLI mode for one-shot and multi-expression evaluation.
- TUI mode with input, output, and history panes.
- Verbose output mode with detailed roll breakdowns.
- Session history persistence with startup reload.
- Platform-aware history paths and color presentation handling.

### Changed
- Parser and evaluation workflow consolidated so CLI and TUI share core behavior.
