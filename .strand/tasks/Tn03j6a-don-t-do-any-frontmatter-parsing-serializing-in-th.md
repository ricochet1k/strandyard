---
type: implement
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T04:42:45.11879Z
date_edited: 2026-02-06T06:43:47.018248Z
owner_approval: false
completed: true
status: done
description: ""
---

# Don't do any frontmatter parsing/serializing in the frontend, do it all in the backend, somewhere in TasksDB

## Summary


## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds

## Completion Report
Moved frontmatter parsing and serialization from the frontend to the backend. Added structured API endpoints for Roles and Templates (/api/roles, /api/role, /api/templates, /api/template) and updated the dashboard to use them. Deleted the redundant frontmatter utility file in the frontend.
