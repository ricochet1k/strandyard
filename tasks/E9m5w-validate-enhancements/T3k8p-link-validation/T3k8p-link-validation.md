---
type: ""
role: developer
priority: ""
parent: E9m5w-validate-enhancements
blockers: []
blocks:
  - T2h9m-simplify-validation
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-28T04:45:47.520113Z
owner_approval: false
completed: true
---

# Add Task Link Validation

## Summary

Add validation to ensure all task references/links in task content point to existing tasks.

## Tasks

- [ ] Scan task content for task links (markdown links, direct references)
- [ ] Extract task IDs from links (format: `[text](path/to/T3k7x-task/file.md)`)
- [ ] Verify each referenced task ID exists in tasks directory
- [ ] Report broken links with clear error messages showing file and line number
- [ ] Handle both relative and absolute paths
- [ ] Add tests for link validation

## Acceptance Criteria

- Detects broken task links
- Reports file path and line number of broken link
- Repairs links to task directories and task files
- Example error: `ERROR: Broken link in T3k7x-example/T3k7x-example.md:15: T9999-missing does not exist`

## Files

- cmd/validate.go
- pkg/metadata/linkchecker.go (new, optional)

## Link Formats to Repair

- Markdown links: `[Task Name](tasks/T3k7x-example/T3k7x-example.md)`
- Direct references: `T3k7x-example`
- Parent field: `parent: E2k7x-metadata-format`
- Blocker field: `blockers: [T5h7w-default-free-task]`

## Implementation Plan

### Architecture Overview

Add a new validation method `repairTaskLinks()` to the existing `Validator` struct in [pkg/task/validate.go](../../pkg/task/validate.go). This maintains the existing validation pattern and keeps all validation logic centralized.

### Implementation Steps

#### 1. Add Link Validation Method to Validator

**File**: [pkg/task/validate.go](../../pkg/task/validate.go)

Add new method to Validator struct:

```go
// repairTaskLinks scans task content for references to other tasks and verifies they exist
func (v *Validator) repairTaskLinks(id string, task *Task) {
    // Implementation details below
}
```

Call this from `Repair()` method after existing validations:

```go
func (v *Validator) Repair() []ValidationError {
    v.errors = []ValidationError{}

    for id, task := range v.tasks {
        v.repairID(id, task)
        v.repairRole(id, task)
        v.repairParent(id, task)
        v.repairBlockers(id, task)
        v.repairTaskLinks(id, task)  // Add this line
    }

    return v.errors
}
```

#### 2. Link Extraction Logic

The `repairTaskLinks` method should:

1. **Split content into lines** for line number tracking
2. **Scan each line** for task references using regex patterns
3. **Extract task IDs** from matched patterns
4. **Verify existence** of each referenced task
5. **Report errors** with file path and line number

**Regex patterns needed**:
- Markdown links: `\[([^\]]+)\]\(([^)]+)\)` - extract path from group 2
- Task ID extraction from path: Extract directory name matching `^[A-Z][0-9a-z]{4}-[a-zA-Z0-9-]+$`

#### 3. Task ID Extraction from Paths

Create helper function:

```go
// extractTaskIDFromPath extracts a task ID from a file or directory path
// Examples:
//   - "tasks/T3k7x-example/T3k7x-example.md" -> "T3k7x-example"
//   - "../T5h7w-task/task.md" -> "T5h7w-task"
//   - "T2k9p-other" -> "T2k9p-other"
func extractTaskIDFromPath(path string) string {
    // Clean path and split by directory separator
    // Scan path components for ID pattern
    // Return first match or empty string
}
```

This function should:
- Handle both Unix and Windows path separators
- Work with relative and absolute paths
- Match against the existing ID pattern in Validator
- Return empty string if no task ID found

#### 4. Error Reporting Format

Use the existing `ValidationError` struct but enhance the error message:

```go
v.errors = append(v.errors, ValidationError{
    TaskID:  id,
    File:    task.FilePath,
    Message: fmt.Sprintf("broken link at line %d: task %s does not exist", lineNum, referencedID),
})
```

This produces output like:
```
ERROR: Task T3k7x-example: broken link at line 15: task T9999-missing does not exist
```

#### 5. Scope Decisions

**Include**:
- Markdown links containing task IDs in the path
- The existing parent and blocker validations already check those fields

**Exclude for MVP**:
- Bare task ID references in content (e.g., "See T3k7x for details") - these are harder to distinguish from false positives like "T3k7x" appearing in code examples
- Can be added in future iteration if needed

**Why this scope**:
- Markdown links are explicit and intentional
- Low false positive rate
- High value - catches broken documentation links
- Parent/blocker fields already repaired

#### 6. Testing Strategy

Add test cases to cover:
- Valid markdown links to existing tasks
- Broken markdown links to non-existent tasks
- Links with various path formats (relative, with tasks/ prefix, etc.)
- Tasks with no links (should not error)
- Edge cases: empty paths, malformed links

### Technical Considerations

**Performance**:
- Regex compilation should happen once (compile at Validator creation or as package-level var)
- Line-by-line scanning is O(n) where n = lines of content across all tasks
- Acceptable for typical project sizes (100s-1000s of tasks)

**Accuracy**:
- Use the same ID pattern already defined in Validator (`v.idPattern`)
- Avoid false positives by only checking markdown link syntax
- Clear error messages help users fix issues quickly

**Maintainability**:
- Keep all validation logic in one file ([pkg/task/validate.go](../../pkg/task/validate.go))
- Reuse existing ValidationError structure
- Follow existing code patterns in Validator

### Alternative Approaches Considered

1. **Using goldmark AST traversal**
   - Pro: More robust markdown parsing
   - Con: Harder to get line numbers, more complex code
   - Decision: REJECTED - regex is simpler for this use case

2. **Separate linkchecker package**
   - Pro: Better separation of concerns
   - Con: Adds indirection, splits validation logic
   - Decision: REJECTED - keep validation centralized

3. **Repair all bare task ID references**
   - Pro: Catches more broken references
   - Con: High false positive rate (IDs in code examples, commit messages, etc.)
   - Decision: DEFER - start with markdown links only, can expand later
