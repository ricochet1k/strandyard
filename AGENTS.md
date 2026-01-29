# AGENTS — memmd

Purpose: record repository-specific knowledge for AI coding agents working on memmd (a Go library + CLI managing a task DB stored as plain Markdown files).

## Overview

- Language: Go (Golang).
- Target: A library + CLI that manages a filesystem-backed task database saved as Markdown.
- Key libraries: `github.com/yuin/goldmark` (Markdown parsing/rendering), `github.com/spf13/cobra` (CLI framework).
- Scaffolding: use `cobra-cli` to bootstrap the CLI and command skeletons.

## CLI Usage

For detailed CLI command documentation, see [CLI.md](CLI.md).

**Quick commands**:
```bash
# Get next task to work on
memmd next

# Mark task as completed
memmd complete <task-id>

# Repair all tasks and update master lists
memmd repair
```

## Quick setup (what to ask the repo owner if missing)

- Desired Go module path (e.g. `github.com/<you>/memmd`) — required for `go mod init`.
- CI/test commands (if any custom ones).

## Recommended developer commands
- Install cobra scaffolder:
```bash
go install github.com/spf13/cobra-cli@latest
```
- Init project (example — replace module path):
```bash
go mod init github.com/yourname/memmd
cobra-cli init --pkg-name github.com/yourname/memmd
go get github.com/yuin/goldmark
go get github.com/spf13/cobra
```
- Build / test:
```bash
go build ./...
go test ./...
```

## Data model and filesystem conventions (authoritative)

- Tasks are stored as directories. Each task directory contains a single markdown file (named `<task-id>.md`, `task.md`, or `README.md`) with YAML frontmatter.
- Task ID is derived from the directory name: must follow format `<PREFIX><4-char-token>-<slug>` (e.g., `T3k7x-example`, `E2k7x-metadata`)
- The directory hierarchy mirrors parent/child lineage. Example:

  ```
  tasks/
    E2k7x-metadata-format/           (root epic)
      E2k7x-metadata-format.md       (epic file)
      T3m9p-add-dep/                 (child task)
        T3m9p-add-dep.md
      T8h4w-update-parser/           (child task)
        T8h4w-update-parser.md
  ```

- Task metadata is stored as YAML frontmatter using `goldmark-frontmatter`:
  - **role**: Role responsible (must match a file in `roles/`)
  - **priority**: Task priority (`high`, `medium`, or `low`; defaults to `medium`)
  - **parent**: Parent task ID (empty for root tasks)
  - **blockers**: Array of task IDs that block this task
  - **blocks**: Array of task IDs this task blocks (optional)
  - **date_created**: ISO 8601 timestamp
  - **date_edited**: ISO 8601 timestamp
  - **completed**: Boolean flag (marks task as done)
  - **owner_approval**: Boolean flag (optional)

- Task file format (follow exactly):

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

# Add goldmark-frontmatter Dependency

## Summary
Add the goldmark-frontmatter library to the project...

## Tasks
- [ ] Run `go get github.com/abhinav/goldmark-frontmatter`
- [ ] Verify dependencies resolve correctly

## Acceptance Criteria
- goldmark-frontmatter is listed in go.mod
- Project builds successfully
```

## Roles

- Each role has a markdown file in `roles/` (e.g. `roles/ai-assistant.md`, `roles/developer.md`). The CLI will treat the role filename (without extension) as the canonical role name. Role files should document responsibilities and any constraints/privileges.

## Master lists
- Two deterministic master list files are kept at `tasks/root-tasks.md` and `tasks/free-tasks.md`.
 - Two deterministic master list files are kept at `tasks/root-tasks.md` and `tasks/free-tasks.md`.
  - `root-tasks.md`: lists all root tasks (no Parent)
  - `free-tasks.md`: lists tasks with no blockers (ready to start), grouped by priority
- These files are updated deterministically by the CLI commands (including `complete` and `repair`); they should not be edited manually except for bootstrapping.

## CLI responsibilities (high level)

- Provide commands to:
  - `init` — initialize repo structure and optional example tasks/roles
  - `scan` — parse the tasks tree and repair references
  - `sync` — update `root-tasks.md` and `free-tasks.md` deterministically
  - `add`/`new` — create a new task directory and `task.md` with provided metadata
  - `assign` — change a task's Role/Assignee
  - `block`/`unblock` — add/remove blockers and update related tasks
  - `render` — render a task (or list) to HTML/terminal via `goldmark`

## Templates

- Task templates: `templates/` (use these for implementable tasks). `ID` and `Parent` are derived from the filesystem; do not include them in templates.
- Document examples: `doc-examples/` (example task outputs, sample documents). Use these as templates for documents in `design-docs/` (for example, `doc-examples/design-alternatives.md`). Task bodies should come from `templates/` at task creation time.
- Task templates must be fully specifiable before work starts; avoid placeholders for results or findings.
- Do not edit task bodies to record outcomes; create follow-up tasks for concerns or deferred decisions.
- The default task template is `templates/task.md`. The `leaf` template/type is deprecated with no backward compatibility.

## Parsing rules & expectations

- Use `goldmark` with `goldmark-frontmatter` extension to parse task files
- Metadata is extracted from YAML frontmatter (between `---` delimiters)
- Task ID is derived from directory name, not from frontmatter
- Markdown content after frontmatter is preserved as task body
- Parsing implementation: see `pkg/task/task.go` for the canonical parser

## Conventions and patterns

- Deterministic ordering: lists (children, blockers, blocks) must be written in sorted order so diffs are stable.
- Use relative paths inside task metadata (repo-relative) to reference other tasks.
- Role names are the lowercase filename (without `.md`) in `roles/`.

## Current repository structure

- **Roles**: `roles/developer.md`, `roles/architect.md`, `roles/designer.md`, `roles/owner.md`, `roles/reviewer*.md` — role documents
- **Master lists**: `tasks/root-tasks.md`, `tasks/free-tasks.md` — auto-generated by `repair` command
- **Task library**: `pkg/task/` — goldmark-based task parser and validator
- **CLI commands**: `cmd/validate.go`, `cmd/next.go`, `cmd/complete.go` — core CLI commands
- **Documentation**: `CLI.md` — CLI usage guide, `AGENTS.md` — this file

## When you need to change behaviour

- If you require changes to parsing or metadata format, update this file and add unit tests exercising the parsing/scanning tools.

## Agent policy

- **Policy**: Agents must not unilaterally choose which alternative to implement. Present clear alternatives with pros/cons and defer the final decision to a human maintainer (mark as "Decision: deferred" in reviews).
- **Execution**: Do the next logical steps yourself (tests, commits, review requests, follow-ups) unless blocked by a missing decision, credentials, or explicit instruction to wait. Do not suggest next steps in responses; execute them or state the specific blocker that prevents execution.
- **Session title**: After `memmd next` returns the task, set the session title to `<role>: <task title>` exactly (lowercase role, task title as shown).
- **Task references in responses**: Use the task ID or `Title (short id)`; avoid full task paths unless explicitly requested.
- **Commit after completion or block**: After completing a task or becoming blocked, commit your changes with a clear message before starting the next task or handing off.
- **Guidance**: When preparing role-based reviews or selecting the next actionable task, run `go run . next` to obtain the canonical role document and next task. Include the full stdout/stderr output from that command in review artifacts and do not assume task selection without running it.
- **Invariant**: The `next` command must print the full role document from `roles/<role>.md` followed by a `---` separator; keep the e2e test in place to prevent regressions.
- **After commit**: Whenever you complete and commit a task, use the `session` tool with `mode: "new"` and `async: true` to start a new session with the exact text `do the next task`, then stop work in the current session.
- **When blocked waiting on other work**: If the current task is blocked on reviews, owner decisions, or other tasks (including when you add wait-only subtasks) or the user says "done for now," use the `session` tool with `mode: "new"` and `async: true` to start a new session with the exact text `do the next task`, then stop work in the current session.
- **When asked "work on the next thing"**: Run `memmd next` and report a summary of the task. Then perform the role's duties on that task.
- **Repair after manual edits**: If you manually edit any task markdown files under `tasks/`, run `go run . repair` immediately afterward to regenerate master lists and confirm consistency. Make sure there's an issue filed to make sure the manual edit can be performed with a command eventually.
- **Complete tasks via CLI**: When a task is done (including planning-only tasks), run `memmd complete <task-id>` rather than editing frontmatter by hand. `memmd complete` should update master lists; if `memmd repair` changes anything afterward, treat it as a bug and file an issue.
- **File issues for manual edits**: If work requires manual edits or repairs outside the CLI, file an issue with repro steps, logs, and affected task IDs.
- **Use add for new tasks/issues**: When asked to add tasks/issues, use `memmd add` instead of creating task files manually.

## Questions to ask the repo owner (useful prompts)
- What Go module path should be used for `go mod init`?
- Any CI or formatting rules (gofmt/pre-commit hooks)?
- Preferred task metadata fields beyond the canonical set above?

If anything in this document is unclear, tell me the desired module path and any additional fields you want included in `task.md` examples and I will update the examples and scaffolding instructions.
