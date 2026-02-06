---
type: issue
role: triage
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T00:09:02.869032Z
date_edited: 2026-02-06T04:20:40.81738Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Allow overriding roles/templates from presets

## Summary
When a preset we depend on evolves, there is no way to pull the new roles or templates into the local storage because the CLI refuses to overwrite existing files. Provide a flag or command that can refresh roles/templates from an updated preset (overwriting local files even when they already exist) while keeping task data safe. This would let us keep roles and templates in sync with upstream tweaks.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds
