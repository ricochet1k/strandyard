# TaskDB API Deletion List (Tb0oq)

## Context and decision records
- `pkg/task/TASKDB_DESIGN.md`
- `design-docs/taskdb-api-implementation-plan.md`
- `design-docs/taskdb-access-control-strategy.md`
- `design-docs/blockers-relationship-management-review.md`
- Parent review: `Ti6zj` (TaskDB API Design Review)

## Reproduction Steps
```bash
go test ./...
go build ./...
go run ./cmd/strand repair
```

## Scope
This list identifies code paths that should be removed once replacement APIs are fully migrated. It also calls out test and example code that currently validates or demonstrates unsafe patterns.

Decision status: deferred to owner review for final keep/deprecate/remove calls.

## Deletion Candidates

| Category | Item | Location | Why remove | Replacement / migration target |
| --- | --- | --- | --- | --- |
| GetOrCreate and similar | `TaskDB.GetOrCreate` | `pkg/task/taskdb.go` | Creates blank tasks and bypasses template-only creation; explicitly conflicts with TaskDB design constraints. | Template-backed creation API only (`strand add`/TaskDB template create entrypoint from `Txvyh`). |
| Test code for deleted functionality | `TestTaskDB_GetOrCreate` | `pkg/task/taskdb_test.go` | Encodes behavior that should no longer be allowed (blank task creation). | Replace with tests asserting creation fails without template input and succeeds through approved creation path. |
| Redundant relationship management | `TaskDB.AddBlocked` | `pkg/task/taskdb.go` | Thin inverse wrapper over `AddBlocker`; adds duplicated verb surface and parameter-order confusion. | Use `AddBlocker(taskID, blockerID)` as the single public directional mutator. |
| Redundant relationship management | `TaskDB.RemoveBlocked` | `pkg/task/taskdb.go` | Thin inverse wrapper over `RemoveBlocker`; same ambiguity and redundant API surface. | Use `RemoveBlocker(taskID, blockerID)` as the single public directional mutator. |
| Example code that demonstrates misuse | Direct relationship field writes (`t.Meta.Blockers = ...`, `t.Meta.Blocks = ...`) in request handling | `pkg/web/handlers.go` | Bypasses TaskDB invariants and normalization flow, making misuse easy in production paths. | Route all relationship edits through TaskDB mutators + reconciliation. |
| Poorly named functions being replaced | `TaskDB.CompleteTask` (bool-oriented completion helper) | `pkg/task/taskdb.go`; used in `cmd/complete.go` | Name implies canonical completion path but delegates to legacy completed-boolean behavior; overlaps newer status-first API. | Standardize on `SetStatusWithReport(..., StatusDone, report)` as canonical completion entrypoint. |
| Poorly named functions being replaced | `TaskDB.SetCompleted` (legacy bool setter) | `pkg/task/taskdb.go` | Encourages legacy `completed`-first mental model and duplicates status transition logic. | Keep internal compatibility only during migration; remove from public surface after call sites move to status APIs. |

## Notes on previously mentioned legacy names
- `UpdateBlockersFromChildren`, `FixBlockerRelationships`, and `SyncBlockersFromChildren` are not present in the current `pkg/task` implementation.
- Their current equivalent behavior is in `ReconcileBlockerRelationships`; this deletion list therefore focuses on APIs that still exist and still create misuse risk.

## Call-Site Impact Summary
- `AddBlocked` / `RemoveBlocked` are currently used in:
  - `cmd/add.go`
  - `cmd/edit.go`
  - `pkg/web/handlers.go`
- `CompleteTask` is currently used in:
  - `cmd/complete.go`
- `GetOrCreate` has no non-test call sites today (removal risk is low once tests are updated).

## Suggested deletion order
1. Remove `GetOrCreate` and `TestTaskDB_GetOrCreate` first (no production call sites).
2. Migrate completion flow to `SetStatusWithReport(..., StatusDone, report)` and then remove `CompleteTask`/`SetCompleted` from public API.
3. Migrate `AddBlocked`/`RemoveBlocked` call sites to `AddBlocker`/`RemoveBlocker` and remove wrapper methods.
4. Remove direct relationship field writes from request handlers and enforce mutator-only relationship edits.
