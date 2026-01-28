# memmd CLI Usage Guide

This document describes how to use the `memmd` command-line tool for managing tasks.

## Prerequisites

```bash
# Build the CLI
go build -o memmd .

# Or run directly with go run
go run . <command>
```

## Core Commands

### `add` - Create tasks from templates

Creates tasks using the appropriate template by type.

```bash
memmd add <type> [title] [flags]

Flags:
  -t, --title string      task title (can also be passed as positional argument)
  -r, --role string       role responsible for the task (defaults from template)
  -p, --parent string     parent task ID (creates task under that directory)
      --priority string   priority: high, medium, or low (defaults from template)
      --blocker strings   blocker task ID(s); can be repeated or comma-separated
      --no-repair       skip repair and master list updates
```

**Example**:
```bash
memmd add leaf "Quick task" --role developer --priority high
```

### `list` - List tasks with filters and output formats

Lists tasks from the tasks tree with optional filtering, sorting, and grouping.

```bash
memmd list [flags]

Flags:
  --scope string         scope of tasks to list: all|root|free (default "all")
  --parent string        list direct children of the given parent task ID
  --path string          list tasks under a subtree path (repo-relative under tasks/)
  --role string          filter by role name
  --priority string      filter by priority: high|medium|low
  --completed            filter by completed status (only when flag is present)
  --blocked              filter by blocked status (has blockers)
  --blocks               filter by blocks status (has blocks)
  --owner-approval       filter by owner approval
  --label string         reserved for future labels support (errors if used)
  --sort string          sort by: id|priority|created|edited|role
  --order string         sort order: asc|desc (default "asc")
  --format string        output format: table|md|json (default "table")
  --columns string       comma-separated list of columns to include
  --group string         group by: none|priority|parent|role (default "none")
  --md-table             use markdown table output (with --format md)
  --use-master-lists     use master lists for root/free scopes when no filters
```

**Examples**:
```bash
# List all tasks (default scope)
memmd list

# List root tasks
memmd list --scope root

# List free tasks grouped by priority in Markdown
memmd list --scope free --format md --group priority

# List children of a parent task
memmd list --parent E2k7x-metadata-format

# List tasks under a subtree path
memmd list --path tasks/E2k7x-metadata-format

# List tasks with filtering and sorting
memmd list --role developer --priority high --sort created --order desc
```

**Notes**:
- `--scope free` cannot be combined with `--parent`, `--path`, or `--group parent`.
- `--parent` and `--path` are mutually exclusive and only valid with `--scope all`.
- `--label` is reserved and will error until labels are implemented.

### `add issue` - Create an issue task

Creates an issue-style task using the issue template and required metadata.

```bash
memmd add issue [title] [flags]

Flags:
  -t, --title string      issue title (can also be passed as positional argument)
  -r, --role string       role responsible for the task (defaults from template)
  -p, --parent string     parent task ID (creates task under that directory)
      --priority string   priority: high, medium, or low (defaults from template)
      --blocker strings   blocker task ID(s); can be repeated or comma-separated
      --no-repair       skip repair and master list updates
```

**Example**:
```bash
memmd add issue "Add issue command" --priority high
```

### `recurring add` - Create recurring task definitions

Creates a recurring task definition that can be materialized into normal tasks.

```bash
memmd recurring add [title] [flags]

Flags:
  -t, --title string           definition title (can also be passed as positional argument)
  -r, --role string            role responsible for generated tasks (defaults from template)
  -p, --parent string          parent task ID (creates definition under that directory)
      --priority string        priority: high, medium, or low (defaults from template)
      --blocker strings        blocker task ID(s); can be repeated or comma-separated
      --interval int           recurrence interval (required)
      --unit string            recurrence unit: days, weeks, months, commits (required)
      --anchor string          anchor date (ISO 8601) or commit hash (required)
      --timezone string        IANA time zone for scheduling (optional)
      --max-instances int      max instances to materialize (optional)
      --no-repair            skip repair and master list updates
```

**Example**:
```bash
memmd recurring add "Quarterly docs review" --interval 3 --unit months --anchor 2026-01-01T00:00:00Z --role reviewer
```

**Resulting definition (example)**:
```markdown
---
type: recurring
role: reviewer
priority: medium
parent:
blockers: []
blocks: []
date_created: 2026-01-01T00:00:00Z
date_edited: 2026-01-01T00:00:00Z
owner_approval: false
completed: false
recurrence_interval: 3
recurrence_unit: months
recurrence_anchor: 2026-01-01T00:00:00Z
recurrence_next_due: 2026-04-01T00:00:00Z
---

# Quarterly docs review
```

**Definition layout (example)**:
```
tasks/
  <parent-dir>/
    <RECURRING_ID>/
      <RECURRING_ID>.md
```

### `recurring materialize` - Materialize due recurring tasks

Generates concrete task instances for any recurring definitions that are due.

```bash
memmd recurring materialize [flags]

Flags:
  --as-of string         override the current time (ISO 8601)
  --path string          path to tasks directory (default "tasks")
  --dry-run              preview materialization without writing files
  --limit int            limit number of instances generated
```

**Example**:
```bash
memmd recurring materialize --as-of 2026-04-01T00:00:00Z
```

### `repair` - Repair task structure

Repairs all tasks and regenerates master lists (`root-tasks.md` and `free-tasks.md`).

```bash
memmd repair [flags]

Flags:
  --path string      path to tasks directory (default "tasks")
  --roots string     path to write root tasks list (default "tasks/root-tasks.md")
  --free string      path to write free tasks list (default "tasks/free-tasks.md")
  --format string    output format: text|json (default "text")
```

**What it repairs**:
- Task IDs match format: `<PREFIX><4-lowercase-alphanumeric>-<slug>` (e.g., `T3k7x-example`)
- Role files exist in `roles/` directory
- Parent tasks exist
- Blocker tasks exist
- Priority is one of: `high`, `medium`, `low` (empty defaults to `medium`)
- YAML frontmatter is valid

**When to run**: After creating, modifying, or completing tasks.

**Example**:
```bash
$ memmd repair
repair: ok
Repaired 25 tasks
Master lists updated: tasks/root-tasks.md, tasks/free-tasks.md
```

### `next` - Get next task to work on

Displays the next free task (tasks with no blockers) with the role document.

```bash
memmd next [flags]

Flags:
  --role string    optional: filter tasks by role
```

**Output format**:
1. Full role document from `roles/<role>.md`
2. Separator line (`---`)
3. Task file content with YAML frontmatter

Invariant: `next` must print the full role document (not just the role name); tests should guard this.

**Example**:
```bash
$ memmd next
# Architect

## Role
Architect (human or senior AI) — breaks accepted designs...

---
---
role: architect
priority: medium
parent:
blockers: []
date_created: 2026-01-27
---

# Implement YAML Frontmatter Metadata Format

## Summary
Replace the current simple field format...
```

**With role filter**:
```bash
$ memmd next --role developer
```

### `complete` - Mark task as completed

Marks a task as completed by setting `completed: true` in the frontmatter and updating `date_edited`.

```bash
memmd complete <task-id>
```

**What it does**:
- Finds task by ID (searches entire task tree)
- Sets `completed: true` in frontmatter
- Updates `date_edited` to current timestamp
- Preserves all other metadata

**After completing**: Run `memmd repair` to update master lists and remove completed task from `free-tasks.md`.

**Example**:
```bash
$ memmd complete T3m9p-add-frontmatter-dep
✓ Task T3m9p-add-frontmatter-dep marked as completed

Run 'memmd repair' to update master lists

$ memmd repair
repair: ok
Repaired 25 tasks
Master lists updated: tasks/root-tasks.md, tasks/free-tasks.md
```

## Typical Workflows

### Working on a task

1. **Find next task**:
   ```bash
   memmd next
   ```

2. **Work on the task** (write code, update files, etc.)

3. **Mark task complete**:
   ```bash
   memmd complete <task-id>
   ```

4. **Update master lists**:
   ```bash
   memmd repair
   ```

### Checking task status

```bash
# See all free tasks (no blockers), grouped by priority
cat tasks/free-tasks.md

# See all root tasks (no parent)
cat tasks/root-tasks.md

# Repair entire task tree
memmd repair
```

### Working with roles

```bash
# Get next task for specific role
memmd next --role developer

# See what roles are available
ls roles/

# View role document
cat roles/developer.md
```

## Task File Format

Tasks use YAML frontmatter for metadata:

```markdown
---
role: developer
priority: medium
parent: E2k7x-metadata-format
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T13:43:58Z
owner_approval: false
completed: false
---

# Task Title

## Summary
Brief description of the task...

## Tasks
- [ ] First step
- [ ] Second step

## Acceptance Criteria
- Criteria for completion
```

### Required Frontmatter Fields

- **role**: Role responsible for this task (must match a file in `roles/`)
- **parent**: Parent task ID (empty string for root tasks)
- **blockers**: Array of task IDs that block this task (empty array if none)
- **date_created**: ISO 8601 timestamp
- **date_edited**: ISO 8601 timestamp

Note: issue tasks (`type: issue`) default to role `triage` and `priority: medium` from the template, but can be overridden with flags.

### Optional Frontmatter Fields

- **blocks**: Array of task IDs this task blocks
- **owner_approval**: Boolean flag for owner approval
- **completed**: Boolean flag marking task as complete
- **priority**: Task priority (`high`, `medium`, or `low`; defaults to `medium`)
- **type**: Task subtype string (e.g., `issue`, `recurring`)

### Recurrence Metadata (for `type: recurring`)

Recurring definitions require additional scheduling fields and are validated by `repair`.

**Required fields**:
- `recurrence_interval` (integer > 0)
- `recurrence_unit` (`days`, `weeks`, `months`, or `commits`)
- `recurrence_anchor` (ISO 8601 date or commit hash)

**Optional fields**:
- `recurrence_next_due` (ISO 8601 date; computed if omitted)
- `recurrence_last_run` (ISO 8601 date)
- `recurrence_timezone` (IANA time zone, e.g., `America/Los_Angeles`)
- `recurrence_max_instances` (positive integer)

**Validation notes**:
- Recurring definitions are excluded from `free-tasks.md` and `root-tasks.md` until materialized.
- Materialized tasks behave like normal tasks and appear in master lists if unblocked.

## Task ID Format

Task IDs must follow this format: `<PREFIX><4-char-token>-<slug>`

- **PREFIX**: Single uppercase letter denoting task type
  - `T` = Task
  - `E` = Epic
  - `D` = Design document
  - `I` = Issue
- **Token**: 4 lowercase alphanumeric characters (base36: 0-9, a-z)
- **Slug**: Human-readable identifier (lowercase, hyphens allowed)

**Valid examples**:
- `T3k7x-implement-parser`
- `E2k7x-metadata-format`
- `D9m2p-api-design`

**Invalid examples**:
- `T123-bad` (only 3 chars in token)
- `TABCD-bad` (uppercase in token)
- `T3k7x_bad` (underscore in slug)

## Directory Structure

```
memmd/
├── tasks/
│   ├── root-tasks.md          # Auto-generated list of root tasks
│   ├── free-tasks.md          # Auto-generated list of free tasks
│   ├── E2k7x-metadata-format/ # Epic directory
│   │   ├── E2k7x-metadata-format.md
│   │   ├── T3m9p-add-dep/     # Child task
│   │   │   └── T3m9p-add-dep.md
│   │   └── T8h4w-update-parser/
│   │       └── T8h4w-update-parser.md
│   └── E6w3m-id-generation/
│       └── E6w3m-id-generation.md
├── roles/
│   ├── developer.md
│   ├── architect.md
│   ├── designer.md
│   └── owner.md
├── templates/
│   └── leaf.md
└── design-docs/
    └── commands-design.md
```

## Environment Variables

- **MEMMD_ROLE**: Default role for `next` command (not currently used, but flag available)

## Error Messages

### "malformed ID: must be <PREFIX><4-lowercase-alphanumeric>-<slug>"
Task directory name doesn't match required format. Rename directory to follow format.

### "role file roles/X.md does not exist"
Create the missing role file in the `roles/` directory.

### "parent task X does not exist"
Referenced parent task doesn't exist. Fix the `parent:` field in frontmatter.

### "blocker task X does not exist"
Referenced blocker task doesn't exist. Fix the `blockers:` array in frontmatter.

### "invalid priority \"X\": must be high, medium, or low"
Set `priority: high|medium|low` or remove the field to default to `medium`.

### "task not found: X"
Task ID doesn't exist in the task tree. Check spelling or use `repair` to see all tasks.

## Tips

1. **Always run repair after changes**: Keeps master lists up-to-date
2. **Use `next` to find work**: Don't manually browse `free-tasks.md`
3. **Complete tasks promptly**: Mark tasks complete as soon as work is done
4. **Check repair errors carefully**: They indicate data integrity issues
5. **Keep role files updated**: Role documents guide AI agents and humans
