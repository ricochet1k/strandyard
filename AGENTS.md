# AGENTS — StrandYard

Purpose: record repository-specific knowledge for AI coding agents working on StrandYard (a Go library + CLI managing a task DB stored as plain Markdown files).

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
strand next --claim

# Mark task as completed
strand complete <task-id> "summary of work"

# Repair all tasks and update master lists
strand repair
```

## Recommended developer commands
- Install cobra scaffolder:
```bash
go install github.com/spf13/cobra-cli@latest
```
- Build / test:
```bash
go build ./...
go test ./...
```

## Data model and filesystem conventions (authoritative)

- Tasks are stored as flat markdown files in the `tasks/` directory (named `<task-id>.md`).
- Task ID is derived from the filename: must follow format `<PREFIX><4-char-token>-<slug>` (e.g., `T3k7x-example`, `E2k7x-metadata`)
- All tasks are stored in the same directory. Example:

  ```
  tasks/
    E2k7x-metadata-format.md       (epic file)
    T3m9p-add-dep.md               (child task)
    T8h4w-update-parser.md         (child task)
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


## Roles

- Each role has a markdown file in `.strand/roles/` (e.g. `.strand/roles/ai-assistant.md`, `.strand/roles/developer.md`). The CLI will treat the role filename (without extension) as the canonical role name. Role files should document responsibilities and any constraints/privileges.

## Master lists
- Two deterministic master list files are kept at `tasks/root-tasks.md` and `tasks/free-tasks.md`.
  - `root-tasks.md`: lists all root tasks (no Parent)
  - `free-tasks.md`: lists tasks with no blockers (ready to start), grouped by priority
- These files are updated deterministically by the CLI commands (including `complete` and `repair`); they should not be edited manually.

## CLI responsibilities (high level)

- Provide commands to:
  - `init` — initialize repo structure and optional example tasks/roles
  - `next` - print out the next task with it's role's full description as the full context an agent needs to complete the task
  - `add`/`new` — create a new task file with provided metadata
  - `edit` — update a task's metadata and description
  - `assign` — change a task's Role/Assignee
  - `block`/`unblock` — add/remove blockers and update related tasks
  - `repair` — parse the tasks tree and repair references

## Templates

- Task templates: run `strand templates` (use these for implementable tasks). `ID` and `Parent` are derived from the filesystem; do not include them in templates.
- Document examples: `doc-examples/` (example task outputs, sample documents). Use these as templates for documents in `design-docs/` (for example, `doc-examples/design-alternatives.md`). Task bodies should come from `templates/` at task creation time.
- Task templates must be fully specifiable before work starts; avoid placeholders for results or findings.
- Do not edit task bodies to record outcomes; tasks are disposable and may be deleted. Use follow-up tasks for open questions/concerns, and record decisions and final rationale in design docs.

## Parsing rules & expectations

- Use `goldmark` with `goldmark-frontmatter` extension to parse task files
- Metadata is extracted from YAML frontmatter (between `---` delimiters)
- Task ID is derived from filename, not from frontmatter
- Markdown content after frontmatter is preserved as task body
- Parsing implementation: see `pkg/task/task.go` for the canonical parser

## Conventions and patterns

- Deterministic ordering: lists (children, blockers, blocks) must be written in sorted order so diffs are stable.
- Subtasks under a parent should preserve creation order by default, support explicit manual reordering via CLI, and free-list generation should preserve that parent-defined subtask order where applicable.
- Use relative paths inside task metadata (repo-relative) to reference other tasks.
- Role names are the lowercase filename (without `.md`) in `roles/`.
- Web dashboard lists projects by name only; do not show local/global storage labels or duplicate entries.

## Current repository structure

- **Roles**: `.strand/roles/developer.md`, `.strand/roles/architect.md`, `.strand/roles/designer.md`, `.strand/roles/owner.md`, `.strand/roles/reviewer*.md` — role documents
- **Master lists**: `tasks/root-tasks.md` and `tasks/free-tasks.md` — auto-generated by `repair` command
- **Recurring tasks**: The `recurring` command has been merged into `add`. Use `strand add --every ...` (or as defined by the current implementation).
- **Task library**: `pkg/task/` — goldmark-based task parser and validator
- **CLI commands**: `cmd/validate.go`, `cmd/next.go`, `cmd/complete.go` — core CLI commands
- **Documentation**: `CLI.md` — CLI usage guide, `AGENTS.md` — this file

## When you need to change behaviour

- If you require changes to parsing or metadata format, update this file and add unit tests exercising the parsing/scanning tools.

## Agent policy

- **Policy**: Agents must not unilaterally choose which alternative to implement. Present clear alternatives with pros/cons and defer the final decision to a human maintainer (mark as "Decision: deferred" in reviews). Once a decision is made, update the entire design doc to reflect the final decision and any user preferences, and remove or condense alternatives.
- **Corrections**: When the user provides a general correction or preference, record it in this file as a durable rule or guideline.
- **TaskDB usage**: All CLI commands should rely on TaskDB for task operations instead of implementing task logic directly.
- **Execution**: Do the next logical steps yourself (tests, commits, review requests, follow-ups) unless blocked by a missing decision, credentials, or explicit instruction to wait. Do not suggest next steps in responses; execute them or state the specific blocker that prevents execution.
- **Next task selection**: `strand next` should respect role metadata that marks roles as opt-in/ignored by default; those roles are only selected when explicitly requested (for example via `strand next --role <role>`), and owner tasks should be handled in task order.
- **Session title**: Always set the session title immediately after receiving any task assignment (including `strand next`). Use `<role>: <task title>` exactly (lowercase role, task title as shown). Do not proceed with task work until the title is set; if you realize it was missed, set it immediately.
- **Task references in responses**: Use the task ID or `Title (short id)`; avoid full task paths unless explicitly requested.
- **Commit after completion or block**: After completing a task or becoming blocked, commit your changes with a clear message before starting the next task or handing off.
- **Guidance**: When preparing role-based reviews or selecting the next actionable task, run `./strand next --claim` to obtain the canonical role document and next task while claiming it. Include the full stdout/stderr output from that command in review artifacts and do not assume task selection without running it.
- **Invariant**: The `next` command must print the full role document from `roles/<role>.md` followed by a `---` separator; keep the e2e test in place to prevent regressions.
- **After commit**: Whenever you complete and commit a task, immediately use the `session` tool with `mode: "new"` and `async: true` to start a new session with the exact text `do the next task and you can commit if complete or blocked`. Only after the new session has been started should you provide a summary to the user. Then stop work in the current session.
- **Questions before handoff**: If you need user input, ask it before committing and before starting a new session. After you commit and start a new session, do not ask more questions in the current session.
- **When blocked waiting on other work**: If the current task is blocked on reviews, owner decisions, or other tasks (including when you add wait-only subtasks) or the user says "done for now," use the `session` tool with `mode: "new"` and `async: true` to start a new session with the exact text `do the next task and you can commit if complete or blocked`, then stop work in the current session.
- **Interactive handoff timing**: If a user is actively interacting (especially after being asked a question), do not hand off or start a new session until the user explicitly says to do so. Only hand off when the user is not intervening or asking follow-ups.
- **When asked "work on the next thing"**: Run `go run ./cmd/strand next --claim` and report a summary of the task. Then perform the role's duties on that task.
- **Repair after manual edits**: If you manually edit any task markdown files under `tasks/`, run `go run ./cmd/strand repair` immediately afterward to regenerate master lists and confirm consistency. Make sure there's an issue filed to make sure the manual edit can be performed with a command eventually.
- **Complete tasks via CLI**: When a task is done (including planning-only tasks), run `strand complete <task-id> "report of what was done"` rather than editing frontmatter by hand. The report should summarize key outcomes, important decisions, and any notable changes. `strand complete` should update master lists; if `strand repair` changes anything afterward, treat it as a bug and file an issue.
- **File issues for manual edits**: If work requires manual edits or repairs outside the CLI, file an issue with repro steps, logs, and affected task IDs.
- **Use add for new tasks/issues**: When asked to add tasks/issues, use `strand add` instead of creating task files manually.

## Questions to ask the repo owner (useful prompts)
- What Go module path should be used for `go mod init`?
- Any CI or formatting rules (gofmt/pre-commit hooks)?
- Preferred task metadata fields beyond the canonical set above?

If anything in this document is unclear, update it to make it more clear.
