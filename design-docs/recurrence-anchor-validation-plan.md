# Implementation Plan — Recurrence Anchor Validation

## Architecture Overview
Enhance the validation of the `--every` flag in `strand add` to ensure that anchors specified in the recurrence rule actually exist and are valid for their respective metrics.

## Specific Files to Modify
- `cmd/add.go`: Update `validateEvery` to perform real-world validation of anchors.
- `pkg/task/recurrence.go`: (If needed) add validation logic that can be shared between `add` and `materialize`.

## Approach
1.  **Date Validation**:
    - Support ISO 8601 and the human-friendly format `Jan 28 2026 09:00 UTC`.
    - Use `time.Parse` to verify the anchor is a valid date.
2.  **Commit Validation**:
    - Use `git rev-parse` to check if the provided anchor (commit hash or `HEAD`) exists in the current repository.
    - If not in a git repo, allow `HEAD` but warn or handle gracefully (currently `CLI.md` says it's a no-op).
3.  **Task Validation**:
    - For `tasks_completed` metric, verify that the anchor (if provided) is a valid task ID in the project.
    - This requires loading the task database, which `validateEvery` currently doesn't do.

## Integration Points
- `validateEvery` is called in `runAdd`.
- Git validation will require calling out to `git`.
- Task validation will require access to the loaded `tasks` map.

## Testing Approach
- Add test cases to `cmd/add_every_test.go` for:
    - Valid and invalid ISO 8601 timestamps.
    - Existing and non-existing commit hashes.
    - Existing and non-existing task IDs.
- Use mocks or temporary git repositories for commit validation tests.

## Decision Rationale
- Early validation prevents misconfigured recurring tasks from being created, which is better than failing during materialization.
- Using `git rev-parse` is the standard way to verify objects in git.
- Task validation should only happen if the task database is available.

## Child Tasks
- [ ] (role: developer) Implement date anchor validation in `validateEvery`.
- [ ] (role: developer) Implement commit anchor validation using `git rev-parse`.
- [ ] (role: developer) Implement task anchor validation for `tasks_completed` metric.
- [ ] (role: master-reviewer) Review recurrence anchor validation implementation.
