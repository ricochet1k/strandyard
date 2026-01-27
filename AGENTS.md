# AGENTS — memmd

Purpose: record repository-specific knowledge for AI coding agents working on memmd (a Go library + CLI managing a task DB stored as plain Markdown files).

## Overview

- Language: Go (Golang).
- Target: A library + CLI that manages a filesystem-backed task database saved as Markdown.
- Key libraries: `github.com/yuin/goldmark` (Markdown parsing/rendering), `github.com/spf13/cobra` (CLI framework).
- Scaffolding: use `cobra-cli` to bootstrap the CLI and command skeletons.

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

- Tasks are stored as directories. Each task directory contains a single `task.md` (or `README.md`) that describes the task and its metadata.
- The directory hierarchy mirrors parent/child lineage. Example:

  ```
  tasks/project-alpha/          (root task)
    task.md                     (metadata + description)
    subtask-a/
      task.md
  ```

- Task metadata (stored as structured markdown sections):
  - Title: human-friendly title header
  - ID: short unique id (slug)
  - Role/Assignee: a named role (must match one of the files in `roles/`)
  - Parent: path to parent task (or blank for root tasks)
  - Blockers: list of paths to tasks that block this task
  - Blocks: list of paths to tasks this task is blocking
  - Children: optional list of child task paths (kept deterministic by CLI)

- Example `task.md` layout (follow exactly to maximize parseability):

```markdown
# Title: Initialize project skeleton

ID: init-project

Role: developer

Parent:

Blockers:
- []

Blocks:
- tasks/project-alpha/task.md

Description:
Add initial Go module, scaffold cobra CLI, and commit.
```

## Roles

- Each role has a markdown file in `roles/` (e.g. `roles/ai-assistant.md`, `roles/developer.md`). The CLI will treat the role filename (without extension) as the canonical role name. Role files should document responsibilities and any constraints/privileges.

## Master lists
- Two deterministic master list files are kept at `tasks/root-tasks.md` and `tasks/free-tasks.md`.
  - `root-tasks.md`: lists all root tasks (no Parent)
  - `free-tasks.md`: lists tasks with no blockers (ready to start)
- These files are updated deterministically by the CLI commands; they should not be edited manually except for bootstrapping.

## CLI responsibilities (high level)

- Provide commands to:
  - `init` — initialize repo structure and optional example tasks/roles
  - `scan` — parse the tasks tree and validate references
  - `sync` — update `root-tasks.md` and `free-tasks.md` deterministically
  - `add`/`new` — create a new task directory and `task.md` with provided metadata
  - `assign` — change a task's Role/Assignee
  - `block`/`unblock` — add/remove blockers and update related tasks
  - `render` — render a task (or list) to HTML/terminal via `goldmark`

## Templates

- Task templates: `templates/task-templates/` (use these for leaf/implementable tasks). `ID` and `Parent` are derived from the filesystem; do not include them in templates.
- Document templates: `templates/doc-templates/` (design alternatives, epic documents, etc.).

## Parsing rules & expectations
- Use `goldmark` to parse markdown bodies, but treat metadata sections (ID/Role/Parent/Blockers/Blocks/Description) as structured text blocks (simple line prefixes and lists). Keep parsing logic robust to minor ordering changes but prefer canonical layout above.

## Conventions and patterns

- Deterministic ordering: lists (children, blockers, blocks) must be written in sorted order so diffs are stable.
- Use relative paths inside task metadata (repo-relative) to reference other tasks.
- Role names are the lowercase filename (without `.md`) in `roles/`.

## Files currently added for bootstrapping

- `roles/ai-assistant.md` and `roles/developer.md` — initial role documents.
- `tasks/root-tasks.md`, `tasks/free-tasks.md` — initial master lists.
- `tasks/project-alpha/task.md` and `tasks/project-alpha/setup-infra/task.md` — example task directories showing parent/child layout and blockers.

## When you need to change behaviour

- If you require changes to parsing or metadata format, update this file and add unit tests exercising the parsing/scanning tools.

## Agent policy

- **Policy**: Agents must not unilaterally choose which alternative to implement. Present clear alternatives with pros/cons and defer the final decision to a human maintainer (mark as "Decision: deferred" in reviews).
- **Guidance**: When preparing role-based reviews or selecting the next actionable task, run `go run . next` to obtain the canonical role document and next task. Include the full stdout/stderr output from that command in review artifacts and do not assume task selection without running it.

## Questions to ask the repo owner (useful prompts)
- What Go module path should be used for `go mod init`?
- Any CI or formatting rules (gofmt/pre-commit hooks)?
- Preferred task metadata fields beyond the canonical set above?

If anything in this document is unclear, tell me the desired module path and any additional fields you want included in `task.md` examples and I will update the examples and scaffolding instructions.
