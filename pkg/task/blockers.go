package task

import (
	"fmt"
	"slices"
	"sort"

	"github.com/google/go-cmp/cmp"
)

// UpdateBlockersFromChildren ensures parent tasks are blocked by incomplete children.
// Returns the number of tasks marked as dirty.
func UpdateBlockersFromChildren(tasks map[string]*Task) (int, error) {
	taskBlockers := map[string]map[string]*Task{}
	for _, t := range tasks {
		if t.Meta.Completed {
			continue
		}

		if t.Meta.Parent != "" {
			blockers := taskBlockers[t.Meta.Parent]
			if blockers == nil {
				blockers = map[string]*Task{}
				taskBlockers[t.Meta.Parent] = blockers
			}
			blockers[t.ID] = t
		}
		for _, blocks := range t.Meta.Blocks {
			blockers := taskBlockers[blocks]
			if blockers == nil {
				blockers = map[string]*Task{}
				taskBlockers[blocks] = blockers
			}
			blockers[t.ID] = t
		}
	}

	updated := 0

	// Update all tasks that should have blockers
	for parentID, foundBlockers := range taskBlockers {
		parent, ok := tasks[parentID]
		if !ok {
			continue
		}
		if parent.Meta.Completed {
			continue
		}

		desired := make([]string, 0, len(foundBlockers))
		for blockerId := range foundBlockers {
			desired = append(desired, blockerId)
		}
		sort.Strings(desired)

		if slices.Equal(parent.Meta.Blockers, desired) {
			continue
		}

		fmt.Printf("UpdateBlockersFromChildren %v diff: %v", parent.FilePath, cmp.Diff(parent.Meta.Blockers, desired))
		fmt.Printf("Blockers %#v\n", parent.Meta.Blockers)
		fmt.Printf("desired %#v\n", desired)

		parent.Meta.Blockers = desired
		parent.MarkDirty()
		updated++
	}

	// Clear blockers from tasks that shouldn't have any
	for _, task := range tasks {
		if task.Meta.Completed {
			continue
		}
		// If task is not in taskBlockers, it should have no blockers
		if _, hasBlockers := taskBlockers[task.ID]; !hasBlockers {
			if len(task.Meta.Blockers) > 0 {
				fmt.Printf("UpdateBlockersFromChildren clearing blockers from %v (had %v)\n", task.FilePath, task.Meta.Blockers)
				task.Meta.Blockers = []string{}
				task.MarkDirty()
				updated++
			}
		}
	}

	return updated, nil
}
