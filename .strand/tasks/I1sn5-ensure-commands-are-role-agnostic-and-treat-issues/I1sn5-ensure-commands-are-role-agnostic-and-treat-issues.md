---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T23:31:16Z
date_edited: 2026-01-28T05:07:24.550261Z
owner_approval: false
completed: true
---

# Ensure commands are role-agnostic and treat issues as tasks

## Summary
The CLI commands treat roles and types inconsistently. The `next` command has role-specific behavior (shows role files, filters by role) while issues and tasks are treated differently based on their type prefix rather than their actual nature.

## Steps to Reproduce
1. Run `strand next` - shows role file and filters by role
2. Run `strand add issue --help` - defaults to "triage" role 
3. Run `strand add leaf --help` - defaults to "developer" role
4. Notice issues get "I" prefix vs "T" prefix for tasks

## Expected Result
Commands should be role-agnostic:
- `next` should show tasks regardless of role without special role file display
- Issues should be treated as regular tasks, just with different metadata/templates
- Role filtering should be optional, not the default behavior

## Actual Result
- `next` command treats roles specially, showing role files and filtering behavior
- Issues vs tasks have different default roles baked into templates
- Role-specific behavior makes commands less flexible

## Acceptance Criteria
- `next` command shows tasks without role file display by default
- Role filtering remains available as a flag but isn't central to behavior  
- Issues and tasks are treated uniformly in commands
- Templates can define default roles but commands don't enforce role-specific behavior
