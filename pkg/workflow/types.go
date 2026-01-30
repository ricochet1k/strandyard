package workflow

import (
	"time"
)

// RoleMetadata represents the frontmatter of a role file
type RoleMetadata struct {
	Role     string         `yaml:"role"`
	Workflow RoleWorkflowMeta `yaml:"workflow,omitempty"`
}

// RoleWorkflowMeta contains workflow-specific metadata from role frontmatter
type RoleWorkflowMeta struct {
	Creates []string `yaml:"creates,omitempty"`
}

// Role represents a complete role definition with metadata and content
type Role struct {
	Name     string
	FilePath string
	Meta     RoleMetadata
	Content  string
}

// TemplateMetadata represents the frontmatter of a template file
type TemplateMetadata struct {
	Role       string   `yaml:"role"`
	Priority   string   `yaml:"priority,omitempty"`
	IDPrefix   string   `yaml:"id_prefix,omitempty"`
}

// Template represents a task template with metadata and content
type Template struct {
	Name         string
	FilePath     string
	Meta         TemplateMetadata
	Content      string
	TodoSequence []TodoStep
}

// TodoStep represents a single TODO item from a template
type TodoStep struct {
	Number      int
	Role        string
	Description string
}

// WorkflowGraph represents the complete workflow extracted from roles and templates
type WorkflowGraph struct {
	Roles     map[string]*Role
	Templates map[string]*Template
	Edges     []WorkflowEdge
}

// WorkflowEdge represents a connection from one role to another via a task type
type WorkflowEdge struct {
	FromRole   string // Role that creates the task
	TaskType   string // Task type created (template name)
	ToRole     string // Role assigned in task template
	Label      string // Description of the transition
	Source     string // Where this edge came from (e.g., "designer.md creates")
}

// ValidationIssue represents a problem found during workflow validation
type ValidationIssue struct {
	Severity string // "error" or "warning"
	Message  string
	Location string // File or role/template name where issue was found
}

// ValidationResult contains all validation issues found
type ValidationResult struct {
	Errors   []ValidationIssue
	Warnings []ValidationIssue
}

// WorkflowStats provides statistics about the workflow
type WorkflowStats struct {
	TotalRoles              int
	TotalTemplates          int
	TotalEdges              int
	RolesWithWorkflow       int
	RolesWithoutWorkflow    int
	OrphanedRoles           []string // Roles defined but never used
	UndefinedRoles          []string // Roles referenced but not defined
	UncreatableTemplates    []string // Templates never created by any role
	TemplatesWithoutRole    []string // Templates with no role assignment
	LastUpdated             time.Time
}
