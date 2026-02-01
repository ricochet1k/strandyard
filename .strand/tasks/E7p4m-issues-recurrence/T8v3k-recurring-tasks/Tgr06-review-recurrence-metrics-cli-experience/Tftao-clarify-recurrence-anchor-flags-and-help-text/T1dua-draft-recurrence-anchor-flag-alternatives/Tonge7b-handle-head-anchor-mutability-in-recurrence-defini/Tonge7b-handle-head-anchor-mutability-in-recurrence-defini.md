---
type: owner-decision
role: owner
priority: high
parent: T1dua-draft-recurrence-anchor-flag-alternatives
blockers: []
blocks: []
date_created: 2026-02-01T20:23:29.697792Z
date_edited: 2026-02-01T20:23:29.697792Z
owner_approval: false
completed: false
description: ""
---

# Handle HEAD anchor mutability in recurrence definitions

## Description
Decide whether HEAD references in recurrence definitions (e.g., `--every "50 commits from HEAD"`) should be:

1. **Resolved at definition time (immutable)**: The HEAD reference is immediately resolved to a specific commit hash, and all recurrences are bound to that fixed anchor.

2. **Resolved at materialization time (mutable)**: Each recurrence instance re-evaluates HEAD at the time of materialization, allowing the recurrence to move forward with the repository.

See design-docs/recurrence-anchor-flags-alternatives.md for context on the `--every` flag syntax (Alternative D was adopted).

The decision will be filled into the design document.

The decision will be filled into the design document.
