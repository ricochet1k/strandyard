---
type: leaf
role: reviewer-usability
priority: medium
parent: T8v3k-recurring-tasks
blockers:
    - T3ebv-define-tasks-completed-scope-and-anchor-semantics
    - Tftao-clarify-recurrence-anchor-flags-and-help-text
    - Tpedm-decide-tasks-completed-storage-strategy
    - Tu1vk-define-lines-changed-scope-flag-and-default
    - Tw6ga-define-recurrence-cli-shape-for-discoverability
    - Twvju-decide-recurrence-schema-option-a-vs-b
blocks:
    - T8v3k-recurring-tasks
date_created: 2026-01-28T17:32:22.794929Z
date_edited: 2026-01-30T22:38:51.494417Z
owner_approval: false
completed: true
---

# Review recurrence metrics CLI experience

## Context
- Design doc: `design-docs/recurrence-metrics.md`
- Recurring task plan: `tasks/E7p4m-issues-recurrence/T8v3k-recurring-tasks/T8v3k-recurring-tasks.md`
- CLI usage guide: `CLI.md` (`recurring add`, `recurring materialize`)

## Usability Review
- Merge `recurring add` into `add` to reduce command surface area. Consider `strand add recurring` or `strand add --recurrence ...` so users discover recurrence alongside existing task creation.
- `recurring add` currently models a single interval; extending to `lines_changed` and `tasks_completed` either adds new `--unit` values (Option A) or needs a multi-trigger CLI shape (Option B). Option A is simpler but can’t express combined triggers; Option B is more flexible but needs heavier CLI input. Decision: deferred.
- `--anchor` accepts date or commit hash today; adding more metrics increases ambiguity. Consider explicit flags like `--anchor-date`/`--anchor-commit` or `--anchor-type` to make help text and validation clearer.
- Lines-changed metrics need a clear counting rule. Suggest a `--lines-scope=added|deleted|total` flag with an explicit default in help and examples.
- Tasks-completed metrics need a scope definition (e.g., only `type: task` vs. all types). Suggest a `--task-scope=tasks|all` flag plus clear guidance on which timestamp/anchor field is used.
- Proposed CLI examples once metrics land (ensure `--help` includes these):
  - `strand recurring add "Repo hygiene" --interval 500 --unit commits --anchor HEAD~500`
  - `strand recurring add "Code churn review" --interval 2000 --unit lines_changed --anchor v1.2.0 --lines-scope total`
  - `strand recurring add "QA audit" --interval 10 --unit tasks_completed --anchor 2026-01-01T00:00:00Z --task-scope tasks`

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] (role: {{ .Role }}) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
- [x] (subtask: T3ebv) Define tasks-completed scope and anchor semantics
- [ ] (subtask: Tftao) Clarify recurrence anchor flags and help text
- [x] (subtask: Tpedm) Decide tasks-completed storage strategy
- [ ] (subtask: Tu1vk) Define lines-changed scope flag and default
- [ ] (subtask: Tw6ga) Define recurrence CLI shape for discoverability
- [ ] (subtask: Twvju) Decide recurrence schema option (A vs B)
