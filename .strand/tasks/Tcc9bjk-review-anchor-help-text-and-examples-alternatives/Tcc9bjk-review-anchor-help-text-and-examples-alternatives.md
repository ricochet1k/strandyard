---
type: review
role: master-reviewer
priority: medium
parent: Tusef-define-anchor-help-text-and-examples
blockers: []
blocks: []
date_created: 2026-02-01T20:20:24.066079Z
date_edited: 2026-02-05T04:07:06.304887Z
owner_approval: false
completed: false
description: ""
---

# Review Anchor Help Text Alternatives

## Summary
Review the design alternatives for `strand add --every` help text and CLI.md documentation in `design-docs/anchor-help-text-and-examples-alternatives.md`.

## Description
The alternatives propose different approaches to help text formatting:
- Alternative A: Compact one-line help with examples
- Alternative B: Grouped help with section headers  
- Alternative C: Minimal help with reference URL

The goal is to balance scannability for users who skim help output with completeness and clarity for anchor types and examples.

## Questions for Review
1. Which alternative best serves users who skim help output?
2. Does the proposed help text format clearly distinguish time-based vs git-based metrics?
3. Are the examples sufficient and deterministic for tests/automation?
4. Is the CLI.md table scannable and comprehensive?

## Acceptance Criteria
- Reviewer provides feedback on which alternative to adopt
- Any usability concerns are documented
- Reviewer approves the selected alternative or requests changes

Delegate concerns to the relevant role via subtasks.
