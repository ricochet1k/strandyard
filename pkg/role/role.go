package role

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ricochet1k/strandyard/pkg/task"
)

// LoadRoles loads all roles from the given directory.
func LoadRoles(rolesDir string) (map[string]*task.Task, error) {
	entries, err := os.ReadDir(rolesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read roles directory: %w", err)
	}

	parser := task.NewParser()
	roles := make(map[string]*task.Task)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		rolePath := filepath.Join(rolesDir, entry.Name())
		r, err := parser.ParseStandaloneFile(rolePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse role file %s: %w", rolePath, err)
		}

		roles[r.ID] = r
	}

	return roles, nil
}

// ValidateRole checks if a role exists in the given directory.
func ValidateRole(rolesDir, role string) error {
	rolePath := filepath.Join(rolesDir, role+".md")
	if _, err := os.Stat(rolePath); err != nil {
		return fmt.Errorf("role file %s does not exist", rolePath)
	}
	return nil
}
