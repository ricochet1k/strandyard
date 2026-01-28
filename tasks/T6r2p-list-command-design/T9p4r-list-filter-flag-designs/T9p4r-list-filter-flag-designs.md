---
type: ""
role: designer
priority: medium
parent: T6r2p-list-command-design
blockers: []
blocks: []
date_created: 2026-01-28T01:52:27Z
date_edited: 2026-01-28T05:20:29.154789Z
owner_approval: false
completed: true
---

# Evaluate list filter flag design alternatives

## Context
Issue: tasks/I5c1s-explore-alternate-list-filter-flag-designs/I5c1s-explore-alternate-list-filter-flag-designs.md
Design doc: design-docs/list-command.md (Filters section)
Current flags: cmd/list.go

## Current Implementation Analysis

The current implementation uses boolean flags with optional `--flag=false` for negation:

```bash
--completed      # filter by completed status
--blocked        # filter by blocked status (has blockers)
--blocks         # filter by blocks status (has blocks)
--owner-approval # filter by owner approval
```

### Problems with current approach:
1. **Non-obvious negation**: Users need to know `--completed=false` to find incomplete tasks
2. **Help clarity**: `--help` doesn't show negation options clearly
3. **Verbosity**: `--completed=false` is longer than dedicated negative flags
4. **Discoverability**: Users might not realize boolean flags accept `false` values

## Alternative Filter Flag Designs

### Alternative 1: Positive/Negative Flag Pairs

```bash
# Instead of: --completed [true|false]
--completed     # show completed tasks
--incomplete    # show incomplete tasks

# Instead of: --blocked [true|false]  
--blocked       # show blocked tasks
--unblocked     # show unblocked tasks

# Instead of: --blocks [true|false]
--blocking      # show tasks that block others
--not-blocking  # show tasks that don't block others

# Instead of: --owner-approval [true|false]
--needs-approval # show tasks needing owner approval
--approved      # show tasks with owner approval
```

**Pros**:
- Explicit and self-documenting
- Clear help text (`--help` shows both positive and negative options)
- Shorter negation (`--incomplete` vs `--completed=false`)
- Intuitive mental model

**Cons**:
- Flag proliferation (doubles the number of boolean flags)
- Potential confusion if both positive and negative flags are used together
- More complex validation logic needed

**Risk**: Low - well-established pattern in CLI tools

### Alternative 2: Single Filter Flag with Key=Value Pairs

```bash
--filter completed=true
--filter completed=false
--filter blocked=true
--filter role=developer
--filter priority=high
--filter owner-approval=true

# Can be repeated for multiple filters:
--filter completed=false --filter role=developer --filter priority=high
```

**Pros**:
- Unified interface for all filters
- Extensible for future filter types
- Consistent syntax
- Easy to parse and repair
- Supports multiple values naturally

**Cons**:
- More verbose for simple cases
- Less discoverable (users need to know available filter keys)
- Breaks from current flag pattern
- Help text becomes more complex

**Risk**: Medium - requires users to learn filter keys

### Alternative 3: Status/State Enumeration Flag

```bash
--status completed|incomplete|blocked|unblocked|ready|waiting
# Can be repeated:
--status incomplete --status unblocked

# or comma-separated:
--status incomplete,unblocked
```

**Pros**:
- Concise for common combinations
- Natural language concepts
- Easy to extend with new status concepts
- Good readability

**Cons**:
- Ambiguous status definitions
- Overlaps with existing filters
- May confuse with `completed` field
- Limited expressiveness for complex combinations

**Risk**: High - introduces new abstraction that may confuse

### Alternative 4: Current Pattern with Helper Aliases

Keep current boolean flags but add negative aliases:

```bash
# Current flags remain:
--completed [true|false]
--blocked [true|false]  
--blocks [true|false]
--owner-approval [true|false]

# Add aliases for common negations:
--incomplete    # alias for --completed=false
--unblocked     # alias for --blocked=false
--not-blocking  # alias for --blocks=false
--approved      # alias for --owner-approval=true
```

**Pros**:
- Full backward compatibility
- Improves discoverability
- Allows gradual migration
- Simple implementation

**Cons**:
- Two ways to do the same thing
- Help text complexity
- Potential user confusion

**Risk**: Low - conservative evolutionary change

## Comparison Matrix

| Alternative | Usability | Composability | Backward Compatibility | Help Clarity | Implementation Complexity |
|-------------|-----------|---------------|----------------------|--------------|--------------------------|
| Current     | Medium    | High          | High                 | Low          | Low                      |
| Flag Pairs  | High      | Medium        | Low                  | High         | Medium                   |
| --filter    | Medium    | High          | Low                  | Medium       | Medium                   |
| --status    | Medium    | Low           | Low                  | Medium       | High                     |
| Aliases     | High      | High          | High                 | High         | Low                      |

## Recommendations

### Primary Recommendation: Alternative 4 (Aliases)

Add negative aliases while preserving current boolean flags:

**Rationale**:
1. **Zero disruption**: Existing scripts and documentation continue working
2. **Improved UX**: Users get clear `--incomplete`, `--unblocked` options
3. **Easy learning**: Both patterns available, users can choose preference
4. **Simple implementation**: Just flag aliases, no complex logic changes

**Implementation approach**:
```bash
# Keep existing flags unchanged
--completed [true|false]
--blocked [true|false]
--blocks [true|false]
--owner-approval [true|false]

# Add these aliases:
--incomplete    # maps to --completed=false
--unblocked     # maps to --blocked=false
--not-blocking  # maps to --blocks=false
--approved      # alias for --owner-approval=true
```

### Secondary Recommendation (Future): Alternative 1 (Flag Pairs)

If breaking changes become acceptable, migrate to explicit positive/negative pairs:

**Rationale**:
- Best long-term UX
- Most discoverable
- Industry standard pattern

**Migration strategy**:
1. Phase 1: Add aliases (primary recommendation)
2. Phase 2: Deprecate `--flag=false` syntax with warnings
3. Phase 3: Remove boolean value support, keep only positive flags
4. Phase 4: Add negative flags as separate options

## Updated Help Text Example

```bash
Filter flags:
  --completed           Show completed tasks
  --incomplete          Show incomplete tasks (alias: --completed=false)
  --blocked             Show tasks with blockers
  --unblocked           Show tasks without blockers (alias: --blocked=false)
  --blocks              Show tasks that block other tasks
  --not-blocking        Show tasks that don't block others (alias: --blocks=false)
  --owner-approval      Show tasks requiring owner approval
  --approved            Show tasks with owner approval (alias: --owner-approval=true)
```

## Tasks
- [x] Draft 2–4 alternative filter flag designs (e.g., `--filter key=value`, `--status`, boolean pairs like `--blocked/--unblocked`).
- [x] Compare each alternative for usability, composability, backward compatibility, and help/UX clarity.
- [x] Identify recommended default and whether to support aliases for transition.
- [ ] Update or propose updates to `design-docs/list-command.md` with the chosen pattern (await owner approval).

## Acceptance Criteria
- [x] Alternatives with pros/cons and risks are documented.
- [x] A recommended pattern is proposed with rationale.
- [ ] Owner approval is requested before any implementation work proceeds.
