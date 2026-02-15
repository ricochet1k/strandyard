# Review: parser.go Task Loading and Parser API

## Context and references
- `pkg/task/parser.go`
- `pkg/task/parser_test.go`
- `pkg/task/task.go`
- `pkg/task/TASKDB_DESIGN.md`
- `design-docs/taskdb-api-implementation-plan.md`
- Parent review task: `Ti6zj` (TaskDB API Design Review)

## Task selection output
```text
Your role is developer. Here's the description of that role:

---
description: "Implements tasks, writes code, and produces working software."
---

# Developer

## Role
Developer (human or AI) — implements tasks, writes code, and produces working software.

## Responsibilities
- Implement tasks assigned by the Architect
- Write clean, maintainable code following project conventions
- Add tests for new functionality
- Document code and update relevant documentation
- Fix bugs and address issues
- Ensure code passes validation and tests before marking tasks complete

## Deliverables
- Working code that meets acceptance criteria
- Tests covering the implemented functionality
- Updated documentation as needed
- Code that passes `go build` and `go test`

## Workflow
1. Read the assigned task and understand acceptance criteria
2. Implement the functionality described in the task
3. Write tests to verify the implementation
4. Run repair: `go build ./...`, `go test ./...`, `strand repair`
5. Update task status and mark as completed when done, including a brief report of what was accomplished: `strand complete <task-id> "Summary of work"`

## Constraints
- Follow existing code patterns and conventions in the codebase
- Ensure all changes are backward compatible unless explicitly noted
- Do not modify task metadata manually - use CLI commands

---

Ancestors:
  Ti6zj: TaskDB API Design Review


Your task is Tdhrq-review-parser-go-and-task-loading. Here's the description of that task:

---
type: task
role: developer
priority: medium
parent: Ti6zj-taskdb-api-design-review
blockers: []
blocks:
    - Ti6zj-taskdb-api-design-review
date_created: 2026-01-31T17:18:46.507223Z
date_edited: 2026-02-15T04:48:41.145067Z
owner_approval: false
completed: false
status: in_progress
description: ""
---

# Review parser.go and task loading

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
Analyze pkg/task/parser.go:
- How are tasks loaded from disk?
- What validation happens during parsing?
- How are relationships read from frontmatter?
- Is there any relationship validation during load?
- Document the Parser API

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.
```

## How tasks are loaded from disk
1. `LoadTasks(tasksRoot string)` walks `tasksRoot` with `filepath.WalkDir`.
2. It skips directories, non-`.md` files, and master list files (`root-tasks.md`, `free-tasks.md`).
3. For each task markdown file, it starts a goroutine (errgroup limit 10) calling `ParseFile(path)`.
4. Parsed tasks are sent through `tasksChan` and collected into `map[string]*Task` keyed by filename-derived ID.
5. Return value is `(tasks map[string]*Task, err error)` where `err` is `errors.Join(walkErr, parseErr)`.

## Validation performed during parsing
- Frontmatter boundary validation:
  - If opening `---` exists without a closing delimiter, parser returns `InvalidFrontmatterError`.
- YAML syntax validation:
  - `yaml.Unmarshal` parses frontmatter into `Metadata`.
  - Malformed YAML returns `FrontmatterParseError` and attempts to extract a YAML line number.
- File-path context enrichment:
  - `ParseFile` and `ParseStandaloneFile` attach the file path to parsing/frontmatter errors when missing.

No semantic validation of task metadata is performed in `parser.go` (for example role existence, parent existence, blocker reciprocity, self-blocking, or cycle checks).

## How relationships are read from frontmatter
- Relationship fields are read by YAML unmarshal into `Metadata`:
  - `parent` -> `Meta.Parent` (string)
  - `blockers` -> `Meta.Blockers` (`[]string`)
  - `blocks` -> `Meta.Blocks` (`[]string`)
- These values are treated as raw parsed data at load time.
- `ParseString` does not reconcile, deduplicate, sort, or cross-check relationship IDs.

## Relationship validation during load
- In current behavior, relationship validation is not done in `LoadTasks` or `Parse*` methods.
- Relationship consistency and normalization are handled later via TaskDB/repair workflows (for example reconciliation and validator paths outside parser loading).

## Parser API documentation

### Construction
- `NewParser() *Parser`: creates a parser instance (stateless today).

### Parse entrypoints
- `(*Parser) ParseFile(filePath string) (*Task, error)`:
  - Reads file, derives task ID from filename, parses markdown/frontmatter, sets `Task.FilePath` and `Task.Dir`.
- `(*Parser) ParseStandaloneFile(filePath string) (*Task, error)`:
  - Same parse flow as `ParseFile`; intended for markdown not constrained to task directory conventions.
- `(*Parser) ParseString(content, id string) (*Task, error)`:
  - Splits frontmatter/body, unmarshals metadata, parses heading sections into title/body/TODOs/subtasks/progress.

### Bulk loading
- `(*Parser) LoadTasks(tasksRoot string) (map[string]*Task, error)`:
  - Concurrently parses task files under a root and returns map by task ID.

### Helpers
- `SplitByHeadings(body string) []Section`: section tokenizer for markdown content.
- `ExtractTitle(content string) string`: returns first H1 heading.

## Design alternatives (relationship validation boundary)

### Option A: Keep parser as syntax-only loader (current)
- Pros:
  - Fast, focused responsibility (I/O + syntax + structural parse).
  - Easy reuse for partial or tooling-only reads.
  - Keeps costly graph checks centralized in TaskDB/repair.
- Cons:
  - Callers can mistakenly assume loaded relationships are valid.

### Option B: Add semantic relationship validation to `LoadTasks`
- Pros:
  - Earlier detection of bad references and malformed graph state.
  - Less risk of callers using invalid tasks before reconciliation.
- Cons:
  - Blurs parser responsibilities and duplicates TaskDB/repair logic.
  - Harder to use parser for permissive reads and migration tooling.

Decision: deferred to owner.

## Reproduction steps
```bash
go test ./pkg/task -run 'Parse|LoadTasks'
go build ./...
go test ./...
go run ./cmd/strand repair
```

## Validation results
```bash
go test ./pkg/task -run 'Parse|LoadTasks'
ok   github.com/ricochet1k/strandyard/pkg/task  0.317s

go build ./...
success

go test ./...
ok   github.com/ricochet1k/strandyard/cmd  1.961s
ok   github.com/ricochet1k/strandyard/pkg/activity  (cached)
ok   github.com/ricochet1k/strandyard/pkg/idgen  (cached)
ok   github.com/ricochet1k/strandyard/pkg/task  (cached)
ok   github.com/ricochet1k/strandyard/pkg/web  (cached)
ok   github.com/ricochet1k/strandyard/test/e2e  (cached)

go run ./cmd/strand repair
repair: ok
Repaired 0 tasks
```
