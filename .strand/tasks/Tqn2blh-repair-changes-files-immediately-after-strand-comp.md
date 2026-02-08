---
type: issue
role: developer
priority: high
parent: Ti6zj-taskdb-api-design-review
blockers: []
blocks:
    - Ti6zj-taskdb-api-design-review
date_created: 2026-02-08T04:10:50.874381Z
date_edited: 2026-02-08T04:11:10.516498Z
owner_approval: false
completed: false
status: ""
description: ""
---

# repair changes files immediately after strand complete

## Repro
1. Run `go run ./cmd/strand complete T0q5n-review-blockers-go-relationship-management "<report>"`.
2. Immediately run `go run ./cmd/strand repair`.

## Observed
- `strand complete` reports success.
- Follow-up `strand repair` reports additional changes: `Repaired 2 tasks`.

## Expected
- `strand complete` should leave the repo in the same repaired state for master lists/relationships.
- Immediate `strand repair` should report `Repaired 0 tasks`.

## Affected task IDs
- T0q5n-review-blockers-go-relationship-management
- Ti6zj-taskdb-api-design-review

## Notes
Policy requires treating post-complete repair deltas as a bug and filing an issue.
