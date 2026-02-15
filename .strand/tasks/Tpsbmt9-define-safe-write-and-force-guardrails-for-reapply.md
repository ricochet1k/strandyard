---
type: issue
role: architect
priority: medium
parent: T448da6-security-review-reapply-template-design
blockers: []
blocks: []
date_created: 2026-02-15T05:18:23.503274Z
date_edited: 2026-02-15T05:18:43.201415Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Define safe write and force guardrails for reapply

## Summary


## Summary
Define secure apply-mode write behavior for `strand reapply` so template reconciliation cannot cause unintended or unsafe file replacement.

## Requirements
- [ ] Specify atomic write strategy (temp file in same dir + fsync + rename) for task rewrites.
- [ ] Require refusal to overwrite non-regular files (symlink/device) at task paths.
- [ ] Define exact semantics for `--force` and `--prune`, including which safeguards are never bypassable.
- [ ] Require dry-run/diff visibility before destructive operations in non-interactive contexts.

## Deliverable
Add these guardrails and invariants to the reapply design with corresponding test expectations.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds
