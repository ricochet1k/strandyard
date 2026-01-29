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
Implementation details and test implementation are out of scope.

## Review Focus
- Deterministic error formatting for automation
- stderr/stdout behavior and exit codes
- Failure-mode coverage for invalid anchors and mismatched units
- Compatibility with future tooling and logging

## Findings
- A stable, two-line error format reduces test churn while keeping guidance actionable for automation.
- Error output should be stderr-only with explicit exit codes to keep scripts deterministic.
- Avoid non-deterministic content in error text or hint examples; use canonical anchors per unit type.
- Reason text should always include the unit/metric and expected anchor type to reduce ambiguous retries.

## Contract Alignment (per decision doc)
- Line 1 prefix: `memmd: error: ` followed by a single-line reason.
- Line 2 (optional): `hint: ` followed by a minimal example.
- `<reason>` should include the unit/metric and expected anchor type (`date` or `commit`) and quote user-provided tokens.
- Emit on stderr only; exit code `2` for `--every` parse/validation failures and `1` for other runtime errors.

## Concerns captured as subtasks
- Decision: define canonical, deterministic anchor examples per unit type for hint lines (T4izq-define-deterministic-anchor-examples-for-every-hin).
- Decision: align CLI-wide exit code conventions with the new `--every` contract and document them (Tc3zv-align-exit-code-conventions-for-every-failures).
