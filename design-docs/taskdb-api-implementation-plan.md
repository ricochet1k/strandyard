# Implementation Plan: TaskDB API Hardening

## Context and decision records
- `pkg/task/TASKDB_DESIGN.md` (problem statement and required outcomes)
- `design-docs/task-go-structure-review.md` (current mutable surface in `Task`/`Metadata`)
- `design-docs/blockers-relationship-management-review.md` (relationship reconciliation behavior and boundary options)
- `pkg/task/doc.go` (current recommended usage and known pitfalls)
- Parent review task: `Ti6zj` (TaskDB API Design Review)

## Objectives
- Make TaskDB the default safe mutation path for task state and relationships.
- Remove or isolate APIs that make misuse easy (`GetOrCreate`, direct relationship writes, overlapping mutators).
- Keep behavior deterministic (sorted relationship fields, stable repair/list output).
- Migrate existing CLI/web call sites without regressions.

## Decision gates (owner decisions required)
The implementation depends on unresolved design tasks. Work should proceed in this order, with explicit owner checkpoints:

1. **Boundary and ownership** (`Thyd1`): confirm TaskDB vs parser/validator responsibilities.
2. **Access-control model** (`T8bgf`): choose hardening approach for `Task`/`Metadata` mutability.
3. **Relationship API shape** (`Tx4jn`): finalize public/private/deleted relationship methods.
4. **Template-only creation contract** (`Txvyh`): finalize creation flow and TaskDB integration.
5. **API audit/deletion list** (`Twcdw`, `Tb0oq`): confirm exact removals and deprecations.

Decision: deferred for all unresolved alternatives until owner approval artifacts exist.

## Risk and dependency ordering

### Highest risk (do first)
1. API boundary decisions (can create breaking public API changes).
2. Task creation contract changes (removal of `GetOrCreate` and similar flows).
3. Relationship API consolidation (must preserve invariants and compatibility).

### Medium risk
4. Call-site migration across CLI/web/tests.
5. Validation/repair alignment with new contracts.

### Lower risk
6. Documentation and examples updates after API stabilizes.

## Breaking changes to plan for
- Remove `TaskDB.GetOrCreate` (currently allows blank task creation).
- Restrict or redesign direct mutable access patterns to relationship/status fields where feasible.
- Rename/remove redundant or confusing relationship APIs once final naming is approved.
- Potentially tighten method preconditions (for example status transitions and task existence rules).

External impact:
- Public `pkg/task` consumers may need code changes if method signatures/exported methods change.
- Internal impact spans `cmd/*`, `pkg/web/*`, and tests that currently rely on direct field edits.

## Migration strategy

### Compatibility window
- Introduce replacement APIs first.
- Keep old APIs as wrappers with explicit deprecation comments for one migration cycle where practical.
- Add compile-time and runtime guidance (errors/docs) for removed unsafe paths.

### Call-site migration approach
1. Migrate CLI commands to new TaskDB entry points.
2. Migrate web handlers to avoid direct relationship field writes.
3. Migrate tests/examples to assert invariants through TaskDB, not direct metadata mutation (except targeted low-level tests).
4. Remove deprecated APIs only after all in-repo call sites are clean and tests pass.

### Data migration
- No on-disk schema rewrite is expected for this task by default.
- Relationship and status invariants remain repairable through existing reconciliation/repair flows.

## Implementation phases

### Phase 0: finalize design decisions
Dependencies: owner review of `Thyd1`, `T8bgf`, `Tx4jn`, `Txvyh`, `Twcdw`, `Tb0oq`, `Tdhrq`, `Trtik`.

Deliverables:
- Approved architecture note describing TaskDB boundaries.
- Approved API table (keep/deprecate/remove/private).
- Approved template-only creation contract.

### Phase 1: API scaffolding and guardrails (incremental)
Dependencies: Phase 0.

Deliverables:
- New/updated TaskDB methods implementing approved safe flows.
- Deprecation markers on legacy unsafe/redundant methods.
- Guardrails in docs/errors for direct misuse patterns.

Files likely touched:
- `pkg/task/taskdb.go`
- `pkg/task/doc.go`
- `pkg/task/taskdb_example_test.go`

### Phase 2: relationship and creation hardening (partly big refactor)
Dependencies: Phase 1.

Deliverables:
- Remove `GetOrCreate` and any equivalent blank creation path.
- Finalize relationship API names and responsibilities.
- Ensure reconciliation/validation semantics are consistent after API changes.

Files likely touched:
- `pkg/task/taskdb.go`
- `pkg/task/blockers.go`
- `pkg/task/parser.go`
- `pkg/task/repair.go`

### Phase 3: call-site migration (incremental)
Dependencies: Phase 2.

Deliverables:
- CLI commands and web handlers use only approved TaskDB operations.
- No direct relationship/status field mutation in command/handler paths.

Files likely touched:
- `cmd/add.go`, `cmd/edit.go`, `cmd/complete.go`, `cmd/repair.go`, `cmd/next.go`
- `pkg/web/handlers.go`

### Phase 4: deletion cleanup and docs convergence (incremental)
Dependencies: Phase 3.

Deliverables:
- Delete deprecated APIs and tests for removed behavior.
- Update GoDoc and design docs to match final API.
- Ensure examples model only approved patterns.

Files likely touched:
- `pkg/task/taskdb.go`
- `pkg/task/taskdb_test.go`
- `pkg/task/taskdb_example_test.go`
- `pkg/task/doc.go`
- `CLI.md`

## Incremental vs big-refactor split
- **Incremental:** call-site migration, docs/examples updates, deprecation shims, test rewrites, validator/repair rule alignment.
- **Big refactor candidates:** changing mutability model of `Task`/`Metadata` (if fields become less directly writable), and any broad exported API contraction.

## Code-change checklist
- [ ] Finalize boundaries and access-control decisions (`Thyd1`, `T8bgf`).
- [ ] Finalize relationship/task-creation API designs (`Tx4jn`, `Txvyh`).
- [ ] Produce complete API misuse inventory and deletion list (`Twcdw`, `Tb0oq`).
- [ ] Review parser/repair integration constraints (`Tdhrq`, `Trtik`).
- [ ] Implement approved API updates in `pkg/task`.
- [ ] Migrate CLI and web call sites to approved TaskDB methods.
- [ ] Remove deprecated APIs and misuse examples once migration completes.
- [ ] Update docs/examples and verify GoDoc guidance.
- [ ] Run `go build ./...`, `go test ./...`, and `go run ./cmd/strand repair`.

## Validation and release criteria
- All in-repo relationship/status mutations go through approved TaskDB APIs.
- No remaining production references to removed/deprecated methods.
- Full test suite passes and repair is stable/no-op on clean data.
- Design docs and GoDoc consistently describe one mental model with no contradictory guidance.
