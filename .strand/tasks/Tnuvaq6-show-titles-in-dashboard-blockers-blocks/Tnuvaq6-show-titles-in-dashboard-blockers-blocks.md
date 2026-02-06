---
type: implement
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T00:28:47.294619Z
date_edited: 2026-02-06T04:31:07.345815Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Show titles in dashboard blockers/blocks

## Summary
- Resolve each blocker/block ID by looking up its title from the cached task list so the relationship list shows `short_id — title` instead of bare IDs.
- Render those entries as links/buttons that load the referenced task, update the browser route with the current project, the selected task, and a `relationship` flag (blocked-by vs blocking), and surface the originating context so the linked task still shows how it relates back to the previous view.
- Keep the blocked-by relationship visible in the editor (e.g., via a label or highlight) so users never lose the context that led them into the new task.

## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds
