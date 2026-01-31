package task

import (
	"slices"
	"sort"
	"time"
)

// UpdateBlockersFromChildren ensures parent tasks are blocked by incomplete children.
// Returns the number of tasks marked as dirty.
func UpdateBlockersFromChildren(tasks map[string]*Task) (int, error) {
	children := map[string][]*Task{}
	for _, t := range tasks {
		if t.Meta.Parent == "" {
			continue
		}
		children[t.Meta.Parent] = append(children[t.Meta.Parent], t)
	}

	updated := 0
	now := time.Now()
	for parentID, kids := range children {
		parent, ok := tasks[parentID]
		if !ok {
			continue
		}
		if parent.Meta.Completed {
			continue
		}

		incomplete := []string{}
		childSet := map[string]struct{}{}
		for _, kid := range kids {
			childSet[kid.ID] = struct{}{}
			if !kid.Meta.Completed {
				incomplete = append(incomplete, kid.ID)
			}
		}

		sort.Strings(incomplete)
		seen := map[string]struct{}{}
		desired := make([]string, 0, len(parent.Meta.Blockers)+len(incomplete))
		for _, blocker := range parent.Meta.Blockers {
			if blocker == "" {
				continue
			}
			if _, isChild := childSet[blocker]; isChild {
				continue
			}
			if _, ok := seen[blocker]; ok {
				continue
			}
			seen[blocker] = struct{}{}
			desired = append(desired, blocker)
		}
		for _, blocker := range incomplete {
			if _, ok := seen[blocker]; ok {
				continue
			}
			seen[blocker] = struct{}{}
			desired = append(desired, blocker)
		}
		sort.Strings(desired)

		if slices.Equal(parent.Meta.Blockers, desired) {
			continue
		}

		parent.Meta.Blockers = desired
		parent.MarkDirty()
		updated++
	}

	return updated, nil
}
