---
type: issue
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-07T00:24:36.28766Z
date_edited: 2026-02-07T00:24:36.28766Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Improve preset refresh help, errors, and validation

## Summary
## Summary
The `strand preset refresh` command needs better help text, error messages, and validation to make it easier to use and debug when things go wrong.

## Current issues
1. Help text doesn't explain what a valid preset is (local directory with roles/ and templates/ subdirs, or git URL)
2. No validation that the preset source contains the expected directory structure before attempting copy
3. Error messages when git clone fails don't explain common causes (invalid URL, network issues, auth required)
4. No clear feedback about what was refreshed or what changed
5. Silent failures when preset is missing expected directories

## Proposed improvements
1. Enhanced help text explaining:
   - What a preset is (directory structure requirements)
   - Examples of valid preset paths (local dir, git URL)
   - What gets overwritten (roles/, templates/) vs preserved (tasks/)
   
2. Better validation:
   - Check preset structure before copying
   - Validate git URLs before attempting clone
   - Provide clear error if preset is missing required directories
   
3. Improved error messages:
   - Git clone failures should suggest common fixes
   - File copy errors should show which file/directory failed
   - Permission errors should be clearly identified
   
4. Better feedback:
   - Show what directories are being refreshed
   - Report number of files copied/updated
   - Confirm successful completion

## Acceptance criteria
- Help text includes examples and explains preset structure
- Command validates preset before attempting copy
- Error messages are actionable and explain how to fix issues
- Success output shows what was changed
- All error paths have tests

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds
