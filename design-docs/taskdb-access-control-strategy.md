# TaskDB Access Control Strategy

## Context and references
- Parent review: `Ti6zj` (TaskDB API Design Review)
- Task: `T8bgf` (Design access control strategy)
- `pkg/task/TASKDB_DESIGN.md`
- `design-docs/task-go-structure-review.md`
- `design-docs/taskdb-api-implementation-plan.md`
- Current mutable surfaces: `pkg/task/task.go` and `pkg/task/taskdb.go`

## Reproduction steps
```bash
go build ./...
go test ./...
go run ./cmd/strand repair
```

## Problem to solve
Current exported `Task` and `Metadata` fields allow:
- manual `Task` creation without template defaults,
- direct field mutation that bypasses relationship/status invariants,
- command/handler code paths that can skip TaskDB safeguards.

The API should be hard to misuse while remaining idiomatic Go.

## Goals
- Prevent ad hoc blank task creation outside approved creation APIs.
- Prevent direct mutation of relationship-critical fields (`Parent`, `Blockers`, `Blocks`, `Status`, `Completed`, role/priority when policy requires validation).
- Make TaskDB the single safe mutation boundary for task graph updates.
- Preserve deterministic write behavior and existing on-disk format.

## Alternatives

### Alternative A: Documentation + linting only
Keep exported structs/fields and rely on docs, examples, and static checks.

Pros:
- Lowest migration cost.
- Minimal public API breakage.

Cons:
- Does not enforce invariants at compile time.
- Third-party consumers can still bypass TaskDB.
- Fails the "hard to use incorrectly" requirement.

### Alternative B: Public read-only views + TaskDB mutators (recommended)
Expose immutable read models and keep mutable task records internal to `pkg/task`.

Shape:
- Introduce exported read-only interface/value (for example `TaskView`) with getters only.
- Make concrete mutable task type unexported (for example `taskRecord`).
- Keep all mutating operations on `TaskDB` methods (`SetParent`, `AddBlocker`, `SetStatus`, `SetRole`, etc.).
- Provide explicit creation entrypoints (`CreateFromTemplate`, `AddTaskFromTemplate`, or equivalent) and remove `GetOrCreate`.
- Return defensive copies for list/slice data (`Blockers`, `Blocks`, recurrence anchors).

Pros:
- Strong compile-time guardrails against direct field writes.
- Keeps Go API straightforward: query via views, mutate via methods.
- Preserves testability and deterministic serialization in one package.

Cons:
- Breaking API change for code reading/writing exported fields directly.
- Requires migration of callers/tests to getters + TaskDB methods.

### Alternative C: Fully opaque handles + transactional edit closures
Expose only opaque IDs/handles and require edits via callback transactions.

Pros:
- Maximum control over mutation boundaries.
- Easier to add future locking/version checks.

Cons:
- Less idiomatic for this repository's current complexity.
- Higher implementation and migration complexity.
- Harder ergonomics for simple command flows.

## Trade-off summary
- Safety: C > B > A
- API simplicity: A >= B > C
- Migration cost: A < B < C
- Fit for current project size: B is best balanced.

## Concrete design proposal

Decision: deferred (owner approval required).

Recommended proposal to approve:

1. **Type boundary**
   - Introduce unexported mutable model (`taskRecord`, `metadataRecord`) inside `pkg/task`.
   - Expose read-only `TaskView`/`MetadataView` interfaces or value snapshots.
   - `TaskDB.Get` and collection APIs return read-only views.

2. **Mutation boundary**
   - Keep all graph/status/content mutation on `TaskDB` methods.
   - Remove (or make internal) mutation-capable surfaces that leak pointers to mutable records.
   - Continue central dirty tracking and date-edited updates inside TaskDB/write path.

3. **Creation boundary**
   - Remove `TaskDB.GetOrCreate`.
   - Add explicit template-backed creation API that requires role/title/description and optional parent.
   - Parse/load paths remain able to read existing files; creation paths must not emit blank metadata.

4. **Compatibility strategy**
   - Phase 1: add read-only views and replacement methods.
   - Phase 2: migrate internal call sites (`cmd/*`, `pkg/web/*`, tests) off direct field access.
   - Phase 3: remove deprecated unsafe APIs and enforce no direct mutable exports.

5. **Validation guardrails**
   - Add tests proving relationship/status invariants can only be changed through TaskDB mutators.
   - Add compile-time tests/examples showing external packages cannot write critical fields.

## API-level acceptance checks (for implementation phase)
- External callers cannot set parent/blockers/blocks/status via field assignment.
- No public API path creates empty tasks without template defaults.
- `go test ./...` includes mutation-boundary tests and passes.
- Existing CLI workflows continue to function after migration.

## Owner decision required
Choose one:
- Approve Alternative B (recommended) for implementation.
- Choose Alternative A (documentation-only) and accept weaker safety.
- Choose Alternative C (opaque transaction model) and accept higher complexity.
