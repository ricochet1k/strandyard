// Package task provides parsing, validation, and graph-safe mutation APIs for
// StrandYard task files.
//
// Task files are markdown documents with YAML frontmatter. Task IDs come from
// filenames, and relationship fields (parent, blockers, blocks) are treated as
// graph edges with deterministic ordering.
//
// # Mental model
//
// Use TaskDB as the single mutation boundary for task relationships and status.
// Load tasks, mutate through TaskDB methods, reconcile relationships when
// needed, then persist dirty tasks.
//
// Pit of success
//
//  1. Create a TaskDB with the tasks directory.
//  2. Resolve user input with ResolveID/GetResolved.
//  3. Mutate through methods such as SetParent/AddBlocker/SetStatus.
//  4. Reconcile with ReconcileBlockerRelationships when relationship state
//     should be canonicalized across the graph.
//  5. Persist with SaveDirty (or SaveAll when explicitly required).
//
// Not allowed (or strongly discouraged)
//
//   - Do not mutate relationship fields directly on returned Task values
//     (Meta.Parent, Meta.Blockers, Meta.Blocks). Doing so bypasses invariants.
//   - Do not create blank tasks ad hoc for normal workflows. Task creation should
//     come from CLI/template flows that initialize required metadata.
//   - Do not assume every mutator auto-reconciles relationships. Some operations
//     intentionally keep reconciliation explicit so callers can batch edits.
//
// Common pitfalls
//
//   - SetParent and SetCompleted do not fully rebuild graph relationships by
//     themselves; call ReconcileBlockerRelationships (or
//     UpdateBlockersAfterCompletion for terminal statuses) before persisting when
//     global consistency is required.
//   - Editing status through SetStatus keeps completed/status fields in sync, but
//     legacy direct writes to Metadata can still produce surprising behavior.
//
// Design references
//
//   - pkg/task/TASKDB_DESIGN.md
//   - design-docs/blockers-relationship-management-review.md
//   - design-docs/status-field-migration.md
package task
