package workflow

import (
	"fmt"
)

// Validate checks the workflow for common issues and returns validation results
func (g *WorkflowGraph) Validate() *ValidationResult {
	result := &ValidationResult{
		Errors:   []ValidationIssue{},
		Warnings: []ValidationIssue{},
	}

	// Check 1: All roles referenced in templates must be defined
	definedRoles := make(map[string]bool)
	for roleName := range g.Roles {
		definedRoles[roleName] = true
	}

	referencedRoles := make(map[string]bool)
	for _, template := range g.Templates {
		if template.Meta.Role != "" {
			referencedRoles[template.Meta.Role] = true
		}
		for _, step := range template.TodoSequence {
			referencedRoles[step.Role] = true
		}
	}

	for roleName := range referencedRoles {
		if !definedRoles[roleName] {
			result.Warnings = append(result.Warnings, ValidationIssue{
				Severity: "warning",
				Message:  fmt.Sprintf("Role '%s' is used in templates but has no role definition file", roleName),
				Location: "templates",
			})
		}
	}

	// Check 2: All task types in role.creates must exist as templates
	for _, role := range g.Roles {
		for _, taskType := range role.Meta.Workflow.Creates {
			if _, exists := g.Templates[taskType]; !exists {
				result.Errors = append(result.Errors, ValidationIssue{
					Severity: "error",
					Message:  fmt.Sprintf("Role '%s' creates task type '%s' but template does not exist", role.Name, taskType),
					Location: role.FilePath,
				})
			}
		}
	}

	// Check 3: Warn about roles without workflow metadata
	for _, role := range g.Roles {
		if len(role.Meta.Workflow.Creates) == 0 {
			result.Warnings = append(result.Warnings, ValidationIssue{
				Severity: "warning",
				Message:  fmt.Sprintf("Role '%s' has no workflow metadata (missing 'creates' field)", role.Name),
				Location: role.FilePath,
			})
		}
	}

	// Check 4: Warn about templates never created by any role
	createdTemplates := make(map[string]bool)
	for _, role := range g.Roles {
		for _, taskType := range role.Meta.Workflow.Creates {
			createdTemplates[taskType] = true
		}
	}

	for templateName := range g.Templates {
		if !createdTemplates[templateName] {
			result.Warnings = append(result.Warnings, ValidationIssue{
				Severity: "warning",
				Message:  fmt.Sprintf("Task type '%s' is never created by any role", templateName),
				Location: g.Templates[templateName].FilePath,
			})
		}
	}

	// Check 5: Error if template has no role assignment and no TODOs
	for _, template := range g.Templates {
		if template.Meta.Role == "" && len(template.TodoSequence) == 0 {
			result.Errors = append(result.Errors, ValidationIssue{
				Severity: "error",
				Message:  fmt.Sprintf("Template '%s' has no role assignment and no TODOs", template.Name),
				Location: template.FilePath,
			})
		}
	}

	// Check 6: Warn about orphaned roles (defined but never used)
	usedRoles := make(map[string]bool)
	for _, template := range g.Templates {
		if template.Meta.Role != "" {
			usedRoles[template.Meta.Role] = true
		}
		for _, step := range template.TodoSequence {
			usedRoles[step.Role] = true
		}
	}

	for roleName := range g.Roles {
		if !usedRoles[roleName] {
			result.Warnings = append(result.Warnings, ValidationIssue{
				Severity: "warning",
				Message:  fmt.Sprintf("Role '%s' is defined but not used in any template", roleName),
				Location: g.Roles[roleName].FilePath,
			})
		}
	}

	return result
}

// HasErrors returns true if validation found any errors
func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// HasWarnings returns true if validation found any warnings
func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// Summary returns a human-readable summary of validation results
func (r *ValidationResult) Summary() string {
	errorCount := len(r.Errors)
	warningCount := len(r.Warnings)

	if errorCount == 0 && warningCount == 0 {
		return "✓ Workflow validation passed with no issues"
	}

	return fmt.Sprintf("%d warnings, %d errors", warningCount, errorCount)
}
