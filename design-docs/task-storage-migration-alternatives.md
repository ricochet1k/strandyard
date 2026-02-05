# Design Alternatives — Task Storage Migration

## Summary
Propose alternatives for migrating task storage between "global" (centralized config) and "local" (`.strand/` directory) modes. The goal is to make it easy for users to switch storage strategies without manual file moves or risk of data loss.

## Context
- StrandYard supports both local and global storage.
- Local storage is preferred for repository-specific tasks that should be tracked in Git.
- Global storage is useful for personal task management across multiple repositories.

## Project Principles
- Keep CLI usage unambiguous and self-documenting.
- Preserve deterministic behavior and data integrity.
- Minimize manual edits required for maintenance.

## Alternatives

### Alternative A — Dedicated `migrate` command
- **Description**: Add `strand migrate --to local|global`.
- **Pros**: Clear intent; specific to migration concerns.
- **Cons**: Adds another top-level command.
- **Effort**: Medium.

### Alternative B — Extend `init` command
- **Description**: If `strand init` is run on a project that is already linked to another storage mode, detect it and offer to migrate.
  - `strand init --storage local` (when currently global) -> "Project 'foo' is currently using global storage. Migrate to local? [y/N]"
- **Pros**: Reuses existing initialization entry point; discoverable.
- **Cons**: Overloads `init` with migration logic.
- **Effort**: Medium.

### Alternative C — Manual Migration Guide
- **Description**: Provide a `strand-migrate.sh` script or a detailed section in `CLI.md` with the necessary `mv` and `git` commands.
- **Pros**: Zero code changes to core CLI.
- **Cons**: Error-prone; poor user experience.
- **Effort**: Small.

## Recommendation
**Adopt Alternative B**. Extending the `init` command is the most natural workflow for users who want to change how their project is stored. It ensures that the initialization process remains the single source of truth for storage configuration.

### Migration Steps (Local to Global):
1. Copy `.strand/{tasks,roles,templates}` to global storage.
2. Update global project map (`config.json`).
3. Remove `.strand/` directory (optionally from Git).

### Migration Steps (Global to Local):
1. Copy files from global storage to `.strand/`.
2. Update global project map to record local path.
3. `git add .strand/`.

## Review Requests
- Request review from: `master-reviewer`, `reviewer-usability`, `reviewer-reliability`.

## Decision
Decision: deferred to Owner.
