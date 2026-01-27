package task

import (
	"os"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// UpdateBlockersFromChildren ensures parent tasks are blocked by incomplete children.
// Returns the number of task files updated.
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
		desired := mergeBlockers(parent.Meta.Blockers, incomplete, childSet)
		if equalStringSlices(parent.Meta.Blockers, desired) {
			continue
		}

		parent.Meta.Blockers = desired
		parent.Meta.DateEdited = now
		if err := writeFrontmatter(parent.FilePath, parent.Meta); err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}

func mergeBlockers(existing []string, childIncomplete []string, childSet map[string]struct{}) []string {
	base := []string{}
	seen := map[string]struct{}{}
	for _, b := range existing {
		if b == "" {
			continue
		}
		if _, isChild := childSet[b]; isChild {
			continue
		}
		if _, ok := seen[b]; ok {
			continue
		}
		seen[b] = struct{}{}
		base = append(base, b)
	}

	for _, b := range childIncomplete {
		if _, ok := seen[b]; ok {
			continue
		}
		seen[b] = struct{}{}
		base = append(base, b)
	}

	sort.Strings(base)
	return base
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func writeFrontmatter(path string, meta Metadata) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return errInvalidFrontmatter(path)
	}
	frontmatterBytes, err := yaml.Marshal(&meta)
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(frontmatterBytes)
	sb.WriteString("---")
	sb.WriteString(parts[2])
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func errInvalidFrontmatter(path string) error {
	return &InvalidFrontmatterError{Path: path}
}

// InvalidFrontmatterError indicates a task file is missing frontmatter delimiters.
type InvalidFrontmatterError struct {
	Path string
}

func (e *InvalidFrontmatterError) Error() string {
	return "invalid task file format: missing frontmatter delimiters in " + e.Path
}
