---
type: task
role: designer
priority: medium
parent: Iff91-add-cli-support-for-task-storage-migration
blockers: []
blocks: []
date_created: 2026-02-05T01:08:56.868412Z
date_edited: 2026-02-05T01:09:42.332022Z
owner_approval: false
completed: true
description: ""
---

# New Task: Design task storage migration

## Description
Design a CLI command or process to migrate task storage between "global" and "local" modes.
The migration should handle:
- Moving `tasks/`, `roles/`, and `templates/` to the new location (e.g., from top-level to `.strand/`).
- Updating the global configuration mapping.
- Handling Git tracking: if migrating to local storage, should the files be `git add`-ed? If migrating from local, should they be removed from Git?
- Updating master lists paths.

Deliverable: A design document in `design-docs/` outlining the migration process and safety checks.

Decide which task template would best fit this task and re-add it with that template and the same parent.

## Completion Report
Triage complete. Re-added as designer task T6y8964 under the same parent.
