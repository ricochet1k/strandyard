---
type: issue
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:38:56.93394Z
date_edited: 2026-02-05T04:38:56.93394Z
owner_approval: false
completed: false
description: ""
---

# Fix YAML serialization spacing and improve parse errors

## Summary
## Summary
Ensure YAML writer produces valid spacing and parser reports useful errors.

## Description
We encountered files with `blockers:[]` (no space) which caused the parser to crash/fail hard.

## Requirements
- Ensure `strand` commands always write `key: value` with a space.
- Improve the parser to either handle missing spaces if valid in stricter YAML, or provide a clear error message pointing to the file and line number without panicking or failing the entire list operation.
