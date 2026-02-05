# Design Alternatives — Delete Task Command

## Summary
Propose alternatives for a CLI command to delete tasks from the StrandYard database. The goal is to provide a safe and reliable way to remove tasks while maintaining the integrity of relationships (parents, blockers) and master lists.

## Context
- Tasks are stored as directories containing markdown files.
- Relationships are bidirectional (parent/child, blockers/blocks).
- Master lists (`root-tasks.md`, `free-tasks.md`) are generated from the task tree.

## Project Principles
- Keep CLI usage unambiguous and self-documenting.
- Preserve deterministic behavior and data integrity.
- Minimize manual edits required for maintenance.

## Alternatives

### Alternative A — Manual Removal + `repair`
- **Description**: Instruct users to use `rm -rf` on the task directory, then run `strand repair` to clean up broken references.
- **Pros**: No new command implementation needed.
- **Cons**: High risk of data corruption if `repair` is forgotten; poor UX.
- **Effort**: None.

### Alternative B — Basic `delete` command
- **Description**: Add a `strand delete <task-id>` command that removes the task directory and then automatically runs the equivalent of `strand repair`.
- **Pros**: Safe; easy to implement; ensures integrity.
- **Cons**: Might leave "orphan" children if not handled.
- **Effort**: Small.

### Alternative C — Smart `delete` with hierarchical options
- **Description**: `strand delete <task-id>` with flags:
  - `--recursive`: Delete all children as well.
  - `--reparent`: Move children to the deleted task's parent.
  - `--force`: Skip confirmation.
- **Pros**: Most flexible; handles hierarchical cleanup cleanly.
- **Cons**: More complex implementation; higher risk of accidental large-scale deletion.
- **Effort**: Medium.

## Recommendation
**Adopt Alternative C** but start with a cautious default:
- `strand delete <task-id>` asks for confirmation.
- Default behavior: refuse to delete if the task has children, unless `--recursive` or `--reparent` is specified.
- Automatically update all references (parent TODOs, blockers) and master lists.

## Review Requests
- Request review from: `master-reviewer`, `reviewer-usability`, `reviewer-reliability`.

## Decision
Decision: deferred to Owner.
