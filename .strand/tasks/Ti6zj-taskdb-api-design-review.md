---
type: review
role: master-reviewer
priority: medium
parent: ""
blockers:
    - T2lt8-review-task-go-structure-and-methods
    - T48or-write-updated-godoc-with-usage-examples
    - T7qkw-create-implementation-plan
    - T8bgf-design-access-control-strategy
    - Tb0oq-identify-code-to-delete
    - Tdhrq-review-parser-go-and-task-loading
    - Thyd1-document-taskdb-responsibilities-and-boundaries
    - Tqn2blh-repair-changes-files-immediately-after-strand-comp
    - Trtik-review-repair-go-validation-logic
    - Tuu6q-update-existing-usage-throughout-codebase
    - Twcdw-audit-api-surface-and-identify-misuse-opportunitie
    - Tx4jn-design-relationship-management-api
    - Txvyh-design-task-creation-api-template-based-only
blocks: []
date_created: 2026-01-31T17:18:35.743126Z
date_edited: 2026-02-08T04:11:10.516938Z
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
- [x] (subtask: T0f98) Review new taskdb.go implementation
- [x] (subtask: T0q5n) Review blockers.go relationship management
- [ ] (subtask: T2lt8) Review task.go structure and methods
- [ ] (subtask: T48or) Write updated godoc with usage examples
- [ ] (subtask: T7qkw) Create implementation plan
- [ ] (subtask: T8bgf) Design access control strategy
- [ ] (subtask: Tb0oq) Identify code to delete
- [ ] (subtask: Tdhrq) Review parser.go and task loading
- [ ] (subtask: Thyd1) Document TaskDB responsibilities and boundaries
- [ ] (subtask: Trtik) Review repair.go validation logic
- [ ] (subtask: Tuu6q) Update existing usage throughout codebase
- [ ] (subtask: Twcdw) Audit API surface and identify misuse opportunities
- [ ] (subtask: Tx4jn) Design relationship management API
- [ ] (subtask: Txvyh) Design task creation API (template-based only)
- [ ] (subtask: Tqn2blh) repair changes files immediately after strand complete
