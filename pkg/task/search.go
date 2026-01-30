package task

import (
	"fmt"
	"strings"
)

// SearchOptions defines parameters for searching tasks.
type SearchOptions struct {
	Query string
	ListOptions
}

// SearchTasks loads tasks and returns those matching the search query.
func SearchTasks(tasksRoot string, opts SearchOptions) ([]*Task, error) {
	query := strings.TrimSpace(opts.Query)
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	parser := NewParser()
	tasks, err := parser.LoadTasks(tasksRoot)
	if err != nil {
		return nil, err
	}

	items, err := filterTasks(tasksRoot, tasks, opts.ListOptions)
	if err != nil {
		return nil, err
	}

	matched := make([]*Task, 0, len(items))
	for _, t := range items {
		ok, err := matchesQuery(t, query)
		if err != nil {
			return nil, err
		}
		if ok {
			matched = append(matched, t)
		}
	}

	sortTasks(matched, opts.ListOptions)
	return matched, nil
}

func matchesQuery(t *Task, query string) (bool, error) {
	text, err := taskSearchText(t)
	if err != nil {
		return false, err
	}
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return false, nil
	}
	return strings.Contains(strings.ToLower(text), q), nil
}

func taskSearchText(t *Task) (string, error) {
	body, err := taskBody(t.Content, t.FilePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(body, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if isSubtaskTodoLine(line) {
			continue
		}
		filtered = append(filtered, line)
	}

	return t.Title() + "\n" + strings.Join(filtered, "\n"), nil
}

func taskBody(content, path string) (string, error) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return "", errInvalidFrontmatter(path)
	}
	return parts[2], nil
}
