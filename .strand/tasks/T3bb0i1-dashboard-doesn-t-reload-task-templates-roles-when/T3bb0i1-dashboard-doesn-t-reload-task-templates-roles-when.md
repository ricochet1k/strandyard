---
type: fix
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T00:23:04.832961Z
date_edited: 2026-02-06T04:18:54.309978Z
owner_approval: false
completed: true
status: done
description: ""
---

# Dashboard doesn't reload task templates/roles when project changes

## Summary


## Acceptance Criteria
- Bug still exists
- Bug is fixed and verified locally
- Tests pass
- Build succeeds

## Completion Report
Switched dashboard fetches to createResource with project-aware roles/templates parsing, updated dashboard E2E coverage for CodeMirror and project switching, and stabilized transition test timing. Ran Playwright, go test ./..., and go build ./...
