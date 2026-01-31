---
type: review-usability
role: reviewer-usability
priority: medium
parent: T1dua-draft-recurrence-anchor-flag-alternatives
blockers:
    - Te6hk-decide-anchor-flag-approach-a-b-c
    - Trlrg-specify-anchor-error-messages
    - Tusef-define-anchor-help-text-and-examples
blocks:
    - T1dua-draft-recurrence-anchor-flag-alternatives
date_created: 2026-01-29T05:45:07.908116Z
date_edited: 2026-01-31T04:41:33.14936Z
owner_approval: false
completed: true
---

# Usability review: recurrence anchor flags

## Artifacts
List the documents, designs, or code paths under review.

## Scope
Clarify what is in and out of scope for this review.

## Primary User Journeys
Describe the key user flows covered by this review.

## Error States and Recovery
List expected errors and recovery paths.

## Review Focus
List the specific usability areas to evaluate.

## Escalation
Create new tasks for concerns or deferred decisions instead of editing this task.

## Checklist
- [ ] Artifacts and scope listed.
- [ ] Journeys and error handling documented.
- [ ] Review focus defined.
- [ ] Concerns captured as subtasks.
- [ ] Decision items deferred to Owner as separate subtasks when needed.

## Artifacts
- design-docs/recurrence-anchor-flags-alternatives.md
- CLI.md (recurring add section)

## Scope
Evaluate how users discover and apply anchor flags for recurring definitions.

## Primary User Journeys
- Define recurring tasks with time-based units
- Define recurring tasks with git-based units
- Interpret help text and resolve anchor errors

## Error States and Recovery
- Missing or malformed anchor value
- Unit/anchor mismatch
- Ambiguous anchor type

## Review Focus
- Flag naming and help text clarity
- Minimizing user confusion across units
- Error messaging expectations

## Subtasks
- [x] (subtask: Te6hk) Decide anchor flag approach (A/B/C)
- [ ] (subtask: Trlrg) Specify anchor error messages
- [ ] (subtask: Tusef) Define anchor help text and examples
