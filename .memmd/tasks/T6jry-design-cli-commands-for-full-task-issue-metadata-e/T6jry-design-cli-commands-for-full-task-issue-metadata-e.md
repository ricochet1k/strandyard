---
type: leaf
role: designer
priority: medium
parent: Ivahc-revamp-task-and-issue-commands-for-full-metadata-e
blockers: []
blocks:
    - Ivahc-revamp-task-and-issue-commands-for-full-metadata-e
date_created: 2026-01-28T05:10:02.740289Z
date_edited: 2026-01-28T05:16:53.663249Z
owner_approval: false
completed: true
---

# Design CLI commands for full task/issue metadata editing

## Context
This task is to design CLI commands that allow creating and editing all metadata fields for tasks and issues without requiring manual markdown file editing. This addresses issue Ivahc-revamp-task-and-issue-commands-for-full-metadata-e.

## Design Considerations
- Current `memmd add` command only supports limited flags (role, priority, parent, blockers)
- Need to support all metadata fields: title, body content, dates, owner approval, etc.
- Need both create and edit workflows
- Should maintain backward compatibility
- UX should be intuitive for CLI users

## Alternatives to Explore
- Option 1: Extend existing `add` and add new `edit` command with rich flag set
- Option 2: Introduce interactive mode with prompts
- Option 3: Subcommand approach: `memmd task add/edit`, `memmd issue add/edit`
- Option 4: Configuration file approach for complex updates

## Research: Current CLI patterns and limitations

### Current `add` command capabilities
- Templates: `leaf` (dev tasks) and `issue` (bug reports)
- Flags supported: `--role`, `--priority`, `--parent`, `--blocker`
- Accepts stdin input for body content
- Auto-generates task IDs and directories
- Supports title via argument or `--title` flag

### Missing metadata fields for full editing
- `--blocks` (tasks this task blocks)
- `--owner-approval` (boolean flag)
- `--completed` (boolean flag)
- Custom dates (created/edited)
- Body content editing (creation only via stdin/template)
- Type-specific metadata

### Current limitations
1. No `edit` command - metadata cannot be changed after creation
2. Limited flag set - cannot set all metadata fields
3. No interactive mode for complex workflows
4. No batch operations
5. Cannot edit body content without opening file

## Design Alternatives

### Option 1: Extend existing commands with rich flag set
**Approach**: Add comprehensive flags to `add` and create new `edit` command

```bash
# Creation
memmd add leaf "Feature X" --role developer --priority high --parent ABC123 \
  --blockers DEF456,GHI789 --blocks JKL012 --owner-approval

# Editing
memmd edit ABC123 --priority low --add-blocker MNO345 --remove-blocker DEF456 \
  --set-owner-approval --toggle-completed
```

**Pros**:
- Familiar CLI pattern (existing users know flag-based approach)
- Scriptable and automation-friendly
- Backward compatible with existing `add` command
- Clear, explicit operations

**Cons**:
- Flag proliferation (`--add-blocker`, `--remove-blocker`, `--set-*`, `--toggle-*`)
- Complex flag combinations hard to discover
- Body editing still awkward via flags
- Long command lines for complex updates

### Option 2: Interactive mode with prompts
**Approach**: Add `--interactive` flag to trigger guided prompts

```bash
memmd add leaf --interactive
memmd edit ABC123 --interactive
```

**Flow**:
1. Show current metadata (for edit)
2. Step through each field with prompts
3. Allow skipping with Enter key
4. Confirm changes before applying

**Pros**:
- Discoverable - users see all available fields
- Reduced command-line complexity
- Good for infrequent users
- Can include validation and help text

**Cons**:
- Not scriptable/automation-friendly
- Slower for experienced users
- Harder to use in CI/automation
- Terminal interaction complexity

### Option 3: Subcommand approach
**Approach**: Separate task/issue subcommands with richer interfaces

```bash
memmd task add leaf "Feature X" --role developer
memmd task edit ABC123 --set-priority high
memmd issue add "Bug description" --severity critical
memmd issue edit I789 --add-blocker T123
```

**Pros**:
- Clear separation of concerns
- Type-specific flags and validation
- Future extensibility per type
- Clean command structure

**Cons**:
- Breaking change from current `memmd add` usage
- More commands to learn
- Type proliferation complexity
- Migration effort

### Option 4: Configuration file approach
**Approach**: Support YAML/JSON files for complex metadata updates

```bash
# Create from config
memmd add leaf --from-config feature-x.yaml

# Edit with config diff
memmd edit ABC123 --apply-changes changes.yaml

# Config file format
changes.yaml:
  priority: high
  add_blockers: [DEF456, GHI789]
  remove_blockers: [OLD123]
  owner_approval: true
```

**Pros**:
- Handles complex operations elegantly
- Version-controllable configurations
- Good for bulk operations
- Clear intent documentation

**Cons**:
- Additional file management overhead
- Overkill for simple changes
- Less immediate than direct commands
- Learning curve for config format

## Recommendation: Hybrid Approach

**Primary Recommendation**: Option 1 + Limited Option 2

1. **Extend `add`** with comprehensive flags for all metadata fields
2. **Add `edit` command** with matching flag set
3. **Add `--interactive`** flag to both for discovery and complex workflows
4. **Keep simple flags** for common operations

This provides:
- Backward compatibility
- Scriptability for automation
- Discoverability for humans
- Clean migration path

## Implementation Phases

### Phase 1: Core flag extensions
- Add missing metadata flags to `add` command
- Create basic `edit` command with same flag set
- Update validation and error handling

### Phase 2: Interactive mode
- Add `--interactive` flag implementation
- Create prompt workflows for both add/edit
- Add help text and validation in prompts

### Phase 3: Type-specific enhancements
- Consider subcommands if type divergence grows
- Add type-specific validation and defaults
- Evaluate config file support for bulk operations

## Risk Assessment

**Low Risk**:
- Backward compatibility (preserve existing behavior)
- Flag extensions (well-established pattern)

**Medium Risk**:
- Interactive implementation complexity
- User experience design

**High Risk**:
- Breaking existing command structure (avoided in recommendation)

## Migration Considerations

- All existing `memmd add` commands continue to work
- New flags are optional - no forced changes
- `edit` command is purely additive
- Interactive mode is opt-in

## TODOs
- [x] (role: designer) Research current CLI patterns and limitations
- [x] (role: designer) Document design alternatives with pros/cons
- [x] (role: designer) Consider user experience for different workflows
- [x] (role: designer) Create mockups/UX flows for each alternative
- [ ] (role: master-reviewer) Review design alternatives
- [ ] (role: owner) Select preferred alternative

## Subtasks
Follow-on implementation tasks will be created based on selected design alternative.

## Acceptance Criteria
- Complete design alternatives document with pros/cons analysis
- User experience flows documented for each alternative
- Risk assessment and migration considerations
- Recommendation ready for owner decision
