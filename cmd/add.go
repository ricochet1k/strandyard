/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
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

	"github.com/ricochet1k/memmd/pkg/idgen"
	"github.com/ricochet1k/memmd/pkg/task"
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
		return runAdd(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addTitle, "title", "t", "", "task title")
	addCmd.Flags().StringVarP(&addRole, "role", "r", "", "role responsible for the task (defaults by type)")
	addCmd.Flags().StringVarP(&addParent, "parent", "p", "", "parent task ID (creates task under that directory)")
	addCmd.Flags().StringVar(&addPriority, "priority", "medium", "priority: high, medium, or low")
	addCmd.Flags().StringSliceVar(&addBlockers, "blocker", nil, "blocker task ID(s); can be repeated or comma-separated")
	addCmd.Flags().BoolVar(&addNoRepair, "no-repair", false, "skip repair and master list updates")
}

var (
	addTitle    string
	addRole     string
	addPriority string
	addParent   string
	addBlockers []string
	addNoRepair bool
)

func runAdd(cmd *cobra.Command, args []string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	tmplName := strings.TrimSpace(args[0])
	if tmplName == "" {
		return fmt.Errorf("type is required")
	}

	templates, err := listTemplateNames(paths.TemplatesDir)
	if err != nil {
		return err
	}
	if len(templates) == 0 {
		return fmt.Errorf("no templates found in %s", paths.TemplatesDir)
	}

	if !containsTemplate(templates, tmplName) {
		return fmt.Errorf("unknown type %q (available: %s)", tmplName, strings.Join(templates, ", "))
	}

	title := strings.TrimSpace(addTitle)
	if title == "" && len(args) > 1 {
		title = strings.TrimSpace(strings.Join(args[1:], " "))
	}
	if title == "" {
		return fmt.Errorf("title is required (use --title or provide it as an argument)")
	}

	templatePath := filepath.Join(paths.TemplatesDir, tmplName+".md")
	templateMeta, templateBody, err := loadTemplate(templatePath)
	if err != nil {
		return err
	}

	role := strings.TrimSpace(addRole)
	if !cmd.Flags().Changed("role") {
		role = strings.TrimSpace(templateMeta.Role)
	}
	if role == "" {
		return fmt.Errorf("role is required (use --role or set role in template frontmatter)")
	}
	if err := validateRole(paths.RolesDir, role); err != nil {
		return err
	}

	priority := task.NormalizePriority(addPriority)
	if !cmd.Flags().Changed("priority") && templateMeta.Priority != "" {
		priority = task.NormalizePriority(templateMeta.Priority)
	}
	if !task.IsValidPriority(priority) {
		return fmt.Errorf("invalid priority: %s", priority)
	}

	parent := strings.TrimSpace(addParent)
	parentDir := paths.TasksDir
	var tasks map[string]*task.Task
	var parser *task.Parser
	if parent != "" {
		parser = task.NewParser()
		loaded, err := parser.LoadTasks(paths.TasksDir)
		if err != nil {
			return err
		}
		parentTask, ok := loaded[parent]
		if !ok {
			return fmt.Errorf("parent task %s does not exist", parent)
		}
		parentDir = parentTask.Dir
		tasks = loaded
	}

	id, err := idgen.GenerateID(taskPrefixForTemplate(tmplName), title)
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

	blockers := normalizeTaskIDs(addBlockers)
	now := time.Now().UTC()
	meta := task.Metadata{
		Type:          tmplName,
		Role:          role,
		Priority:      priority,
		Parent:        parent,
		Blockers:      blockers,
		Blocks:        []string{},
		DateCreated:   now,
		DateEdited:    now,
		OwnerApproval: false,
		Completed:     false,
	}

	stdinBody, err := readStdin()
	if err != nil {
		return err
	}

	body := renderTemplateBody(templateBody, map[string]string{
		"Title":               title,
		"SuggestedSubtaskDir": fmt.Sprintf("%s-subtask", id),
		"Body":                stdinBody,
	})
	if stdinBody != "" && !strings.Contains(templateBody, "{{ .Body }}") {
		if strings.TrimSpace(body) != "" {
			body += "\n\n"
		}
		body += stdinBody
	}
	taskFile := filepath.Join(taskDir, id+".md")
	if err := writeTaskFile(taskFile, meta, body); err != nil {
		return err
	}

	fmt.Printf("✓ Task created: %s\n", filepath.ToSlash(taskFile))

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

	if !addNoRepair {
		if err := runRepair(paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
			return err
		}
	}

	return nil
}

func listTemplateNames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".md") {
			names = append(names, strings.TrimSuffix(name, ".md"))
		}
	}
	sort.Strings(names)
	return names, nil
}

func containsTemplate(names []string, name string) bool {
	for _, item := range names {
		if item == name {
			return true
		}
	}
	return false
}

func taskPrefixForTemplate(name string) string {
	switch name {
	case "issue":
		return "I"
	default:
		return "T"
	}
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

type templateDefaults struct {
	Role     string `yaml:"role"`
	Priority string `yaml:"priority"`
}

func loadTemplate(path string) (templateDefaults, string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return templateDefaults{}, "", fmt.Errorf("failed to read template %s: %w", path, err)
	}
	text := string(content)
	if !strings.HasPrefix(text, "---") {
		return templateDefaults{}, "", fmt.Errorf("template %s missing frontmatter", path)
	}

	parts := strings.SplitN(text, "---", 3)
	if len(parts) < 3 {
		return templateDefaults{}, "", fmt.Errorf("template %s frontmatter delimiter missing", path)
	}

	var meta templateDefaults
	if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
		return templateDefaults{}, "", fmt.Errorf("failed to parse template frontmatter %s: %w", path, err)
	}

	body := strings.TrimLeft(parts[2], "\r\n")
	return meta, body, nil
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

func validateRole(rolesDir, role string) error {
	rolePath := filepath.Join(rolesDir, role+".md")
	if _, err := os.Stat(rolePath); err != nil {
		return fmt.Errorf("role file %s does not exist", rolePath)
	}
	return nil
}
