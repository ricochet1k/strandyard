---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T00:08:48.216444Z
date_edited: 2026-02-06T00:29:20.564126Z
owner_approval: false
completed: true
status: done
description: ""
---

# Show titles in dashboard blockers/blocks

## Summary
Dashboard relationship lists currently only show IDs, which makes it hard to know what the blockers or blocks actually are. Display the task titles and make each entry a link so users can jump directly to the related task. The dashboard should place the project and task into the route (and surface the blocked-by relationship) when a user clicks a blocker/blocks entry, keeping the blocked task context visible.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds

## Completion Report
Verified that the dashboard still renders blockers/blocks as bare IDs in the editor and the relationships are not linked; created Tnuvaq6 to render titles with links, push the project+task into the route, and keep the blocked-by context visible.
