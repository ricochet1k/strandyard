---
type: issue
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-31T17:29:31.04212Z
date_edited: 2026-01-31T17:29:31.04212Z
owner_approval: false
completed: false
---

# Add auth + file write safeguards for dashboard

## Summary
The web dashboard API has no authentication or authorization controls. Anyone who can access the dashboard port (default 8686) can read and write all task/role/template files across all projects. The `/api/file` PUT endpoint allows unrestricted file writes.

## Context
- The `strand web` command starts a dashboard server at `http://localhost:8686`
- API endpoints include `/api/file` with PUT support for writing files directly to disk
- CORS is set to allow all origins (`Access-Control-Allow-Origin: *`)
- No authentication mechanism exists
- No authorization checks for read/write operations
- File writes are validated only for path containment (within tasks/roles/templates) but not for content or permissions

## Impact
**Severity: High**

- Anyone on the local network (or who can access the port via port forwarding) can modify task data, roles, and templates
- No audit trail for who made changes
- Potential for data corruption or malicious edits
- In shared environments, users can interfere with each other's work

## Acceptance Criteria
- Add authentication mechanism to web dashboard (e.g., password/token flag to `strand web`)
- Add CORS configuration to restrict origins in production
- Add read-only mode option for dashboard
- Ensure file write operations validate content type and structure
- Add authorization checks for write operations
- Document security model and usage

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.
