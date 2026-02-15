---
type: review
role: reviewer-security
priority: medium
parent: Tk259ie-design-allow-reapplying-templates-to-existing-task
blockers: []
blocks: []
date_created: 2026-02-15T05:16:43.102152Z
date_edited: 2026-02-15T05:18:43.20133Z
owner_approval: false
completed: true
status: done
description: ""
---

# Security review: reapply-template design

## Summary


## Summary
Review template/task identifier handling, path-safety assumptions, and destructive-operation guardrails for the proposed reapply command.

## Instructions
Focus on path traversal resistance, safe write strategy, and conflict/force safeguards.

## Instructions
Delegate concerns to the relevant role via subtasks. Mark complete when review is finished.

## Subtasks
- [ ] (subtask: Tr4krqs) Define path-safe reapply identifier resolution
- [ ] (subtask: Tpsbmt9) Define safe write and force guardrails for reapply

## Completion Report
Security review complete. Concerns: added Tr4krqs for identifier/path-safety invariants and Tpsbmt9 for atomic write + non-bypassable force/prune guardrails; review notes captured in design-docs/reapply-template-existing-tasks-security-review.md.
