## Design Alternatives Review: CLI command alignment

**Designer**: ai-assistant
**Date**: 2026-01-27
**Status**: Complete

## Summary

This review compares the implemented CLI (`cmd/*.go`) against the intended behaviour in [design-docs/commands-design.md](design-docs/commands-design.md) and proposes concrete alternatives for each divergence.

## Test Results

**Tests run**: 2026-01-27
```bash
# go test ./...
?   	github.com/ricochet1k/memmd	[no test files]
?   	github.com/ricochet1k/memmd/cmd	[no test files]

# go run . repair --path tasks --format text
ERROR: missing role file roles/developer.md for task T000002-setup-infra
ERROR: missing Role in tasks/D000001-review-design/task.md
ERROR: missing role file roles/developer.md for task T000001-project-alpha

# go build ./...
SUCCESS (builds without errors)
```

## Command Implementation Status

| Command | Specified in Design | Implementation Status |
|---------|--------------------|-----------------------|
| `repair` | ✓ | **COMPLETE** - Fully functional |
| `next` | ✓ | **COMPLETE** - Fully functional |
| `init` | ✓ | **STUB** - Only prints "init called" |
| `add` / `new` | ✓ | **STUB** - Only prints "add called" |
| `assign` | ✓ | **STUB** - Only prints "assign called" |
| `block` | ✓ (with subcommands) | **STUB** - Only prints "block called", no subcommands |
| `templates list` | ✓ | **STUB** - Only prints "templates called", no list subcommand |

## Identified Discrepancies

### 1) **CRITICAL: Task Metadata Format Inconsistency**

**Observation**: There are THREE different metadata formats in use across the codebase:
- **[AGENTS.md](AGENTS.md:51-70)** specifies simple field format: `Role: developer`, `Blockers:`, `Blocks:`, `Parent:`
- **[templates/leaf.md](templates/leaf.md:9-13)** uses markdown heading format: `## Role`, `## Track`
- **[tasks/T000001-project-alpha/task.md](tasks/T000001-project-alpha/task.md:4-14)** uses simple field format (matches AGENTS.md)
- **[tasks/D000001-review-design/task.md](tasks/D000001-review-design/task.md:3-4)** uses heading format (matches template)
- **Parser in [cmd/validate.go:214-225](cmd/validate.go:214-225)** expects `Role:` with colon (simple field format only)

This causes the validation error: `ERROR: missing Role in tasks/D000001-review-design/task.md`

**Alternatives**:

**A) Standardize on Simple Field Format** (matches AGENTS.md, parser expects this)
- Change template from `## Role` to `Role:`
- Update D000001 task to use `Role: designer` instead of `## Role` heading
- Keep parser as-is (simplest)
- **Pros**: Matches existing AGENTS.md spec; parser already works for this format; simpler parsing (no markdown AST needed)
- **Cons**: Less visually structured in rendered markdown; fields blend with description text
- **Impact**: Update 1 template file, 1 task file

**B) Standardize on Markdown Heading Format** (template uses this)
- Change AGENTS.md examples to use `## Role`, `## Blockers`, `## Blocks`, `## Parent`
- Update parser to recognize headings (requires markdown AST parsing or regex for `## Field`)
- Update T000001 task to use heading format
- **Pros**: More visually structured; cleaner separation in rendered docs; aligns with modern markdown practices
- **Cons**: More complex parsing; larger change to codebase and docs
- **Impact**: Update AGENTS.md, parser code, T000001 task file

**C) Support Both Formats** (parser accepts either)
- Enhance parser to check for both `Role:` and `## Role`
- Keep templates and docs as-is (mixed)
- **Pros**: Backward compatible; flexible for users
- **Cons**: Ambiguous standard; potential for confusion; more complex parser; harder to maintain consistency
- **Impact**: Update parser only

**Owner Decision Required**: Choose between A (minimal change, simpler), B (more structured, complex), or C (flexible, ambiguous).

### 2) Task ID Format Mismatch

**Observation**:
- **Design doc ([design-docs/commands-design.md:20-27](design-docs/commands-design.md:20-27))** specifies: `<prefix><token>-<mini>` where token is 6-char random alphanumeric (base36), example: `T4k3a1-init`
- **Actual tasks** use sequential format: `T000001-project-alpha`, `T000002-setup-infra`, `D000001-review-design`
- **Parser ([cmd/validate.go:108](cmd/validate.go:108))** enforces: `^[A-Z][0-9A-Za-z]{6}-[a-zA-Z0-9-]{1,}$` (7 alphanumeric chars total after prefix)

Current task IDs like `T000001` have 6 digits, which matches the regex (total 7 chars including prefix).

**Alternatives**:

**A) Keep Sequential Numbering** (current implementation)
- Document that IDs use sequential format `<prefix>NNNNNN-<slug>` (6-digit counter)
- Update design doc to match actual implementation
- **Pros**: Predictable IDs; easy to implement; sequential ordering; no collision risks
- **Cons**: Reveals total task count; sequential IDs might encode information over time
- **Impact**: Update design-docs/commands-design.md only

**B) Switch to Random Token Format** (design doc specification)
- Implement crypto/rand base36 token generator as specified in design doc
- Migrate existing tasks (rename directories, update references)
- **Pros**: Matches design doc; non-sequential (no info leak); cryptographically strong
- **Cons**: Requires migration tool; ID collisions theoretically possible (though extremely unlikely with 6-char base36 = 2B+ combinations); harder to remember
- **Impact**: Implement ID generator, migration tool, update 3 existing tasks

**Owner Decision Required**: Choose between A (keep sequential, update docs) or B (switch to random tokens, matches original design spec).

### 3) `next` role-filtering behaviour
- Observation: Implementation requires a role (flag or `MEMMD_ROLE`) and selects tasks by exact `Role:` match first, then unassigned tasks, then any task. Design doc suggests `memmd next` should emit role doc then next task and not require filters.

Alternatives
- Minimal change: keep current behaviour (require role) but document it clearly in the design doc and role files.
  - Pros: simple, minimal code change; explicit control for agent role selection.
  - Cons: deviates from design doc; agents must set env/flag.

- Larger change: make `next` default to printing the role doc for a role inferred from environment or a configured default, and if no role provided, print a generic role doc + first free task.
  - Pros: matches design doc intent (no required flags); better UX for interactive humans and agents.
  - Cons: requires changes to `next` selection logic and documentation; potential ambiguity when multiple roles present.

**Owner Decision Required**: Choose between minimal change (keep current behavior, document it) or larger change (optional role with defaults).

### 4) `repair` behaviour and failure semantics
**Observation**: `repair` regenerates master lists and exits non-zero on errors - this matches design doc. However, validation currently fails due to the metadata format issue (see Discrepancy #1) and references to missing role files.

**Current behavior**:
- Parses tasks correctly when using `Role:` field format
- Repairs role files exist
- Regenerates master lists deterministically (sorted)
- Exits with error code if any validation fails

**Alternatives**:

**A) Keep Current Strict Validation**
- Maintain current behavior: fail fast on any error
- Fix metadata format issues (see Discrepancy #1)
- **Pros**: Catches errors early; enforces data quality; predictable behavior
- **Cons**: Can block other operations even for minor issues
- **Impact**: No code changes needed, just fix task metadata

**B) Add Validation Levels** (`--strict` / `--lenient`)
- `--strict`: Current behavior (default)
- `--lenient`: Warn but don't fail; still regenerate master lists
- **Pros**: Flexibility for gradual migration or legacy repos
- **Cons**: More complex; users might ignore warnings
- **Impact**: Add flag parsing and conditional error handling

**Owner Decision Required**: Choose between A (keep strict validation) or B (add validation levels with --strict/--lenient flags).

### 5) Master lists path and determinism

**Observation**: `repair` writes `tasks/root-tasks.md` and `tasks/free-tasks.md` deterministically with sorted entries - behavior matches design doc exactly. ✓ **No change recommended.**

### 6) Unimplemented Commands

**Observation**: Five commands are specified in design doc but only have boilerplate stubs:

| Command | Design Spec | Current State |
|---------|-------------|---------------|
| `init` | Bootstrap repo structure, roles, examples | Prints "init called" |
| `add`/`new` | Create task from template with ID generation | Prints "add called" |
| `assign` | Change task Role field | Prints "assign called" |
| `block` | Manage blocker relationships (add/remove/list subcommands) | Prints "block called", no subcommands |
| `templates` | List available templates | Prints "templates called", no list subcommand |

**Alternatives**:

**A) Implement All Commands** (comprehensive)
- Implement full functionality for all 5 commands per design spec
- **Pros**: Complete feature set; matches design doc; full CLI functionality
- **Cons**: Significant implementation effort (estimate 20-40 hours)
- **Impact**: Implement 5 commands with tests

**B) Prioritize Core Workflow** (incremental)
- Phase 1: Implement `add` (create tasks) and `assign` (change ownership) - enables basic task management
- Phase 2: Implement `block` subcommands - enables dependency management
- Phase 3: Implement `templates list` and `init --force` - nice-to-have features
- **Pros**: Delivers value incrementally; focuses on most-used commands first
- **Cons**: Partial functionality for a period; need to plan phases
- **Impact**: 3 phased epics

**Owner Decision Required**: Choose between A (implement all commands at once) or B (incremental phased approach).

### 7) Test Coverage

**Observation**:
- **Current state**: No test files exist (`[no test files]`)
- **Design doc requirement**: Not explicitly specified but good practice expected
- **Risk**: No automated validation of parsing logic, ID validation, master list generation

**Alternatives**:

**A) Add Comprehensive Test Suite**
- Unit tests for parsers (`parseRole`, `parseBlockers`)
- Integration tests for `repair` and `next` commands
- Test fixtures with sample tasks (valid and malformed)
- **Pros**: Catches regressions; documents expected behavior; enables confident refactoring
- **Cons**: Initial effort to write tests; ongoing maintenance
- **Impact**: Add `cmd/repair_test.go`, `cmd/next_test.go`, `testdata/` fixtures

**B) Minimal Test Coverage** (critical paths only)
- Tests for ID regex validation
- Tests for master list generation (deterministic ordering)
- Tests for role file existence checks
- **Pros**: Faster to implement; covers highest-risk areas
- **Cons**: May miss edge cases; less comprehensive
- **Impact**: Single test file `cmd/cmd_test.go`

**Owner Decision Required**: Choose between A (comprehensive test suite) or B (minimal test coverage for critical paths only).

### 8) Template Organization

**Observation**:
- **Design doc** mentions `templates/task-templates/` and `templates/doc-templates/`
- **AGENTS.md** mentions `templates/task-templates/` and `templates/doc-templates/`
- **Actual structure**: `templates/leaf.md` (single file), `doc-examples/` directory exists with different purpose
- **[templates/leaf.md](templates/leaf.md:1-6)** has Go template syntax `{{ .Title }}` but no template execution code exists

**Alternatives**:

**A) Reorganize to Match Design Doc**
- Create `templates/task-templates/` and `templates/doc-templates/` directories
- Move `leaf.md` to `templates/task-templates/leaf.md`
- Add additional templates (epic, design-doc, etc.)
- **Pros**: Matches design doc; organized by template type; room for growth
- **Cons**: Need to update any hardcoded paths; migration needed
- **Impact**: Create directories, move files, update references

**B) Keep Flat Structure, Update Docs**
- Keep `templates/` flat with `leaf.md`, `epic.md`, etc.
- Update design doc and AGENTS.md to match
- **Pros**: Simpler; less migration; fine for small number of templates
- **Cons**: Doesn't scale well; mixes different template types
- **Impact**: Update 2 doc files

**Owner Decision Required**: Choose between A (reorganize to match design doc structure) or B (keep flat structure, update docs to match).

## Possible Prioritization (for Owner consideration)

One possible priority ordering based on impact and effort analysis:

| Suggested Priority | Discrepancy | Impact | Effort | Notes |
|----------|-------------|--------|--------|-----------|
| High | #1 - Metadata format | High | Low | Blocks current validation; quick fix |
| High | #7 - Test coverage | High | Medium | Prevents regressions during implementation |
| Medium | #6 - Unimplemented commands (Phase 1: add, assign) | High | High | Core functionality for task management |
| Medium | #8 - Template organization | Medium | Low | Prerequisite for `add` command |
| Lower | #2 - ID format mismatch | Low | Low | Documentation issue only |
| Lower | #3 - `next` role filtering | Medium | Low | UX improvement |
| Lower | #6 - Unimplemented commands (Phase 2: block) | Medium | Medium | Enables dependency management |
| Lower | #4 - Repair strictness | Low | Low | Current behavior is acceptable |
| Lower | #6 - Unimplemented commands (Phase 3: init, templates) | Low | Medium | Nice-to-have features |

**Owner**: Please review and establish actual priority ordering based on project goals.

## Example Epics & Tasks (if Owner approves alternatives)

### Epic E1: Fix Metadata Format and Add Tests (P0)
**Owner**: developer
**Estimated effort**: 4-8 hours

- **Task E1-T1**: Standardize task metadata format to simple field format
  - Update [templates/leaf.md](templates/leaf.md) to use `Role:` instead of `## Role`
  - Update [tasks/D000001-review-design/task.md](tasks/D000001-review-design/task.md) to use simple field format
  - Verify `repair` passes for all tasks
  - **Files**: templates/leaf.md, tasks/D000001-review-design/task.md
  - **Acceptance**: `go run . repair` succeeds with no errors

- **Task E1-T2**: Add comprehensive test suite
  - Create `cmd/repair_test.go` with tests for `parseRole`, `parseBlockers`, ID validation
  - Create `cmd/next_test.go` with tests for role selection logic
  - Create `testdata/` with valid and malformed task fixtures
  - **Files**: cmd/repair_test.go, cmd/next_test.go, testdata/*
  - **Acceptance**: `go test ./...` passes with >80% coverage of repair and next commands

### Epic E2: Template Organization (P1)
**Owner**: developer
**Estimated effort**: 2-4 hours

- **Task E2-T1**: Reorganize templates directory
  - Create `templates/task-templates/` and `templates/doc-templates/`
  - Move `leaf.md` to `templates/task-templates/`
  - Add `templates/task-templates/epic.md` template
  - Add `templates/doc-templates/design-alternatives.md` template
  - **Files**: templates/* (restructure)
  - **Acceptance**: Template directories exist and contain appropriate templates

- **Task E2-T2**: Update documentation references
  - Update AGENTS.md to reflect new template paths
  - Update design-docs/commands-design.md template references
  - **Files**: AGENTS.md, design-docs/commands-design.md
  - **Acceptance**: All doc references to templates point to correct paths

### Epic E3: Implement Core Task Management Commands (P1)
**Owner**: developer
**Estimated effort**: 16-24 hours

- **Task E3-T1**: Implement ID generation for `add` command
  - Create `pkg/idgen/` package with sequential ID generator
  - Support format: `<prefix>NNNNNN-<slug>` where NNNNNN is zero-padded counter
  - Load last ID from tasks directory, increment
  - Add slugify function for title → mini-title conversion
  - **Files**: pkg/idgen/generator.go, pkg/idgen/generator_test.go
  - **Acceptance**: ID generator creates valid sequential IDs; tests pass

- **Task E3-T2**: Implement template expansion engine
  - Create `pkg/templates/` package to load and execute Go templates
  - Support template variables: Title, Role, Track, SuggestedSubtaskDir
  - Load templates from `templates/task-templates/`
  - **Files**: pkg/templates/expand.go, pkg/templates/expand_test.go
  - **Acceptance**: Template expansion works with test fixtures

- **Task E3-T3**: Implement `add` command
  - Parse flags: `--title`, `--role`, `--parent`, `--template`
  - Generate ID using idgen package
  - Create task directory: `tasks/<parent-path>/<task-id>/`
  - Expand template and write `<task-id>.md`
  - Run `repair` to update master lists
  - **Files**: cmd/add.go
  - **Acceptance**: `memmd add --title "Test task" --role developer` creates valid task

- **Task E3-T4**: Implement `assign` command
  - Parse flags: `<task-id>` (positional), `--role <role>`
  - Read task file, find `Role:` line, replace value
  - Write atomically (temp file + rename)
  - Repair role file exists before updating
  - **Files**: cmd/assign.go
  - **Acceptance**: `memmd assign T000001-test --role owner` updates role correctly

### Epic E4: Implement Blocker Management (P2)
**Owner**: developer
**Estimated effort**: 12-16 hours

- **Task E4-T1**: Implement `block add` subcommand
  - Parse flags: `--task <task-id>` `--blocks <blocker-id>`
  - Update both tasks: add to `Blockers:` list of task, add to `Blocks:` list of blocker
  - Maintain sorted order
  - Run `repair` to update free-tasks list
  - **Files**: cmd/block.go
  - **Acceptance**: `memmd block add --task T000002 --blocks T000001` updates both tasks

- **Task E4-T2**: Implement `block remove` subcommand
  - Parse flags: `--task <task-id>` `--blocks <blocker-id>`
  - Remove from both tasks' lists
  - **Files**: cmd/block.go
  - **Acceptance**: `memmd block remove --task T000002 --blocks T000001` removes blocker

- **Task E4-T3**: Implement `block list` subcommand
  - Parse flags: `--task <task-id>`
  - Display blockers and blocks for the task
  - **Files**: cmd/block.go
  - **Acceptance**: `memmd block list --task T000002` shows blockers and blocks

### Epic E5: Polish & Nice-to-Have Features (P3)
**Owner**: developer
**Estimated effort**: 8-12 hours

- **Task E5-T1**: Implement `init` command
  - Create directory structure: `tasks/`, `roles/`, `templates/`, `doc-examples/`
  - Copy example roles if `--examples` flag provided
  - Create sample tasks if `--examples` flag provided
  - **Files**: cmd/init.go, embedded example files
  - **Acceptance**: `memmd init --examples` bootstraps a working repository

- **Task E5-T2**: Implement `templates list` subcommand
  - List all templates in `templates/task-templates/` and `templates/doc-templates/`
  - Show template name and first line of description
  - **Files**: cmd/templates.go
  - **Acceptance**: `memmd templates list` shows available templates

- **Task E5-T3**: Improve `next` role selection UX
  - Make `--role` optional
  - Check `MEMMD_DEFAULT_ROLE` env var if `MEMMD_ROLE` not set
  - Provide helpful error message if no role configured
  - **Files**: cmd/next.go
  - **Acceptance**: `next` works without explicit role when default configured

- **Task E5-T4**: Update design docs to match implementation
  - Update ID format specification to sequential format in design-docs/commands-design.md
  - Document all command flags and examples
  - **Files**: design-docs/commands-design.md
  - **Acceptance**: Design doc accurately reflects implementation

## Reviewer Checklist

This design review should be reviewed by:

- **Master Reviewer** (roles/reviewer.md): Overall review for completeness, consistency, and architectural soundness
- **Reliability Reviewer** (roles/reviewer-reliability.md): Review error handling, validation logic, data consistency guarantees
- **Usability Reviewer** (roles/reviewer-usability.md): Review CLI UX, command naming, help text, error messages
- **Security Reviewer** (roles/reviewer-security.md): Review file operations (atomic writes, permissions), path traversal risks

## Files Produced

1. **doc-examples/design-alternatives-review.md** (this file) - Complete alternatives analysis with options for Owner decision

## Next Steps

1. **Owner**: Review alternatives and make decisions on each discrepancy (choose A, B, C, etc. for each)
2. **Owner** (or delegate to Designer): Document selected alternatives in a design recommendation document
3. **Architect**: Create detailed task files for approved work in `tasks/` directory
4. **Developer**: Begin implementation according to Owner's prioritization
