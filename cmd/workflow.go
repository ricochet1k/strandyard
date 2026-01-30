package cmd

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ricochet1k/strandyard/pkg/workflow"
	"github.com/spf13/cobra"
)

var (
	workflowFormat   string
	workflowTemplate string
	workflowRole     string
	workflowValidate bool
	workflowStats    bool
)

// workflowCmd represents the workflow command
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Extract and visualize role-based workflows",
	Long: `Extract workflow information from roles and templates.

The workflow command analyzes role definitions and task templates to generate
workflow diagrams, validate workflow completeness, and show role usage.

Examples:
  # Generate Mermaid diagram of the complete workflow
  strand workflow --format mermaid

  # Generate JSON output with statistics
  strand workflow --format json --stats

  # Show workflow for a specific template
  strand workflow --template task --format mermaid

  # Show usage information for a role
  strand workflow --role developer

  # Validate workflow completeness
  strand workflow --validate`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWorkflow(cmd.OutOrStdout(), cmd.ErrOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(workflowCmd)

	workflowCmd.Flags().StringVar(&workflowFormat, "format", "mermaid", "Output format: mermaid, json")
	workflowCmd.Flags().StringVar(&workflowTemplate, "template", "", "Show workflow for specific template")
	workflowCmd.Flags().StringVar(&workflowRole, "role", "", "Show usage information for specific role")
	workflowCmd.Flags().BoolVar(&workflowValidate, "validate", false, "Validate workflow completeness")
	workflowCmd.Flags().BoolVar(&workflowStats, "stats", false, "Include statistics in output")
}

func runWorkflow(w io.Writer, errW io.Writer) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	// Load roles and templates
	parser := workflow.NewParser()

	roles, err := parser.LoadRoles(paths.RolesDir)
	if err != nil {
		return fmt.Errorf("failed to load roles: %w", err)
	}

	templates, err := parser.LoadTemplates(paths.TemplatesDir)
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Build workflow graph
	graph := workflow.BuildGraph(roles, templates)

	// Handle --validate flag
	if workflowValidate {
		return runWorkflowValidation(graph, w, errW)
	}

	// Handle --role flag
	if workflowRole != "" {
		return runWorkflowRoleUsage(graph, workflowRole, w)
	}

	// Handle --template flag with diagram
	if workflowTemplate != "" {
		return runWorkflowTemplateDiagram(graph, workflowTemplate, w)
	}

	// Generate output based on format
	switch workflowFormat {
	case "mermaid":
		output := graph.GenerateMermaid()
		fmt.Fprintln(w, output)
		return nil

	case "json":
		data, err := graph.ToJSON(workflowStats)
		if err != nil {
			return fmt.Errorf("failed to generate JSON: %w", err)
		}
		fmt.Fprintln(w, string(data))
		return nil

	default:
		return fmt.Errorf("unsupported format: %s (use 'mermaid' or 'json')", workflowFormat)
	}
}

func runWorkflowValidation(graph *workflow.WorkflowGraph, w io.Writer, errW io.Writer) error {
	result := graph.Validate()

	// Print errors
	if len(result.Errors) > 0 {
		for _, issue := range result.Errors {
			fmt.Fprintf(errW, "✗ Error: %s\n", issue.Message)
			if issue.Location != "" {
				fmt.Fprintf(errW, "  Location: %s\n", issue.Location)
			}
		}
		fmt.Fprintln(errW)
	}

	// Print warnings
	if len(result.Warnings) > 0 {
		for _, issue := range result.Warnings {
			fmt.Fprintf(w, "⚠ Warning: %s\n", issue.Message)
			if issue.Location != "" {
				fmt.Fprintf(w, "  Location: %s\n", issue.Location)
			}
		}
		fmt.Fprintln(w)
	}

	// Print summary
	if !result.HasErrors() && !result.HasWarnings() {
		fmt.Fprintln(w, "✓ All roles referenced in templates have definitions")
		fmt.Fprintln(w, "✓ All task types referenced in role.creates exist as templates")
	}

	fmt.Fprintln(w, result.Summary())

	if result.HasErrors() {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func runWorkflowRoleUsage(graph *workflow.WorkflowGraph, roleName string, w io.Writer) error {
	usage := graph.GetRoleUsage(roleName)
	if usage == nil {
		return fmt.Errorf("role '%s' not found", roleName)
	}

	fmt.Fprintf(w, "Role: %s\n\n", usage.RoleName)

	// Assigned templates
	if len(usage.AssignedTemplates) > 0 {
		fmt.Fprintln(w, "Assigned to templates:")
		sort.Strings(usage.AssignedTemplates)
		for _, tmpl := range usage.AssignedTemplates {
			fmt.Fprintf(w, "  - %s (primary role)\n", tmpl)
		}
		fmt.Fprintln(w)
	}

	// Used in template TODOs
	if len(usage.UsedInTodos) > 0 {
		fmt.Fprintln(w, "Used in template TODOs:")
		templateNames := make([]string, 0, len(usage.UsedInTodos))
		for name := range usage.UsedInTodos {
			templateNames = append(templateNames, name)
		}
		sort.Strings(templateNames)
		for _, name := range templateNames {
			todos := usage.UsedInTodos[name]
			fmt.Fprintf(w, "  - %s (TODOs: %s)\n", name, formatTodoNumbers(todos))
		}
		fmt.Fprintln(w)
	}

	// Receives work via task types
	if len(usage.ReceivesVia) > 0 {
		fmt.Fprintln(w, "Receives work via task types:")
		sort.Strings(usage.ReceivesVia)
		for _, taskType := range usage.ReceivesVia {
			fmt.Fprintf(w, "  - %s\n", taskType)
		}
		fmt.Fprintln(w)
	}

	// Creates task types
	if len(usage.Creates) > 0 {
		fmt.Fprintln(w, "Creates task types:")
		for _, taskType := range usage.Creates {
			fmt.Fprintf(w, "  - %s\n", taskType)
		}
	} else {
		fmt.Fprintln(w, "Creates task types:")
		fmt.Fprintln(w, "  (none in workflow metadata)")
	}

	return nil
}

func runWorkflowTemplateDiagram(graph *workflow.WorkflowGraph, templateName string, w io.Writer) error {
	diagram, err := graph.GenerateMermaidForTemplate(templateName)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, diagram)
	return nil
}

func formatTodoNumbers(numbers []int) string {
	if len(numbers) == 0 {
		return ""
	}

	sort.Ints(numbers)
	parts := make([]string, len(numbers))
	for i, num := range numbers {
		parts[i] = fmt.Sprintf("%d", num)
	}
	return strings.Join(parts, ", ")
}
