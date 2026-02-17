---
type: review
role: master-reviewer
priority: medium
parent: ""
blockers:
    - Trtik-review-repair-go-validation-logic
    - Tuu6q-update-existing-usage-throughout-codebase
    - Twcdw-audit-api-surface-and-identify-misuse-opportunitie
    - Tx4jn-design-relationship-management-api
    - Txvyh-design-task-creation-api-template-based-only
blocks: []
date_created: 2026-01-31T17:18:35.743126Z
date_edited: 2026-02-15T05:12:33.908807Z
owner_approval: false
completed: false
status: ""
description: ""
---

# TaskDB API Design Review

## Artifacts
List the documents, designs, or code paths under review.

## Scope
Clarify what is in and out of scope for this review.

## Review Focus
List the specific areas to evaluate (e.g., usability, API ergonomics, error handling).

## Escalation
Create new tasks for concerns or open questions instead of editing this task. Record decisions and final rationale in design docs.

## Checklist
- [ ] Artifacts and scope listed.
- [ ] Review focus defined.
- [ ] Concerns captured as subtasks.
- [ ] Decision items deferred to Owner as separate subtasks when needed.


Review the current TaskDB implementation against repository principles and design a proper API that makes it hard to misuse.

Key concerns:
- Tasks should only be created from templates
- Manual Task creation and field manipulation must be prevented
- Relationship integrity must be enforced automatically
- Clear API surface with well-named, single-responsibility functions
- Remove code that violates core principles (e.g., GetOrCreate)

Reference: pkg/task/TASKDB_DESIGN.md

## Subtasks
- [x] (subtask: T06ubsf) Consolidate blocker relationship repair
  Consolidated blocker reconciliation into TaskDB.ReconcileBlockerRelationships and removed duplicate Sync/Fix paths. Updated repair command and TaskDB/blocker tests/examples to use the unified bidirectional reconciliation flow.
- [x] (subtask: T0f98) Review new taskdb.go implementation
  Reviewed taskdb.go exported API and relationship methods; identified redundant blocker reconciliation paths and naming issues; captured decision to keep AddBlocked/RemoveBlocked and logged follow-up consolidation task T06ubsf.
- [x] (subtask: T0q5n) Review blockers.go relationship management
  Reviewed blocker relationship logic in pkg/task/blockers.go and documented behavior, strengths, and risks in design-docs/blockers-relationship-management-review.md. Captured boundary alternatives (TaskDB vs separate engine) with decision deferred to Owner, ran build/test/repair, and filed follow-up issue T09easy for reconciliation invariant gaps.
- [x] (subtask: T2lt8) Review task.go structure and methods
  Reviewed pkg/task/task.go structure, documented fields/methods/dirty tracking in design-docs/task-go-structure-review.md, added unit tests for Task setters/content and write helpers.
- [x] (subtask: T48or) Write updated godoc with usage examples
  Added comprehensive pkg/task package godoc with TaskDB mental model, pitfalls, and design references; rewrote TaskDB examples for recommended workflow, reconciliation semantics, status lifecycle, validation/repair, and misuse guidance. Verified with go build ./..., go test ./..., and strand repair.
- [x] (subtask: T7qkw) Create implementation plan
  Created design-docs/taskdb-api-implementation-plan.md with phased TaskDB hardening plan covering dependency/risk ordering, breaking changes, migration strategy, and incremental-vs-refactor split. Included decision gates mapped to unresolved Ti6zj subtasks and a concrete code-change checklist tied to pkg/task, cmd, and web call sites.
- [x] (subtask: T8bgf) Design access control strategy
  Created design-docs/taskdb-access-control-strategy.md with access-control alternatives (A/B/C), trade-offs, and a concrete recommended design that keeps decision deferred for owner approval. Linked the new strategy doc from design-docs/taskdb-api-implementation-plan.md, removed no APIs yet, and validated repository health with go build ./..., go test ./..., and go run ./cmd/strand repair.
- [x] (subtask: Tb0oq) Identify code to delete
  Added design-docs/taskdb-api-deletion-list.md with a concrete removal inventory covering GetOrCreate, redundant relationship wrappers, misuse call sites, and legacy completion helpers, including rationale, migration targets, and deletion order. Verified with go test ./..., go build ./..., and go run ./cmd/strand repair.
- [x] (subtask: Tdhrq) Review parser.go and task loading
  Added design-docs/parser-go-task-loading-review.md documenting parser load flow, parse-time validation, relationship parsing behavior, and current no-semantic-validation boundary. Included parser API reference, owner-deferred boundary alternatives, task selection transcript, and verification outputs from parser tests, full test suite, build, and repair (0 changes).
- [x] (subtask: Thyd1) Document TaskDB responsibilities and boundaries
  Documented TaskDB boundaries across TaskDB, Parser, Validator, and command layers; added architecture diagram and operation ownership map in design-docs/taskdb-responsibilities-and-boundaries.md with alternatives/pros-cons and deferred decision notes; validated with targeted pkg/task tests, full go test, build, and repair.
- [ ] (subtask: Trtik) Review repair.go validation logic
- [ ] (subtask: Tuu6q) Update existing usage throughout codebase
- [ ] (subtask: Twcdw) Audit API surface and identify misuse opportunities
- [ ] (subtask: Tx4jn) Design relationship management API
- [ ] (subtask: Txvyh) Design task creation API (template-based only)
- [x] (subtask: Tqn2blh) repair changes files immediately after strand complete
  Added blocker reconciliation to strand complete paths and a regression test ensuring immediate strand repair reports Repaired 0 tasks after completion.
