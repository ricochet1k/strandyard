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

### `validate` - Validate task structure

Validates all tasks and regenerates master lists (`root-tasks.md` and `free-tasks.md`).

```bash
memmd validate [flags]

Flags:
  --path string      path to tasks directory (default "tasks")
  --roots string     path to write root tasks list (default "tasks/root-tasks.md")
  --free string      path to write free tasks list (default "tasks/free-tasks.md")
  --format string    output format: text|json (default "text")
```

**What it validates**:
- Task IDs match format: `<PREFIX><4-lowercase-alphanumeric>-<slug>` (e.g., `T3k7x-example`)
- Role files exist in `roles/` directory
- Parent tasks exist
- Blocker tasks exist
- YAML frontmatter is valid

**When to run**: After creating, modifying, or completing tasks.

**Example**:
```bash
$ memmd validate
validate: ok
Validated 25 tasks
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

**Example**:
```bash
$ memmd next
# Architect

## Role
Architect (human or senior AI) — breaks accepted designs...

---
---
role: architect
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

**After completing**: Run `memmd validate` to update master lists and remove completed task from `free-tasks.md`.

**Example**:
```bash
$ memmd complete T3m9p-add-frontmatter-dep
✓ Task T3m9p-add-frontmatter-dep marked as completed

Run 'memmd validate' to update master lists

$ memmd validate
validate: ok
Validated 25 tasks
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
   memmd validate
   ```

### Checking task status

```bash
# See all free tasks (no blockers)
cat tasks/free-tasks.md

# See all root tasks (no parent)
cat tasks/root-tasks.md

# Validate entire task tree
memmd validate
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

### Optional Frontmatter Fields

- **blocks**: Array of task IDs this task blocks
- **owner_approval**: Boolean flag for owner approval
- **completed**: Boolean flag marking task as complete

## Task ID Format

Task IDs must follow this format: `<PREFIX><4-char-token>-<slug>`

- **PREFIX**: Single uppercase letter denoting task type
  - `T` = Task
  - `E` = Epic
  - `D` = Design document
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

### "task not found: X"
Task ID doesn't exist in the task tree. Check spelling or use `validate` to see all tasks.

## Tips

1. **Always run validate after changes**: Keeps master lists up-to-date
2. **Use `next` to find work**: Don't manually browse `free-tasks.md`
3. **Complete tasks promptly**: Mark tasks complete as soon as work is done
4. **Check validation errors carefully**: They indicate data integrity issues
5. **Keep role files updated**: Role documents guide AI agents and humans
