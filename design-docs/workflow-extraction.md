# Workflow Extraction Design

## Problem

The workflow is currently embedded implicitly in role definitions and template TODOs. This makes it:
- Hard to visualize the complete workflow at a glance
- Difficult to validate that workflows are complete and consistent
- Impossible to generate documentation automatically
- Challenging to compare different preset workflows

## Goals

1. Extract workflow information from existing roles and templates without duplicating data
2. Provide a `strand workflow` command that outputs workflows in multiple formats
3. Enable workflow validation (detect missing handoffs, orphaned roles, etc.)
4. Support multiple presets with different workflows

## Non-Goals

- Separate workflow definition files (workflow stays embedded in roles/templates)
- Runtime workflow enforcement (roles already define their own constraints)
- Workflow versioning or migration tools

## Proposed Solution

### Phase 1: Add Workflow Metadata to Role Files

Extend role file frontmatter with optional workflow metadata:

```markdown
---
role: designer
workflow:
  creates: [design-task, alternatives-task, review-task, owner-decision]
---

# Designer

## Role
Designer (human or senior AI) — explores alternatives and produces design artifacts.
...
```

**Frontmatter Fields:**
- `creates`: Array of task type names this role typically creates

This field is **optional** and **informational**. It documents which task types this role creates, and the workflow is inferred from the task templates (which have embedded role assignments).

**Key Insight:** Roles don't hand off directly to other roles. They create specific task types, and each task type has a role embedded in its template. For example:
- Designer creates a `review-task` → template has `role: master-reviewer` → Master Reviewer picks it up
- Designer creates an `owner-decision` → template has `role: owner` → Owner picks it up
- Architect creates a `task` → template has `role: developer` → Developer picks it up

This means the workflow is: **Role → Creates Task Type → (Task Type has Role) → Next Role**

### Phase 2: Extract Workflow from Template TODOs

Parse template files to extract the sequence of roles:

```markdown
## TODOs
1. [ ] (role: developer) Implement the behavior described in Context.
2. [ ] (role: developer) Add unit and integration tests covering the main flows.
3. [ ] (role: tester) Execute test-suite and report failures.
4. [ ] (role: master-reviewer) Coordinate required reviews.
5. [ ] (role: documentation) Update user-facing docs and examples.
```

**Extraction Logic:**
- Parse `(role: X)` annotations from TODO items
- Determine sequence: developer → developer → tester → master-reviewer → documentation
- Identify parallelizable steps (multiple roles at same TODO level)
- Track which templates use which role sequences

### Phase 3: Workflow Command

Implement `strand workflow` command with multiple output formats:

```bash
# Generate Mermaid diagram (like in README)
strand workflow --format mermaid

# Generate structured JSON for tooling
strand workflow --format json

# Generate DOT format for Graphviz
strand workflow --format dot

# Show workflow for specific template
strand workflow --template task --format mermaid

# Show all templates that use a specific role
strand workflow --role developer --show templates

# Validate workflow completeness
strand workflow --validate
```

### Phase 4: Workflow Analysis & Validation

The `--validate` flag checks for:

1. **Orphaned roles** - Roles defined but never referenced in templates
2. **Missing handoffs** - Role A hands off to Role B, but Role B doesn't expect work from Role A
3. **Incomplete templates** - Templates missing common workflow steps (e.g., no testing step)
4. **Circular dependencies** - Role A → B → A cycles without Owner intervention
5. **Missing role definitions** - Templates reference roles that don't exist

## Implementation Details

### Data Structures

```go
type RoleWorkflow struct {
    RoleName string
    Creates  []string  // Task types this role can create
}

type TemplateWorkflow struct {
    TemplateName string
    RoleSequence []RoleStep
}

type RoleStep struct {
    Role        string
    Description string
    Parallel    bool  // Can execute in parallel with next step
}

type WorkflowGraph struct {
    Roles     map[string]*RoleWorkflow
    Templates map[string]*TemplateWorkflow
    Edges     []WorkflowEdge
}

type WorkflowEdge struct {
    FromRole   string  // Role name that creates the task
    TaskType   string  // Task type created (e.g., "review-task")
    ToRole     string  // Role name from task template
    Label      string  // e.g., "creates review-task for"
    Source     string  // e.g., "designer.md frontmatter" or "task.md template"
}
```

### Mermaid Generation Algorithm

```go
func GenerateMermaid(graph *WorkflowGraph) string {
    // 1. For each role, look at what task types it creates (from frontmatter)
    // 2. For each task type, look at what role is assigned (from template frontmatter)
    // 3. Create edges: FromRole --[creates TaskType]--> ToRole
    // 4. Extract TODO sequences from templates for within-task flow
    // 5. Generate Mermaid graph with:
    //    - Nodes for each role
    //    - Edges labeled with task types (e.g., "creates review-task")
    //    - Colors based on role type (owner=green, reviewer=purple, etc.)
    //    - Subgraphs for within-task TODO sequences
}
```

### JSON Output Format

```json
{
  "roles": {
    "designer": {
      "creates": ["design-task", "alternatives-task", "review-task", "owner-decision"]
    },
    "architect": {
      "creates": ["task", "epic"]
    }
  },
  "templates": {
    "task": {
      "role": "developer",
      "todo_sequence": [
        {"role": "developer", "description": "Implement the behavior"},
        {"role": "developer", "description": "Add tests"},
        {"role": "tester", "description": "Execute test-suite"},
        {"role": "master-reviewer", "description": "Coordinate reviews"},
        {"role": "documentation", "description": "Update docs"}
      ]
    },
    "review-task": {
      "role": "master-reviewer",
      "todo_sequence": [
        {"role": "master-reviewer", "description": "Coordinate specialized reviewers"},
        {"role": "reviewer-security", "description": "Security review"},
        {"role": "reviewer-reliability", "description": "Reliability review"}
      ]
    }
  },
  "edges": [
    {
      "from_role": "designer",
      "task_type": "review-task",
      "to_role": "master-reviewer",
      "source": "designer.md creates + review-task.md role"
    },
    {
      "from_role": "architect",
      "task_type": "task",
      "to_role": "developer",
      "source": "architect.md creates + task.md role"
    }
  ]
}
```

## File Format Changes

### Role Files (Backward Compatible)

**Before:**
```markdown
# Designer

## Role
Designer (human or senior AI) — explores alternatives and produces design artifacts.
```

**After (optional workflow metadata):**
```markdown
---
role: designer
workflow:
  creates: [design-task, alternatives-task, review-task, owner-decision]
---

# Designer

## Role
Designer (human or senior AI) — explores alternatives and produces design artifacts.
```

**Migration:** Existing role files without frontmatter continue to work. Workflow extraction can scan task creation commands in role documentation as fallback, or emit warnings about missing metadata.

### Template Files (No Changes Required)

Templates already contain `(role: X)` annotations in TODOs. No changes needed, but we formalize the parsing:

```markdown
## TODOs
Check this off one at a time with `strand complete <task_id> --role <my_given_role> --todo <num> "report"` only if your Role matches the todo's role.
1. [ ] (role: developer) Implement the behavior described in Context.
2. [ ] (role: developer) Add unit and integration tests covering the main flows.
3. [ ] (role: tester) Execute test-suite and report failures.
4. [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
5. [ ] (role: documentation) Update user-facing docs and examples.
```

**Parsing Rules:**
- Extract `(role: X)` from each TODO line
- If multiple roles mentioned in description (e.g., `reviewer-reliability`), treat as parallel sub-tasks
- Sequence is determined by TODO order
- Empty lines or non-role TODOs are skipped

## Examples

### Generate Workflow Diagram for README

```bash
$ strand workflow --format mermaid --template task
graph TD
    Developer1[Developer: Implement] --> Developer2[Developer: Add Tests]
    Developer2 --> Tester[Tester: Execute Tests]
    Tester --> MasterReviewer[Master Reviewer: Coordinate Reviews]
    MasterReviewer --> Security[Security Reviewer]
    MasterReviewer --> Reliability[Reliability Reviewer]
    MasterReviewer --> Usability[Usability Reviewer]
    Security --> Documentation
    Reliability --> Documentation
    Usability --> Documentation
    Documentation[Documentation: Update Docs]
```

### Validate Workflow Consistency

```bash
$ strand workflow --validate
✓ All roles referenced in templates have definitions
✓ All task types referenced in role.creates exist as templates
⚠ Warning: Role 'designer' has no workflow metadata (missing 'creates' field)
⚠ Warning: Role 'tester' is used in templates but has no role definition file
⚠ Warning: Task type 'review-task' is never created by any role
✗ Error: Template 'epic' has no testing step (missing 'tester' role in TODOs)

3 warnings, 1 error
```

### Show Role Usage

```bash
$ strand workflow --role developer --show templates
Role: developer

Assigned to templates:
  - task.md (primary role)

Used in template TODOs:
  - task.md (TODOs: 1, 2)
  - epic.md (TODOs: 1)

Receives work via task types:
  - task (created by architect)
  - issue (created by triage, owner)

Creates task types:
  - (none in workflow metadata)
```

## Migration Plan

### Step 1: Implement Workflow Extraction (No Breaking Changes)
- Add workflow frontmatter parsing (optional, backward compatible)
- Implement template TODO parsing
- Build in-memory workflow graph

### Step 2: Implement `strand workflow` Command
- Add `--format mermaid` output
- Add `--format json` output
- Add `--format dot` output
- Add `--template` filter

### Step 3: Add Workflow Validation
- Implement `--validate` flag
- Add validation checks (orphaned roles, missing handoffs, etc.)
- Add `--show templates` filter

### Step 4: Add Workflow Metadata to Built-in Roles
- Update all built-in role files with frontmatter
- Validate that metadata matches actual template usage
- Update documentation

### Step 5: Auto-generate README Workflow Diagram
- Add script to generate Mermaid diagram from `strand workflow`
- Update README build process to include workflow diagram
- Document how custom presets can generate their own diagrams

## Open Questions

1. **Should workflow metadata be required or optional?**
   - ✓ Decision: Optional, with validation warnings if missing

2. **How to handle role name variations?**
   - ✓ Decision: No variations - each role has exactly one name. `reviewer`, `master-reviewer`, `reviewer-security` are distinct roles.

3. **Should we auto-generate workflow metadata from templates?**
   - Proposal: Yes, as a `strand workflow --update-roles` command that scans role documentation for task creation patterns

4. **How to handle conditional workflows (e.g., security review only for certain tasks)?**
   - Proposal: Document as separate templates (e.g., `task.md` vs `task-security.md`) or as optional TODOs in templates

5. **Should we show both inter-task workflow (role → task type → role) AND intra-task workflow (TODO sequences)?**
   - Proposal: Yes, with different diagram types:
     - `--format mermaid` shows inter-task workflow (like README diagram)
     - `--format mermaid --template task` shows intra-task TODO sequence for specific template

## Future Enhancements

- `strand workflow --compare preset1 preset2` - Compare workflows between presets
- `strand workflow --live` - Show current project's workflow with active tasks highlighted
- Interactive workflow editor (TUI/web dashboard)
- Workflow metrics (average time per role, bottleneck detection)
- Export workflows for external tools (Jira, Linear, etc.)
