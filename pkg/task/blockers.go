package task

import (
	"slices"
	"sort"
)

// ReconcileBlockerRelationships repairs blockers/blocks relationships in one pass.
//
// Expected behavior:
// - Parent tasks are blocked by every incomplete child task.
// - Incomplete tasks listed in Blockers or Blocks are treated as blocker edges.
// - Completed tasks do not block other tasks.
// - Blockers and Blocks are always rewritten as bidirectional, sorted, unique sets.
//
// Returns the number of tasks marked dirty.
func ReconcileBlockerRelationships(tasks map[string]*Task) (int, error) {
	desiredBlockers := make(map[string]map[string]struct{})

	addEdge := func(blockedID, blockerID string) {
		blocked, blockedOK := tasks[blockedID]
		blocker, blockerOK := tasks[blockerID]
		if !blockedOK || !blockerOK {
			return
		}
		if blocked.Meta.Completed || blocker.Meta.Completed {
			return
		}
		if desiredBlockers[blockedID] == nil {
			desiredBlockers[blockedID] = make(map[string]struct{})
		}
		desiredBlockers[blockedID][blockerID] = struct{}{}
	}

	for _, current := range tasks {
		if current.Meta.Completed {
			continue
		}

		if current.Meta.Parent != "" {
			addEdge(current.Meta.Parent, current.ID)
		}

		for _, blockedID := range current.Meta.Blocks {
			addEdge(blockedID, current.ID)
		}

		for _, blockerID := range current.Meta.Blockers {
			addEdge(current.ID, blockerID)
		}
	}

	updated := 0
	desiredBlocks := make(map[string]map[string]struct{})

	for taskID, current := range tasks {
		desired := sortedKeys(desiredBlockers[taskID])
		if !slices.Equal(current.Meta.Blockers, desired) {
			current.Meta.Blockers = desired
			current.MarkDirty()
			updated++
		}

		for _, blockerID := range desired {
			if desiredBlocks[blockerID] == nil {
				desiredBlocks[blockerID] = make(map[string]struct{})
			}
			desiredBlocks[blockerID][taskID] = struct{}{}
		}
	}

	for taskID, current := range tasks {
		desired := sortedKeys(desiredBlocks[taskID])
		if !slices.Equal(current.Meta.Blocks, desired) {
			current.Meta.Blocks = desired
			current.MarkDirty()
			updated++
		}
	}

	return updated, nil
}

func sortedKeys(items map[string]struct{}) []string {
	if len(items) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(items))
	for key := range items {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
