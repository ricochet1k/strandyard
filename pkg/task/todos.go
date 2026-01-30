package task

import (
	"regexp"
	"strings"
)

var todoItemPattern = regexp.MustCompile(`^(?:\d+\.|-)\s*\[([ xX])\]\s*(?:\(role:\s*([^)]+)\)\s*)?(.*)$`)
var subtaskItemPattern = regexp.MustCompile(`^-\s*\[([ xX])\]\s*\(subtask:\s*([^)]+)\)\s*(.*)$`)

// TodoItem represents an entry in the TODOs section.
type TodoItem struct {
	Index   int    `json:"index"`
	Checked bool   `json:"checked"`
	Role    string `json:"role,omitempty"`
	Text    string `json:"text"`
	Raw     string `json:"raw"`
}

// SubtaskItem represents a subtask entry in the Tasks section.
type SubtaskItem struct {
	ID      string `json:"id"`
	Checked bool   `json:"checked"`
	Title   string `json:"title"`
	Raw     string `json:"raw"`
}

// ParseTodoItems parses the ## TODOs section into structured items.
func ParseTodoItems(content, path string) ([]TodoItem, error) {
	body, err := taskBody(content, path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(body, "\n")
	start := findHeaderLine(lines, "## TODOs")
	if start == -1 {
		return nil, nil
	}
	end := findNextHeader(lines, start+1)

	items := []TodoItem{}
	index := 1
	for _, line := range lines[start+1 : end] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		match := todoItemPattern.FindStringSubmatch(trimmed)
		if match == nil {
			continue
		}
		checked := strings.ToLower(match[1]) == "x"
		role := strings.TrimSpace(match[2])
		text := strings.TrimSpace(match[3])
		items = append(items, TodoItem{
			Index:   index,
			Checked: checked,
			Role:    role,
			Text:    text,
			Raw:     trimmed,
		})
		index++
	}

	return items, nil
}

// ParseSubtaskItems parses subtask entries in the ## Tasks section.
func ParseSubtaskItems(content, path string) ([]SubtaskItem, error) {
	body, err := taskBody(content, path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(body, "\n")
	start := findHeaderLine(lines, "## Tasks")
	if start == -1 {
		return nil, nil
	}
	end := findNextHeader(lines, start+1)

	items := []SubtaskItem{}
	for _, line := range lines[start+1 : end] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		match := subtaskItemPattern.FindStringSubmatch(trimmed)
		if match == nil {
			continue
		}
		checked := strings.ToLower(match[1]) == "x"
		id := strings.TrimSpace(match[2])
		title := strings.TrimSpace(match[3])
		items = append(items, SubtaskItem{
			ID:      id,
			Checked: checked,
			Title:   title,
			Raw:     trimmed,
		})
	}

	return items, nil
}
