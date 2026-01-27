{{/*
  Task template (for leaf/implementable tasks).
  Notes:
  - `ID` and `Parent` are derived from the filesystem path and set by the CLI; do not include them here.
  - Blockers and Blocks are managed by the CLI and should not be populated in templates.
*/}}
# {{ .Title }}

## Role
{{ .Role }}

## Track
{{ .Track }}

## Context
Provide links to relevant design documents, diagrams, and decision records.

## TODOs
- [ ] (role: {{ .Role }}) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
Use subtasks for work that should be tracked as separate task directories. List them here when useful:
- tasks/{{ .SuggestedSubtaskDir }}/task.md — short description of subtask

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.
