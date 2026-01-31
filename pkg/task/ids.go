package task

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	shortIDPattern = regexp.MustCompile(`^[A-Z][0-9a-z]{4,6}$`)
	fullIDPattern  = regexp.MustCompile(`^([A-Z][0-9a-z]{4,6})-[a-zA-Z0-9-]+$`)
)

// ShortID returns the short form of a task ID (prefix + token).
func ShortID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	if matches := fullIDPattern.FindStringSubmatch(id); len(matches) == 2 {
		return matches[1]
	}
	if shortIDPattern.MatchString(id) {
		return id
	}
	return id
}

// ResolveTaskID resolves a short or full task ID to the full ID.
func ResolveTaskID(tasks map[string]*Task, input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}
	if _, ok := tasks[input]; ok {
		return input, nil
	}

	if extracted := extractTaskIDFromPath(input); extracted != "" {
		if _, ok := tasks[extracted]; ok {
			return extracted, nil
		}
	}

	// Try exact short ID match first
	if shortIDPattern.MatchString(input) {
		matches := make([]string, 0, 2)
		for id := range tasks {
			if strings.HasPrefix(id, input+"-") {
				matches = append(matches, id)
			}
		}
		if len(matches) == 1 {
			return matches[0], nil
		}
		if len(matches) > 1 {
			sort.Strings(matches)
			return "", fmt.Errorf("short id %s is ambiguous: %s", input, strings.Join(matches, ", "))
		}
	}

	// Try partial prefix match for any valid prefix
	matches := make([]string, 0, 2)
	for id := range tasks {
		if strings.HasPrefix(id, input) {
			matches = append(matches, id)
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		sort.Strings(matches)
		return "", fmt.Errorf("prefix %s is ambiguous: %s", input, strings.Join(matches, ", "))
	}

	return "", fmt.Errorf("task not found: %s", input)
}
