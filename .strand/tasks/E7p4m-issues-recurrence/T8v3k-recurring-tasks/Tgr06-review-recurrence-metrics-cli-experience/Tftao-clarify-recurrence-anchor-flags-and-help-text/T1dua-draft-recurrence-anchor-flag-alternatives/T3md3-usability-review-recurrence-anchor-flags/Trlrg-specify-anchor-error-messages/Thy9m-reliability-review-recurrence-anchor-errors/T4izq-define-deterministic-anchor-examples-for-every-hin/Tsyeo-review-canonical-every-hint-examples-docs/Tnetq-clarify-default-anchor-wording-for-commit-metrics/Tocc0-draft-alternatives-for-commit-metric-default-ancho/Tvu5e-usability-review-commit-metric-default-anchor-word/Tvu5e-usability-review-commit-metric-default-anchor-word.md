---
type: review-usability
role: reviewer-usability
priority: medium
parent: Tocc0-draft-alternatives-for-commit-metric-default-ancho
blockers:
    - T9eo473-decision-needed-clarify-invalid-head-behavior-for
blocks:
    - Tocc0-draft-alternatives-for-commit-metric-default-ancho
date_created: 2026-01-29T20:00:34.4364Z
date_edited: 2026-02-01T09:04:50.758779Z
owner_approval: false
completed: true
description: ""
---

# Usability review: commit-metric default anchor wording

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
Create new tasks for concerns or open questions instead of editing this task. Record decisions and final rationale in design docs.

## Checklist
- [ ] Artifacts and scope listed.
- [ ] Journeys and error handling documented.
- [ ] Review focus defined.
- [ ] Concerns captured as subtasks.
- [ ] Decision items deferred to Owner as separate subtasks when needed.

## Artifacts
- design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-alternatives.md

## Scope
User-facing wording for default anchors in commit-based metrics, including hint examples and docs.

## Primary User Journeys
- Read docs to learn default anchor behavior for commits/lines_changed.
- Recover from invalid anchor errors using hint examples.

## Error States and Recovery
- Missing anchor for commit metrics and unclear defaults.
- Misinterpreting "now" as time-based for commit metrics.

## Review Focus
- Clarity and ambiguity reduction in wording
- Consistency with other recurrence metric docs
- Brevity vs. explicitness in hints

## Subtasks
- [x] (subtask: T9eo473) Decision needed: Clarify invalid HEAD behavior for commit metrics
