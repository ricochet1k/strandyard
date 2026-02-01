# Error Messages for --every Anchor Parsing

## Summary
Define the complete error message contract for `strand recurring add --every` anchor parsing, including specific error messages for each failure mode, recovery hints, and canonical examples for tests.

## Context
- design-docs/recurrence-anchor-error-messages-alternatives.md
- design-docs/recurrence-anchor-error-messages-reliability-review.md
- design-docs/recurrence-anchor-hint-examples.md
- design-docs/recurrence-anchor-flags-alternatives.md
- design-docs/recurrence-metrics.md
- CLI.md (recurring add section)

## Decision
Adopt Alternative B unified error format with structured reason + hint line. Error output goes to stderr only, with exit code 2 for --every parse/validation failures and 1 for other runtime errors.

## Error Format Contract

### Output Format
- **Line 1 (required)**: `strand: error: <reason>`
- **Line 2 (optional)**: `hint: <example>`
- **Channel**: stderr only
- **Exit codes**: `2` for `--every` parse/validation failures, `1` for other runtime errors

### Reason Format
The `<reason>` must include:
- The metric/unit being parsed
- The expected anchor type (date or commit)
- The problematic user-provided token (quoted)

## Specific Error Messages

### Missing Anchor (time-based metrics)

**Metrics**: `days`, `weeks`, `months`, `tasks_completed`

**Error**:
```
strand: error: invalid --every value: expected date anchor after 'from', got nothing
hint: --every "10 days"
```

**Notes**: Default to "from now" behavior; hint omits explicit anchor.

### Missing Anchor (commit-based metrics)

**Metrics**: `commits`, `lines_changed`

**Error**:
```
strand: error: invalid --every value: expected commit anchor after 'from', got nothing
hint: --every "50 commits from HEAD"
```

**Notes**: Default to "from HEAD"; hint uses canonical commit anchor.

### Malformed Date Anchor

**Error**:
```
strand: error: invalid date anchor '2026-13-01': invalid month
hint: --every "10 days from Jan 28 2026 09:00 UTC"
```

**Notes**: Use human-friendly canonical date in hint; specific parse error in reason.

### Malformed Commit Anchor

**Error**:
```
strand: error: invalid commit anchor 'not-a-hash': expected 40-character hex string
hint: --every "50 commits from HEAD"
```

**Notes**: Use `HEAD` in hint; explain expected format in reason.

### Unit/Anchor Mismatch (time metric with commit anchor)

**Error**:
```
strand: error: metric 'days' expects a date anchor, got commit hash '0123456789abcdef'
hint: --every "10 days from Jan 28 2026 09:00 UTC"
```

### Unit/Anchor Mismatch (commit metric with date anchor)

**Error**:
```
strand: error: metric 'commits' expects a commit anchor, got date '2026-01-28'
hint: --every "50 commits from HEAD"
```

### Ambiguous Anchor Type

**Error**:
```
strand: error: metric 'tasks_completed' requires a date anchor
hint: --every "20 tasks_completed from Jan 28 2026 09:00 UTC"
```

**Notes**: `tasks_completed` uses completion time (date anchor) not commit count.

### Invalid Metric

**Error**:
```
strand: error: invalid metric 'bogus': supported metrics are days, weeks, months, commits, lines_changed, tasks_completed
hint: --every "10 days"
```

### Invalid Interval Amount

**Error**:
```
strand: error: invalid interval 'zero': must be a positive integer
hint: --every "10 days"
```

### Missing Preposition

**Error**:
```
strand: error: invalid --every value: expected 'from' after anchor value, got '10 days Jan 28'
hint: --every "10 days from Jan 28 2026 09:00 UTC"
```

## Canonical Examples for Hint Lines

Use these stable examples in hint lines (from `recurrence-anchor-hint-examples.md`):

- `days`/`weeks`/`months` (default anchor): `--every "10 days"`
- `commits` (default anchor): `--every "50 commits from HEAD"`
- `lines_changed` (default anchor): `--every "500 lines_changed from HEAD"`
- `tasks_completed` (default anchor): `--every "20 tasks_completed"`
- Explicit date anchor: `--every "10 days from Jan 28 2026 09:00 UTC"`

## Anchor Type by Metric

| Metric | Expected Anchor Type | Default Anchor |
|--------|---------------------|----------------|
| `days` | date | now (implicit) |
| `weeks` | date | now (implicit) |
| `months` | date | now (implicit) |
| `commits` | commit | HEAD |
| `lines_changed` | commit | HEAD |
| `tasks_completed` | date | now (implicit) |

## Implementation Notes

### Parsing Behavior
- Time-based anchors: Use `github.com/olebedev/when` with UTC base time, full-string matches only, English rules (`en` + `common`).
- Commit anchors: Validate as 40-character hex string or `HEAD`.
- Missing anchors: Default to "now" (time) or "HEAD" (commit) for metrics that support defaults.

### Determinism
- All hint examples use deterministic strings (no timestamps, no locale-dependent formatting).
- Error messages quote user-provided tokens exactly.
- Avoid non-deterministic content in error text or hint examples.

### Testing
- Tests should assert exact string matches for both error and hint lines.
- Use canonical examples from this document for hint validation.
- Cover all failure modes listed above.

## Open Concerns (Captured for Resolution)

1. **CLI.md Outdated**: `CLI.md` still documents `--anchor` flag while Alternative D (`--every`) was adopted. File issue to update CLI.md to reflect the new flag format.

2. **Default Behavior Documentation**: The design doc states defaults apply (now/HEAD), but this was from Alternative C (which was not adopted for error messages). Clarify whether defaults apply for missing anchors or if anchors are always required. This is a decision for the Owner.

3. **Exit Code Convention**: New exit code `2` for parse/validation; align with CLI-wide exit code conventions (see subtask `Tc3zv-align-exit-code-conventions-for-every-failures`).

## Local Verification Steps

1. Add a recurring task with a malformed date anchor:
   ```bash
   strand add "Test" --every "10 days from 2026-13-01"
   ```
   Expected output to stderr:
   ```
   strand: error: invalid date anchor '2026-13-01': invalid month
   hint: --every "10 days from Jan 28 2026 09:00 UTC"
   ```

2. Add a recurring task with unit/anchor mismatch:
   ```bash
   strand add "Test" --every "10 commits from 2026-01-01"
   ```
   Expected output to stderr:
   ```
   strand: error: metric 'commits' expects a commit anchor, got date '2026-01-01'
   hint: --every "50 commits from HEAD"
   ```

3. Verify exit code:
   ```bash
   strand add "Test" --every "10 days from 2026-13-01"
   echo $?
   ```
   Expected: `2`
