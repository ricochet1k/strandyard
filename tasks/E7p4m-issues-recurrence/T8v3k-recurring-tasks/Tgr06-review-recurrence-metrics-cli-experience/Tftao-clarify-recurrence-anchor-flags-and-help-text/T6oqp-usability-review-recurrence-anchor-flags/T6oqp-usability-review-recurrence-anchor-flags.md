---
type: review-usability
role: reviewer-usability
priority: medium
parent: Tftao-clarify-recurrence-anchor-flags-and-help-text
blockers: []
blocks:
    - Tftao-clarify-recurrence-anchor-flags-and-help-text
date_created: 2026-01-29T05:19:58.866092Z
date_edited: 2026-01-28T22:19:58.875005-07:00
owner_approval: false
completed: false
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
- design-docs/recurrence-metrics.md
- CLI.md (recurring add section)

## Scope
- User experience for selecting anchors in `memmd recurring add`
- Help text and examples for time-based vs git-based metrics

## Primary User Journeys
- Create a time-based recurring task with ISO 8601 anchor
- Create a commit-based recurring task with commit hash anchor

## Error States and Recovery
- Missing `--anchor`/anchor field for selected unit
- Invalid ISO 8601 timestamp or invalid commit hash
- Ambiguous anchor format and guidance to correct it

## Review Focus
- Help text readability and concision
- Examples that reduce user error
- Consistency with other CLI flag patterns
