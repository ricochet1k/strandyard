package task

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var subtaskTodoPattern = regexp.MustCompile(`^- \[[ xX]\] \(subtask: ([^)]+)\)`)

// UpdateParentTodoEntries syncs the parent's ## Tasks section with its subtasks.
// Returns true if the parent task content was updated.
func UpdateParentTodoEntries(tasks map[string]*Task, parentID string) (bool, error) {
	parent, ok := tasks[parentID]
	if !ok {
		return false, fmt.Errorf("parent task not found: %s", parentID)
	}

	entries := buildSubtaskTodoEntries(tasks, parentID)
	updatedContent, changed, err := updateTasksSection(parent.Content, entries)
	if err != nil {
		return false, err
	}
	if !changed {
		return false, nil
	}

	parent.Content = updatedContent
	parent.Meta.DateEdited = time.Now().UTC()
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

func buildSubtaskTodoEntries(tasks map[string]*Task, parentID string) []string {
	subtasks := []*Task{}
	for _, t := range tasks {
		if t.Meta.Parent == parentID {
			subtasks = append(subtasks, t)
		}
	}
	sort.Slice(subtasks, func(i, j int) bool {
		return subtasks[i].ID < subtasks[j].ID
	})

	entries := make([]string, 0, len(subtasks))
	for _, sub := range subtasks {
		title := strings.TrimSpace(sub.Title())
		if title == "" {
			title = sub.ID
		}
		status := " "
		if sub.Meta.Completed {
			status = "x"
		}
		entries = append(entries, fmt.Sprintf("- [%s] (subtask: %s) %s", status, sub.ID, title))
	}
	return entries
}

func updateTasksSection(content string, entries []string) (string, bool, error) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return "", false, errInvalidFrontmatter("task content")
	}

	body := parts[2]
	updatedBody, changed := updateTasksSectionBody(body, entries)
	if !changed {
		return content, false, nil
	}

	updated := parts[0] + "---" + parts[1] + "---" + updatedBody
	return updated, true, nil
}

func updateTasksSectionBody(body string, entries []string) (string, bool) {
	lines := strings.Split(body, "\n")
	start := findHeaderLine(lines, "## Tasks")
	if start == -1 {
		return insertTasksSection(lines, entries)
	}

	end := findNextHeader(lines, start+1)
	section := lines[start+1 : end]
	filtered := make([]string, 0, len(section))
	for _, line := range section {
		if isSubtaskTodoLine(line) {
			continue
		}
		filtered = append(filtered, line)
	}
	filtered = trimBlankLines(filtered)

	sectionLines := buildTasksSectionLines(filtered, entries)
	if end < len(lines) && len(sectionLines) > 0 {
		if strings.TrimSpace(sectionLines[len(sectionLines)-1]) != "" && strings.TrimSpace(lines[end]) != "" {
			sectionLines = append(sectionLines, "")
		}
	}
	updated := make([]string, 0, len(lines)-len(section)+len(sectionLines))
	updated = append(updated, lines[:start]...)
	updated = append(updated, sectionLines...)
	updated = append(updated, lines[end:]...)

	updatedBody := strings.Join(updated, "\n")
	return updatedBody, updatedBody != body
}

func insertTasksSection(lines []string, entries []string) (string, bool) {
	sectionLines := buildTasksSectionLines(nil, entries)
	insertAt := findHeaderLine(lines, "## Acceptance Criteria")
	if insertAt == -1 {
		insertAt = len(lines)
	}

	updated := make([]string, 0, len(lines)+len(sectionLines)+2)
	updated = append(updated, lines[:insertAt]...)

	if insertAt > 0 && strings.TrimSpace(updated[len(updated)-1]) != "" {
		updated = append(updated, "")
	}
	updated = append(updated, sectionLines...)
	if insertAt < len(lines) && strings.TrimSpace(lines[insertAt]) != "" {
		updated = append(updated, "")
	}
	updated = append(updated, lines[insertAt:]...)

	updatedBody := strings.Join(updated, "\n")
	return updatedBody, updatedBody != strings.Join(lines, "\n")
}

func buildTasksSectionLines(existing []string, entries []string) []string {
	section := []string{"## Tasks"}
	content := append([]string{}, existing...)

	if len(entries) > 0 {
		if len(content) > 0 && strings.TrimSpace(content[len(content)-1]) != "" {
			content = append(content, "")
		}
		content = append(content, entries...)
	}

	if len(content) == 0 {
		content = append(content, "")
	}
	if len(content) > 0 && strings.TrimSpace(content[0]) != "" {
		section = append(section, "")
	}
	section = append(section, content...)
	return section
}

func findHeaderLine(lines []string, header string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == header {
			return i
		}
	}
	return -1
}

func findNextHeader(lines []string, start int) int {
	for i := start; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "# ") {
			return i
		}
	}
	return len(lines)
}

func isSubtaskTodoLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return subtaskTodoPattern.MatchString(trimmed)
}

func trimBlankLines(lines []string) []string {
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	return lines[start:end]
}
