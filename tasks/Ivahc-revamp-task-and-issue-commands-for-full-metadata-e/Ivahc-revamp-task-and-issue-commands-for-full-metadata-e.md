---
kind: issue
role: triage
priority: medium
parent:
blockers: []
blocks: []
date_created: 2026-01-27T23:25:07Z
date_edited: 2026-01-27T23:25:07Z
owner_approval: false
completed: false
---

# Revamp task and issue commands for full metadata editing

## Summary
Current task/issue creation flows require editing Markdown files to fill in fields or fix mistakes. We need CLI commands that can create and update all metadata and content fields for tasks and issues, so end-to-end task management can be done without opening the Markdown files directly.

## Steps to Reproduce
1. Create a task or issue via CLI (e.g., `memmd add` or `memmd issue add`).
2. Attempt to set or edit metadata fields like title, priority, blockers, or body content via CLI.
3. Notice that updates require manual edits to the Markdown file.

## Expected Result
CLI supports creating and editing tasks/issues, including full frontmatter and body content, without manual file edits.

## Actual Result
CLI can create tasks/issues, but editing metadata or body content requires manual Markdown edits.

## Acceptance Criteria
- CLI provides commands to create tasks/issues with all metadata fields supplied via flags or prompts.
- CLI provides commands to edit metadata and body content for existing tasks/issues.
- No manual editing of task Markdown files is required for standard workflows.
- `memmd validate` passes after creating or editing tasks/issues.
