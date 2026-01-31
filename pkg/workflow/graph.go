package workflow

import (
	"fmt"
	"time"
)

// BuildGraph constructs a workflow graph from roles and templates
func BuildGraph(roles map[string]*Role, templates map[string]*Template) *WorkflowGraph {
	graph := &WorkflowGraph{
		Roles:     roles,
		Templates: templates,
		Edges:     []WorkflowEdge{},
	}

	// Build edges from role.creates -> template.role
	for _, role := range roles {
		for _, taskType := range role.Meta.Workflow.Creates {
			// Find the template for this task type
			template, ok := templates[taskType]
			if !ok {
				// Template doesn't exist - we'll catch this in validation
				continue
			}

			// Create edge from role -> template's role
			toRole := template.Meta.Role
			if toRole == "" {
				// No role in template, try to infer from first TODO
				if len(template.TodoSequence) > 0 {
					toRole = template.TodoSequence[0].Role
				}
			}

			if toRole != "" {
				edge := WorkflowEdge{
					FromRole: role.Name,
					TaskType: taskType,
					ToRole:   toRole,
					Label:    fmt.Sprintf("creates %s", taskType),
					Source:   fmt.Sprintf("%s.md workflow.creates", role.Name),
				}
				graph.Edges = append(graph.Edges, edge)
			}
		}
	}

	return graph
}

// GetStats calculates statistics about the workflow
func (g *WorkflowGraph) GetStats() WorkflowStats {
	stats := WorkflowStats{
		TotalRoles:     len(g.Roles),
		TotalTemplates: len(g.Templates),
		TotalEdges:     len(g.Edges),
		LastUpdated:    time.Now(),
	}

	// Count roles with/without workflow metadata
	for _, role := range g.Roles {
		if len(role.Meta.Workflow.Creates) > 0 {
			stats.RolesWithWorkflow++
		} else {
			stats.RolesWithoutWorkflow++
		}
	}

	// Find roles that are defined but never used (orphaned)
	usedRoles := make(map[string]bool)

	// Mark roles used in templates (primary role assignment)
	for _, template := range g.Templates {
		if template.Meta.Role != "" {
			usedRoles[template.Meta.Role] = true
		}
		// Also mark roles used in TODOs
		for _, step := range template.TodoSequence {
			usedRoles[step.Role] = true
		}
	}

	// Check which defined roles are never used
	for roleName := range g.Roles {
		if !usedRoles[roleName] {
			stats.OrphanedRoles = append(stats.OrphanedRoles, roleName)
		}
	}

	// Find roles that are referenced but not defined
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
			stats.UndefinedRoles = append(stats.UndefinedRoles, roleName)
		}
	}

	// Find templates that are never created by any role
	createdTemplates := make(map[string]bool)
	for _, role := range g.Roles {
		for _, taskType := range role.Meta.Workflow.Creates {
			createdTemplates[taskType] = true
		}
	}

	for templateName := range g.Templates {
		if !createdTemplates[templateName] {
			stats.UncreatableTemplates = append(stats.UncreatableTemplates, templateName)
		}
	}

	// Find templates without role assignment
	for _, template := range g.Templates {
		if template.Meta.Role == "" && len(template.TodoSequence) == 0 {
			stats.TemplatesWithoutRole = append(stats.TemplatesWithoutRole, template.Name)
		}
	}

	return stats
}

// GetRoleUsage returns information about how a role is used in the workflow
func (g *WorkflowGraph) GetRoleUsage(roleName string) *RoleUsageInfo {
	role, exists := g.Roles[roleName]
	if !exists {
		return nil
	}

	info := &RoleUsageInfo{
		RoleName:          roleName,
		AssignedTemplates: []string{},
		UsedInTodos:       make(map[string][]int),
		ReceivesVia:       []string{},
		Creates:           role.Meta.Workflow.Creates,
	}

	// Find templates where this role is the primary assigned role
	for _, template := range g.Templates {
		if template.Meta.Role == roleName {
			info.AssignedTemplates = append(info.AssignedTemplates, template.Name)
		}

		// Find TODOs that use this role
		for i, step := range template.TodoSequence {
			if step.Role == roleName {
				info.UsedInTodos[template.Name] = append(info.UsedInTodos[template.Name], i+1)
			}
		}
	}

	// Find which task types lead to this role
	for _, edge := range g.Edges {
		if edge.ToRole == roleName {
			info.ReceivesVia = append(info.ReceivesVia, edge.TaskType)
		}
	}

	return info
}

// RoleUsageInfo contains information about how a role is used
type RoleUsageInfo struct {
	RoleName          string
	AssignedTemplates []string         // Templates where this is the primary role
	UsedInTodos       map[string][]int // Template name -> TODO numbers
	ReceivesVia       []string         // Task types that assign to this role
	Creates           []string         // Task types this role creates
}
