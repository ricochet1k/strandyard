---
role: reviewer-usability
priority: medium
---

# Usability Review: Canonical --every Hint Examples (Docs)

## Artifacts
- design-docs/recurrence-anchor-hint-examples.md
- design-docs/recurrence-anchor-error-messages-alternatives.md
- CLI.md (recurring add section; contextual reference)
- Task selection output (go run . next):
```text
Your role is reviewer-usability. Here's the description of that role:

# Usability Reviewer

## Role
Usability Reviewer — review designs and plans for human-facing usability and clarity.

## Responsibilities
- Evaluate UX flows, documentation clarity, and user-facing error handling.
- Do not wait for interactive responses; capture concerns as tasks.
- Use `templates/review-usability.md` for usability reviews.
- Avoid editing review tasks to record outcomes; file new tasks for concerns or open questions. Record decisions and final rationale in design docs.

## Escalation
- For obvious concerns, create a new subtask under the current task and assign it to Architect for technical/design documents or Designer for UX/documentation artifacts.
- For decisions needing maintainer input, create a new subtask assigned to the Owner role and note that the decision should be recorded in design docs.

---

Your task is Tsyeo-review-canonical-every-hint-examples-docs. Here's the description of that task:

---
type: review
role: reviewer-usability
priority: medium
parent: T4izq-define-deterministic-anchor-examples-for-every-hin
blockers: []
blocks:
    - T4izq-define-deterministic-anchor-examples-for-every-hin
    - Tm6qi-document-canonical-every-hint-examples
date_created: 2026-01-29T19:24:40.775932Z
date_edited: 2026-01-29T12:24:50.68882-07:00
owner_approval: false
completed: false
---

# Review canonical --every hint examples (docs)

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
Review documentation-facing examples for clarity and consistency with error message contract.

## Tasks
- [ ] Validate wording and clarity of canonical examples in docs
- [ ] Check consistency with CLI.md and design docs

## Acceptance Criteria
- Usability review feedback captured as tasks or approval noted
```

## Scope
Review documentation-facing canonical hint examples for `memmd recurring add --every`, focusing on clarity of examples and recovery hints. Implementation behavior or parser details are out of scope.

## Primary User Journeys
- Read docs to understand how to format `--every` without anchors.
- Recover from invalid anchor errors using hint examples.

## Error States and Recovery
- Missing anchor with default behavior; user needs to understand what the default means per metric.
- Malformed anchor value; user needs example format.
- Unit/anchor mismatch; user needs guidance to correct anchor type.

## Review Focus
- Clarity of example strings and explanations of defaults.
- Mapping between metric type and anchor expectation.
- Consistency with error message contract and hint line phrasing.

## Findings
- Canonical examples are short and deterministic; human-friendly date anchors reduce formatting ambiguity.
- Default-anchor examples are appropriate for most hints, but commit-based metrics would be clearer if the docs/hints state the default anchor as `HEAD` rather than generic "now" wording.

## Escalation
- Tnetq-clarify-default-anchor-wording-for-commit-metrics (role: designer).

## Checklist
- [x] Artifacts and scope listed.
- [x] Journeys and error handling documented.
- [x] Review focus defined.
- [x] Concerns captured as subtasks.
- [x] Decision items deferred to Owner as separate subtasks when needed.
