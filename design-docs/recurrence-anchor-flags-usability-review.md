---
role: reviewer-usability
priority: medium
---

# Usability Review: Recurrence Anchor Flags

## Artifacts
- design-docs/recurrence-anchor-flags-alternatives.md
- CLI.md (recurring add section)
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
- Avoid editing review tasks to record outcomes; file new tasks for concerns or decisions.

## Escalation
- For obvious concerns, create a new subtask under the current task and assign it to Architect for technical/design documents or Designer for UX/documentation artifacts.
- For decisions needing maintainer input, create a new subtask assigned to the Owner role.

---

Your task is T3md3-usability-review-recurrence-anchor-flags. Here's the description of that task:

---
type: review-usability
role: reviewer-usability
priority: medium
parent: T1dua-draft-recurrence-anchor-flag-alternatives
blockers: []
blocks:
    - T1dua-draft-recurrence-anchor-flag-alternatives
date_created: 2026-01-29T05:45:07.908116Z
date_edited: 2026-01-28T22:45:07.917246-07:00
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
- CLI.md (recurring add section)

## Scope
Evaluate how users discover and apply anchor flags for recurring definitions.

## Primary User Journeys
- Define recurring tasks with time-based units.
- Define recurring tasks with git-based units.
- Interpret help text and resolve anchor errors.

## Error States and Recovery
- Missing or malformed anchor value.
- Unit/anchor mismatch.
- Ambiguous anchor type.

## Review Focus
- Flag naming and help text clarity.
- Minimizing user confusion across units.
- Error messaging expectations.
```

## Scope
Evaluate how users discover the correct anchor value and recover from errors when using `memmd recurring add` for time-based and git-based units.

## Primary User Journeys
- Define a recurring task using a time-based unit (days/weeks/months).
- Define a recurring task using a git-based unit (commits/lines_changed).
- Resolve errors after providing an anchor of the wrong type or format.

## Error States and Recovery
- Missing anchor: show required format and a concrete example for the selected unit.
- Malformed anchor: show expected format and hint about where to find the value.
- Unit/anchor mismatch: explain mismatch and suggest the correct flag or format.
- Ambiguous anchor type: explain how to disambiguate (explicit flag or help entry).

## Review Focus
- Flag naming and help text clarity across units.
- Discoverability of anchor format for users who skim help output.
- Error messages that provide recovery actions and examples.

## Findings
- The current CLI help (`CLI.md`) describes `--anchor` as date or commit hash but does not map formats to specific units; this becomes more confusing as units expand (tasks_completed, lines_changed).
- Alternative A is minimal but relies heavily on help text; without a per-unit mapping and examples, users may select the wrong anchor format.
- Alternative B improves clarity, but migration needs explicit guidance and aliasing to avoid user confusion.
- Alternative C adds a second axis (`--anchor-type`) and is likely to increase cognitive load unless errors guide users to the correct combination.

## Decision
- Decision: deferred to Owner (see Te6hk-decide-anchor-flag-approach-a-b-c).

## Escalation
- Tusef-define-anchor-help-text-and-examples (role: designer).
- Trlrg-specify-anchor-error-messages (role: designer).
- Te6hk-decide-anchor-flag-approach-a-b-c (role: owner).

## Checklist
- [x] Artifacts and scope listed.
- [x] Journeys and error handling documented.
- [x] Review focus defined.
- [x] Concerns captured as subtasks.
- [x] Decision items deferred to Owner as separate subtasks when needed.
