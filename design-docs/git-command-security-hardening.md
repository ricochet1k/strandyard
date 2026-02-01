# Implementation Plan — Git Command Security Hardening

## Architecture Overview
Harden all calls to the `git` executable to prevent flag injection vulnerabilities. This ensures that user-provided strings (anchors, commits, presets) are always treated as positional arguments and never interpreted as flags for the `git` command or its subcommands.

## Specific Files to Modify
- `pkg/task/recurrence.go`: Update all `exec.Command` calls to use `--` separator where appropriate.
- `cmd/init.go`: Update `git clone` to use `--` separator.
- `cmd/projects.go`: Ensure `git rev-parse --show-toplevel` is safe (it doesn't currently take user input, but good to check).
- `pkg/task/recurrence_test.go`: Add regression tests for flag injection.

## Approach
1.  **Harden `pkg/task/recurrence.go`**:
    - For `git rev-list`, `git diff`, `git rev-parse`, and `git show`, insert the `--` separator before user-controlled revision/commit arguments.
    - Example: `exec.Command("git", "rev-list", "--count", "--", fmt.Sprintf("%s..HEAD", anchor))`
2.  **Harden `cmd/init.go`**:
    - Update `applyPreset` to use `exec.Command("git", "clone", "--depth", "1", "--", preset, tempDir)`.
3.  **Harden `pkg/task/recurrence_test.go`**:
    - Add a test case that attempts to use a flag-like string (e.g., `--help`) as an anchor and verify it is treated as an invalid revision rather than triggering help output or other unexpected behavior.

## Integration Points
- This hardening is transparent to users but protects against malicious task configurations or CLI inputs.

## Testing Approach
- Unit tests in `pkg/task/recurrence_test.go` using mock/temp git repos.
- Verification that `strand init --preset="--help"` results in a "failed to clone" error (treats it as a URL/path) rather than printing git clone help.

## Decision Rationale
- Standard security practice for CLI tools calling other CLI tools.
- Go's `exec.Command` prevents shell injection but not flag injection; the `--` separator is the standard way to fix this in most POSIX-compliant tools.

## Child Tasks
- [ ] (role: developer) Harden git command calls in `pkg/task/recurrence.go`.
- [ ] (role: developer) Harden git clone in `cmd/init.go`.
- [ ] (role: developer) Add regression tests for git flag injection in recurrence.
- [ ] (role: master-reviewer) Review git security hardening.
