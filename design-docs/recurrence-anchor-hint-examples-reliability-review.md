---
role: reviewer-reliability
priority: medium
---

# Review: Canonical --every hint examples (Reliability)

## Artifacts
- design-docs/recurrence-anchor-hint-examples.md
- design-docs/recurrence-anchor-error-messages-alternatives.md
- design-docs/recurrence-anchor-error-messages-reliability-review.md

## Task Selection Output
```
Your role is reviewer-reliability. Here's the description of that role:

# Reliability Reviewer

## Role
Reliability Reviewer — review designs and plans for operational reliability, SLOs, and runbook needs.

## Responsibilities
- Evaluate operational impact and failure modes.
- Suggest SLOs, monitoring, and runbook items.
- Do not wait for interactive responses; capture concerns as tasks.
- Use `templates/review.md` unless a more specific template applies.
- Avoid editing review tasks to record outcomes; file new tasks for concerns or open questions. Record decisions and final rationale in design docs.

## Escalation
- For obvious concerns, create a new subtask under the current task and assign it to Architect for technical/design documents or Designer for UX/documentation artifacts.
- For decisions needing maintainer input, create a new subtask assigned to the Owner role and note that the decision should be recorded in design docs.

---

Your task is Tm2sq-review-canonical-every-hint-examples-implementatio. Here's the description of that task:

---
type: review
role: reviewer-reliability
priority: medium
parent: T4izq-define-deterministic-anchor-examples-for-every-hin
blockers: []
blocks:
    - T4izq-define-deterministic-anchor-examples-for-every-hin
    - Tv4cw-implement-deterministic-every-hint-examples
date_created: 2026-01-29T19:24:33.271976Z
date_edited: 2026-01-29T12:24:46.400087-07:00
owner_approval: false
completed: false
---

# Review canonical --every hint examples (implementation)

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


## Summary
Review the canonical --every hint examples for determinism and automation stability.

## Tasks
- [ ] Confirm examples in design-docs/recurrence-anchor-hint-examples.md are stable and deterministic
- [ ] Validate hint strings are suitable for tests and error output contracts

## Acceptance Criteria
- Review notes recorded and actionable follow-ups captured as tasks
```

## Scope
Review deterministic hint examples for `strand recurring add --every` anchor parsing.
Implementation details for parsing or formatting are out of scope.

## Review Focus
- Determinism and stability of example strings for automation
- Suitability for error output contracts and test fixtures
- Avoidance of locale/time-dependent formatting

## Findings
- Examples are stable, time-independent constants that avoid runtime values.
- Date anchors are explicitly UTC and human-friendly; ISO 8601 is limited to validation contexts.
- Commit anchor uses `HEAD` with a fixed placeholder for hash-only tests, avoiding repo-dependent values.

## Concerns captured as subtasks
- None.
