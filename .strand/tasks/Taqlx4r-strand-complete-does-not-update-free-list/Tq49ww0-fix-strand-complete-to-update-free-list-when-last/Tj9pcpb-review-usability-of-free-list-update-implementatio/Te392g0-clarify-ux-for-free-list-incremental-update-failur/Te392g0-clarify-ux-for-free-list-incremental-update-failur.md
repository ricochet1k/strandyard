---
type: task
role: designer
priority: medium
parent: Tj9pcpb-review-usability-of-free-list-update-implementatio
blockers: []
blocks: []
date_created: 2026-02-05T12:09:49.743763Z
date_edited: 2026-02-05T12:09:49.743763Z
owner_approval: false
completed: false
description: ""
---

# New Task: Clarify UX for free-list incremental update failures

## Description
## Context
When `UpdateFreeListIncrementally` fails, users see warning message:
```
⚠️  Incremental update failed, falling back to full repair: %v
```

This requires clarification:
1. When/why does incremental update fail?
2. What does "full repair" mean in user language?
3. Why should the user care?

## Questions
- Should this be a warning to the user or just transparent?
- Is the fallback automatic or does user need to do anything?
- Are there scenarios where fallback repair also fails?

## Acceptance Criteria
- UX messaging explains the failure reason in plain language
- User knows whether action is required
- Documentation reflects the behavior

Decide which task template would best fit this task and re-add it with that template and the same parent.
