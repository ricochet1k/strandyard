package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ricochet1k/strandyard/pkg/idgen"
	rPkg "github.com/ricochet1k/strandyard/pkg/role"
	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/ricochet1k/strandyard/pkg/template"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// addCmd groups task creation commands.
var addCmd = &cobra.Command{
	Use:   "add <type> [title]",
	Short: "Create tasks from templates",
	Long:  "Create a task using a template in templates/. Types correspond to template filenames (without .md). Templates define default roles and priorities. Provide a detailed body on stdin (pipe or heredoc); it will be inserted where the template uses {{ .Body }} or appended to the end.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := readStdin()
		if err != nil {
			return err
		}
		opts, err := addOptionsFromFlags(cmd, args, body)
		if err != nil {
			return err
		}
		return runAdd(cmd.OutOrStdout(), opts)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addTitle, "title", "t", "", "task title")
	addCmd.Flags().StringVarP(&addRole, "role", "r", "", "role responsible for the task (defaults by type)")
	addCmd.Flags().StringVarP(&addParent, "parent", "p", "", "parent task ID (creates task under that directory)")
	addCmd.Flags().StringVar(&addPriority, "priority", "medium", "priority: high, medium, or low")
	addCmd.Flags().StringSliceVar(&addBlockers, "blocker", nil, "blocker task ID(s); can be repeated or comma-separated")
}

var (
	addTitle    string
	addRole     string
	addPriority string
	addParent   string
	addBlockers []string
)

type addOptions struct {
	ProjectName       string
	TemplateName      string
	Title             string
	Role              string
	Priority          string
	Parent            string
	Blockers          []string
	RoleSpecified     bool
	PrioritySpecified bool
	Body              string
}

func addOptionsFromFlags(cmd *cobra.Command, args []string, body string) (addOptions, error) {
	if len(args) == 0 {
		return addOptions{}, fmt.Errorf("type is required")
	}
	title := strings.TrimSpace(addTitle)
	if title == "" && len(args) > 1 {
		title = strings.TrimSpace(strings.Join(args[1:], " "))
	}
	return addOptions{
		ProjectName:       projectName,
		TemplateName:      strings.TrimSpace(args[0]),
		Title:             title,
		Role:              strings.TrimSpace(addRole),
		Priority:          strings.TrimSpace(addPriority),
		Parent:            strings.TrimSpace(addParent),
		Blockers:          addBlockers,
		RoleSpecified:     cmd.Flags().Changed("role"),
		PrioritySpecified: cmd.Flags().Changed("priority"),
		Body:              body,
	}, nil
}

func runAdd(w io.Writer, opts addOptions) error {
	paths, err := resolveProjectPaths(opts.ProjectName)
	if err != nil {
		return err
	}

	tmplName := strings.TrimSpace(opts.TemplateName)
	if tmplName == "" {
		return fmt.Errorf("type is required")
	}

	templates, err := template.LoadTemplates(paths.TemplatesDir)
	if err != nil {
		return err
	}

	tmpl, ok := templates[tmplName]
	if !ok {
		fmt.Fprintln(w, "Unknown type. Available templates:")
		var names []string
		for name := range templates {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			desc := templates[name].Meta.Description
			if desc == "" {
				desc = "(no description found)"
			}
			fmt.Fprintf(w, "  %-15s %s\n", name, desc)
		}
		return fmt.Errorf("unknown type %q", tmplName)
	}

	title := strings.TrimSpace(opts.Title)
	if title == "" {
		return fmt.Errorf("title is required (use --title or provide it as an argument)")
	}

	roleName := strings.TrimSpace(opts.Role)
	if !opts.RoleSpecified {
		roleName = strings.TrimSpace(tmpl.Meta.Role)
	}
	if roleName == "" {
		return fmt.Errorf("role is required (use --role or set role in template frontmatter)")
	}

	roles, err := rPkg.LoadRoles(paths.RolesDir)
	if err != nil {
		return err
	}

	if _, ok := roles[roleName]; !ok {
		fmt.Fprintln(w, "Invalid role. Available roles:")
		var names []string
		for name := range roles {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			desc := roles[name].Meta.Description
			if desc == "" {
				desc = "(no description found)"
			}
			fmt.Fprintf(w, "  %-15s %s\n", name, desc)
		}
		return fmt.Errorf("invalid role %q", roleName)
	}

	priority := task.NormalizePriority(opts.Priority)
	if !opts.PrioritySpecified {
		if pStr, ok := tmpl.Meta.Priority.(string); ok && pStr != "" {
			priority = task.NormalizePriority(pStr)
		}
	}
	if !task.IsValidPriority(priority) {
		return fmt.Errorf("invalid priority: %s", priority)
	}

	parent := strings.TrimSpace(opts.Parent)
	parentDir := paths.TasksDir
	var tasks map[string]*task.Task
	var parser *task.Parser
	if parent != "" || len(opts.Blockers) > 0 {
		parser = task.NewParser()
		loaded, err := parser.LoadTasks(paths.TasksDir)
		if err != nil {
			return err
		}
		tasks = loaded
	}
	if parent != "" {
		resolvedParent, err := task.ResolveTaskID(tasks, parent)
		if err != nil {
			return fmt.Errorf("parent task %s does not exist: %w", parent, err)
		}
		parent = resolvedParent
		parentTask, ok := tasks[parent]
		if !ok {
			return fmt.Errorf("parent task %s does not exist", parent)
		}
		parentDir = parentTask.Dir
	}

	// We don't have id_prefix in task.Metadata yet, but it was in templateDefaults.
	// For now, default to "T" or "E" if the template name suggests it.
	prefix := "T"
	if strings.Contains(strings.ToLower(tmplName), "epic") {
		prefix = "E"
	}

	id, err := idgen.GenerateID(prefix, title)
	if err != nil {
		return err
	}

	taskDir := filepath.Join(parentDir, id)
	if _, err := os.Stat(taskDir); err == nil {
		return fmt.Errorf("task directory already exists: %s", taskDir)
	}
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		return fmt.Errorf("failed to create task directory: %w", err)
	}

	blockers, err := resolveTaskIDs(tasks, normalizeTaskIDs(opts.Blockers))
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	meta := task.Metadata{
		Type:          tmplName,
		Role:          roleName,
		Priority:      priority,
		Parent:        parent,
		Blockers:      blockers,
		Blocks:        []string{},
		DateCreated:   now,
		DateEdited:    now,
		OwnerApproval: false,
		Completed:     false,
	}

	body := renderTemplateBody(tmpl.BodyContent, map[string]string{
		"Title":               title,
		"SuggestedSubtaskDir": fmt.Sprintf("%s-subtask", id),
		"Body":                opts.Body,
	})
	if opts.Body != "" && !strings.Contains(tmpl.BodyContent, "{{ .Body }}") {
		if strings.TrimSpace(body) != "" {
			body += "\n\n"
		}
		body += opts.Body
	}
	taskFile := filepath.Join(taskDir, id+".md")
	if err := writeTaskFile(taskFile, meta, body); err != nil {
		return err
	}

	fmt.Fprintf(w, "✓ Task created: %s\n", id)

	if parent != "" {
		newTask, err := parser.ParseFile(taskFile)
		if err != nil {
			return fmt.Errorf("failed to parse new task: %w", err)
		}
		tasks[newTask.ID] = newTask
		if _, err := task.UpdateParentTodoEntries(tasks, parent); err != nil {
			return fmt.Errorf("failed to update parent task TODO entries: %w", err)
		}
		if _, err := task.WriteDirtyTasks(tasks); err != nil {
			return fmt.Errorf("failed to write parent task updates: %w", err)
		}
	}

	// TODO: This should not be necessary
	if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
		return err
	}

	return nil
}

func normalizeTaskIDs(items []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, item := range items {
		parts := strings.Split(item, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			if _, ok := seen[trimmed]; ok {
				continue
			}
			seen[trimmed] = struct{}{}
			out = append(out, trimmed)
		}
	}
	sort.Strings(out)
	return out
}

func resolveTaskIDs(tasks map[string]*task.Task, inputs []string) ([]string, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	if tasks == nil {
		return nil, fmt.Errorf("failed to resolve task IDs: no tasks loaded")
	}
	seen := map[string]struct{}{}
	resolved := make([]string, 0, len(inputs))
	for _, input := range inputs {
		id, err := task.ResolveTaskID(tasks, input)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		resolved = append(resolved, id)
	}
	sort.Strings(resolved)
	return resolved, nil
}

func renderTemplateBody(body string, data map[string]string) string {
	out := body
	for key, value := range data {
		out = strings.ReplaceAll(out, "{{ ."+key+" }}", value)
	}
	return out
}

func writeTaskFile(path string, meta task.Metadata, body string) error {
	frontmatterBytes, err := yaml.Marshal(&meta)
	if err != nil {
		return fmt.Errorf("failed to marshal frontmatter: %w", err)
	}
	frontmatterBytes = bytes.TrimSpace(frontmatterBytes)

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(frontmatterBytes)
	sb.WriteString("\n---\n\n")
	sb.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		sb.WriteString("\n")
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func readStdin() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat stdin: %w", err)
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		return "", nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}
	return strings.TrimRight(string(data), "\r\n"), nil
}
