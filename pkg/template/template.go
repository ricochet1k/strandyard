package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// TemplateMetadata represents the YAML frontmatter of a template,
// which may contain Go template expressions.
type TemplateMetadata struct {
	Type          string      `yaml:"type"`
	Role          string      `yaml:"role"`
	Priority      interface{} `yaml:"priority"`
	Parent        string      `yaml:"parent"`
	Blockers      []string    `yaml:"blockers"`
	Blocks        []string    `yaml:"blocks"`
	DateCreated   interface{} `yaml:"date_created"`
	DateEdited    interface{} `yaml:"date_edited"`
	OwnerApproval bool        `yaml:"owner_approval"`
	Completed     bool        `yaml:"completed"`
	Description   string      `yaml:"description"`
	IDPrefix      string      `yaml:"id_prefix"`
}

// Template represents a parsed task template.
type Template struct {
	ID          string
	Meta        TemplateMetadata
	BodyContent string
}

// LoadTemplates loads all templates from the given directory.
func LoadTemplates(templatesDir string) (map[string]*Template, error) {
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	templates := make(map[string]*Template)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		templatePath := filepath.Join(templatesDir, entry.Name())
		data, err := os.ReadFile(templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read template file %s: %w", templatePath, err)
		}

		content := string(data)
		parts := strings.SplitN(content, "---", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid template file format: missing frontmatter delimiters in %s", templatePath)
		}

		var meta TemplateMetadata
		if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
			return nil, fmt.Errorf("failed to decode frontmatter in %s: %w", templatePath, err)
		}

		id := strings.TrimSuffix(entry.Name(), ".md")
		templates[id] = &Template{
			ID:          id,
			Meta:        meta,
			BodyContent: strings.TrimSpace(parts[2]),
		}
	}

	return templates, nil
}
