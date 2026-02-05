package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ricochet1k/strandyard/pkg/activity"
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
	addCmd.Flags().StringSliceVar(&addBlocks, "blocks", nil, "task ID(s) this task blocks; can be repeated or comma-separated")
	addCmd.Flags().StringSliceVar(&addEvery, "every", nil, "recurrence rule (e.g., \"10 days\", \"50 commits from HEAD\")")
}

var (
	addTitle    string
	addRole     string
	addPriority string
	addParent   string
	addBlockers []string
	addBlocks   []string
	addEvery    []string
)

type addOptions struct {
	ProjectName       string
	TemplateName      string
	Title             string
	Role              string
	Priority          string
	Parent            string
	Blockers          []string
	Blocks            []string
	Every             []string
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
		Blocks:            addBlocks,
		Every:             addEvery,
		RoleSpecified:     cmd.Flags().Changed("role"),
		PrioritySpecified: cmd.Flags().Changed("priority"),
		Body:              body,
	}, nil
}

// validateEvery validates --every flag values and returns resolved recurrence rules
func validateEvery(every []string, repoPath string, tasks map[string]*task.Task) ([]string, error) {
	if len(every) == 0 {
		return nil, nil // --every is optional
	}

	resolvedEvery := make([]string, 0, len(every))
	for _, value := range every {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		// Parse format: <amount> <metric> [from <anchor>]
		parts := strings.Fields(value)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "strand: error: invalid --every value: expected format \"<amount> <metric> [from <anchor>]\"\n")
			fmt.Fprintf(os.Stderr, "hint: --every \"10 days\"\n")
			return nil, fmt.Errorf("invalid --every format")
		}

		// Validate amount (must be integer)
		amount := parts[0]
		if _, err := strconv.Atoi(amount); err != nil {
			fmt.Fprintf(os.Stderr, "strand: error: invalid --every value: amount must be an integer\n")
			fmt.Fprintf(os.Stderr, "hint: --every \"10 days\"\n")
			return nil, fmt.Errorf("invalid --every amount")
		}

		// Validate metric
		metric := parts[1]
		validMetrics := map[string]bool{
			"days":            true,
			"weeks":           true,
			"months":          true,
			"commits":         true,
			"lines_changed":   true,
			"tasks_completed": true,
		}
		if !validMetrics[metric] {
			fmt.Fprintf(os.Stderr, "strand: error: invalid --every value: unsupported metric \"%s\"\n", metric)
			fmt.Fprintf(os.Stderr, "hint: --every \"10 days\"\n")
			return nil, fmt.Errorf("invalid --every metric")
		}

		// Validate anchor if present
		if len(parts) >= 4 && parts[2] == "from" {
			anchor := strings.Join(parts[3:], " ")
			resolved, err := task.ValidateAnchor(metric, anchor, repoPath, tasks)
			if err != nil {
				fmt.Fprintf(os.Stderr, "strand: error: invalid --every value: %v\n", err)
				if metric == "commits" || metric == "lines_changed" {
					fmt.Fprintf(os.Stderr, "hint: --every \"50 commits from HEAD\"\n")
				} else if metric == "days" || metric == "weeks" || metric == "months" {
					fmt.Fprintf(os.Stderr, "hint: --every \"10 days from Jan 28 2026 09:00 UTC\"\n")
				}
				return nil, err
			}
			// Update rule with resolved anchor
			value = fmt.Sprintf("%s %s from %s", amount, metric, resolved)
		}
		resolvedEvery = append(resolvedEvery, value)
	}

	return resolvedEvery, nil // Validation passed
}

func runAdd(w io.Writer, opts addOptions) error {
	paths, err := resolveProjectPaths(opts.ProjectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	if err := db.LoadAllIfEmpty(); err != nil {
		return err
	}

	resolvedEvery, err := validateEvery(opts.Every, paths.BaseDir, db.GetAll())
	if err != nil {
		os.Exit(2)
	}
	opts.Every = resolvedEvery

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
	if parent != "" {
		resolvedParent, err := db.ResolveID(parent)
		if err != nil {
			return fmt.Errorf("parent task %s does not exist: %w", parent, err)
		}
		parent = resolvedParent
		parentTask, err := db.Get(parent)
		if err != nil {
			return fmt.Errorf("parent task %s does not exist: %w", parent, err)
		}
		parentDir = parentTask.Dir
	}

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

	blockers, err := db.ResolveIDs(normalizeTaskIDs(opts.Blockers))
	if err != nil {
		return err
	}
	blocks, err := db.ResolveIDs(normalizeTaskIDs(opts.Blocks))
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	meta := task.Metadata{
		Type:          tmplName,
		Role:          roleName,
		Priority:      priority,
		Parent:        parent,
		Blockers:      []string{},
		Blocks:        []string{},
		DateCreated:   now,
		DateEdited:    now,
		OwnerApproval: false,
		Completed:     false,
		Every:         opts.Every,
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

	// Load the new task and set up blocker/blocks relationships via TaskDB
	if len(blockers) > 0 || len(blocks) > 0 {
		if _, err := db.Load(id); err != nil {
			return fmt.Errorf("failed to load new task: %w", err)
		}
		for _, blockerID := range blockers {
			if err := db.AddBlocker(id, blockerID); err != nil {
				return fmt.Errorf("failed to add blocker %s: %w", blockerID, err)
			}
		}
		for _, blockedID := range blocks {
			if err := db.AddBlocked(id, blockedID); err != nil {
				return fmt.Errorf("failed to add blocked %s: %w", blockedID, err)
			}
		}
		if _, err := db.SaveDirty(); err != nil {
			return fmt.Errorf("failed to write blocker updates: %w", err)
		}
	}

	if len(opts.Every) > 0 {
		activeLog, err := activity.Open(paths.BaseDir)
		if err == nil {
			defer activeLog.Close()
			for _, rule := range opts.Every {
				parts := strings.Fields(rule)
				if len(parts) >= 2 {
					metric := parts[1]
					anchor := ""
					if len(parts) >= 4 && parts[2] == "from" {
						anchor = strings.Join(parts[3:], " ")
					}

					if metric == "commits" || metric == "lines_changed" {
						if anchor == "HEAD" || anchor == "" {
							if resolved, err := task.ResolveGitHash(paths.BaseDir, "HEAD"); err == nil {
								_ = activeLog.WriteRecurrenceAnchorResolution(id, "HEAD", resolved)
							}
						}
					} else {
						if anchor == "now" || anchor == "" {
							_ = activeLog.WriteRecurrenceAnchorResolution(id, "now", now.Format("Jan 2 2006 15:04 MST"))
						}
					}
				}
			}
		}
	}

	if parent != "" {
		// Load the newly created task into the DB
		if _, err := db.Load(id); err != nil {
			return fmt.Errorf("failed to load new task: %w", err)
		}
		if _, err := db.UpdateParentTodos(parent); err != nil {
			return fmt.Errorf("failed to update parent task TODO entries: %w", err)
		}
		if _, err := db.SaveDirty(); err != nil {
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
