---
type: task
role: designer
priority: medium
parent: I8lew-allow-reapplying-templates-to-existing-tasks
blockers:
    - T39dxo3-owner-decision-reapply-template-command-direction
    - Tfo4cg6-reliability-review-reapply-template-design
    - Tmzzb0o-usability-review-reapply-template-design
blocks: []
date_created: 2026-02-06T04:48:43.012277Z
date_edited: 2026-02-15T05:18:43.201406Z
owner_approval: false
completed: false
status: in_progress
description: ""
---

# Design: Allow reapplying templates to existing tasks

## Summary


## Summary
Design the mechanism for reapplying templates to existing tasks. 
Consider:
- How to detect changes in templates.
- How to merge template structure with existing task content without data loss.
- Command syntax (e.g., `strand edit --template <type>` or a new command `strand reapply`).
- Handling of frontmatter vs body.
- Dry run and interactive merge options.

## Acceptance Criteria
- Design document created in `design-docs/`.
- Alternatives considered and pros/cons listed.
- Implementation plan (epics/tasks) defined.

## Instructions
Decide which task template would best fit this task and re-add it with that template and the same parent.

## Subtasks
- [x] (subtask: Te0b2ww) Review reapply-template design alternatives
  Verdict: coordination complete. Confirmed existing usability/reliability/owner-decision subtasks and added security review subtask T448da6 for path-safety and force-guard concerns; final command-direction decision remains deferred to owner task T39dxo3.
- [ ] (subtask: Tmzzb0o) Usability review: reapply-template design
- [ ] (subtask: Tfo4cg6) Reliability review: reapply-template design
- [ ] (subtask: T39dxo3) Owner decision: reapply-template command direction
- [x] (subtask: T448da6) Security review: reapply-template design
  Security review complete. Concerns: added Tr4krqs for identifier/path-safety invariants and Tpsbmt9 for atomic write + non-bypassable force/prune guardrails; review notes captured in design-docs/reapply-template-existing-tasks-security-review.md.
