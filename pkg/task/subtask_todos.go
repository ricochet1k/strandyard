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

	// fmt.Printf("ParentSubtaskTodoItems diff: %v", cmp.Diff(parent.SubsItems, newSubs))

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
	parent, ok := tasks[parentID]
	if !ok {
		return nil
	}

	subtasks := []*Task{}
	remainingByID := map[string]*Task{}
	for _, t := range tasks {
		if t.Meta.Parent == parentID {
			subtasks = append(subtasks, t)
			remainingByID[t.ID] = t
		}
	}

	ordered := make([]*Task, 0, len(subtasks))
	for _, item := range parent.SubsItems {
		if item.SubtaskID == "" {
			continue
		}
		fullID := resolveChildSubtaskID(parentID, item.SubtaskID, remainingByID)
		if fullID == "" {
			continue
		}
		ordered = append(ordered, remainingByID[fullID])
		delete(remainingByID, fullID)
	}

	remaining := make([]*Task, 0, len(remainingByID))
	for _, sub := range remainingByID {
		remaining = append(remaining, sub)
	}

	// New subtasks are appended in creation order.
	sort.Slice(remaining, func(i, j int) bool {
		if !remaining[i].Meta.DateCreated.Equal(remaining[j].Meta.DateCreated) {
			return remaining[i].Meta.DateCreated.Before(remaining[j].Meta.DateCreated)
		}
		return remaining[i].ID < remaining[j].ID
	})
	ordered = append(ordered, remaining...)

	items := make([]TaskItem, 0, len(ordered))
	for _, sub := range ordered {
		title := sub.Title()
		if title == "" {
			title = sub.ID
		}
		items = append(items, TaskItem{
			Checked:   sub.Meta.Completed,
			SubtaskID: ShortID(sub.ID),
			Text:      title,
		})
	}
	return items
}

func resolveChildSubtaskID(parentID, input string, subtasks map[string]*Task) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	if t, ok := subtasks[input]; ok && t.Meta.Parent == parentID {
		return input
	}

	short := ShortID(input)
	match := ""
	for id, t := range subtasks {
		if t.Meta.Parent != parentID {
			continue
		}
		if ShortID(id) != short {
			continue
		}
		if match != "" {
			return ""
		}
		match = id
	}

	return match
}
