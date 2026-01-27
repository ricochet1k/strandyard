---
role: architect
parent:
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
---

# Update Next Command Behavior

## Summary

Update the `next` command to default to reading the first free task and printing that task's role (or the role from its first TODO), then the task description.

## Context

**Owner Decision**: `next` should default to reading the first free task, print that task's role (or the role of its first TODO), then print the task description. No role filtering required by default.

**Current state**: `next` requires `--role` flag or `MEMMD_ROLE` env var, filters tasks by role, prints role doc + task.

**Target state**: `next` reads first free task, extracts role from task metadata or first TODO, prints minimal output (task role + task content).

## Subtasks

1. [T5h7w-default-free-task](T5h7w-default-free-task/T5h7w-default-free-task.md) - Update next to default to first free task
2. [T8n2m-print-task-role](T8n2m-print-task-role/T8n2m-print-task-role.md) - Update next to print task's role
3. [T6p4k-todo-role-detection](T6p4k-todo-role-detection/T6p4k-todo-role-detection.md) - Add TODO-based role detection

## Acceptance Criteria

- `memmd next` works without any flags
- Reads first task from free-tasks.md
- Prints task's role or role from first TODO
- Prints task content
- No role filtering by default

## References

- Current implementation: cmd/next.go
