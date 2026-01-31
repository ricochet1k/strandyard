---
type: review
role: reviewer-reliability
priority: medium
parent: Trlrg-specify-anchor-error-messages
blockers:
    - T4izq-define-deterministic-anchor-examples-for-every-hin
    - Tc3zv-align-exit-code-conventions-for-every-failures
blocks:
    - Trlrg-specify-anchor-error-messages
date_created: 2026-01-29T15:19:02.019232Z
date_edited: 2026-01-30T22:38:51.494406Z
owner_approval: false
completed: true
---

# Reliability review: recurrence anchor errors

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

## Artifacts
- design-docs/recurrence-anchor-error-messages-alternatives.md

## Scope
Determinism, testability, and backward compatibility of error messaging.

## Review Focus
- Stability of error strings for tests
- Failure mode coverage and determinism
- Impact on automation workflows

## Subtasks
- [x] (subtask: T4izq) Define deterministic anchor examples for --every hints
- [ ] (subtask: Tc3zv) Align exit code conventions for --every failures
