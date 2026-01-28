---
type: ""
role: developer
priority: ""
parent: E9m5w-validate-enhancements
blockers: []
blocks:
  - T2h9m-simplify-validation
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-28T04:59:23.327776Z
owner_approval: false
completed: true
---

# Add Blocker Status Validation

## Summary

Repair that blocker relationships are bidirectional and consistent, and that free-tasks.md accurately reflects tasks with no blockers.

## Tasks

- [ ] Repair bidirectional blocker relationships:
  - If task A blocks task B, then task B should list A in blockers
  - If task A has blocker B, then task B should list A in blocks
- [ ] Repair free-tasks.md only contains tasks with empty blockers array
- [ ] Repair tasks in free-tasks.md actually exist
- [ ] Report inconsistent blocker relationships with clear errors
- [ ] Suggest fixes for broken relationships

## Acceptance Criteria

- Detects missing bidirectional blocker links
- Detects tasks in free-tasks.md that have blockers
- Detects tasks not in free-tasks.md that should be (no blockers)
- Example errors:
  - `ERROR: Task T3k7x has blocker T5h7w, but T5h7w doesn't list T3k7x in blocks`
  - `ERROR: Task T8n2m is in free-tasks.md but has blocker T6p4k`
  - `ERROR: Task T2h9m has no blockers but is not in free-tasks.md`

## Files

- cmd/validate.go

## Bidirectional Validation Example

Task A:
```yaml
blockers:
  - T5h7w-example
```

Task B (T5h7w-example):
```yaml
blocks:
  - Ta3k7x-taskname  # Should list Task A
```

## Implementation Plan

### Architecture Overview

Add bidirectional blocker validation to the existing `Validator` struct and update the `GenerateMasterLists` validation logic. This ensures blocker relationships are consistent and that the free-tasks.md file accurately reflects unblocked tasks.

### Implementation Steps

#### 1. Add Bidirectional Blocker Validation

**File**: [pkg/task/validate.go](../../pkg/task/validate.go)

Add new validation method:

```go
// repairBidirectionalBlockers ensures blocker relationships are bidirectional
func (v *Validator) repairBidirectionalBlockers(id string, task *Task) {
    // For each blocker B that task A lists:
    //   - Verify B.blocks contains A
    // For each task C in A.blocks:
    //   - Verify C.blockers contains A
}
```

Call from `Repair()` method:

```go
func (v *Validator) Repair() []ValidationError {
    v.errors = []ValidationError{}

    for id, task := range v.tasks {
        v.repairID(id, task)
        v.repairRole(id, task)
        v.repairParent(id, task)
        v.repairBlockers(id, task)
        v.repairTaskLinks(id, task)
        v.repairBidirectionalBlockers(id, task)  // Add this
    }

    return v.errors
}
```

#### 2. Blocker Relationship Validation Logic

The validation should check:

**A. For each blocker listed by a task:**
```go
for _, blockerID := range task.Meta.Blockers {
    blocker, exists := v.tasks[blockerID]
    if !exists {
        // Already caught by repairBlockers
        continue
    }

    // Check if blocker lists this task in its blocks field
    if !contains(blocker.Meta.Blocks, id) {
        v.errors = append(v.errors, ValidationError{
            TaskID:  id,
            File:    task.FilePath,
            Message: fmt.Sprintf(
                "task has blocker %s, but %s doesn't list this task in blocks field",
                blockerID, blockerID,
            ),
        })
    }
}
```

**B. For each task this one blocks:**
```go
for _, blockedID := range task.Meta.Blocks {
    blocked, exists := v.tasks[blockedID]
    if !exists {
        v.errors = append(v.errors, ValidationError{
            TaskID:  id,
            File:    task.FilePath,
            Message: fmt.Sprintf("blocks non-existent task %s", blockedID),
        })
        continue
    }

    // Check if blocked task lists this in its blockers field
    if !contains(blocked.Meta.Blockers, id) {
        v.errors = append(v.errors, ValidationError{
            TaskID:  id,
            File:    task.FilePath,
            Message: fmt.Sprintf(
                "task blocks %s, but %s doesn't list this task in blockers field",
                blockedID, blockedID,
            ),
        })
    }
}
```

**Helper function:**
```go
// contains checks if a string slice contains a value
func contains(slice []string, val string) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}
```

#### 3. Free Tasks List Validation

Currently, `GenerateMasterLists` generates free-tasks.md but doesn't repair it. We need to add validation that runs during the repair command.

**Two approaches considered:**

**Option A: Repair the generated file matches reality**
- After generating free-tasks.md, read it back
- Parse the task list
- Compare with actual free tasks
- Report discrepancies

**Option B: Don't repair the file, just ensure it's regenerated**
- The repair command already regenerates free-tasks.md
- If it's out of sync, running repair fixes it automatically
- This is the "one right way" - repair regenerates truth

**Decision: Option B (simpler)**
- The repair command already regenerates both master lists
- If free-tasks.md is out of date, repair fixes it
- No need for additional validation of the file itself
- Focus validation on task metadata correctness

**What we DO repair:**
- Blocker relationships are bidirectional (covered above)
- Tasks with no blockers can be correctly identified (already works)
- Master lists are regenerated on every repair run (already works)

#### 4. Error Message Format

Bidirectional blocker errors:

```
ERROR: Task T3k7x-example: task has blocker T5h7w-task, but T5h7w-task doesn't list this task in blocks field
```

```
ERROR: Task T5h7w-task: task blocks T3k7x-example, but T3k7x-example doesn't list this task in blockers field
```

Non-existent blocked task:

```
ERROR: Task T5h7w-task: blocks non-existent task T9999-missing
```

These are actionable - they tell the user exactly what's wrong and what fields need to be updated.

#### 5. Handling Completed Tasks

**Question**: Should completed tasks be repaired for blocker relationships?

**Answer**: YES
- Completed tasks may still have blocker/blocks fields
- Validation ensures historical consistency
- If relationships are wrong, they should be fixed even for completed tasks
- The `completed` field only affects master list generation, not validation

#### 6. Testing Strategy

Test cases:
- Task A blocks B, B lists A as blocker ✓ (valid)
- Task A blocks B, B doesn't list A as blocker ✗ (error)
- Task A has blocker B, B doesn't list A in blocks ✗ (error)
- Task A blocks non-existent task ✗ (error)
- Completed tasks with invalid blocker relationships ✗ (error)
- Empty blocker/blocks fields ✓ (valid)
- Circular blockers (A blocks B, B blocks A) ✓ (valid, but unusual)

### Technical Considerations

**Performance**:
- For each task, iterate through its blockers and blocks arrays
- Lookup other tasks in map (O(1))
- Overall O(n*m) where n=tasks, m=average blocker count
- Acceptable for typical projects

**Order Independence**:
- Validation may report the same relationship error twice (once from each side)
- This is acceptable - it makes errors clearer
- Alternative: deduplicate errors, but adds complexity
- Decision: Keep it simple, allow duplicate errors

**Data Model Clarification**:
Looking at the Metadata struct:
```go
type Metadata struct {
    Blockers []string `yaml:"blockers"`  // Tasks that block this one
    Blocks   []string `yaml:"blocks"`    // Tasks that this one blocks
}
```

**Relationship semantics**:
- If task A has `blockers: [B]`, then B must be completed before A can start
- If task A has `blocks: [C]`, then A must be completed before C can start
- These should be bidirectional:
  - A.blockers contains B ⟺ B.blocks contains A
  - A.blocks contains C ⟺ C.blockers contains A

### Alternative Approaches Considered

1. **Auto-fix blocker relationships**
   - Pro: Convenient for users
   - Con: Magic behavior, unclear which side is "correct"
   - Decision: REJECTED - validation should not modify files

2. **Warn instead of error for blocker mismatches**
   - Pro: Less strict
   - Con: Violates "one right way" principle
   - Decision: REJECTED - must be an error to fix

3. **Repair free-tasks.md file content**
   - Pro: Catches stale files
   - Con: File is auto-generated by repair, adds redundant checking
   - Decision: REJECTED - regeneration is sufficient

### Implementation Order

1. Add `repairBidirectionalBlockers` method with blocker→blocks validation
2. Add blocks→blocker validation in same method
3. Add non-existent blocked task validation
4. Add helper `contains` function
5. Wire into `Repair()` method
6. Add tests
