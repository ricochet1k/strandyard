package workflow

import (
	"fmt"
	"sort"
	"strings"
)

// GenerateMermaid generates a Mermaid diagram of the workflow
func (g *WorkflowGraph) GenerateMermaid() string {
	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("graph TD\n")

	// Build edges with deduplication
	edgeMap := make(map[string]WorkflowEdge)
	for _, edge := range g.Edges {
		key := fmt.Sprintf("%s->%s->%s", edge.FromRole, edge.TaskType, edge.ToRole)
		edgeMap[key] = edge
	}

	// Sort edges for consistent output
	var edges []WorkflowEdge
	for _, edge := range edgeMap {
		edges = append(edges, edge)
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].FromRole != edges[j].FromRole {
			return edges[i].FromRole < edges[j].FromRole
		}
		if edges[i].TaskType != edges[j].TaskType {
			return edges[i].TaskType < edges[j].TaskType
		}
		return edges[i].ToRole < edges[j].ToRole
	})

	// Generate edges
	for _, edge := range edges {
		fromNode := sanitizeMermaidID(edge.FromRole)
		toNode := sanitizeMermaidID(edge.ToRole)
		label := edge.TaskType

		sb.WriteString(fmt.Sprintf("    %s[\"%s\"] -->|%s| %s[\"%s\"]\n",
			fromNode, edge.FromRole, label, toNode, edge.ToRole))
	}

	// Add styling
	sb.WriteString("\n")
	roleColors := getRoleColors()
	for roleName := range g.Roles {
		if color, ok := roleColors[roleName]; ok {
			nodeID := sanitizeMermaidID(roleName)
			sb.WriteString(fmt.Sprintf("    style %s fill:%s\n", nodeID, color))
		} else {
			// Default color for roles not in predefined map
			nodeID := sanitizeMermaidID(roleName)
			sb.WriteString(fmt.Sprintf("    style %s fill:#e0e0e0\n", nodeID))
		}
	}

	sb.WriteString("```")
	return sb.String()
}

// GenerateMermaidForTemplate generates a Mermaid diagram showing the TODO sequence within a template
func (g *WorkflowGraph) GenerateMermaidForTemplate(templateName string) (string, error) {
	template, ok := g.Templates[templateName]
	if !ok {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var sb strings.Builder
	sb.WriteString("```mermaid\n")
	sb.WriteString("graph TD\n")

	if len(template.TodoSequence) == 0 {
		sb.WriteString("    NoTodos[\"No TODOs in template\"]\n")
		sb.WriteString("```")
		return sb.String(), nil
	}

	// Generate nodes and edges for TODO sequence
	for i, step := range template.TodoSequence {
		nodeID := fmt.Sprintf("Todo%d", i+1)
		label := fmt.Sprintf("%s: %s", step.Role, truncate(step.Description, 40))

		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", nodeID, label))

		// Connect to next step
		if i < len(template.TodoSequence)-1 {
			nextID := fmt.Sprintf("Todo%d", i+1)
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", nodeID, nextID))
		}
	}

	// Add styling by role
	sb.WriteString("\n")
	roleColors := getRoleColors()
	for i, step := range template.TodoSequence {
		nodeID := fmt.Sprintf("Todo%d", i+1)
		if color, ok := roleColors[step.Role]; ok {
			sb.WriteString(fmt.Sprintf("    style %s fill:%s\n", nodeID, color))
		} else {
			sb.WriteString(fmt.Sprintf("    style %s fill:#e0e0e0\n", nodeID))
		}
	}

	sb.WriteString("```")
	return sb.String(), nil
}

// sanitizeMermaidID converts a role name to a valid Mermaid node ID
func sanitizeMermaidID(name string) string {
	// Replace hyphens and other special characters with underscores
	result := strings.ReplaceAll(name, "-", "_")
	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ReplaceAll(result, ".", "_")
	return result
}

// getRoleColors returns a map of role names to colors for styling
func getRoleColors() map[string]string {
	return map[string]string{
		"owner":                "#99ff99",
		"designer":             "#ffff99",
		"architect":            "#99ffff",
		"developer":            "#9999ff",
		"tester":               "#ff99cc",
		"documentation":        "#ccff99",
		"triage":               "#ffcc99",
		"reviewer":             "#cc99ff",
		"master-reviewer":      "#cc99ff",
		"reviewer-security":    "#ff99ff",
		"reviewer-reliability": "#ff99ff",
		"reviewer-usability":   "#ff99ff",
	}
}

// truncate truncates a string to maxLen characters, adding "..." if truncated
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
