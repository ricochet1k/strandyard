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
	Report    string `json:"report,omitempty"`
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
			// If it doesn't match the pattern, it might be an indented report line for the last item
			if len(items) > 0 && (strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t")) {
				if items[len(items)-1].Report != "" {
					items[len(items)-1].Report += "\n"
				}
				items[len(items)-1].Report += strings.TrimSpace(line)
			}
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

// FormatTodoItems formats a list of items as a bulleted list.
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
		if item.Report != "" {
			reportLines := strings.Split(item.Report, "\n")
			for _, rl := range reportLines {
				sb.WriteString("\n  ")
				sb.WriteString(rl)
			}
		}
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
		if item.Report != "" {
			reportLines := strings.Split(item.Report, "\n")
			for _, rl := range reportLines {
				sb.WriteString("\n  ")
				sb.WriteString(rl)
			}
		}
	}
	return sb.String()
}
