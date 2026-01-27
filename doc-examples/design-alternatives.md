# Design Alternatives — Task storage layout

## Summary
This document compares two approaches for storing the task database that the CLI and library will operate on. The choices affect discoverability, editability by humans/agents, and ease of implementing deterministic updates.

## Project Principles
- Keep data human-editable and version-control friendly.
- Prefer deterministic on-disk layout for stable diffs and reproducible CLI behaviour.
- Minimize required central infrastructure — aim for filesystem-first solutions.

## Alternatives

### Alternative A — Single-index file (monolithic)
- Description: Store all task metadata in a single repository-level file (JSON/YAML/TOML) that indexes tasks, relationships and metadata.
- Assumptions: Central index simplifies lookups and global queries.
- Pros:
  - Fast global queries and simple serialization.
  - Easier to validate referential integrity in one pass.
- Cons:
  - Large merge conflicts for active repositories.
  - Less friendly to ad-hoc human edits and simple agent-driven changes.
  - Harder to allow per-task directory attachments (notes, patches, artifacts).
- Risks:
  - Single point of contention for parallel edits; requires lock/merge tooling.

### Alternative B — Filesystem-per-task (chosen)
- Description: Each task is a directory containing a single `task.md` (and optional assets). Directory hierarchy represents parent/child lineage.
- Assumptions: Humans and agents will edit individual task files; CLI will maintain master lists and relations deterministically.
- Pros:
  - Optimized for git workflows and small, focused edits.
  - Easy for agents to create/update single tasks without large merges.
  - Supports attachments and sub-documents alongside a task.
- Cons:
  - Global queries require scanning files or a lightweight index updated by the CLI.
  - Slightly more work to implement efficient cross-task operations.
- Risks:
  - Need deterministic ordering and master-list maintenance to avoid drift; must be enforced by CLI commands.

## Recommendation
- Preferred alternative: Filesystem-per-task.
- Rationale: This project is explicitly designed for AI agent interaction and human editing; per-task directories minimize merge conflicts, encourage small updates, and support ancillary documents. The CLI will provide deterministic master lists (`tasks/root-tasks.md`, `tasks/free-tasks.md`) and scanning commands to make global queries reliable.

## Review Requests
- Request review from: `master-reviewer`, `reviewer-reliability`, `reviewer-security` as appropriate for implementation details (format, on-disk path choices, validation rules).

## Example TODOs
- [ ] (role: designer) Confirm on-disk layout for attachments and sub-documents.
- [ ] (role: architect) Draft CLI sync algorithm for `root-tasks.md` and `free-tasks.md`.
- [ ] (role: developer) Implement scan/sync commands and add unit tests for deterministic ordering.
