package cmd

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

// mcpCmd represents the MCP server command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run an MCP server for strand",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMCP()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

type addArgs struct {
	Project  string   `json:"project,omitempty" jsonschema_description:"Project name (equivalent to --project)"`
	Type     string   `json:"type" jsonschema:"required" jsonschema_description:"Template type name"`
	Title    string   `json:"title" jsonschema:"required" jsonschema_description:"Task title"`
	Role     string   `json:"role,omitempty" jsonschema_description:"Role responsible for the task"`
	Priority string   `json:"priority,omitempty" jsonschema:"enum=high,enum=medium,enum=low" jsonschema_description:"Task priority"`
	Parent   string   `json:"parent,omitempty" jsonschema_description:"Parent task ID"`
	Blockers []string `json:"blockers,omitempty" jsonschema_description:"Blocker task IDs"`
	NoRepair bool     `json:"no_repair,omitempty" jsonschema_description:"Skip repair and master list updates"`
	Body     string   `json:"body,omitempty" jsonschema_description:"Task body content"`
}

type nextArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"Project name (equivalent to --project)"`
	Role    string `json:"role,omitempty" jsonschema_description:"Filter by role"`
}

type completeArgs struct {
	Project string `json:"project,omitempty" jsonschema_description:"Project name (equivalent to --project)"`
	TaskID  string `json:"task_id" jsonschema:"required" jsonschema_description:"Task ID or short ID"`
}

type initArgs struct {
	ProjectName string `json:"project_name,omitempty" jsonschema_description:"Project name (defaults to git root name)"`
	Storage     string `json:"storage,omitempty" jsonschema:"enum=global,enum=local" jsonschema_description:"Storage mode"`
	Preset      string `json:"preset,omitempty" jsonschema_description:"Preset directory or git repo"`
}

type repairArgs struct {
	Project   string `json:"project,omitempty" jsonschema_description:"Project name (equivalent to --project)"`
	TasksRoot string `json:"tasks_root,omitempty" jsonschema_description:"Tasks root directory"`
	RootsFile string `json:"roots_file,omitempty" jsonschema_description:"Root tasks list path"`
	FreeFile  string `json:"free_file,omitempty" jsonschema_description:"Free tasks list path"`
	Format    string `json:"format,omitempty" jsonschema:"enum=text,enum=json" jsonschema_description:"Output format"`
}

type listArgs struct {
	Project        string   `json:"project,omitempty" jsonschema_description:"Project name (equivalent to --project)"`
	Scope          string   `json:"scope,omitempty" jsonschema:"enum=all,enum=root,enum=free" jsonschema_description:"Scope of tasks to list"`
	Children       string   `json:"children,omitempty" jsonschema_description:"List direct children of the given task ID"`
	Role           string   `json:"role,omitempty" jsonschema_description:"Filter by role"`
	Priority       string   `json:"priority,omitempty" jsonschema:"enum=high,enum=medium,enum=low" jsonschema_description:"Filter by priority"`
	Completed      *bool    `json:"completed,omitempty" jsonschema_description:"Filter by completed status"`
	Blocked        *bool    `json:"blocked,omitempty" jsonschema_description:"Filter by blocked status"`
	Blocks         *bool    `json:"blocks,omitempty" jsonschema_description:"Filter by blocks status"`
	OwnerApproval  *bool    `json:"owner_approval,omitempty" jsonschema_description:"Filter by owner approval"`
	Label          string   `json:"label,omitempty" jsonschema_description:"Reserved for future labels support"`
	Sort           string   `json:"sort,omitempty" jsonschema:"enum=id,enum=priority,enum=created,enum=edited,enum=role" jsonschema_description:"Sort field"`
	Order          string   `json:"order,omitempty" jsonschema:"enum=asc,enum=desc" jsonschema_description:"Sort order"`
	Format         string   `json:"format,omitempty" jsonschema:"enum=table,enum=md,enum=json" jsonschema_description:"Output format"`
	Columns        []string `json:"columns,omitempty" jsonschema_description:"Columns to include"`
	Group          string   `json:"group,omitempty" jsonschema:"enum=none,enum=priority,enum=parent,enum=role" jsonschema_description:"Group by"`
	MdTable        bool     `json:"md_table,omitempty" jsonschema_description:"Use markdown table output"`
	UseMasterLists bool     `json:"use_master_lists,omitempty" jsonschema_description:"Use master lists for root/free"`
}

type searchArgs struct {
	Project string   `json:"project,omitempty" jsonschema_description:"Project name (equivalent to --project)"`
	Query   string   `json:"query" jsonschema:"required" jsonschema_description:"Search query"`
	Sort    string   `json:"sort,omitempty" jsonschema:"enum=id,enum=priority,enum=created,enum=edited,enum=role" jsonschema_description:"Sort field"`
	Order   string   `json:"order,omitempty" jsonschema:"enum=asc,enum=desc" jsonschema_description:"Sort order"`
	Format  string   `json:"format,omitempty" jsonschema:"enum=table,enum=md,enum=json" jsonschema_description:"Output format"`
	Columns []string `json:"columns,omitempty" jsonschema_description:"Columns to include"`
	Group   string   `json:"group,omitempty" jsonschema:"enum=none,enum=priority,enum=parent,enum=role" jsonschema_description:"Group by"`
	MdTable bool     `json:"md_table,omitempty" jsonschema_description:"Use markdown table output"`
}

func runMCP() error {
	s := server.NewMCPServer("strand", "dev", server.WithToolCapabilities(true))
	registerMCPTools(s)
	return server.ServeStdio(s)
}

func registerMCPTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("strand_add",
			mcp.WithDescription("Create tasks from templates"),
			mcp.WithInputSchema[addArgs](),
		),
		mcp.NewTypedToolHandler(handleMCPAdd),
	)

	s.AddTool(
		mcp.NewTool("strand_next",
			mcp.WithDescription("Print the next free task"),
			mcp.WithInputSchema[nextArgs](),
		),
		mcp.NewTypedToolHandler(handleMCPNext),
	)

	s.AddTool(
		mcp.NewTool("strand_complete",
			mcp.WithDescription("Mark a task as completed"),
			mcp.WithInputSchema[completeArgs](),
		),
		mcp.NewTypedToolHandler(handleMCPComplete),
	)

	s.AddTool(
		mcp.NewTool("strand_init",
			mcp.WithDescription("Initialize strand storage"),
			mcp.WithInputSchema[initArgs](),
		),
		mcp.NewTypedToolHandler(handleMCPInit),
	)

	s.AddTool(
		mcp.NewTool("strand_repair",
			mcp.WithDescription("Repair task tree and regenerate master lists"),
			mcp.WithInputSchema[repairArgs](),
		),
		mcp.NewTypedToolHandler(handleMCPRepair),
	)

	s.AddTool(
		mcp.NewTool("strand_list",
			mcp.WithDescription("List tasks with filtering and formatting options"),
			mcp.WithInputSchema[listArgs](),
		),
		mcp.NewTypedToolHandler(handleMCPList),
	)

	s.AddTool(
		mcp.NewTool("strand_search",
			mcp.WithDescription("Search tasks by title, description, and todos"),
			mcp.WithInputSchema[searchArgs](),
		),
		mcp.NewTypedToolHandler(handleMCPSearch),
	)

	s.AddTool(
		mcp.NewTool("strand_agents",
			mcp.WithDescription("Print portable agent instructions"),
			mcp.WithInputSchema[struct{}](),
		),
		mcp.NewTypedToolHandler(handleMCPAgents),
	)

	s.AddTool(
		mcp.NewTool("strand_templates",
			mcp.WithDescription("Describe available templates"),
			mcp.WithInputSchema[struct{}](),
		),
		mcp.NewTypedToolHandler(handleMCPTemplates),
	)

	s.AddTool(
		mcp.NewTool("strand_assign",
			mcp.WithDescription("Assign a task to a role"),
			mcp.WithInputSchema[struct{}](),
		),
		mcp.NewTypedToolHandler(handleMCPAssign),
	)

	s.AddTool(
		mcp.NewTool("strand_block",
			mcp.WithDescription("Block or unblock tasks"),
			mcp.WithInputSchema[struct{}](),
		),
		mcp.NewTypedToolHandler(handleMCPBlock),
	)
}

func handleMCPAdd(ctx context.Context, request mcp.CallToolRequest, args addArgs) (*mcp.CallToolResult, error) {
	opts := addOptions{
		ProjectName:       strings.TrimSpace(args.Project),
		TemplateName:      strings.TrimSpace(args.Type),
		Title:             strings.TrimSpace(args.Title),
		Role:              strings.TrimSpace(args.Role),
		Priority:          strings.TrimSpace(args.Priority),
		Parent:            strings.TrimSpace(args.Parent),
		Blockers:          args.Blockers,
		NoRepair:          args.NoRepair,
		RoleSpecified:     strings.TrimSpace(args.Role) != "",
		PrioritySpecified: strings.TrimSpace(args.Priority) != "",
		Body:              args.Body,
	}
	return runWithOutput(func(w io.Writer) error {
		return runAdd(w, opts)
	})
}

func handleMCPNext(ctx context.Context, request mcp.CallToolRequest, args nextArgs) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		return runNext(w, strings.TrimSpace(args.Project), strings.TrimSpace(args.Role))
	})
}

func handleMCPComplete(ctx context.Context, request mcp.CallToolRequest, args completeArgs) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		return runComplete(w, strings.TrimSpace(args.Project), strings.TrimSpace(args.TaskID), 0)
	})
}

func handleMCPInit(ctx context.Context, request mcp.CallToolRequest, args initArgs) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		return runInit(w, initOptions{
			ProjectName: strings.TrimSpace(args.ProjectName),
			StorageMode: strings.TrimSpace(args.Storage),
			Preset:      strings.TrimSpace(args.Preset),
		})
	})
}

func handleMCPRepair(ctx context.Context, request mcp.CallToolRequest, args repairArgs) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		format := strings.ToLower(strings.TrimSpace(args.Format))
		if format == "" {
			format = "text"
		}
		paths, err := resolveProjectPaths(strings.TrimSpace(args.Project))
		if err != nil {
			return err
		}
		tasksRoot := strings.TrimSpace(args.TasksRoot)
		if tasksRoot == "" {
			tasksRoot = paths.TasksDir
		}
		rootsFile := strings.TrimSpace(args.RootsFile)
		if rootsFile == "" {
			rootsFile = paths.RootTasksFile
		}
		freeFile := strings.TrimSpace(args.FreeFile)
		if freeFile == "" {
			freeFile = paths.FreeTasksFile
		}
		return runRepair(w, tasksRoot, rootsFile, freeFile, format)
	})
}

func handleMCPList(ctx context.Context, request mcp.CallToolRequest, args listArgs) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		opts := task.ListOptions{
			Scope:          normalizeEnum(args.Scope, "all"),
			Parent:         strings.TrimSpace(args.Children),
			Role:           strings.TrimSpace(args.Role),
			Priority:       normalizeEnum(args.Priority, ""),
			Label:          strings.TrimSpace(args.Label),
			Sort:           normalizeEnum(args.Sort, ""),
			Order:          normalizeEnum(args.Order, "asc"),
			Format:         normalizeEnum(args.Format, "table"),
			Group:          normalizeEnum(args.Group, "none"),
			MdTable:        args.MdTable,
			UseMasterLists: args.UseMasterLists,
			Color:          false,
		}
		if args.Completed == nil {
			opts.Completed = boolPtr(false)
		} else {
			opts.Completed = boolPtr(*args.Completed)
		}
		if args.Blocked != nil {
			opts.Blocked = boolPtr(*args.Blocked)
		}
		if args.Blocks != nil {
			opts.Blocks = boolPtr(*args.Blocks)
		}
		if args.OwnerApproval != nil {
			opts.OwnerApproval = boolPtr(*args.OwnerApproval)
		}
		if len(args.Columns) > 0 {
			opts.Columns = normalizeColumns(args.Columns)
		}
		return runListWithProject(w, strings.TrimSpace(args.Project), opts)
	})
}

func handleMCPSearch(ctx context.Context, request mcp.CallToolRequest, args searchArgs) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		query := strings.TrimSpace(args.Query)
		if query == "" {
			return mcpError("search query cannot be empty")
		}
		opts := task.SearchOptions{
			Query: query,
			ListOptions: task.ListOptions{
				Sort:    normalizeEnum(args.Sort, ""),
				Order:   normalizeEnum(args.Order, "asc"),
				Format:  normalizeEnum(args.Format, "table"),
				Group:   normalizeEnum(args.Group, "none"),
				MdTable: args.MdTable,
				Color:   false,
			},
		}
		if len(args.Columns) > 0 {
			opts.Columns = normalizeColumns(args.Columns)
		}
		return runSearchWithProject(w, strings.TrimSpace(args.Project), opts)
	})
}

func handleMCPAgents(ctx context.Context, request mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		return runAgents(w)
	})
}

func handleMCPTemplates(ctx context.Context, request mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		return runTemplates(w)
	})
}

func handleMCPAssign(ctx context.Context, request mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		return runAssign(w)
	})
}

func handleMCPBlock(ctx context.Context, request mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, error) {
	return runWithOutput(func(w io.Writer) error {
		return runBlock(w)
	})
}

func runWithOutput(run func(w io.Writer) error) (*mcp.CallToolResult, error) {
	var buf bytes.Buffer
	if err := run(&buf); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(buf.String()), nil
}

func normalizeEnum(value, fallback string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func normalizeColumns(columns []string) []string {
	out := make([]string, 0, len(columns))
	for _, col := range columns {
		col = strings.ToLower(strings.TrimSpace(col))
		if col == "" {
			continue
		}
		out = append(out, col)
	}
	return out
}

type mcpError string

func (e mcpError) Error() string {
	return string(e)
}
