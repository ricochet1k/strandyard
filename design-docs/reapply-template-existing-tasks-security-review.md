# Security Review: Reapply templates to existing tasks

## Artifacts
- `design-docs/reapply-template-existing-tasks-alternatives.md`

## Scope
Security-focused review for template reapplication behavior, with emphasis on untrusted input boundaries, path safety, and destructive-operation safeguards.

## Task selection output
```text
Your role is reviewer-security. Here's the description of that role:

---
description: "Reviews designs and plans for security concerns."
---

# Security Reviewer

## Role
Security Reviewer — review designs and plans for security concerns, threat models and mitigations.

## Responsibilities
- Identify real security vulnerabilities, not defensive paranoia against your own codebase.
- Focus on untrusted input boundaries: where user input or external data enters the system and whether it's properly validated.
- Evaluate threat models, data handling, and access control implications.
- Do not wait for interactive responses; capture concerns as tasks.
- Avoid editing review tasks to record outcomes; file new tasks for concerns or open questions. Record decisions and final rationale in design docs.

## Review Focus
- **Untrusted input:** CLI arguments, file contents, network data. Is it validated before use?
- **Trust boundaries:** Where does data cross from untrusted to trusted? Are there gaps?
- **Real threats:** Only flag actual vulnerabilities (injection, auth bypass, data leaks). Skip hypothetical defenses against competent code.
- **Not in scope:** Defensive coding practices within trusted code, internal API robustness, or defensive assumptions that other code is malicious.

## Escalation
- For obvious concerns, create a new subtask under the current task and assign it to Architect for technical/design documents or Designer for UX/documentation artifacts.
- For decisions needing maintainer input, create a new subtask assigned to the Owner role and note that the decision should be recorded in design docs.

## Workflow
1. Review the design or plan for security.
2. Capture any concerns as follow-up tasks.
3. Mark the review task as completed: `strand complete <task-id> "Security review complete. Concerns: ..."`

---

Ancestors:
  Tk259ie: Design: Allow reapplying templates to existing tasks
  I8lew: Allow reapplying templates to existing tasks


Your task is T448da6-security-review-reapply-template-design. Here's the description of that task:

---
type: review
role: reviewer-security
priority: medium
parent: Tk259ie-design-allow-reapplying-templates-to-existing-task
blockers: []
blocks:
    - Tk259ie-design-allow-reapplying-templates-to-existing-task
date_created: 2026-02-15T05:16:43.102152Z
date_edited: 2026-02-15T05:17:38.156567Z
owner_approval: false
completed: false
status: in_progress
description: ""
---

# Security review: reapply-template design

## Summary


## Summary
Review template/task identifier handling, path-safety assumptions, and destructive-operation guardrails for the proposed reapply command.

## Instructions
Focus on path traversal resistance, safe write strategy, and conflict/force safeguards.

## Instructions
Delegate concerns to the relevant role via subtasks. Mark complete when review is finished.
```

## Findings
- The recommended command shape (`strand reapply`) is safer than overloading `edit`, but the design does not yet define mandatory input-to-path security invariants for task/template identifiers.
- Current merge/apply guidance does not yet require an explicit atomic-write strategy or symlink refusal policy for file replacement.
- `--force`/`--prune` are proposed but missing a hard list of safeguards that remain non-bypassable.

## Concerns captured as subtasks
- `Tr4krqs`: Define path-safe reapply identifier resolution (architect)
- `Tpsbmt9`: Define safe write and force guardrails for reapply (architect)

## Decision
- Decision: deferred pending architect follow-up tasks above.
