---
type: review-usability
role: reviewer-usability
priority: medium
parent: Trlrg-specify-anchor-error-messages
blockers: []
blocks:
    - Trlrg-specify-anchor-error-messages
date_created: 2026-01-29T15:19:01.830901Z
date_edited: 2026-01-29T08:19:01.838171-07:00
owner_approval: false
completed: false
---

# Usability review: recurrence anchor errors

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
- design-docs/recurrence-anchor-error-messages-alternatives.md
- CLI.md (recurring add section)

## Scope
User-facing error strings and recovery hints for missing/malformed anchors.

## Primary User Journeys
- Add a recurring task with time-based metrics
- Add a recurring task with commit-based metrics

## Error States and Recovery
- Missing anchor
- Malformed date anchor
- Malformed commit anchor
- Unit/anchor mismatch
- Ambiguous anchor type

## Review Focus
- Clarity and actionability of hints
- Consistency with CLI language patterns
- Potential confusion around defaults
