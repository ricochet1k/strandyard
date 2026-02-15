# Review: Reapply templates to existing tasks (Alternatives)

## Artifacts
- `design-docs/reapply-template-existing-tasks-alternatives.md`

## Scope
Master review coordination for the reapply-template design alternatives.

## Review Focus
- Decision framing and tradeoffs for command shape
- Safety and determinism of merge/apply behavior
- Delegation to specialized reviewers and owner decision capture

## Task selection output
```text
Your role is master-reviewer. Here's the description of that role:

---
description: "Coordinates specialized reviewers and consolidates feedback."
---

# Reviewer (master)

## Role
Master Reviewer — central review role that coordinates specialized reviewers.

## Responsibilities
- Accept review requests and delegate to specialized reviewers (Reliability, Security, Usability, etc.) by adding TODOs.
- Consolidate feedback and return a single review verdict to the requestor.
- Do not wait for interactive responses; capture concerns as subtasks.
- Avoid editing review tasks to record outcomes; file new tasks for concerns or open questions. Record decisions and final rationale in design docs.

## Escalation
- For obvious concerns, create a new subtask under the current task and assign it to the appropriate reviewer role.
- For decisions that require maintainer input, create a new `owner-decision` subtask.

## Workflow
1. Receive review request.
2. Delegate to specialized reviewers if needed.
3. Consolidate feedback.
4. Mark the review task as completed: `strand complete <task-id> "Verdict: ..."`

---

Ancestors:
  Tk259ie: Design: Allow reapplying templates to existing tasks
  I8lew: Allow reapplying templates to existing tasks


Your task is Te0b2ww-review-reapply-template-design-alternatives. Here's the description of that task:

---
type: review
role: master-reviewer
priority: medium
parent: Tk259ie-design-allow-reapplying-templates-to-existing-task
blockers: []
blocks:
    - Tk259ie-design-allow-reapplying-templates-to-existing-task
date_created: 2026-02-15T05:14:56.889849Z
date_edited: 2026-02-15T05:15:46.764182Z
owner_approval: false
completed: false
status: in_progress
description: ""
---

# Review reapply-template design alternatives

## Summary


## Instructions
Delegate concerns to the relevant role via subtasks. Mark complete when review is finished.
```

## Findings
- Alternative B (`strand reapply`) keeps structural template updates explicit and separates risky operations from ordinary task edits.
- Proposed merge rules are mostly deterministic and preserve user-authored task state by default.
- The design correctly defers command-direction approval to owner decision.

## Concerns captured as subtasks
- Existing subtask: `Tmzzb0o` (usability review)
- Existing subtask: `Tfo4cg6` (reliability review)
- Existing subtask: `T39dxo3` (owner decision)
- Added subtask: `T448da6` (security review) for path-safety and force/overwrite safeguards in template reapply flows.

## Decision
- Decision: deferred to Owner (`T39dxo3`).
