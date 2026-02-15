# Design Alternatives - Reapply Templates to Existing Tasks

## Summary
This document compares approaches for applying updated templates to already-existing tasks without losing task-specific edits. The goal is to support deterministic updates, safe merges, and clear CLI behavior.

## Context
- Templates in `.strand/templates/` evolve over time (new sections, wording, acceptance criteria).
- Existing task files often include user-authored frontmatter overrides and task body edits.
- Current workflow supports template use on creation (`strand add`) but not reapplication.

## Project Principles
- Preserve user-authored task data by default.
- Keep behavior deterministic so repeated runs produce stable output.
- Prefer explicit CLI operations over implicit mutation.
- Keep repair/list workflows compatible with TaskDB invariants.

## Alternatives

### Alternative A - Extend `strand edit` with `--template`
- Description: Add `strand edit <task-id> --template <type>` that reapplies the named template to the target task.
- Assumptions: Users already reach for `edit` when changing an existing task.
- Pros:
  - Minimal command surface area.
  - Keeps "modify task" behavior centralized in one command.
- Cons:
  - `edit` becomes overloaded (metadata edits + structural merge logic).
  - Harder to expose specialized merge flags cleanly.
  - Higher risk of surprising behavior for users expecting simple edits.
- Risks:
  - Confusing UX around precedence between edit flags and template fields.
- Rough effort estimate: Medium.

### Alternative B - New `strand reapply` command
- Description: Add `strand reapply <task-id> --template <type>` dedicated to template reconciliation.
- Assumptions: Reapplication is conceptually distinct from metadata edits.
- Pros:
  - Clear intent and safer UX for a potentially destructive operation.
  - Clean place for dry-run, diff, and interactive conflict options.
  - Easier to document precedence rules and merge strategy.
- Cons:
  - Adds a new top-level command to maintain.
  - Slightly more discoverability burden than extending `edit`.
- Risks:
  - Command proliferation if similar niche operations become standalone commands.
- Rough effort estimate: Medium.

### Alternative C - `strand repair --sync-templates`
- Description: Add a repair mode that detects template drift and rewrites tasks to match current templates.
- Assumptions: Repair already scans all tasks and can apply bulk changes.
- Pros:
  - Convenient for large-scale consistency updates.
  - Reuses existing global scan plumbing.
- Cons:
  - Violates expectation that `repair` is primarily integrity-focused.
  - High blast radius; easy to rewrite many tasks unintentionally.
  - Difficult to support granular interactive merges per task.
- Risks:
  - Accidental broad mutations in CI or routine maintenance workflows.
- Rough effort estimate: Large.

## Merge Strategy (applies to selected alternative)
- Frontmatter handling:
  - Template defaults are candidates only; existing explicit task values win.
  - Preserve task-managed fields (`date_created`, `date_edited`, `completed`, `status`, relations).
  - Add newly introduced template fields only if not already present.
- Body handling:
  - Parse template into structured sections.
  - Preserve existing section content when headings match.
  - Insert newly introduced template sections with default text.
  - Retain unmatched user-authored sections unless `--prune` is explicitly set.
- Conflict handling:
  - When both template and task changed same canonical section scaffold, mark conflict in dry-run output and require `--interactive` or `--force`.

## Recommended Direction
- Recommended alternative: **Alternative B (`strand reapply`)**.
- Why: It keeps risky structural updates explicit and gives room for safe merge UX (`--dry-run`, `--interactive`, `--output-diff`) without complicating `edit`.

## Proposed CLI Shape
- `strand reapply <task-id> --template <type>`
- Optional flags:
  - `--dry-run`: Show proposed changes, do not write.
  - `--output-diff`: Print markdown/frontmatter diff.
  - `--interactive`: Resolve section-level conflicts interactively.
  - `--prune`: Remove template sections that no longer exist.
  - `--no-repair`: Skip post-write repair pass.

## Implementation Plan (epics/tasks)
1. Add a TaskDB template-reapply API that accepts task ID, template type, and merge options; return a structured plan for dry-run and apply modes.
2. Implement template/task normalization and section-aware merge logic with deterministic ordering for frontmatter keys and section insertion.
3. Add the `strand reapply` CLI command and wire dry-run/diff/interactive flows.
4. Add test coverage:
   - Unit tests for frontmatter precedence and section merge behavior.
   - Golden tests for deterministic output.
   - E2E tests for dry-run, apply, and conflict paths.
5. Update CLI documentation with examples and migration guidance from manual copy/paste workflows.

## Review Requests
- Request review from: `master-reviewer`, `reviewer-usability`, `reviewer-reliability`.

## Decision
- Decision: deferred to Owner.
