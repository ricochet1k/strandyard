# Review: blockers.go Relationship Management

## Artifacts
- `pkg/task/blockers.go`
- `pkg/task/taskdb.go`
- `pkg/task/blockers_test.go`
- `pkg/task/taskdb_test.go`
- `pkg/task/TASKDB_DESIGN.md`
- `.strand/tasks/T06ubsf-consolidate-blocker-relationship-repair.md`

## Scope
Analyze the current blocker relationship reconciliation behavior (including the previous `UpdateBlockersFromChildren` intent), identify strengths and issues, and frame whether this logic belongs in TaskDB or outside of it.

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


Your task is T0q5n-review-blockers-go-relationship-management. Here's the description of that task:

---
type: task
role: developer
priority: medium
parent: Ti6zj-taskdb-api-design-review
blockers: []
blocks:
    - Ti6zj-taskdb-api-design-review
date_created: 2026-01-31T17:18:48.599142Z
date_edited: 2026-01-31T17:18:48.617275Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Review blockers.go relationship management

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
Analyze pkg/task/blockers.go:
- Document UpdateBlockersFromChildren behavior and purpose
- Identify what it does well (original design)
- Note any issues with current implementation
- Understand the relationship computation logic
- Determine if this belongs in TaskDB or is a separate concern

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.
```

## Behavior and purpose (including legacy naming)
- The previous `UpdateBlockersFromChildren` behavior is now represented by `ReconcileBlockerRelationships`.
- Purpose: compute canonical blocker edges from three sources, then rewrite `blockers` and `blocks` as sorted, unique, bidirectional sets.
- Sources of edges:
  - Parent-child rule: each incomplete child blocks its parent.
  - Explicit `blocks` edges.
  - Explicit `blockers` edges.
- Completed tasks are excluded as blockers and as blocked nodes in computed edges.

## Relationship computation logic
1. Build `desiredBlockers[blockedID][blockerID]` by scanning all tasks and adding valid edges.
2. Rewrite every task's `Meta.Blockers` to sorted keys of `desiredBlockers[taskID]`.
3. Derive inverse map `desiredBlocks[blockerID][blockedID]` from rewritten blockers.
4. Rewrite every task's `Meta.Blocks` from `desiredBlocks`.
5. Mark tasks dirty only when slices changed; return count of modified tasks.

## What the current design does well
- Single-pass canonicalization eliminates drift between `blockers` and `blocks`.
- Deterministic sorting and deduplication keeps diffs stable.
- Parent-derived blocking is centrally enforced instead of duplicated in callers.
- Ignores dangling references safely during reconciliation (while validator can flag/fix data quality separately).

## Issues and risks in current implementation
- `UpdateBlockersAfterCompletion` is incremental and can temporarily diverge from full reconciliation semantics (it removes inbound blockers from blocked tasks but does not rebuild all inverse edges).
- Missing-task references are silently dropped during reconciliation; if caller code expects strict failure here, it must run validation explicitly.
- Naming history (`UpdateBlockersFromChildren` -> `ReconcileBlockerRelationships`) is not visible in code comments, which makes older design discussions harder to map.

Follow-up concern task: `T09easy` (unify completion cleanup with reconciliation invariants).

## Boundary analysis: TaskDB vs separate concern

### Option A: Keep reconciliation in TaskDB (recommended)
- Pros:
  - Matches repo policy that relationship mutation should go through TaskDB.
  - Keeps integrity logic close to relationship mutators (`SetParent`, `AddBlocker`, `SetCompleted`).
  - Easier to guarantee invariants before persistence.
- Cons:
  - TaskDB API surface remains broad unless reconciliation entry points are tightly defined.

### Option B: Move reconciliation to a separate relationship engine
- Pros:
  - Clearer separation between storage/cache and graph normalization logic.
  - Could simplify isolated testing and future reuse.
- Cons:
  - More orchestration burden on callers and easier misuse if callers bypass the engine.
  - Conflicts with current project direction to centralize task operations in TaskDB.

Decision: deferred to Owner.

## Repro steps
```bash
go test ./pkg/task -run ReconcileBlockerRelationships
go test ./pkg/task -run UpdateBlockersAfterCompletion
go run ./cmd/strand repair
```

## Test and validation results
```bash
go test ./pkg/task -run ReconcileBlockerRelationships
ok  github.com/ricochet1k/strandyard/pkg/task  0.296s

go test ./pkg/task -run UpdateBlockersAfterCompletion
ok  github.com/ricochet1k/strandyard/pkg/task  0.182s

go build ./...
success

go test ./...
ok across all packages

go run ./cmd/strand repair
repair: ok
Repaired 0 tasks
```
