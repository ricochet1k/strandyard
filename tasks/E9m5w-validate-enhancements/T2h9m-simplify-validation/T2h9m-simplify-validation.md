---
role: developer
parent: E9m5w-validate-enhancements
blockers:
  - T3k8p-link-validation
  - T7w4n-blocker-validation
date_created: 2026-01-27
date_edited: 2026-01-27
---

# Ensure Single Right Way Validation

## Summary

Remove any optional validation modes and ensure there's one right way to structure tasks. Keep validation simple and strict.

## Tasks

- [ ] Remove any lenient/strict mode options (if they exist or were planned)
- [ ] Ensure validation always fails fast on first error category
- [ ] Provide clear, actionable error messages
- [ ] No warnings - only errors that must be fixed
- [ ] Document the canonical task format in AGENTS.md
- [ ] Update error messages to be specific and helpful

## Acceptance Criteria

- No `--strict` or `--lenient` flags
- Validation either passes completely or fails with errors
- Each error message clearly states what's wrong and how to fix it
- No ambiguity in what constitutes a valid task
- Documentation clearly describes the one right way

## Files

- cmd/validate.go
- AGENTS.md (update canonical format)

## Error Message Examples

Good error messages:
- `ERROR: Task T3k7x-example missing required frontmatter field 'role'`
- `ERROR: Invalid ID format 'T123-bad' in tasks/T123-bad/ - must be <PREFIX><4-lowercase-alphanumeric>-<slug>`

Bad error messages:
- `WARNING: Task might have issues`
- `ERROR: Invalid format`

## Implementation Plan

### Architecture Overview

Review the existing validation implementation to ensure it follows "one right way" principles. The validation should be strict, clear, and fail-fast with no optional modes or warnings.

### Current State Analysis

**Existing implementation in [cmd/validate.go](../../cmd/validate.go)**:

Current flags:
- `--path`: tasks directory path (configuration, not validation mode) ✓
- `--roots`: output file for root tasks (configuration) ✓
- `--free`: output file for free tasks (configuration) ✓
- `--format`: output format text|json (output format, not validation mode) ✓

**Good**: No --strict or --lenient flags exist

**Current behavior**:
- Collects all validation errors
- Reports all errors at end
- Returns error exit code if any validation fails
- Regenerates master lists even when validation fails

### Issues to Address

#### Issue 1: Master Lists Generated Even on Validation Failure

**Current code** ([cmd/validate.go:51-53](../../cmd/validate.go#L51-L53)):
```go
// Generate master lists
if err := task.GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
    return fmt.Errorf("failed to generate master lists: %w", err)
}
```

This happens BEFORE checking if validation passed. This means invalid task trees can have their master lists updated.

**Problem**: Violates fail-fast principle. If tasks are invalid, we shouldn't update generated files.

**Fix**: Move master list generation AFTER validation check:

```go
func runValidate(tasksRoot, rootsFile, freeFile, outFormat string) error {
    // Parse all tasks
    parser := task.NewParser()
    tasks, err := parser.LoadTasks(tasksRoot)
    if err != nil {
        return fmt.Errorf("failed to load tasks: %w", err)
    }

    // Validate tasks
    validator := task.NewValidator(tasks)
    errors := validator.Validate()

    // Report errors and fail fast
    if len(errors) > 0 {
        reportErrors(errors, outFormat)
        return fmt.Errorf("validation failed: %d error(s)", len(errors))
    }

    // Only generate master lists if validation passed
    if err := task.GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
        return fmt.Errorf("failed to generate master lists: %w", err)
    }

    reportSuccess(tasks, rootsFile, freeFile, outFormat)
    return nil
}
```

**Rationale**:
- Fail fast: don't update anything if validation fails
- Clear behavior: validation must pass before generating files
- Prevents invalid state from being written to master lists

#### Issue 2: Error Messages Could Be More Specific

Current error messages are good but can be enhanced:

**Current**:
```go
Message: fmt.Sprintf("malformed ID: must be <PREFIX><4-lowercase-alphanumeric>-<slug> (e.g., T3k7x-example)")
```

**Enhanced**: Include the actual invalid ID in the message:
```go
Message: fmt.Sprintf("malformed ID '%s': must be <PREFIX><4-lowercase-alphanumeric>-<slug> (e.g., T3k7x-example)", id)
```

**Current**:
```go
Message: "missing role in frontmatter and no role found in first TODO"
```

**Enhanced**: Add suggestion for how to fix:
```go
Message: "missing role: add 'role: <rolename>' to frontmatter or '- [ ] (role: <rolename>)' as first TODO"
```

#### Issue 3: No Clear Documentation of Canonical Format

The task mentions updating AGENTS.md with canonical format documentation.

**Action items**:
1. Document required frontmatter fields
2. Document optional frontmatter fields
3. Document ID format requirements
4. Document blocker relationship requirements
5. Document role file requirements

### Implementation Steps

#### 1. Refactor validate command to fail fast

**File**: [cmd/validate.go](../../cmd/validate.go)

Changes:
- Extract error reporting to helper function `reportErrors()`
- Extract success reporting to helper function `reportSuccess()`
- Move master list generation after validation
- Only generate if validation passed

```go
func reportErrors(errors []error, outFormat string) {
    if outFormat == "json" {
        errMsgs := make([]string, len(errors))
        for i, e := range errors {
            errMsgs[i] = e.Error()
        }
        b, _ := json.MarshalIndent(map[string]interface{}{"errors": errMsgs}, "", "  ")
        fmt.Println(string(b))
    } else {
        for _, e := range errors {
            fmt.Println("ERROR:", e.Error())
        }
    }
}

func reportSuccess(tasks map[string]*Task, rootsFile, freeFile, outFormat string) {
    if outFormat == "json" {
        roots := []string{}
        free := []string{}
        for id, t := range tasks {
            if t.Meta.Parent == "" {
                roots = append(roots, id)
            }
            if len(t.Meta.Blockers) == 0 {
                free = append(free, id)
            }
        }
        b, _ := json.MarshalIndent(map[string]interface{}{"roots": roots, "free": free}, "", "  ")
        fmt.Println(string(b))
    } else {
        fmt.Println("validate: ok")
        fmt.Printf("Validated %d tasks\n", len(tasks))
        fmt.Printf("Master lists updated: %s, %s\n", rootsFile, freeFile)
    }
}
```

#### 2. Improve error messages

**File**: [pkg/task/validate.go](../../pkg/task/validate.go)

Update error messages to be more specific and actionable:

**ID validation** ([pkg/task/validate.go:64-70](../../pkg/task/validate.go#L64-L70)):
```go
func (v *Validator) validateID(id string, task *Task) {
    if !v.idPattern.MatchString(id) {
        v.errors = append(v.errors, ValidationError{
            TaskID:  id,
            File:    task.FilePath,
            Message: fmt.Sprintf(
                "malformed ID '%s': must match <PREFIX><4-base36>-<slug> (e.g., T3k7x-example, E2h9w-epic-name)",
                id,
            ),
        })
    }
}
```

**Role validation** ([pkg/task/validate.go:73-94](../../pkg/task/validate.go#L73-L94)):
```go
func (v *Validator) validateRole(id string, task *Task) {
    role := task.GetEffectiveRole()
    if role == "" {
        v.errors = append(v.errors, ValidationError{
            TaskID:  id,
            File:    task.FilePath,
            Message: "missing role: add 'role: developer' to frontmatter or '- [ ] (role: developer)' as first TODO",
        })
        return
    }

    // Check if role file exists
    rolePath := filepath.Join("roles", role+".md")
    if _, err := os.Stat(rolePath); os.IsNotExist(err) {
        v.errors = append(v.errors, ValidationError{
            TaskID:  id,
            File:    task.FilePath,
            Message: fmt.Sprintf(
                "role file %s does not exist - create it or use an existing role",
                rolePath,
            ),
        })
    }
}
```

**Parent validation** ([pkg/task/validate.go:96-109](../../pkg/task/validate.go#L96-L109)):
```go
func (v *Validator) validateParent(id string, task *Task) {
    if task.Meta.Parent == "" {
        return // Root task, no parent to validate
    }

    if _, exists := v.tasks[task.Meta.Parent]; !exists {
        v.errors = append(v.errors, ValidationError{
            TaskID:  id,
            File:    task.FilePath,
            Message: fmt.Sprintf(
                "parent task '%s' does not exist - check parent field in frontmatter",
                task.Meta.Parent,
            ),
        })
    }
}
```

**Blocker validation** ([pkg/task/validate.go:111-126](../../pkg/task/validate.go#L111-L126)):
```go
func (v *Validator) validateBlockers(id string, task *Task) {
    for _, blocker := range task.Meta.Blockers {
        if blocker == "" {
            continue
        }

        if _, exists := v.tasks[blocker]; !exists {
            v.errors = append(v.errors, ValidationError{
                TaskID:  id,
                File:    task.FilePath,
                Message: fmt.Sprintf(
                    "blocker task '%s' does not exist - check blockers field in frontmatter",
                    blocker,
                ),
            })
        }
    }
}
```

#### 3. Document canonical format in AGENTS.md

**File**: [AGENTS.md](../../AGENTS.md)

Add section documenting the canonical task format (if not already present):

```markdown
## Task Format Specification

### Required Frontmatter Fields

All tasks must have YAML frontmatter with these fields:

- `role`: The role that should work on this task (must match a file in roles/)
  - Can be specified in frontmatter OR as first TODO: `- [ ] (role: developer)`
- `parent`: Parent task ID (empty string "" for root tasks)
- `blockers`: Array of task IDs that must be completed first (empty array [] if none)

### Optional Frontmatter Fields

- `blocks`: Array of task IDs that this task blocks
- `date_created`: ISO 8601 timestamp
- `date_edited`: ISO 8601 timestamp
- `owner_approval`: Boolean indicating owner approved this task
- `completed`: Boolean indicating task is done

### Task ID Format

Task IDs must follow this pattern: `<PREFIX><TOKEN>-<SLUG>`

- PREFIX: Single uppercase letter (T, E, D, etc.)
  - T: Task
  - E: Epic
  - D: Design doc
- TOKEN: Exactly 4 base36 characters (0-9, a-z lowercase)
  - Example: 3k7x, 2h9w, 5n2q
- SLUG: Descriptive name using lowercase letters, numbers, and hyphens
  - Example: add-validation, update-docs, create-user

Examples: `T3k7x-example`, `E2h9w-metadata-format`, `D5n2q-architecture`

### Blocker Relationships

Blocker relationships must be bidirectional:

- If task A has `blockers: [B]`, then task B must have `blocks: [A]`
- If task A has `blocks: [C]`, then task C must have `blockers: [A]`

### File Structure

Tasks can have their markdown file named in order of preference:
1. `<task-id>/<task-id>.md` (e.g., `T3k7x-example/T3k7x-example.md`)
2. `<task-id>/task.md`
3. `<task-id>/README.md`
```

#### 4. Ensure no warnings, only errors

**Current state**: The codebase only uses ValidationError, no warnings ✓

**Maintain this**: Never add warnings. If something is wrong, it's an error.

### Testing Strategy

Test the fail-fast behavior:
1. Create tasks with validation errors
2. Run validate command
3. Verify it returns error exit code
4. Verify master lists are NOT updated
5. Fix errors
6. Run validate again
7. Verify it succeeds
8. Verify master lists ARE updated

Test error messages:
1. Verify all error messages include specific details
2. Verify all error messages suggest how to fix
3. Verify error messages reference the correct file paths

### Technical Considerations

**Backward Compatibility**:
- Moving master list generation after validation could break workflows that depend on lists being updated even when validation fails
- This is acceptable because:
  - It's the correct behavior (fail-fast)
  - Invalid state shouldn't be persisted
  - Users should fix validation errors, then lists will update

**Exit Codes**:
- Current implementation returns error if validation fails ✓
- This is correct for scripts/CI integration
- Maintain this behavior

**Output Format Flag**:
- The `--format` flag is acceptable
- It controls output format, not validation behavior
- Both text and json report the same errors
- Keep this flag

### Alternative Approaches Considered

1. **Add --fix flag to auto-correct errors**
   - Pro: Convenient
   - Con: Violates "one right way" - unclear what's correct
   - Decision: REJECTED - user must fix manually

2. **Add warnings for potential issues**
   - Pro: More gradual migration
   - Con: Ambiguity about what needs fixing
   - Decision: REJECTED - everything is error or valid

3. **Update master lists even on failure**
   - Pro: Files stay in sync
   - Con: Persists invalid state
   - Decision: REJECTED - fail fast is better

### Implementation Priority

This task is blocked on T3k8p and T7w4n because:
- We're adding new validations (link and blocker)
- Those should be integrated first
- Then we refactor the overall validation flow
- This ensures we don't have to refactor twice

**Implementation order**:
1. Wait for T3k8p and T7w4n to be implemented
2. Refactor validate command to fail fast
3. Improve all error messages
4. Document canonical format in AGENTS.md
5. Test fail-fast behavior
