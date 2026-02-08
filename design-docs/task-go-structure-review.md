# Review: pkg/task/task.go Structure and Methods

## Context
- design-docs/status-field-migration.md
- design-docs/commands-design.md
- pkg/task/task.go

## Reproduction Steps
```bash
go test ./...
go build ./...
go run ./cmd/strand repair
```

## Task Struct Inventory

### Metadata (YAML frontmatter)
- Type (string): task type label.
- Role (string): assigned role; used for selection and role documents.
- Priority (string): priority label.
- Parent (string): parent task ID; drives subtask ordering and hierarchy.
- Blockers ([]string): task IDs that block this task.
- Blocks ([]string): task IDs this task blocks.
- DateCreated (time.Time): creation timestamp.
- DateEdited (time.Time): last edit timestamp, updated by MarkDirty.
- OwnerApproval (bool): optional approval gate.
- Completed (bool): completion flag.
- Status (string): status field (current data model includes both Completed and Status).
- Every ([]string): recurrence anchors.
- Description (string): short summary or description.

### Task (full task model)
- ID (string): full task ID derived from filename.
- Dir (string): directory containing the task file.
- FilePath (string): absolute path to task file.
- Meta (Metadata): parsed frontmatter.
- TitleContent (string): H1 title text.
- BodyContent (string): markdown content excluding title and special sections.
- TodoItems ([]TaskItem): parsed entries from "## TODOs".
- SubsItems ([]TaskItem): parsed entries from "## Subtasks".
- ProgressContent (string): "## Progress" section body.
- OtherContent (string): remaining markdown not captured by the above sections.
- Dirty (bool): dirty flag used by write helpers.

## Methods on *Task
- SetTitle(newTitle string): update TitleContent and mark dirty when changed.
- SetBody(newBody string): replace body while stripping title and reserved sections (TODOs, Subtasks, Progress), then mark dirty when changed.
- MarkDirty(): set Dirty and update Meta.DateEdited once per dirty cycle.
- Write() error: serialize and write Content() to FilePath, clear Dirty.
- Title() string: return TitleContent.
- Content() string: render YAML frontmatter and markdown sections for file output.
- GetEffectiveRole() string: prefer first unchecked TODO role, then Meta.Role.

## Exported Functions Operating on Tasks (pkg/task/task.go)
- WriteAllTasks(tasks map[string]*Task) (int, error): write all tasks regardless of Dirty.
- WriteDirtyTasks(tasks map[string]*Task) (int, error): write only Dirty tasks.

## Manual Field Edits That Can Break Relationships
- Meta.Parent: changing this re-parents a task; affects subtask lists and ordering.
- Meta.Blockers / Meta.Blocks: inconsistent edits can break dependency symmetry.
- Meta.Role: affects role-based selection and workflow.
- Meta.Completed / Meta.Status: affects free list and completion logic.
- Meta.Every: affects recurrence generation and anchoring.
- TitleContent and Subtasks/TODOs sections: manual edits can desync Subtasks list from actual children.

## Methods That Modify State
- SetTitle, SetBody: mutate content and mark dirty.
- MarkDirty: updates Meta.DateEdited and Dirty flag.
- Write: writes to disk and clears Dirty.

## Dirty Tracking Mechanism
- Dirty is a boolean flag on Task.
- MarkDirty sets Meta.DateEdited to time.Now().UTC() when transitioning from clean to dirty.
- Write clears Dirty after successfully persisting Content().
- SetTitle and SetBody call MarkDirty when they detect changes.
