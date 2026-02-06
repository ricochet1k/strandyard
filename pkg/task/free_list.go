package task

import (
	"path/filepath"
	"regexp"
	"strings"
)

// FreeListParse represents parsed free-tasks.md data.
type FreeListParse struct {
	Title   string
	TaskIDs []string
}

// ParseFreeList parses free-tasks.md content and resolves task IDs.
// It tolerates malformed list entries by extracting IDs from paths when possible.
func ParseFreeList(content string, tasks map[string]*Task) FreeListParse {
	title := "Free tasks"
	ids := []string{}
	seen := make(map[string]bool)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			parsedTitle := strings.TrimSpace(strings.TrimPrefix(line, "# "))
			if parsedTitle != "" {
				title = parsedTitle
			}
			continue
		}
		if !strings.HasPrefix(line, "- ") {
			continue
		}

		entry := strings.TrimSpace(strings.TrimPrefix(line, "- "))
		path := parseListPath(entry)
		if path == "" {
			path = entry
		}

		taskID := resolveTaskIDFromListPath(path, tasks)
		if taskID == "" {
			taskID = resolveTaskIDFromListPath(entry, tasks)
		}

		if taskID != "" && !seen[taskID] {
			seen[taskID] = true
			ids = append(ids, taskID)
		}
	}

	return FreeListParse{Title: title, TaskIDs: ids}
}

func resolveTaskIDFromListPath(path string, tasks map[string]*Task) string {
	cleaned := filepath.ToSlash(strings.TrimSpace(path))
	for id, task := range tasks {
		taskPath := filepath.ToSlash(task.FilePath)
		if cleaned == taskPath || strings.HasSuffix(cleaned, taskPath) || strings.HasSuffix(taskPath, cleaned) {
			return id
		}
	}

	return extractTaskIDFromPathLast(cleaned)
}

func extractTaskIDFromPathLast(path string) string {
	path = filepath.Clean(path)
	parts := strings.Split(filepath.ToSlash(path), "/")
	idPattern := regexp.MustCompile(`^[A-Z][0-9a-z]{4,6}-[a-zA-Z0-9-]+$`)
	match := ""
	for _, part := range parts {
		if idPattern.MatchString(part) {
			match = part
		}
	}
	return match
}

func parseListPath(entry string) string {
	if !strings.HasPrefix(entry, "[") {
		return entry
	}
	open := strings.Index(entry, "](")
	close := strings.LastIndex(entry, ")")
	if open == -1 || close == -1 || close <= open+2 {
		return ""
	}
	return strings.TrimSpace(entry[open+2 : close])
}
