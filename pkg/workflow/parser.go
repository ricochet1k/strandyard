package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

// Parser handles parsing role and template files
type Parser struct {
	md goldmark.Markdown
}

// NewParser creates a new workflow parser
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
	)
	return &Parser{md: md}
}

// ParseRole parses a role file and returns a Role
func (p *Parser) ParseRole(filePath string) (*Role, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read role file %s: %w", filePath, err)
	}

	// Parse the markdown with frontmatter
	var meta RoleMetadata
	ctx := parser.NewContext()
	_ = p.md.Parser().Parse(text.NewReader(data), parser.WithContext(ctx))

	// Extract frontmatter (optional for roles)
	fm := frontmatter.Get(ctx)
	if fm != nil {
		if err := fm.Decode(&meta); err != nil {
			return nil, fmt.Errorf("failed to decode frontmatter in %s: %w", filePath, err)
		}
	}

	// Extract role name from filename if not in frontmatter
	name := filepath.Base(filePath)
	name = strings.TrimSuffix(name, ".md")
	if meta.Role == "" {
		meta.Role = name
	}

	role := &Role{
		Name:     name,
		FilePath: filePath,
		Meta:     meta,
		Content:  string(data),
	}

	return role, nil
}

// ParseTemplate parses a template file and returns a Template
func (p *Parser) ParseTemplate(filePath string) (*Template, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file %s: %w", filePath, err)
	}

	// Parse the markdown with frontmatter
	var meta TemplateMetadata
	ctx := parser.NewContext()
	_ = p.md.Parser().Parse(text.NewReader(data), parser.WithContext(ctx))

	// Extract frontmatter
	fm := frontmatter.Get(ctx)
	if fm != nil {
		if err := fm.Decode(&meta); err != nil {
			return nil, fmt.Errorf("failed to decode frontmatter in %s: %w", filePath, err)
		}
	}

	// Extract template name from filename
	name := filepath.Base(filePath)
	name = strings.TrimSuffix(name, ".md")

	// Extract TODO sequence
	todoSequence := extractTodoSequence(string(data))

	template := &Template{
		Name:         name,
		FilePath:     filePath,
		Meta:         meta,
		Content:      string(data),
		TodoSequence: todoSequence,
	}

	return template, nil
}

// LoadRoles loads all role files from a directory
func (p *Parser) LoadRoles(rolesDir string) (map[string]*Role, error) {
	roles := make(map[string]*Role)

	entries, err := os.ReadDir(rolesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return roles, nil // Empty roles map if directory doesn't exist
		}
		return nil, fmt.Errorf("failed to read roles directory %s: %w", rolesDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(rolesDir, entry.Name())
		role, err := p.ParseRole(filePath)
		if err != nil {
			return nil, err
		}

		roles[role.Name] = role
	}

	return roles, nil
}

// LoadTemplates loads all template files from a directory
func (p *Parser) LoadTemplates(templatesDir string) (map[string]*Template, error) {
	templates := make(map[string]*Template)

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return templates, nil // Empty templates map if directory doesn't exist
		}
		return nil, fmt.Errorf("failed to read templates directory %s: %w", templatesDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(templatesDir, entry.Name())
		template, err := p.ParseTemplate(filePath)
		if err != nil {
			return nil, err
		}

		templates[template.Name] = template
	}

	return templates, nil
}

// extractTodoSequence extracts TODO items with role annotations from content
// Format: 1. [ ] (role: developer) Do something
var todoRolePattern = regexp.MustCompile(`^\s*\d+\.\s*\[[ x]\]\s*\(role:\s*([^)]+)\)\s*(.*)`)

func extractTodoSequence(content string) []TodoStep {
	var steps []TodoStep
	lines := strings.Split(content, "\n")
	todoNum := 0

	for _, line := range lines {
		matches := todoRolePattern.FindStringSubmatch(line)
		if matches != nil && len(matches) >= 3 {
			todoNum++
			role := strings.TrimSpace(matches[1])
			description := strings.TrimSpace(matches[2])

			steps = append(steps, TodoStep{
				Number:      todoNum,
				Role:        role,
				Description: description,
			})
		}
	}

	return steps
}
