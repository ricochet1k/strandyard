---
type: review
role: master-reviewer
priority: medium
parent: Trlrg-specify-anchor-error-messages
blockers:
    - T14az-define-error-message-format-contract-for-every-anc
    - T4kdz-decide-recurrence-anchor-error-message-strategy-a
blocks:
    - Trlrg-specify-anchor-error-messages
date_created: 2026-01-29T15:19:01.741649Z
date_edited: 2026-01-31T17:29:31.078613Z
owner_approval: false
completed: true
---

# Review alternatives: recurrence anchor error messages

## Artifacts
List the documents, designs, or code paths under review.

## Scope
Clarify what is in and out of scope for this review.

## Review Focus
List the specific areas to evaluate (e.g., usability, API ergonomics, error handling).

## Escalation
Create new tasks for concerns or deferred decisions instead of editing this task.

## Checklist
- [ ] Artifacts and scope listed.
- [ ] Review focus defined.
- [ ] Concerns captured as subtasks.
- [ ] Decision items deferred to Owner as separate subtasks when needed.

## Artifacts
- design-docs/recurrence-anchor-error-messages-alternatives.md

## Scope
Error message strategy for recurrence anchors in `strand recurring add --every`.

## Review Focus
- Decision framing and tradeoffs
- Alignment with project principles
- Testability and determinism

## Subtasks
- [x] (subtask: T14az) Define error message format contract for --every anchor parsing
- [x] (subtask: T4kdz) Decide recurrence anchor error message strategy (A/B/C)
