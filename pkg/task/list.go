package task

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

// ListOptions defines filters and output parameters for listing tasks.
type ListOptions struct {
	Scope          string
	Parent         string
	Path           string
	Role           string
	Priority       string
	Completed      *bool
	Blocked        *bool
	Blocks         *bool
	OwnerApproval  *bool
	Label          string
	Sort           string
	Order          string
	Format         string
	Columns        []string
	Group          string
	MdTable        bool
	UseMasterLists bool
}

// ListTasks loads tasks and returns a filtered, deterministically sorted list.
func ListTasks(tasksRoot string, opts ListOptions) ([]*Task, error) {
	parser := NewParser()
	tasks, err := parser.LoadTasks(tasksRoot)
	if err != nil {
		return nil, err
	}

	items, err := filterTasks(tasksRoot, tasks, opts)
	if err != nil {
		return nil, err
	}

	sortTasks(items, opts)
	return items, nil
}

// FormatList formats tasks according to the requested output format.
func FormatList(tasks []*Task, opts ListOptions) (string, error) {
	switch opts.Format {
	case "table", "":
		return formatTable(tasks, opts)
	case "md":
		return formatMarkdown(tasks, opts)
	case "json":
		return formatJSON(tasks, opts)
	default:
		return "", fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

func filterTasks(tasksRoot string, tasks map[string]*Task, opts ListOptions) ([]*Task, error) {
	items := make([]*Task, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, t)
	}

	if opts.Parent != "" {
		if _, ok := tasks[opts.Parent]; !ok {
			return nil, fmt.Errorf("parent task not found: %s", opts.Parent)
		}
	}

	basePath := ""
	if strings.TrimSpace(opts.Path) != "" {
		path := filepath.Clean(opts.Path)
		if strings.HasPrefix(path, "tasks"+string(filepath.Separator)) || path == "tasks" {
			basePath = filepath.Join(filepath.Dir(tasksRoot), path)
		} else {
			basePath = filepath.Join(tasksRoot, path)
		}
	}

	filtered := make([]*Task, 0, len(items))
	for _, t := range items {
		if !matchesScope(t, opts) {
			continue
		}
		if opts.Parent != "" && t.Meta.Parent != opts.Parent {
			continue
		}
		if basePath != "" && !isWithinPath(t.Dir, basePath) {
			continue
		}
		if opts.Role != "" && strings.ToLower(t.GetEffectiveRole()) != strings.ToLower(opts.Role) {
			continue
		}
		if opts.Priority != "" && NormalizePriority(t.Meta.Priority) != NormalizePriority(opts.Priority) {
			continue
		}
		if opts.Completed != nil && t.Meta.Completed != *opts.Completed {
			continue
		}
		if opts.Blocked != nil {
			hasBlockers := len(t.Meta.Blockers) > 0
			if hasBlockers != *opts.Blocked {
				continue
			}
		}
		if opts.Blocks != nil {
			hasBlocks := len(t.Meta.Blocks) > 0
			if hasBlocks != *opts.Blocks {
				continue
			}
		}
		if opts.OwnerApproval != nil && t.Meta.OwnerApproval != *opts.OwnerApproval {
			continue
		}
		filtered = append(filtered, t)
	}

	return filtered, nil
}

func matchesScope(t *Task, opts ListOptions) bool {
	switch opts.Scope {
	case "", "all":
		return true
	case "root":
		return strings.TrimSpace(t.Meta.Parent) == ""
	case "free":
		return len(t.Meta.Blockers) == 0
	default:
		return true
	}
}

func isWithinPath(dir, base string) bool {
	rel, err := filepath.Rel(base, dir)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".."
}

func sortTasks(items []*Task, opts ListOptions) {
	sortKey := strings.ToLower(strings.TrimSpace(opts.Sort))
	order := strings.ToLower(strings.TrimSpace(opts.Order))
	desc := order == "desc"

	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		less := compareTasks(a, b, sortKey)
		if desc {
			return !less
		}
		return less
	})
}

func compareTasks(a, b *Task, sortKey string) bool {
	switch sortKey {
	case "id":
		return a.ID < b.ID
	case "priority":
		if PriorityRank(a.Meta.Priority) != PriorityRank(b.Meta.Priority) {
			return PriorityRank(a.Meta.Priority) < PriorityRank(b.Meta.Priority)
		}
		return a.ID < b.ID
	case "created":
		return compareTime(a.Meta.DateCreated, b.Meta.DateCreated, a.ID, b.ID)
	case "edited":
		return compareTime(a.Meta.DateEdited, b.Meta.DateEdited, a.ID, b.ID)
	case "role":
		roleA := strings.ToLower(a.GetEffectiveRole())
		roleB := strings.ToLower(b.GetEffectiveRole())
		if roleA != roleB {
			return roleA < roleB
		}
		return a.ID < b.ID
	default:
		// Default sort: priority, completed, ID.
		if PriorityRank(a.Meta.Priority) != PriorityRank(b.Meta.Priority) {
			return PriorityRank(a.Meta.Priority) < PriorityRank(b.Meta.Priority)
		}
		if a.Meta.Completed != b.Meta.Completed {
			return !a.Meta.Completed && b.Meta.Completed
		}
		return a.ID < b.ID
	}
}

func compareTime(a, b time.Time, ida, idb string) bool {
	if !a.Equal(b) {
		return a.Before(b)
	}
	return ida < idb
}

type listRow struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Role        string   `json:"role"`
	Priority    string   `json:"priority"`
	Parent      string   `json:"parent"`
	Completed   bool     `json:"completed"`
	Blockers    []string `json:"blockers"`
	Blocks      []string `json:"blocks"`
	Path        string   `json:"path"`
	DateCreated string   `json:"date_created"`
	DateEdited  string   `json:"date_edited"`
}

func toListRows(tasks []*Task) []listRow {
	rows := make([]listRow, 0, len(tasks))
	for _, t := range tasks {
		rows = append(rows, listRow{
			ID:          t.ID,
			Title:       t.Title(),
			Role:        t.GetEffectiveRole(),
			Priority:    NormalizePriority(t.Meta.Priority),
			Parent:      t.Meta.Parent,
			Completed:   t.Meta.Completed,
			Blockers:    append([]string{}, t.Meta.Blockers...),
			Blocks:      append([]string{}, t.Meta.Blocks...),
			Path:        filepath.ToSlash(t.FilePath),
			DateCreated: t.Meta.DateCreated.Format(time.RFC3339),
			DateEdited:  t.Meta.DateEdited.Format(time.RFC3339),
		})
	}
	return rows
}

func formatTable(tasks []*Task, opts ListOptions) (string, error) {
	columns := normalizeColumns(opts.Columns, []string{"id", "title", "priority", "role", "completed", "blockers"})
	rows := toListRows(tasks)

	builder := &strings.Builder{}
	writer := tabwriter.NewWriter(builder, 0, 4, 2, ' ', 0)
	fmt.Fprintln(writer, strings.ToUpper(strings.Join(columns, "\t")))
	for _, row := range rows {
		fmt.Fprintln(writer, strings.Join(rowValues(row, columns, true), "\t"))
	}
	_ = writer.Flush()
	return strings.TrimRight(builder.String(), "\n"), nil
}

func formatMarkdown(tasks []*Task, opts ListOptions) (string, error) {
	rows := toListRows(tasks)
	if opts.Group != "" && opts.Group != "none" {
		grouped := groupRows(rows, opts.Group)
		return formatMarkdownGrouped(grouped, opts)
	}

	if opts.MdTable {
		return formatMarkdownTable(rows, opts)
	}
	return formatMarkdownList(rows), nil
}

func formatMarkdownGrouped(grouped map[string][]listRow, opts ListOptions) (string, error) {
	keys := sortedGroupKeys(grouped)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		title := groupTitle(opts.Group, key)
		part := &strings.Builder{}
		fmt.Fprintf(part, "## %s\n\n", title)
		if opts.MdTable {
			table, _ := formatMarkdownTable(grouped[key], opts)
			fmt.Fprintln(part, table)
		} else {
			list := formatMarkdownList(grouped[key])
			fmt.Fprintln(part, list)
		}
		parts = append(parts, strings.TrimRight(part.String(), "\n"))
	}
	return strings.Join(parts, "\n\n"), nil
}

func formatMarkdownTable(rows []listRow, opts ListOptions) (string, error) {
	columns := normalizeColumns(opts.Columns, []string{"id", "title", "priority", "role", "completed", "blockers"})
	builder := &strings.Builder{}
	fmt.Fprintln(builder, "| "+strings.Join(columns, " | ")+" |")
	separators := make([]string, 0, len(columns))
	for range columns {
		separators = append(separators, "---")
	}
	fmt.Fprintln(builder, "| "+strings.Join(separators, " | ")+" |")
	for _, row := range rows {
		fmt.Fprintln(builder, "| "+strings.Join(rowValues(row, columns, false), " | ")+" |")
	}
	return strings.TrimRight(builder.String(), "\n"), nil
}

func formatMarkdownList(rows []listRow) string {
	if len(rows) == 0 {
		return ""
	}
	builder := &strings.Builder{}
	for _, row := range rows {
		fmt.Fprintf(builder, "- %s — %s (priority: %s, role: %s, completed: %t)\n",
			row.ID, row.Title, row.Priority, row.Role, row.Completed)
	}
	return strings.TrimRight(builder.String(), "\n")
}

func formatJSON(tasks []*Task, opts ListOptions) (string, error) {
	rows := toListRows(tasks)
	if opts.Group != "" && opts.Group != "none" {
		grouped := groupRows(rows, opts.Group)
		ordered := map[string][]listRow{}
		for _, key := range sortedGroupKeys(grouped) {
			ordered[key] = grouped[key]
		}
		b, err := json.MarshalIndent(ordered, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	b, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func normalizeColumns(cols, defaults []string) []string {
	if len(cols) == 0 {
		return defaults
	}
	out := make([]string, 0, len(cols))
	for _, col := range cols {
		col = strings.ToLower(strings.TrimSpace(col))
		if col == "" {
			continue
		}
		out = append(out, col)
	}
	if len(out) == 0 {
		return defaults
	}
	return out
}

func rowValues(row listRow, columns []string, numericForCounts bool) []string {
	values := make([]string, 0, len(columns))
	for _, col := range columns {
		switch col {
		case "id":
			values = append(values, row.ID)
		case "title":
			values = append(values, row.Title)
		case "priority":
			values = append(values, row.Priority)
		case "role":
			values = append(values, row.Role)
		case "parent":
			values = append(values, row.Parent)
		case "completed":
			values = append(values, fmt.Sprintf("%t", row.Completed))
		case "blockers":
			values = append(values, formatListValue(row.Blockers, numericForCounts))
		case "blocks":
			values = append(values, formatListValue(row.Blocks, numericForCounts))
		case "path":
			values = append(values, row.Path)
		case "date_created":
			values = append(values, row.DateCreated)
		case "date_edited":
			values = append(values, row.DateEdited)
		default:
			values = append(values, "")
		}
	}
	return values
}

func formatListValue(values []string, numeric bool) string {
	if numeric {
		return fmt.Sprintf("%d", len(values))
	}
	if len(values) == 0 {
		return ""
	}
	return strings.Join(values, ",")
}

func groupRows(rows []listRow, group string) map[string][]listRow {
	grouped := make(map[string][]listRow)
	for _, row := range rows {
		key := ""
		switch group {
		case "priority":
			key = NormalizePriority(row.Priority)
		case "parent":
			key = row.Parent
			if strings.TrimSpace(key) == "" {
				key = "(root)"
			}
		case "role":
			key = strings.ToLower(strings.TrimSpace(row.Role))
			if key == "" {
				key = "(unassigned)"
			}
		default:
			key = "(other)"
		}
		grouped[key] = append(grouped[key], row)
	}
	return grouped
}

func sortedGroupKeys(grouped map[string][]listRow) []string {
	keys := make([]string, 0, len(grouped))
	for key := range grouped {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func groupTitle(group string, key string) string {
	switch group {
	case "priority":
		switch NormalizePriority(key) {
		case PriorityHigh:
			return "High"
		case PriorityMedium:
			return "Medium"
		case PriorityLow:
			return "Low"
		default:
			return key
		}
	case "parent":
		if key == "(root)" {
			return "Root Tasks"
		}
		return key
	case "role":
		if key == "(unassigned)" {
			return "Unassigned"
		}
		return strings.Title(key)
	default:
		return key
	}
}
