---
type: ""
role: architect
priority: ""
parent: ""
blockers: []
blocks:
    - E5w8m-e2e-tests
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-30T22:38:51.494414Z
owner_approval: false
completed: true
---

# Update Next Command Behavior

## Summary
Update the `next` command to default to reading the first free task and printing that task's role (or the role from its first TODO), then the task description.

## Context
**Owner Decision**: `next` should default to reading the first free task, print that task's role (or the role of its first TODO), then print the task description. No role filtering required by default.

**Current state**: `next` requires `--role` flag or `MEMMD_ROLE` env var, filters tasks by role, prints role doc + task.

**Target state**: `next` reads first free task, extracts role from task metadata or first TODO, prints minimal output (task role + task content).

## Acceptance Criteria
- `strand next` works without any flags
- Reads first task from free-tasks.md
- Prints task's role or role from first TODO
- Prints task content
- No role filtering by default

## References
- Current implementation: cmd/next.go

## TODOs
- [ ] [T5h7w-default-free-task](T5h7w-default-free-task/T5h7w-default-free-task.md) - Update next to default to first free task
- [ ] [T8n2m-print-task-role](T8n2m-print-task-role/T8n2m-print-task-role.md) - Update next to print task's role
- [ ] [T6p4k-todo-role-detection](T6p4k-todo-role-detection/T6p4k-todo-role-detection.md) - Add TODO-based role detection

## Subtasks
- [x] (subtask: T5h7w) Update Next to Default to First Free Task
- [x] (subtask: T6p4k) Some Task
- [x] (subtask: T8n2m) Task Title
