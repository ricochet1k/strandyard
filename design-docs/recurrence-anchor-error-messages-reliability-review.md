---
role: reviewer-reliability
priority: medium
---

# Review: Error message format contract for --every anchor parsing (Reliability)

## Artifacts
- design-docs/recurrence-anchor-error-messages-alternatives.md
- design-docs/recurrence-metrics.md
- CLI.md (recurring add section)

## Scope
Define the error output contract for parsing `--every` anchors in `memmd recurring add`.
Implementation details, tests, and the final alternative selection are out of scope.

## Review Focus
- Deterministic error formatting for automation
- stderr/stdout behavior and exit codes
- Failure-mode coverage for invalid anchors and mismatched units
- Compatibility with future tooling and logging

## Findings
- A stable, two-line error format reduces test churn while keeping guidance actionable.
- Errors should be emitted on stderr with a non-zero exit code and no stdout output to keep scripts deterministic.
- Avoid non-deterministic content in error text; examples should be static per unit type.
- Reason text should always include the unit/metric and expected anchor type to reduce ambiguous retries.

## Proposed Contract (non-binding)
- Line 1: `invalid --every value: <reason>` (single line, lowercase prefix).
- Line 2 (optional): `hint: --every "<interval> <unit> from <anchor>"`.
- `<reason>` uses single quotes around user-provided tokens and names the expected anchor type (`date` or `commit`).
- Emit on stderr only; exit code `1` (or a documented constant) for parse/validation failures.

## Concerns captured as subtasks
- Decision: finalize error line prefix, stderr/stdout behavior, and exit code contract.
- Decision: define canonical example anchors for hint lines per unit.
