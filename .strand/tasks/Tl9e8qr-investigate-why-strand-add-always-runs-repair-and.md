---
type: task
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-15T04:57:46.029164Z
date_edited: 2026-02-15T04:57:46.029164Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Investigate why strand add always runs repair and outputs long IDs

## Summary
## Summary
Triaged from Ti1ugbe. `strand add` still prints `repair: ok` and `Repaired 0 tasks` after creating a task, and output emphasizes long IDs rather than concise human-friendly task details.

## Repro
1. In a clean temporary repo, run `strand init --storage local --preset /Users/matt/mycode/strandyard/.strand`.
2. Run `strand add issue "triage repro add output" --priority medium`.
3. Observe output:
   - `✓ Task created: Tgvodvf-triage-repro-add-output`
   - `repair: ok`
   - `Repaired 0 tasks`

## Investigation Goals
- Confirm whether `add` should skip full repair by default and only do incremental TaskDB/list updates.
- Determine whether repair invocation can be removed safely without regressions.
- Propose/implement more human-friendly `add` output (short id/title/priority) while keeping machine-parseable behavior stable where required.
- Add tests that lock expected behavior.

## Instructions
Decide which task template would best fit this task and re-add it with that template and the same parent.
