# StrandYard CLI Usage Guide

This document describes how to use the `strand` command-line tool for managing tasks.

## Prerequisites

```bash
# Build the CLI
go build -o strand .

# Or run directly with go run
go run ./cmd/strand <command>
```

## Project storage and selection

strand stores `tasks/`, `roles/`, and `templates/` either in a local `.strand/` directory at the git root or in a global project directory under `~/.config/strand/projects/<project_name>` (default).

### `init` - Initialize strand storage

Initialize the strand project storage for the current repository.

```bash
strand init [project_name] [--storage global|local] [--preset <dir-or-git-url>]
```

Flags:
- `--storage`: choose `global` (default) or `local` (`.strand/` at git root)
- `--preset`: path to a directory or a git repo containing `tasks/`, `roles/`, and `templates/` to copy into the project. If a git URL is provided, it will be cloned securely.

### `preset refresh` - Refresh roles and templates from a preset

Refresh roles and templates from a preset source (local directory or git URL).

A preset is a directory containing:
- `roles/` - role documents
- `templates/` - task templates

This command will:
- ✓ Overwrite existing role and template files
- ✓ Preserve your `tasks/` directory (tasks are never touched)
- ✓ Validate preset structure before copying
- ✓ Run `repair` automatically after refreshing

```bash
strand preset refresh <preset>
```

**The preset source can be**:
- Local directory path: `/path/to/my-preset`
- Git HTTPS URL: `https://github.com/user/strand-preset.git`
- Git SSH URL: `git@github.com:user/strand-preset.git`

**Examples**:
```bash
# Refresh from local directory
strand preset refresh /path/to/my-preset

# Refresh from GitHub repository
strand preset refresh https://github.com/example/strand-presets.git

# Refresh using SSH (requires configured keys)
strand preset refresh git@github.com:user/strand-preset.git
```

**Output**:
The command provides verbose feedback about what's happening:
```
Using local preset directory: /path/to/preset
Validating preset structure...
✓ Preset structure validated

Refreshing roles/...
  Refreshing roles/developer.md
  Refreshing roles/architect.md
Refreshing templates/...
  Refreshing templates/task.md
  Refreshing templates/epic.md
✓ Refresh complete. Running repair...
repair: ok
```

**Common errors**:
- **"preset is missing required directories"**: The preset doesn't have `roles/` and/or `templates/` subdirectories
- **"failed to clone preset: repository not found"**: Git URL is incorrect or repository is inaccessible
- **"failed to clone preset: authentication required"**: Private repository requires SSH keys or access tokens
- **"project not initialized"**: Run `strand init` first

**Notes**:
- Fails if the project is not already initialized.
- Validates preset structure before making any changes.
- Shows exactly which files are being refreshed.

## Core Commands


### `add` - Create tasks from templates

Creates tasks using the appropriate template by type.

```bash
strand add <type> [title] [flags]

Flags:
  -t, --title string      task title (can also be passed as positional argument)
  -r, --role string       role responsible for the task (defaults from template)
  -p, --parent string     parent task ID
      --priority string   priority: high, medium, or low (defaults from template)
      --blocker strings   blocker task ID(s); can be repeated or comma-separated
      --no-repair       skip repair and master list updates
```

**Example**:
```bash
strand add task "Quick task" --role developer --priority high
```

**Detailed body via stdin**:
```bash
# Pipe from a file
strand add task "Incident followup" --role developer < ./notes.md

# Heredoc
strand add task "Incident followup" --role developer <<'EOF'
## Summary
Capture findings from the incident review.

## Tasks
- [ ] Draft timeline
- [ ] Identify owners
EOF
```

Notes:
- Stdin content is inserted where the template uses `{{ .Body }}` or appended to the end.

### `edit` - Edit a task

Edits a task's metadata and description.

```bash
strand edit <task-id> [flags]

Flags:
  -t, --title string      task title
  -r, --role string       role responsible for the task
  -p, --parent string     parent task ID
      --priority string   priority: high, medium, or low
      --blocker strings   blocker task ID(s); can be repeated or comma-separated
      --no-repair       skip repair and master list updates
```

**Example**:
```bash
strand edit T3k7x --priority high --role architect
```

**Edit description via stdin (heredoc)**:
```bash
strand edit T3k7x <<'EOF'
# Updated Title
New description goes here.
EOF
```

Notes:
- If stdin is a terminal, the description remains unchanged.
- If both `--title` and stdin are provided, the title flag overrides any H1 in the stdin content.

### `list` - List tasks with filters and output formats

Lists tasks from the tasks tree with optional filtering, sorting, and grouping.

```bash
strand list [flags]

Flags:
  --scope string         scope of tasks to list: all|root|free (default "all")
  --children string      list direct children of the given task ID
  --role string          filter by role name
  --priority string      filter by priority: high|medium|low
  --completed            list only completed tasks (default: uncompleted)
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
# List uncompleted tasks (default)
strand list

# List completed tasks
strand list --completed

# List root tasks
strand list --scope root

# List free tasks grouped by priority in Markdown
strand list --scope free --format md --group priority

# List children of a task (short IDs supported)
strand list --children E2k7x

# List tasks with filtering and sorting
strand list --role developer --priority high --sort created --order desc
```

**Notes**:
- `--scope free` cannot be combined with `--children` or `--group parent`.
- `--children` is only valid with `--scope all`.
- `--label` is reserved and will error until labels are implemented.

### `search` - Search tasks by content

Searches task title, description, and todos. Subtask names are excluded from search results.

```bash
strand search <query> [flags]

Flags:
  --sort string          sort by: id|priority|created|edited|role
  --order string         sort order: asc|desc (default "asc")
  --format string        output format: table|md|json (default "table")
  --columns string       comma-separated list of columns to include
  --group string         group by: none|priority|parent|role (default "none")
  --md-table             use markdown table output (with --format md)
```

**Examples**:
```bash
# Search by keyword
strand search "frontmatter"

# Search and output Markdown
strand search "owner approval" --format md --group priority
```
- Task ID flags accept short IDs like `T3k7x` (prefix + token).

### `add issue` - Create an issue task

Creates an issue-style task using the issue template and required metadata.

```bash
strand add issue [title] [flags]

Flags:
  -t, --title string      issue title (can also be passed as positional argument)
  -r, --role string       role responsible for the task (defaults from template)
  -p, --parent string     parent task ID
      --priority string   priority: high, medium, or low (defaults from template)
      --blocker strings   blocker task ID(s); can be repeated or comma-separated
      --no-repair       skip repair and master list updates
```

**Example**:
```bash
strand add issue "Add issue command" --priority high
```

### `add` - Create recurring task definitions

Creates a task with a recurrence rule that can be materialized into subsequent tasks.

```bash
strand add <type> [title] --every "<interval> <unit> [from|after <anchor>]"
```

Flags:
  - `--every`: Recurrence rule (e.g., "10 days", "50 commits from HEAD", "20 tasks_completed"). Can be repeated for multiple rules. Supports `from <anchor>` (start at anchor) or `after <anchor>` (start one interval after anchor).

**Example**:
```bash
strand add task "Quarterly docs review" --every "3 months after Jan 1 2026 00:00 UTC" --role reviewer
```

**Resulting task metadata (example)**:
```yaml
---
role: master-reviewer
priority: medium
parent:
every:
  - 3 months from Jan 1 2026 00:00 UTC
date_created: 2026-01-01T00:00:00Z
date_edited: 2026-01-01T00:00:00Z
completed: false
---
```

### Recurrence Audit Logging

When a recurring task is created or materialized using a dynamic anchor (`HEAD`, `now`, or an empty anchor), the system automatically logs the resolved value (the specific commit hash or timestamp) to the activity log for auditability.

These entries have the type `recurrence_anchor_resolved` and include:
- `timestamp`: When the resolution occurred
- `task_id`: The ID of the task being created or evaluated
- `metadata.original`: The original anchor string (e.g., "HEAD", "now")
- `metadata.resolved`: The resolved value (e.g., a 40-char commit hash or a formatted timestamp)

The activity log is stored at `.strand/activity.log`.

### `repair` - Repair task structure

Repairs all tasks and regenerates master lists (`root-tasks.md` and `free-tasks.md`).

```bash
strand repair [flags]

Flags:
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
$ strand repair
repair: ok
Repaired 25 tasks
Master lists updated: tasks/root-tasks.md, tasks/free-tasks.md
```

### `next` - Get next task to work on

Displays the next free task (tasks with no blockers) with the role document.

```bash
strand next [flags]

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
$ strand next
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
$ strand next --role developer
```

### `complete` - Mark task as completed

Marks a task as completed by setting `completed: true` in the frontmatter and updating `date_edited`.

```bash
strand complete <task-id> [report]
```

**What it does**:
- Finds task by ID
- Sets `completed: true` in frontmatter
- Updates `date_edited` to current timestamp
- Appends the report to the task body if provided
- Preserves all other metadata

**After completing**: `strand complete` should update master lists and remove the completed task from `free-tasks.md`. If `strand repair` changes anything afterward, treat it as a bug.

**Example**:
```bash
$ strand complete T3m9p-add-frontmatter-dep "Added goldmark-frontmatter and updated go.mod"
✓ Task T3m9p-add-frontmatter-dep marked as completed
```

## Typical Workflows

### Working on a task

1. **Find next task**:
   ```bash
   strand next
   ```

2. **Work on the task** (write code, update files, etc.)

3. **Mark task complete**:
   ```bash
   strand complete <task-id> "report of what was done"
   ```

4. **Optional verification** (should be a no-op):
   ```bash
   strand repair
   ```

### Checking task status

```bash
# See all free tasks (no blockers), grouped by priority
cat tasks/free-tasks.md

# See all root tasks (no parent)
cat tasks/root-tasks.md

# Repair entire task tree
strand repair
```

### Working with roles

```bash
# Get next task for specific role
strand next --role developer

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

### Recurrence Metadata

Tasks can include recurrence rules using the `every` field in the frontmatter.

**Example**:
```yaml
every:
  - 10 days
  - 50 commits from HEAD
```

**Supported units**:
- `days`, `weeks`, `months`
- `commits`, `lines_changed` (git-based)
- `tasks_completed` (activity-log-based)

**Anchors**:
- If no anchor is specified (e.g., `10 days`), the anchor defaults to `now` for time-based rules or `HEAD` for git-based rules.
- Explicit anchors can be provided using `from <anchor>` (start at anchor) or `after <anchor>` (start one interval after anchor).
- Date anchors support ISO 8601 (e.g., `2026-01-28T09:00:00Z`) and the human-friendly format `Jan 2 2006 15:04 MST`.
- `tasks_completed` anchors can be a task ID (short or full) or a date/time.

**Special considerations for git-based recurrence**:
- When a rule uses `HEAD` as an anchor, it indicates the latest commit.
- Invalid or "unborn" HEAD states are treated as a no-op (no tasks are materialized) rather than an error.

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
strand/
├── tasks/
│   ├── root-tasks.md          # Auto-generated list of root tasks
│   ├── free-tasks.md          # Auto-generated list of free tasks
│   ├── E2k7x-metadata-format.md
│   ├── T3m9p-add-dep.md
│   ├── T8h4w-update-parser.md
│   └── E6w3m-id-generation.md
├── roles/
│   ├── developer.md
│   ├── architect.md
│   ├── designer.md
│   └── owner.md
├── templates/
│   └── task.md
└── design-docs/
    └── commands-design.md
```

## Environment Variables

- **MEMMD_ROLE**: Default role for `next` command (not currently used, but flag available)

## Error Messages

### "malformed ID: must be <PREFIX><4-lowercase-alphanumeric>-<slug>"
Task filename doesn't match required format. Rename the file to follow format.

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

### `add --every` hint examples
Hint lines for `strand add --every` use deterministic examples so automation and tests remain stable. Canonical examples and anchor guidance live in `design-docs/recurrence-anchor-hint-examples.md`.

Default anchor examples:
- `--every "10 days"`
- `--every "50 commits from HEAD"`
- `--every "500 lines_changed from HEAD"`
- `--every "20 tasks_completed"`

Explicit anchor examples:
- `--every "10 days from Jan 28 2026 09:00 UTC"`
- `--every "1 month after Jan 1 2026 00:00 UTC"`

Anchor guidance:
- Use the human-friendly date anchor above for explicit date/time examples.
- Use `HEAD` for commit-based defaults and explicit commit anchors.
- Task ID anchors (short or full) can be used for `tasks_completed` (e.g., `from T1a1a`).
- ISO 8601 anchors are fully supported (for example, `2026-01-28T09:00:00Z`).

## Tips

1. **Run repair after manual edits**: It should be a no-op after `complete` and `add`
2. **Use `next` to find work**: Don't manually browse `free-tasks.md`
3. **Complete tasks promptly**: Mark tasks complete as soon as work is done
4. **Check repair errors carefully**: They indicate data integrity issues
5. **Keep role files updated**: Role documents guide AI agents and humans
