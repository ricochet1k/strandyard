package task

import (
	"fmt"
	"regexp"
	"strings"
)

var itemPattern = regexp.MustCompile(`^(?:\d+\.|-)\s*(?:\[([ xX])\]\s*)?(?:\(role:\s*([^)]+)\)\s*)?(?:\(subtask:\s*([^)]+)\)\s*)?(.*)$`)

// TaskItem represents an entry in a task list (TODOs or Subtasks).
type TaskItem struct {
	Checked   bool   `json:"checked"`
	Role      string `json:"role,omitempty"`
	SubtaskID string `json:"subtask_id,omitempty"`
	Text      string `json:"text"`
}

// ParseTaskItems parses a section's content into structured items.
func ParseTaskItems(content string) []TaskItem {
	lines := strings.Split(content, "\n")
	var items []TaskItem

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		match := itemPattern.FindStringSubmatch(trimmed)
		if match == nil {
			continue
		}
		checked := strings.ToLower(match[1]) == "x"
		role := strings.TrimSpace(match[2])
		subtaskID := strings.TrimSpace(match[3])
		text := strings.TrimSpace(match[4])

		items = append(items, TaskItem{
			Checked:   checked,
			Role:      role,
			SubtaskID: subtaskID,
			Text:      text,
		})
	}

	return items
}

// FormatTodoItems formats a list of items as a numbered list.
func FormatTodoItems(items []TaskItem) string {
	var sb strings.Builder
	for i, item := range items {
		if i > 0 {
			sb.WriteString("\n")
		}
		status := " "
		if item.Checked {
			status = "x"
		}
		sb.WriteString(fmt.Sprintf("- [%s] ", status))
		if item.Role != "" {
			sb.WriteString(fmt.Sprintf("(role: %s) ", item.Role))
		}
		if item.SubtaskID != "" {
			sb.WriteString(fmt.Sprintf("(subtask: %s) ", ShortID(item.SubtaskID)))
		}
		sb.WriteString(item.Text)
	}
	return sb.String()
}

// FormatSubtaskItems formats a list of items as a bulleted list.
func FormatSubtaskItems(items []TaskItem) string {
	var sb strings.Builder
	for i, item := range items {
		if i > 0 {
			sb.WriteString("\n")
		}
		status := " "
		if item.Checked {
			status = "x"
		}
		sb.WriteString(fmt.Sprintf("- [%s] ", status))
		if item.Role != "" {
			sb.WriteString(fmt.Sprintf("(role: %s) ", item.Role))
		}
		if item.SubtaskID != "" {
			sb.WriteString(fmt.Sprintf("(subtask: %s) ", ShortID(item.SubtaskID)))
		}
		sb.WriteString(item.Text)
	}
	return sb.String()
}
