package workflow

import (
	"encoding/json"
)

// JSONOutput represents the JSON output format for workflow data
type JSONOutput struct {
	Roles     map[string]JSONRole     `json:"roles"`
	Templates map[string]JSONTemplate `json:"templates"`
	Edges     []WorkflowEdge          `json:"edges"`
	Stats     *WorkflowStats          `json:"stats,omitempty"`
}

// JSONRole represents a role in JSON format
type JSONRole struct {
	Name    string   `json:"name"`
	Creates []string `json:"creates,omitempty"`
}

// JSONTemplate represents a template in JSON format
type JSONTemplate struct {
	Name         string     `json:"name"`
	Role         string     `json:"role,omitempty"`
	TodoSequence []TodoStep `json:"todo_sequence,omitempty"`
}

// ToJSON converts the workflow graph to JSON output format
func (g *WorkflowGraph) ToJSON(includeStats bool) ([]byte, error) {
	output := JSONOutput{
		Roles:     make(map[string]JSONRole),
		Templates: make(map[string]JSONTemplate),
		Edges:     g.Edges,
	}

	// Convert roles
	for name, role := range g.Roles {
		output.Roles[name] = JSONRole{
			Name:    role.Name,
			Creates: role.Meta.Workflow.Creates,
		}
	}

	// Convert templates
	for name, template := range g.Templates {
		output.Templates[name] = JSONTemplate{
			Name:         template.Name,
			Role:         template.Meta.Role,
			TodoSequence: template.TodoSequence,
		}
	}

	// Include stats if requested
	if includeStats {
		stats := g.GetStats()
		output.Stats = &stats
	}

	return json.MarshalIndent(output, "", "  ")
}
