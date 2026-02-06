---
type: implement
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T04:45:04.24278Z
date_edited: 2026-02-06T04:45:04.24278Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Add strand preset refresh command

## Summary
## Context
- The preset cloning logic in `cmd/init.go` only runs during `strand init`, and once a project dir exists the command refuses to rerun, so roles/templates never refresh from updated presets.
- For downstream workflows we still want to pull template/role updates from presets without touching existing task data, which means rerunning `applyPreset` against a project while leaving `tasks/` untouched.
- The triage issue `Iz1zd` captures the need to overwrite stale role/template files safely, so the implementation should reuse the preset download/copy helpers and log what was refreshed so triage/developers can verify.

## Deliverables
- Add a `strand preset refresh <preset>` command that clones or reads a preset (supporting local dirs and git URLs) and copies `roles/` and `templates/` into the current project, overwriting files while keeping `tasks/` untouched.
- Ensure the command detects when the target project is initialized (local/global) and fails fast with a clear message when not.
- Share context notes about what files changed after running the command and run `strand repair` to keep root/free lists deterministic.
- Extend tests to cover refreshing from both local preset dirs and remote git presets, verifying that templates and roles are overwritten and that task data remains unaffected.

## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds
