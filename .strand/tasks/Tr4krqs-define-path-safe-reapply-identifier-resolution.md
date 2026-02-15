---
type: issue
role: architect
priority: medium
parent: T448da6-security-review-reapply-template-design
blockers: []
blocks: []
date_created: 2026-02-15T05:18:19.190616Z
date_edited: 2026-02-15T05:18:43.201412Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Define path-safe reapply identifier resolution

## Summary


## Summary
Define the security requirements for resolving `strand reapply` task IDs and template names so user input cannot escape the repository task/template roots.

## Requirements
- [ ] Specify canonical validation for `<task-id>` and `--template` values before any filesystem path construction.
- [ ] Require path normalization and repository-root containment checks after join/clean.
- [ ] Require rejection of symlink traversal for task/template lookups used by reapply.
- [ ] Define expected error messages for invalid identifiers versus missing resources.

## Deliverable
Update the reapply design doc with explicit path-safety invariants and test cases.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds
