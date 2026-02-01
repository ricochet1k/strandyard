---
type: review-usability
role: reviewer-usability
priority: medium
parent: Trlrg-specify-anchor-error-messages
blockers: []
blocks:
    - Trlrg-specify-anchor-error-messages
date_created: 2026-01-29T15:19:01.830901Z
date_edited: 2026-02-01T20:16:41.149143Z
owner_approval: false
completed: true
description: ""
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

## Completion Report
Usability review complete. Concerns captured as subtasks: (1) Anchor flag inconsistency between CLI.md (--anchor) and design doc examples (--every), (2) Missing anchor default behavior unclear - design doc says defaults to 'now'/'HEAD' but Alternative C (which had defaults) was not adopted, (3) Ambiguous anchor type error listed but no concrete message/hint defined, (4) 'after now' vs 'from now' distinction needs clearer documentation. Overall error message format (Alternative B) with structured reason + hint line is actionable and consistent.
