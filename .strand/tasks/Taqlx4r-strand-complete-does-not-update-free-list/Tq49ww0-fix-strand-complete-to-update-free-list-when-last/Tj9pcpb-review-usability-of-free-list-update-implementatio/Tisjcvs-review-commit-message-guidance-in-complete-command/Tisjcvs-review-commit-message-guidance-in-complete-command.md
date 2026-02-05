---
type: task
role: designer
priority: medium
parent: Tj9pcpb-review-usability-of-free-list-update-implementatio
blockers: []
blocks: []
date_created: 2026-02-05T12:09:53.599105Z
date_edited: 2026-02-05T12:09:53.599105Z
owner_approval: false
completed: false
description: ""
---

# New Task: Review commit message guidance in complete command

## Description
## Context
When completing tasks via last TODO, the commit message suggestion is:
```
git add -A && git commit -m "complete: %s"
```

When completing individual todos, the suggestion is:
```
git add -A && git commit -m "%v (%v) check off %v"
```

## Issues
1. Inconsistent message formats across different completion paths
2. For incomplete todos, the format is verbose and may not match repo conventions
3. Users might not know what message format is expected

## Acceptance Criteria
- Consistent commit message guidance across all paths
- Message format aligns with documented repo conventions
- Users clearly understand what the message represents

Decide which task template would best fit this task and re-add it with that template and the same parent.
