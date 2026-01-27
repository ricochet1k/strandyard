# Design: memmd list Command

## Overview
Define the CLI contract and output formats for a `memmd list` command, including filters, sorting, and groupings aligned with master lists and task metadata.

## Command and scopes
Primary command: `memmd list`

Scopes (mutually exclusive):
- `--scope all` (default): list all tasks discovered in the tasks tree.
- `--scope root`: list only root tasks (parent empty).
- `--scope free`: list only tasks with no blockers.
- `--parent <TASK-ID>`: list direct children of the given parent task ID.
- `--path <REL-PATH>`: list tasks under a subtree path (repo-relative under `tasks/`), useful for epics.

## Filters
All filters are additive (AND).
- `--role <role>`: filter by role name.
- `--priority <high|medium|low>`: filter by priority.
- `--completed <true|false>`: filter by completion state.
- `--blocked <true|false>`: `true` means blockers length > 0, `false` means zero blockers.
- `--blocks <true|false>`: `true` means blocks length > 0.
- `--owner-approval <true|false>`: if present in frontmatter.
- `--label <name>`: reserved for future use; return "unsupported field" error if invoked before labels exist.

## Sorting and determinism
Default sort order (stable):
1. Priority (high, medium, low)
2. Completed (false before true)
3. Task ID (lexicographic)

Optional overrides:
- `--sort id|priority|created|edited|role`
- `--order asc|desc` (default `asc`)

When `--sort created|edited` is used, fall back to Task ID on ties.

Determinism requirements:
- When reading and writing lists, all outputs are stable given the same tasks tree.
- Any derived lists (markdown/table) use the same stable ordering rules.

## Output formats
Default: `table` for terminal.

Formats:
- `--format table`: aligned columns for terminal.
- `--format md`: Markdown list or table (see below).
- `--format json`: array of objects; intended for automation.

Common schema fields across formats:
`id`, `title`, `role`, `priority`, `parent`, `completed`, `blockers`, `blocks`, `path`, `date_created`, `date_edited`.

Table columns (default, with optional `--columns`):
`id`, `title`, `priority`, `role`, `completed`, `blockers`.

Markdown output:
- If `--group priority`, output grouped sections (High/Medium/Low) each as list.
- Otherwise output a single list or table (configurable via `--md-table` boolean).

JSON output (shape example):
```
[
  {
    "id": "T3m9p-add-dep",
    "title": "Add goldmark-frontmatter Dependency",
    "role": "developer",
    "priority": "medium",
    "parent": "E2k7x-metadata-format",
    "completed": false,
    "blockers": [],
    "blocks": [],
    "path": "tasks/E2k7x-metadata-format/T3m9p-add-dep/T3m9p-add-dep.md",
    "date_created": "2026-01-27T00:00:00Z",
    "date_edited": "2026-01-27T13:43:58Z"
  }
]
```

## Grouping options
- `--group none|priority|parent|role` (default `none`)
- Grouping applies to `table` and `md`; for `json` add a top-level object keyed by group if `--group` is set.

## Data sources
Default data source: scan the tasks tree (`scan`-equivalent) to ensure accurate filtering.

Optional optimization:
- `--use-master-lists` (boolean): when `--scope root` or `--scope free` and no additional filters, use `tasks/root-tasks.md` or `tasks/free-tasks.md` for speed. For any filters or non-root/free scopes, use scan.

## CLI flags summary
`memmd list [--scope <all|root|free>] [--parent <TASK-ID>] [--path <REL-PATH>] [--role <role>] [--priority <p>] [--completed <bool>] [--blocked <bool>] [--blocks <bool>] [--owner-approval <bool>] [--sort <key>] [--order <asc|desc>] [--format <table|md|json>] [--columns <list>] [--group <none|priority|parent|role>] [--md-table] [--use-master-lists]`

Errors:
- Invalid flag combinations: `--scope free` with `--parent`, `--path`, or `--group parent` should error with helpful message.
- Unknown parent ID should return "not found" with non-zero exit code.

## Implementation notes
- `cmd/list.go`: new Cobra command, flag parsing, output dispatch.
- `pkg/task/scan.go` (or existing scan module): add `ListOptions` and filtering helpers.
- `pkg/task/format.go`: output formatting (table/md/json) to keep command thin.
- `tasks/root-tasks.md`, `tasks/free-tasks.md`: treat as read-only sources when `--use-master-lists` is set.
- `CLI.md`: add `list` command docs.

## Testing strategy
- Unit tests for filter combinations and sorting stability in `pkg/task`.
- Golden tests for output formatting (`table`, `md`, `json`).
- CLI tests to ensure flags map to options correctly.

## Alternatives considered
1. Separate subcommands (`list root`, `list free`) instead of `--scope`.
   - Pros: clearer help/usage.
   - Cons: more commands to maintain; more duplication.
   - Decision: prefer `--scope` for a stable single interface.
2. Default to master lists for root/free scopes always.
   - Pros: faster.
   - Cons: can be stale if user edits tasks without validate.
   - Decision: default to scan; allow `--use-master-lists` opt-in.
