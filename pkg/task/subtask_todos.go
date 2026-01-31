package task

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// UpdateParentTodoEntries syncs the parent's ## Subtasks section with its subtasks.
// Returns true if the parent task content was updated.
func UpdateParentTodoEntries(tasks map[string]*Task, parentID string) (bool, error) {
	parent, ok := tasks[parentID]
	if !ok {
		return false, fmt.Errorf("parent task not found: %s", parentID)
	}

	newSubs := buildSubtaskTodoItems(tasks, parentID)
	if slices.Equal(parent.SubsItems, newSubs) {
		return false, nil
	}

	parent.SubsItems = newSubs
	parent.MarkDirty()
	return true, nil
}

// UpdateAllParentTodoEntries syncs all parent tasks' ## Tasks sections with subtasks.
// Returns the number of parents updated.
func UpdateAllParentTodoEntries(tasks map[string]*Task) (int, error) {
	parents := map[string]struct{}{}
	for _, t := range tasks {
		if strings.TrimSpace(t.Meta.Parent) == "" {
			continue
		}
		parents[t.Meta.Parent] = struct{}{}
	}

	updated := 0
	for parentID := range parents {
		if _, ok := tasks[parentID]; !ok {
			continue
		}
		changed, err := UpdateParentTodoEntries(tasks, parentID)
		if err != nil {
			return updated, err
		}
		if changed {
			updated++
		}
	}
	return updated, nil
}

func buildSubtaskTodoItems(tasks map[string]*Task, parentID string) []TaskItem {
	subtasks := []*Task{}
	for _, t := range tasks {
		if t.Meta.Parent == parentID {
			subtasks = append(subtasks, t)
		}
	}
	sort.Slice(subtasks, func(i, j int) bool {
		return subtasks[i].ID < subtasks[j].ID
	})

	items := make([]TaskItem, 0, len(subtasks))
	for _, sub := range subtasks {
		title := sub.Title()
		if title == "" {
			title = sub.ID
		}
		items = append(items, TaskItem{
			Checked:   sub.Meta.Completed,
			SubtaskID: sub.ID,
			Text:      title,
		})
	}
	return items
}
