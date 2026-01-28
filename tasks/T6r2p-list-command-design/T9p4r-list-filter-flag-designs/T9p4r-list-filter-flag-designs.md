---
kind: ""
role: designer
priority: medium
parent: T6r2p-list-command-design
blockers: []
blocks: []
date_created: 2026-01-28T01:52:27Z
date_edited: 2026-01-28T01:52:27Z
owner_approval: false
completed: false
---

# Evaluate list filter flag design alternatives

## Context
Issue: tasks/I5c1s-explore-alternate-list-filter-flag-designs/I5c1s-explore-alternate-list-filter-flag-designs.md
Design doc: design-docs/list-command.md (Filters section)
Current flags: cmd/list.go

## Tasks
- [ ] Draft 2–4 alternative filter flag designs (e.g., `--filter key=value`, `--status`, boolean pairs like `--blocked/--unblocked`).
- [ ] Compare each alternative for usability, composability, backward compatibility, and help/UX clarity.
- [ ] Identify recommended default and whether to support aliases for transition.
- [ ] Update or propose updates to `design-docs/list-command.md` with the chosen pattern (await owner approval).

## Acceptance Criteria
- Alternatives with pros/cons and risks are documented.
- A recommended pattern is proposed with rationale.
- Owner approval is requested before any implementation work proceeds.
