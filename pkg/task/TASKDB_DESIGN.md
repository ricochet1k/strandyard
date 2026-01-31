# TaskDB Design Review

## Problem Statement

The current TaskDB implementation has fundamental design flaws that violate the repository's core principles and create opportunities for misuse.

## Core Repository Principles (That Were Violated)

1. **Tasks should ONLY be created from templates** - never completely blank
2. **Parent/title/description should already be filled in** when creating tasks
3. **GetOrCreate should NOT exist** - it creates blank tasks, violating principle #1
4. **The library should be hard to use incorrectly** - manual field manipulation breaks relationships
5. **GoDoc-visible APIs must be either:**
   - Hard to misuse, OR
   - Clearly document pitfalls

## Current Implementation Issues

### 1. Manual Task Creation/Modification
- Direct `*Task` creation is possible
- Manual field editing (e.g., `task.Meta.Parent = "foo"`) bypasses relationship integrity
- This breaks the "hard to use incorrectly" principle

### 2. GetOrCreate Method
- **SHOULD NOT EXIST**
- Creates blank tasks with empty metadata
- Violates template-based task creation principle
- Example of code that should never have been written

### 3. Relationship Management Confusion

#### Good Examples (Keep These Patterns)
- `AddBlocker` / `RemoveBlocker` - Lazy, maintains bidirectional relationships automatically
- These are the RIGHT approach

#### Problematic Areas
- `UpdateBlockersFromChildren` - Now poorly named after changes
- `FixBlockerRelationships` - Duplicates/copies behavior from UBFC
- `SyncBlockersFromChildren` - Just calls UBFC, adding another layer
- Confusion about which function does what

### 4. Code That Was Copied vs. Designed
- Parts of the implementation copied existing code without understanding the design
- Lost sight of the original request during implementation
- Need to distinguish between:
  - Original, well-designed code (UpdateBlockersFromChildren)
  - New code that should integrate properly
  - Code that shouldn't exist at all (GetOrCreate)

## What Needs to Happen

### Phase 1: Complete Code Review
Review ALL existing task management code to understand:
- What exists and why
- What follows good patterns
- What violates principles
- What creates opportunities for misuse

### Phase 2: Design Proper API Surface
Determine what should be:
- Public (exported, in GoDoc)
- Private (internal use only)
- Not exist at all

### Phase 3: Enforce Relationship Integrity
Design mechanisms that make it:
- Impossible to manually break relationships
- Clear when you're doing something dangerous
- Natural to do the right thing

### Phase 4: Clean Implementation
Remove, refactor, and properly integrate the code

## Key Questions to Answer

1. How do we prevent manual `*Task` creation?
2. How do we prevent direct field manipulation?
3. What's the proper naming and responsibility split for blocker syncing?
4. What belongs in TaskDB vs. what belongs elsewhere?
5. How do tasks get created (only from templates)?
6. What operations should TaskDB support?
7. What should be impossible to do?

## Success Criteria

- No way to create tasks except through templates
- No way to break relationships through the API
- Clear, well-named functions with single responsibilities
- GoDoc that makes correct usage obvious
- Remove all code that violates principles
