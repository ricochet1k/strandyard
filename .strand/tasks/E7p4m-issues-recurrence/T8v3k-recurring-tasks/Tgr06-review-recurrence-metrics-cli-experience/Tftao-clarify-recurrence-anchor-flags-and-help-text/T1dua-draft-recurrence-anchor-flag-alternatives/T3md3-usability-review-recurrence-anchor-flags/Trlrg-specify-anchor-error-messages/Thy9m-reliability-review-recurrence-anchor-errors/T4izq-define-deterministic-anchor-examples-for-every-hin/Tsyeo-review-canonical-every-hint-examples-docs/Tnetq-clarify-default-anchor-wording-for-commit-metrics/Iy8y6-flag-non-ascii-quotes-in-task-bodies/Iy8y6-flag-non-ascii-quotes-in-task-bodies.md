---
type: issue
role: triage
priority: medium
parent: Tnetq-clarify-default-anchor-wording-for-commit-metrics
blockers: []
blocks:
    - Tnetq-clarify-default-anchor-wording-for-commit-metrics
date_created: 2026-01-29T20:01:12.782966Z
date_edited: 2026-01-29T20:04:13.264215Z
owner_approval: false
completed: true
---

# Flag non-ASCII quotes in task bodies

## Summary
## Summary
Manual edit needed to replace a smart quote in a review task body. Consider adding a lint/repair rule or add command validation to flag non-ASCII quotes so manual edits are avoidable.

## Repro Steps
1. Run `strand add review-usability "Usability review: commit-metric default anchor wording" --parent Tocc0-draft-alternatives-for-commit-metric-default-ancho` with body text containing smart quotes.
2. Observe the resulting task file contains non-ASCII quotes.
3. Manual edit was required to normalize the quote to ASCII.

## Logs
- Manual edit after task creation; `go run . repair` succeeded.

## Affected Task IDs
- Tvu5e-usability-review-commit-metric-default-anchor-word

## Context
Provide relevant logs, links, and environment details.

## Impact
Describe severity and who/what is affected.

## Acceptance Criteria
- Define what fixes or mitigations are required.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.
