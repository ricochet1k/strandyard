# Graceful Error Handling for Invalid Anchors During Materialization

## Summary
Design and specify error handling behavior for when recurrence anchors become invalid during task materialization (when recurring task instances are being created). This includes detection, logging, error messages, and fallback strategies.

## Context
- design-docs/recurrence-metrics.md (defines recurrence metrics and anchors)
- design-docs/recurrence-anchor-flags-alternatives.md (anchor CLI design)
- design-docs/recurrence-anchor-error-messages.md (error format for parsing/validation phase)
- pkg/task/recurrence.go (current evaluateGitMetric implementation)
- Task: Tjlzm08-add-graceful-error-handling-for-invalid-anchors-du

## Problem Statement

When recurring tasks are materialized (new instances created based on recurrence rules), the stored anchors may become invalid due to:
- Git operations (rebasing, force-pushing, garbage collection) removing commit references
- Repository state changes (commits deleted, branches renamed)
- Invalid configuration stored in task metadata

The current implementation in `pkg/task/recurrence.go`:
- Returns (0, nil) for "unknown revision" and "ambiguous argument" errors (silent no-op)
- Returns (0, nil) for invalid HEAD (silent no-op)
- Returns an error for other git command failures (non-deterministic message format)

**Issues with current behavior:**
1. Silent no-op may confuse users - recurrence stops triggering without clear indication
2. No logging of anchor validation failures during materialization
3. Error messages for non-no-op cases don't follow the unified error format
4. No audit trail of when anchors become invalid

## Design Goals
1. **Detectability**: Clearly identify when anchors become invalid during materialization
2. **Recoverability**: Provide actionable guidance for fixing invalid anchors
3. **Auditability**: Log anchor validation failures for debugging and auditing
4. **Consistency**: Follow the unified error format defined in recurrence-anchor-error-messages.md
5. **Graceful degradation**: Allow system to continue operating with clear warnings when possible

## What Constitutes an Invalid Anchor

### Git-Based Anchors (commits, lines_changed)
- Commit hash that doesn't exist in repository (`unknown revision`)
- HEAD reference is invalid/unborn in the repository
- Ambiguous commit reference (rare with full hashes)
- Git command execution errors (e.g., corrupt repository)

### Date-Based Anchors (days, weeks, months, tasks_completed)
- Invalid date format in stored metadata (should be caught by validation at creation time)
- Date parsing errors during materialization (unlikely with ISO 8601 format)
- Edge case: future dates that prevent materialization (log as warning, don't error)

## Detection Strategy

### Git-Based Anchor Validation
1. **Check HEAD validity first** using `git rev-parse --verify HEAD`
   - If invalid: This is a repository state issue
   - Action: Return error with clear message, log to stderr
2. **Verify anchor exists** using `git rev-parse --verify <anchor>`
   - If unknown revision: Anchor was deleted or never existed
   - Action: Return error with actionable message, log to audit trail
3. **Execute metric command** with proper error capture
   - Catch all errors, not just "unknown revision" and "ambiguous argument"
   - Action: Return unified error format with hint lines

### Date-Based Anchor Validation
1. Parse date using ISO 8601 parser
2. Validate date is not in the future (warning only)
3. Anchor validation is primarily done at task creation time, not materialization

## Error Handling Approaches

### Option A — Fail-Fast with Clear Errors (Recommended)
- **Description**: Materialization fails when anchor is invalid; provide clear error messages with recovery hints
- **Behavior**:
  - Exit code 2 for anchor validation failures (matches parse/validation errors)
  - Error format: `strand: error: <reason>` followed by `hint: <example>`
  - Log to stderr only (no silent no-op)
- **Pros**:
  - Prevents silent failures - user knows immediately when recurrence stops working
  - Follows the unified error format from parsing phase
  - Forces user to fix the invalid anchor (better than broken recurrence)
- **Cons**:
  - More disruptive - recurring task creation stops entirely
  - Requires manual intervention to fix anchors
- **Implementation Effort**: Low to Medium

### Option B — Graceful Degradation with Warnings
- **Description**: Materialization continues with zero metric value; emit warning messages
- **Behavior**:
  - Return (0, nil) for invalid anchors (current behavior)
  - Emit warning to stderr: `strand: warning: <reason>`
  - Log to audit trail for visibility
- **Pros**:
  - Less disruptive - materialization continues
  - System remains functional even with invalid anchors
- **Cons**:
  - Silent behavior may be missed in automated workflows
  - Warnings may be ignored
  - Harder to detect when recurrence stops working as expected
- **Implementation Effort**: Low

### Option C — Configurable Failure Mode
- **Description**: Add a configuration option to choose between fail-fast or graceful degradation
- **Behavior**:
  - Default to fail-fast (Option A)
  - Allow configuration for graceful degradation (Option B) via flag or config file
  - Per-task or per-repository setting
- **Pros**:
  - Flexibility for different use cases
  - Allows fail-fast by default but permits graceful degradation when needed
- **Cons**:
  - Adds complexity
  - Requires configuration management
  - May lead to inconsistent behavior across environments
- **Implementation Effort**: Medium to High

## Decision

**Adopt Option A — Fail-Fast with Clear Errors**

Rationale:
- Silent failures (current behavior) are worse than explicit failures
- Follows the unified error format established for parsing/validation
- Forces users to fix invalid anchors rather than operating with broken recurrence
- Aligns with the principle of keeping behavior deterministic and auditable
- Future extensibility: can add graceful degradation mode if needed via config

## Implementation Plan

### Files to Modify

1. **pkg/task/recurrence.go** - Update `evaluateGitMetric` function
   - Remove silent no-op behavior for unknown revision errors
   - Return proper errors following unified format
   - Add detailed error messages with recovery hints

2. **pkg/task/errors.go** (new file) - Define error types and messages
   - Create error types for different anchor validation failures
   - Implement error formatting that follows unified format
   - Provide helper functions for constructing error messages with hints

3. **pkg/task/logger.go** (new file) - Add structured logging for audit trail
   - Log anchor validation failures with context (task ID, anchor value, metric type)
   - Use structured logging format for easy parsing and analysis
   - Log to stderr or file based on configuration

### Code Structure

```
pkg/task/
├── recurrence.go (modified)
│   ├── evaluateGitMetric (return errors instead of no-op)
│   ├── isHeadValid (keep as-is)
│   └── tempGitRepo (keep as-is)
├── errors.go (new)
│   ├── type AnchorValidationError struct
│   ├── type InvalidCommitAnchorError struct
│   ├── type InvalidHeadError struct
│   ├── type GitCommandError struct
│   ├── func NewAnchorValidationError(...)
│   ├── func NewInvalidCommitAnchorError(...)
│   └── func (e *Error) Error() string // returns unified format
└── logger.go (new)
    ├── func LogAnchorValidationFailure(taskID, metricType, anchor, error)
    └── func SetLoggerDestination(io.Writer)
```

### Error Message Format

Follow the unified format from recurrence-anchor-error-messages.md:

```go
type AnchorValidationError struct {
    MetricType  string // "commits", "lines_changed", "days", etc.
    AnchorValue string // the problematic anchor
    Reason      string // specific validation failure
}

func (e *AnchorValidationError) Error() string {
    return fmt.Sprintf("strand: error: metric '%s' anchor '%s' is invalid: %s",
        e.MetricType, e.AnchorValue, e.Reason)
}

func (e *AnchorValidationError) Hint() string {
    // Return recovery hint based on metric type and anchor
}
```

### Specific Error Messages

#### Invalid HEAD
```
strand: error: cannot evaluate recurrence: HEAD is invalid or unborn
hint: Run 'git init' or ensure you are in a valid git repository
```

#### Unknown Revision (Git Anchor)
```
strand: error: metric 'commits' anchor 'abc123def' does not exist in repository
hint: Update anchor in task metadata to a valid commit hash, or use 'HEAD'
hint: Run 'git log' to find valid commit hashes
```

#### Git Command Error
```
strand: error: failed to evaluate metric 'commits': git command failed
hint: Ensure repository is not corrupted and git is working correctly
hint: Try running 'git fsck' to check repository integrity
```

### Integration Points

1. **Recurring task materialization** (main integration point)
   - Call `evaluateGitMetric` during materialization
   - Check for errors and propagate to caller
   - Handle errors at CLI level (emit to stderr, exit with code 2)

2. **Task metadata validation** (prevention)
   - Validate anchors exist when task is created/edited
   - This is covered by other tasks (Tg96jgm-add-validation-for-anchor-existence-at-recurrence)

3. **Audit trail** (observability)
   - Log anchor validation failures with task ID and context
   - Enable post-mortem analysis of failed materializations
   - Support debugging in production environments

### Testing Strategy

#### Unit Tests
1. Test `evaluateGitMetric` with invalid HEAD
   - Expect error, not (0, nil)
   - Verify error message format
2. Test `evaluateGitMetric` with unknown commit hash
   - Expect error, not (0, nil)
   - Verify error message includes anchor value
3. Test `evaluateGitMetric` with ambiguous argument
   - Expect error, not (0, nil)
   - Verify error message format
4. Test error message formatting
   - Assert exact string format matches unified error format
   - Verify hint lines are provided
5. Test logger output
   - Verify structured logging format
   - Check that task ID, metric type, and anchor are logged

#### Integration Tests
1. Create a recurring task with commit-based anchor
2. Delete the commit from repository
3. Attempt materialization
4. Verify error is returned and logged correctly

#### Edge Cases
1. Repository with unborn HEAD (fresh git init)
2. Corrupted git repository
3. Anchor references a tag that was deleted
4. Simultaneous git operations during materialization

## Decision Rationale

### Why Fail-Fast (Option A) Over Graceful Degradation (Option B)?
- **Visibility**: Silent failures are worse than explicit failures
- **Recoverability**: Clear errors with hints enable faster recovery
- **Auditability**: Failures are logged and tracked
- **Principle**: Aligns with project's emphasis on deterministic, auditable behavior

### Why Not Configurable Mode (Option C)?
- **Simplicity**: One clear behavior is easier to understand and maintain
- **Principle of least surprise**: Consistent behavior across all environments
- **Future-proof**: Can add config later if needed without breaking changes

### Why Follow Unified Error Format?
- **Consistency**: Same format as parsing/validation errors
- **Testability**: Tests can assert exact string matches
- **Documentation**: Users see familiar error format across all error conditions

## Trade-Offs

| Aspect | Fail-Fast | Graceful Degradation |
|--------|-----------|----------------------|
| User Experience | Explicit, actionable errors | Warnings may be missed |
| System Availability | Materialization stops | Materialization continues |
| Debugging | Clear failure point | Harder to detect issues |
| Recovery Path | Immediate action required | May ignore warnings |
| Auditability | All failures logged | Warnings logged, may be ignored |

## Open Questions

1. **Log destination**: Should audit logs go to stderr, file, or both?
   - Proposal: stderr for now (matches CLI conventions); add file logging later if needed

2. **Retry logic**: Should materialization retry on transient git errors?
   - Proposal: No - fail-fast on all anchor errors; retries add complexity without clear benefit

3. **Exit code**: Should anchor validation failures use exit code 2 (same as parse errors) or a different code?
   - Proposal: Exit code 2 for all anchor validation errors (consistent with parsing errors)

## Acceptance Criteria

- [ ] `evaluateGitMetric` returns errors instead of (0, nil) for invalid anchors
- [ ] All errors follow unified format: `strand: error: <reason>` + `hint: <example>`
- [ ] Error messages include specific anchor value and metric type
- [ ] Recovery hints provide actionable guidance
- [ ] Audit trail logs anchor validation failures with context
- [ ] Unit tests cover all error scenarios (invalid HEAD, unknown revision, git errors)
- [ ] Integration tests verify materialization fails with invalid anchors
- [ ] Exit code 2 is returned for anchor validation failures
- [ ] Error messages are deterministic (no timestamps, locale-dependent formatting)

## Dependencies

- **Pre-existing**: design-docs/recurrence-anchor-error-messages.md (error format specification)
- **Pre-existing**: pkg/task/recurrence.go (current implementation)
- **Blocking**: None - this can be implemented independently
- **Follow-on tasks**:
  - Tg96jgm-add-validation-for-anchor-existence-at-recurrence (validation at creation time)
  - Tk1cdj4-add-audit-logging-for-default-anchor-values (audit logging enhancement)

## Implementation Tasks (for Developer Role)

This architect task creates the plan. Implementation tasks for developer role:

1. Create `pkg/task/errors.go` with error types and unified formatting
2. Create `pkg/task/logger.go` for audit trail logging
3. Modify `evaluateGitMetric` in `pkg/task/recurrence.go` to return errors
4. Update CLI command that materializes recurring tasks to handle errors properly
5. Add unit tests for all error scenarios
6. Add integration tests for materialization failures
7. Update CLI.md documentation for new error messages
8. Run full test suite and fix any failures

## Review Checklist

- [ ] Error messages follow unified format from recurrence-anchor-error-messages.md
- [ ] All error scenarios are covered (invalid HEAD, unknown revision, git errors)
- [ ] Recovery hints are actionable and helpful
- [ ] Audit trail provides sufficient context for debugging
- [ ] Tests assert exact string matches for error messages
- [ ] Exit codes are consistent (2 for validation failures)
- [ ] Documentation is updated (CLI.md, AGENTS.md if needed)
- [ ] No silent no-op behavior remains
