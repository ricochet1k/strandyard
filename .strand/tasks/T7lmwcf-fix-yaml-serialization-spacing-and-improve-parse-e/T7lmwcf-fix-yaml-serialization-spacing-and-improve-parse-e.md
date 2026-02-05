---
type: issue
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:38:56.93394Z
date_edited: 2026-02-05T11:55:47.323102Z
owner_approval: false
completed: true
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

## Completion Report
Implemented YAML parser improvements and spacing fixes

- Added FrontmatterParseError type with file path and line number tracking
- Enhanced parser to extract and report line numbers from YAML errors
- Verified YAML writer produces correct spacing (key: value with space)
- Added extractLineNumberFromYAMLError helper to parse YAML error format
- Created comprehensive test suite covering malformed YAML, valid parsing, and file I/O
- All tests pass and build succeeds with no regressions

Users now get clear, actionable error messages when encountering malformed YAML files,
pointing to the exact file and line with the problem instead of cryptic parse failures.
