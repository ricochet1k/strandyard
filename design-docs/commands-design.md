# CLI Design: Commands & Examples

This document describes the planned CLI surface for `memmd`. The implementation uses `cobra` for commands and `goldmark` for parsing markdown task bodies. IDs and Parent are derived from task directory/file names; Blockers/Blocks are maintained by the CLI and should not be edited by hand.

Command summary (revised)
- `memmd init` — bootstrap repository structure and optional examples.
- `memmd validate` — parse tasks, validate metadata/links, and automatically update deterministic master lists (`tasks/root-tasks.md`, `tasks/free-tasks.md`).
- `memmd add` / `memmd new` — create a new task directory from a task template.
- `memmd assign` — change a task's `Role`/assignee.
- `memmd block` — add/remove blocker relationships between tasks.
- `memmd templates list` — list available templates.
- `memmd next` — primary agent-facing command; emits role doc followed by the next task document for that role.

High-level behaviour and guarantees
- Filesystem-first: each task is a file `<task-id>/<task-id>.md` (plus attachments). Directory hierarchy represents parent/child lineage.
- Deterministic outputs: `validate` writes `tasks/root-tasks.md` and `tasks/free-tasks.md` with sorted entries to keep diffs stable; this happens automatically.
- CLI manages relationships (Blocks/Blockers) — prefer CLI commands over manual edits.
- IDs and Parents: `ID` is derived from the task directory name and the CLI enforces a canonical ID format (see below). Templates must not declare `ID`/`Parent`.

Task ID format
- Each task is identified by a short ID composed of:
	- a single-character type prefix (uppercase) denoting task type (e.g. `T`=task, `D`=design doc, `E`=epic),
	- a short random alphanumeric token (4 chars, base36),
	- a short mini-title (1-2 words hyphen-joined) used only for human-readability.

- Format: `<prefix><token>-<mini>`
- Example: `T3k7x-init` or `D9x2b-api`.

Primary commands (details)

1) `memmd init`
- Purpose: create initial directories and bootstrap templates, roles and example tasks if missing.
- Example:

```bash
memmd init --force
```

2) `memmd validate`
- Purpose: parse the `tasks/` tree, validate referential integrity (roles exist, parent links valid), enforce formatting rules, and automatically regenerate deterministic master lists.
- Behaviour:
	- Parse all `<task-id>.md` files, extract metadata sections, and validate references.
	- Regenerate `tasks/root-tasks.md` and `tasks/free-tasks.md` deterministically (sorted).
	- Exit non-zero if validation errors exist (missing roles, broken links, malformed IDs).
- Example:

```bash
memmd validate --path tasks --format json
```

3) `memmd add` / `memmd new`
- Purpose: create a new task directory and `<task-id>.md` from a template. CLI will generate the canonical ID and directory name.
- Flags: `--title`, `--role`, `--parent` (task id or "root"), `--template`.
- Example:

```bash
memmd add --title "Implement validation" --role developer --parent scaffold-cli --template task_template
```

- Effect: CLI generates an ID (e.g. `T3k7x-impl`), creates `tasks/T3k7x-impl/<task-id>.md` with template-filled fields. `validate` runs implicitly (or will be run later) to update master lists.

4) `memmd assign`
- Purpose: change the `Role`/assignee of a task.
- Example:

```bash
memmd assign T3k7x-impl --role owner
```

- Effect: updates the `role` field in the YAML frontmatter of `<task-id>.md` atomically.

5) `memmd block` (subcommands: `add`, `remove`, `list`)
- Purpose: manage blocker relationships using canonical IDs.
- Examples:

```bash
memmd block add --task T3k7x-impl --blocks T9x2b-design
memmd block list --task T3k7x-impl
```

- Effect: updates the `blockers` and `blocks` fields in the YAML frontmatter of involved tasks deterministically and ensures sorted order.

6) `memmd templates list`
- Purpose: list available templates in `templates/`.

7) `memmd next` (primary agent-facing command)
- Purpose: provide an agent with all context it needs to start work: first the Role document, then the task document for the next actionable item for that role.
- Behaviour:
	- It selects the highest-priority/first `free` task for that role from `tasks/free-tasks.md` and prints the task's `<task-id>.md` content.
	- It reads `roles/<role>.md` for the first todo in the task and prints it to stdout.
	- The output is plain concatenated markdown (role doc then task doc) so an agent can consume it in one pass.
- Example:

```bash
memmd next
```

Implementation notes
- Use `github.com/yuin/goldmark` with `goldmark-frontmatter` extension to parse YAML frontmatter and Markdown bodies.
- The CLI must ensure deterministic ordering for master lists to avoid churn.
- ID generation: provide a deterministic, cryptographically-strong short token generator (e.g., `crypto/rand` → base36 4 chars) and a slugify function for mini-title.

Developer commands (examples)

```bash
# Run unit tests
go test ./... 

# Build CLI
go build -o bin/memmd ./

# Quick local test (init + add + validate + next)
memmd init --force
memmd add --title "Quick task" --role developer
memmd validate
MEMMD_ROLE=developer memmd next
```

Files and templates
- Task templates: `templates/` (must not include ID/Parent/Blockers fields).
- Document examples: `doc-examples/` (examples of alternatives and epics).
- Roles: `roles/` (each role is a markdown file; filename without extension is canonical role name).

Open implementation questions
- Slug generation: allow explicit `--id` override? Recommend slugify by default and allow `--id` for rare overrides. No --id.
- Claiming semantics: should `--claim` create a claim file or directly modify `Role:`? No --claim.

If this matches your expectations, I'll implement `validate` (parse + auto-sync) next, then `next`.
