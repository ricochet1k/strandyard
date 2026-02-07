package task

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
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
	Status         string
	Label          string
	Sort           string
	Order          string
	Format         string
	Columns        []string
	Group          string
	MdTable        bool
	UseMasterLists bool
	Color          bool
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
		resolved, err := ResolveTaskID(tasks, opts.Parent)
		if err != nil {
			return nil, fmt.Errorf("parent task not found: %w", err)
		}
		opts.Parent = resolved
	}

	pathFilter := strings.TrimSpace(opts.Path)
	var pathRoot string
	if pathFilter != "" {
		pathRoot = filepath.Clean(filepath.Join(tasksRoot, pathFilter))
	}

	filtered := make([]*Task, 0, len(items))
	for _, t := range items {
		if pathRoot != "" && !isUnderPath(t.Dir, pathRoot) {
			continue
		}
		if !matchesScope(t, opts) {
			continue
		}
		if opts.Parent != "" && t.Meta.Parent != opts.Parent {
			continue
		}
		if opts.Role != "" && !strings.EqualFold(t.GetEffectiveRole(), opts.Role) {
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
		if opts.Status != "" && !strings.EqualFold(t.Meta.Status, opts.Status) {
			continue
		}
		filtered = append(filtered, t)
	}

	return filtered, nil
}

func isUnderPath(path, root string) bool {
	path = filepath.Clean(path)
	root = filepath.Clean(root)
	if path == root {
		return true
	}
	if !strings.HasSuffix(root, string(filepath.Separator)) {
		root += string(filepath.Separator)
	}
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}
	return strings.HasPrefix(path, root)
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
	Status      string   `json:"status"`
	Blockers    []string `json:"blockers"`
	Blocks      []string `json:"blocks"`
	Path        string   `json:"path"`
	DateCreated string   `json:"date_created"`
	DateEdited  string   `json:"date_edited"`
}

func toListRows(tasks []*Task) []listRow {
	rows := make([]listRow, 0, len(tasks))
	for _, t := range tasks {
		shortParent := ShortID(t.Meta.Parent)
		shortBlockers := shortenTaskIDs(t.Meta.Blockers)
		shortBlocks := shortenTaskIDs(t.Meta.Blocks)
		rows = append(rows, listRow{
			ID:          ShortID(t.ID),
			Title:       t.Title(),
			Role:        t.GetEffectiveRole(),
			Priority:    NormalizePriority(t.Meta.Priority),
			Parent:      shortParent,
			Completed:   t.Meta.Completed,
			Status:      t.Meta.Status,
			Blockers:    shortBlockers,
			Blocks:      shortBlocks,
			Path:        filepath.ToSlash(t.FilePath),
			DateCreated: t.Meta.DateCreated.Format(time.RFC3339),
			DateEdited:  t.Meta.DateEdited.Format(time.RFC3339),
		})
	}
	return rows
}

func shortenTaskIDs(ids []string) []string {
	if len(ids) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, ShortID(id))
	}
	return out
}

func formatTable(tasks []*Task, opts ListOptions) (string, error) {
	rows := toListRows(tasks)
	columns := defaultColumnsForRows(opts, rows, []string{"id", "title", "priority", "role", "status", "completed", "blockers"}, true)

	if len(rows) == 0 {
		return "", nil
	}

	// Calculate widths based on non-colored values
	widths := make(map[string]int)
	for _, col := range columns {
		widths[col] = len(col)
	}
	for _, row := range rows {
		for _, col := range columns {
			val := columnValue(row, col, true)
			if len(val) > widths[col] {
				widths[col] = len(val)
			}
		}
	}
	if w, ok := widths["title"]; ok {
		widths["title"] = min(w, 50)
	}

	builder := &strings.Builder{}
	// Print header
	for i, col := range columns {
		header := strings.ToUpper(col)
		builder.WriteString(header)
		if i < len(columns)-1 {
			padding := widths[col] - len(header) + 2 // 2 spaces minimal padding
			builder.WriteString(strings.Repeat(" ", padding))
		}
	}
	builder.WriteString("\n")

	// Print rows
	for _, row := range rows {
		for i, col := range columns {
			val := columnValue(row, col, true)
			width := widths[col]
			vallen := len(val)
			if vallen > width {
				// len("…") == 3 even though it's only 1 cell wide
				val = val[:width-1] + "…"
				vallen = width
			}
			colored := colorizeValue(row, col, val, opts)
			builder.WriteString(colored)
			if i < len(columns)-1 {
				padding := width - vallen + 2
				builder.WriteString(strings.Repeat(" ", padding))
			}
		}
		builder.WriteString("\n")
	}

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
	return formatMarkdownList(rows, opts), nil
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
			list := formatMarkdownList(grouped[key], opts)
			fmt.Fprintln(part, list)
		}
		parts = append(parts, strings.TrimRight(part.String(), "\n"))
	}
	return strings.Join(parts, "\n\n"), nil
}

func formatMarkdownTable(rows []listRow, opts ListOptions) (string, error) {
	columns := defaultColumnsForRows(opts, rows, []string{"id", "title", "priority", "role", "status", "completed", "blockers"}, false)
	builder := &strings.Builder{}
	fmt.Fprintln(builder, "| "+strings.Join(columns, " | ")+" |")
	separators := make([]string, 0, len(columns))
	for range columns {
		separators = append(separators, "---")
	}
	fmt.Fprintln(builder, "| "+strings.Join(separators, " | ")+" |")
	for _, row := range rows {
		fmt.Fprintln(builder, "| "+strings.Join(rowValues(row, columns, false, opts), " | ")+" |")
	}
	return strings.TrimRight(builder.String(), "\n"), nil
}

func formatMarkdownList(rows []listRow, opts ListOptions) string {
	if len(rows) == 0 {
		return ""
	}
	columns := defaultColumnsForRows(opts, rows, []string{"id", "title", "priority", "role", "status", "completed"}, false)
	builder := &strings.Builder{}
	for _, row := range rows {
		id := colorizeValue(row, "id", row.ID, opts)
		title := colorizeValue(row, "title", row.Title, opts)
		parts := listMetadataParts(row, columns, opts)
		if len(parts) == 0 {
			fmt.Fprintf(builder, "- %s — %s\n", id, title)
			continue
		}
		fmt.Fprintf(builder, "- %s — %s (%s)\n", id, title, strings.Join(parts, ", "))
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

func rowValues(row listRow, columns []string, numericForCounts bool, opts ListOptions) []string {
	values := make([]string, 0, len(columns))
	for _, col := range columns {
		value := columnValue(row, col, numericForCounts)
		values = append(values, colorizeValue(row, col, value, opts))
	}
	return values
}

func defaultColumnsForRows(opts ListOptions, rows []listRow, defaults []string, numericForCounts bool) []string {
	if len(opts.Columns) > 0 {
		return normalizeColumns(opts.Columns, defaults)
	}
	columns := normalizeColumns(nil, defaults)
	return filterConstantColumns(rows, columns, numericForCounts)
}

func filterConstantColumns(rows []listRow, columns []string, numericForCounts bool) []string {
	if len(rows) == 0 {
		return columns
	}
	alwaysKeep := map[string]bool{"id": true, "title": true}
	out := make([]string, 0, len(columns))
	for _, col := range columns {
		if alwaysKeep[col] {
			out = append(out, col)
			continue
		}
		baseline := columnValue(rows[0], col, numericForCounts)
		allSame := true
		for _, row := range rows[1:] {
			if columnValue(row, col, numericForCounts) != baseline {
				allSame = false
				break
			}
		}
		if allSame {
			continue
		}
		out = append(out, col)
	}
	return out
}

func columnValue(row listRow, col string, numericForCounts bool) string {
	switch col {
	case "id":
		return row.ID
	case "title":
		return row.Title
	case "priority":
		return row.Priority
	case "role":
		return row.Role
	case "parent":
		return row.Parent
	case "completed":
		return fmt.Sprintf("%t", row.Completed)
	case "status":
		return row.Status
	case "blockers":
		return formatListValue(row.Blockers, numericForCounts)
	case "blocks":
		return formatListValue(row.Blocks, numericForCounts)
	case "path":
		return row.Path
	case "date_created":
		return row.DateCreated
	case "date_edited":
		return row.DateEdited
	default:
		return ""
	}
}

func colorizeValue(row listRow, col string, value string, opts ListOptions) string {
	if !opts.Color || value == "" {
		return value
	}
	switch col {
	case "priority":
		return colorizePriority(value)
	case "id", "title":
		if len(row.Blockers) > 0 {
			return colorBlocked(value)
		}
		if NormalizePriority(row.Priority) == PriorityHigh {
			return colorBold(value)
		}
	}
	return value
}

func colorizePriority(priority string) string {
	switch NormalizePriority(priority) {
	case PriorityHigh:
		return colorHigh(priority)
	case PriorityMedium:
		return colorMedium(priority)
	case PriorityLow:
		return colorLow(priority)
	default:
		return colorOther(priority)
	}
}

func listMetadataParts(row listRow, columns []string, opts ListOptions) []string {
	parts := []string{}
	for _, col := range columns {
		if col == "id" || col == "title" {
			continue
		}
		value := columnValue(row, col, false)
		if value == "" {
			continue
		}
		value = colorizeValue(row, col, value, opts)
		label := columnLabel(col)
		if label == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s: %s", label, value))
	}
	return parts
}

func columnLabel(col string) string {
	switch col {
	case "priority":
		return "priority"
	case "role":
		return "role"
	case "completed":
		return "completed"
	case "status":
		return "status"
	case "blockers":
		return "blockers"
	case "blocks":
		return "blocks"
	case "parent":
		return "parent"
	case "path":
		return "path"
	case "date_created":
		return "created"
	case "date_edited":
		return "edited"
	default:
		return ""
	}
}

var (
	colorHigh    = color.New(color.FgRed, color.Bold).SprintFunc()
	colorMedium  = color.New(color.FgYellow).SprintFunc()
	colorLow     = color.New(color.FgGreen).SprintFunc()
	colorOther   = color.New(color.FgCyan).SprintFunc()
	colorBlocked = color.New(color.FgHiBlack).SprintFunc()
	colorBold    = color.New(color.Bold).SprintFunc()
)

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
