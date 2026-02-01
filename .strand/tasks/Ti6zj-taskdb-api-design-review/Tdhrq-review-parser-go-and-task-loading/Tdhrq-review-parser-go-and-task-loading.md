---
type: task
role: developer
priority: medium
parent: Ti6zj-taskdb-api-design-review
blockers: []
blocks:
    - Ti6zj-taskdb-api-design-review
date_created: 2026-01-31T17:18:46.507223Z
date_edited: 2026-01-31T17:18:46.528688Z
owner_approval: false
completed: false
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

