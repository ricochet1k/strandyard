package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ricochet1k/strandyard/pkg/idgen"
	rPkg "github.com/ricochet1k/strandyard/pkg/role"
	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/ricochet1k/strandyard/pkg/template"
	"gopkg.in/yaml.v3"
)

type fileEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Kind string `json:"kind"`
}

type taskListItem struct {
	ID          string   `json:"id"`
	ShortID     string   `json:"short_id"`
	Title       string   `json:"title"`
	Role        string   `json:"role"`
	Priority    string   `json:"priority"`
	Completed   bool     `json:"completed"`
	Parent      string   `json:"parent"`
	Blockers    []string `json:"blockers"`
	Blocks      []string `json:"blocks"`
	Path        string   `json:"path"`
	DateCreated string   `json:"date_created"`
	DateEdited  string   `json:"date_edited"`
}

type filePayload struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type taskUpdateRequest struct {
	Title     *string   `json:"title,omitempty"`
	Role      *string   `json:"role,omitempty"`
	Priority  *string   `json:"priority,omitempty"`
	Completed *bool     `json:"completed,omitempty"`
	Parent    *string   `json:"parent,omitempty"`
	Blockers  *[]string `json:"blockers,omitempty"`
	Blocks    *[]string `json:"blocks,omitempty"`
	Body      *string   `json:"body,omitempty"`
}

type taskCreateRequest struct {
	TemplateName string   `json:"template_name"`
	Title        string   `json:"title"`
	Role         string   `json:"role,omitempty"`
	Priority     string   `json:"priority,omitempty"`
	Parent       string   `json:"parent,omitempty"`
	Blockers     []string `json:"blockers,omitempty"`
	Blocks       []string `json:"blocks,omitempty"`
	Body         string   `json:"body,omitempty"`
}

type taskDetailResponse struct {
	ID          string   `json:"id"`
	ShortID     string   `json:"short_id"`
	Title       string   `json:"title"`
	Role        string   `json:"role"`
	Priority    string   `json:"priority"`
	Completed   bool     `json:"completed"`
	Parent      string   `json:"parent"`
	Blockers    []string `json:"blockers"`
	Blocks      []string `json:"blocks"`
	Path        string   `json:"path"`
	DateCreated string   `json:"date_created"`
	DateEdited  string   `json:"date_edited"`
	Body        string   `json:"body"`
}

type projectResponse struct {
	Name          string `json:"name"`
	StorageRoot   string `json:"storage_root"`
	TasksRoot     string `json:"tasks_root"`
	RolesRoot     string `json:"roles_root"`
	TemplatesRoot string `json:"templates_root"`
	GitRoot       string `json:"git_root"`
	Storage       string `json:"storage"`
}

type roleDetailResponse struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

type roleUpdateRequest struct {
	Description *string `json:"description,omitempty"`
	Body        *string `json:"body,omitempty"`
}

type templateDetailResponse struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Role        string `json:"role"`
	Priority    string `json:"priority"`
	Description string `json:"description"`
	IDPrefix    string `json:"id_prefix"`
	Body        string `json:"body"`
}

type templateUpdateRequest struct {
	Role        *string `json:"role,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	Description *string `json:"description,omitempty"`
	IDPrefix    *string `json:"id_prefix,omitempty"`
	Body        *string `json:"body,omitempty"`
}

func (s *Server) getProject(r *http.Request) (*ProjectInfo, error) {
	projectName := strings.TrimSpace(r.URL.Query().Get("project"))
	if projectName == "" {
		// Default to current project if set
		if s.config.CurrentProject != "" {
			projectName = s.config.CurrentProject
		} else if len(s.config.Projects) > 0 {
			projectName = s.config.Projects[0].Name
		} else {
			return nil, fmt.Errorf("no projects available")
		}
	}

	proj, ok := s.projects[projectName]
	if !ok {
		return nil, fmt.Errorf("project not found: %s", projectName)
	}

	return proj, nil
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	projects := make([]projectResponse, 0, len(s.config.Projects))
	for _, proj := range s.config.Projects {
		projects = append(projects, projectResponse{
			Name:          proj.Name,
			StorageRoot:   proj.StorageRoot,
			TasksRoot:     proj.TasksRoot,
			RolesRoot:     proj.RolesRoot,
			TemplatesRoot: proj.TemplatesRoot,
			GitRoot:       proj.GitRoot,
			Storage:       proj.Storage,
		})
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"projects": projects,
		"current":  s.config.CurrentProject,
	})
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"project":        proj.Name,
		"storage_root":   proj.StorageRoot,
		"tasks_root":     proj.TasksRoot,
		"roles_root":     proj.RolesRoot,
		"templates_root": proj.TemplatesRoot,
		"git_root":       proj.GitRoot,
		"storage":        proj.Storage,
	})
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	items, err := s.listTasks(proj)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (s *Server) handleRoles(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	roles, err := rPkg.LoadRoles(proj.RolesRoot)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	items := make([]roleDetailResponse, 0, len(roles))
	for _, t := range roles {
		items = append(items, roleDetailResponse{
			Name:        t.ID,
			Path:        makeRelative(proj.StorageRoot, t.FilePath),
			Description: t.Meta.Description,
			Body:        t.BodyContent,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	respondJSON(w, http.StatusOK, items)
}

func (s *Server) handleTemplates(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	templates, err := template.LoadTemplates(proj.TemplatesRoot)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	items := make([]templateDetailResponse, 0, len(templates))
	for _, t := range templates {
		priorityStr := ""
		if p, ok := t.Meta.Priority.(string); ok {
			priorityStr = p
		}

		items = append(items, templateDetailResponse{
			Name:        t.ID,
			Path:        makeRelative(proj.StorageRoot, filepath.Join(proj.TemplatesRoot, t.ID+".md")),
			Role:        t.Meta.Role,
			Priority:    priorityStr,
			Description: t.Meta.Description,
			IDPrefix:    t.Meta.IDPrefix,
			Body:        t.BodyContent,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	respondJSON(w, http.StatusOK, items)
}

func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.handleTaskCreate(w, r, proj)
	case http.MethodGet, http.MethodPatch:
		s.handleTaskGetOrUpdate(w, r, proj)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRole(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		respondError(w, http.StatusBadRequest, fmt.Errorf("missing path"))
		return
	}

	resolved, err := s.resolvePath(proj, path)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		p := task.NewParser()
		t, err := p.ParseStandaloneFile(resolved)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respondJSON(w, http.StatusOK, roleDetailResponse{
			Name:        t.ID,
			Path:        path,
			Description: t.Meta.Description,
			Body:        t.BodyContent,
		})

	case http.MethodPut:
		if s.config.ReadOnly {
			respondError(w, http.StatusForbidden, fmt.Errorf("server is in read-only mode"))
			return
		}

		var req roleUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		p := task.NewParser()
		t, err := p.ParseStandaloneFile(resolved)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		if req.Description != nil {
			t.Meta.Description = *req.Description
		}
		if req.Body != nil {
			t.BodyContent = *req.Body
		}

		if err := t.Write(); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respondJSON(w, http.StatusOK, roleDetailResponse{
			Name:        t.ID,
			Path:        path,
			Description: t.Meta.Description,
			Body:        t.BodyContent,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTemplate(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		respondError(w, http.StatusBadRequest, fmt.Errorf("missing path"))
		return
	}

	resolved, err := s.resolvePath(proj, path)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		data, err := os.ReadFile(resolved)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		content := string(data)
		parts := strings.SplitN(content, "---", 3)
		if len(parts) < 3 {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("invalid template format"))
			return
		}

		var meta template.TemplateMetadata
		if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		priorityStr := ""
		if p, ok := meta.Priority.(string); ok {
			priorityStr = p
		}

		respondJSON(w, http.StatusOK, templateDetailResponse{
			Name:        strings.TrimSuffix(filepath.Base(path), ".md"),
			Path:        path,
			Role:        meta.Role,
			Priority:    priorityStr,
			Description: meta.Description,
			IDPrefix:    meta.IDPrefix,
			Body:        strings.TrimSpace(parts[2]),
		})

	case http.MethodPut:
		if s.config.ReadOnly {
			respondError(w, http.StatusForbidden, fmt.Errorf("server is in read-only mode"))
			return
		}

		var req templateUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		data, err := os.ReadFile(resolved)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		content := string(data)
		parts := strings.SplitN(content, "---", 3)
		if len(parts) < 3 {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("invalid template format"))
			return
		}

		var meta template.TemplateMetadata
		if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		if req.Role != nil {
			meta.Role = *req.Role
		}
		if req.Priority != nil {
			meta.Priority = *req.Priority
		}
		if req.Description != nil {
			meta.Description = *req.Description
		}
		if req.IDPrefix != nil {
			meta.IDPrefix = *req.IDPrefix
		}

		body := strings.TrimSpace(parts[2])
		if req.Body != nil {
			body = *req.Body
		}

		frontmatterBytes, err := yaml.Marshal(&meta)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		var sb strings.Builder
		sb.WriteString("---\n")
		sb.Write(frontmatterBytes)
		sb.WriteString("---\n\n")
		sb.WriteString(body)
		sb.WriteString("\n")

		if err := os.WriteFile(resolved, []byte(sb.String()), 0o644); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		priorityStr := ""
		if p, ok := meta.Priority.(string); ok {
			priorityStr = p
		}

		respondJSON(w, http.StatusOK, templateDetailResponse{
			Name:        strings.TrimSuffix(filepath.Base(path), ".md"),
			Path:        path,
			Role:        meta.Role,
			Priority:    priorityStr,
			Description: meta.Description,
			IDPrefix:    meta.IDPrefix,
			Body:        body,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTaskGetOrUpdate(w http.ResponseWriter, r *http.Request, proj *ProjectInfo) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		respondError(w, http.StatusBadRequest, fmt.Errorf("missing task id"))
		return
	}

	db := task.NewTaskDB(proj.TasksRoot)
	if err := db.LoadAll(); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	t, err := db.Get(taskID)
	if err != nil {
		respondError(w, http.StatusNotFound, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		snapshot, err := taskToSnapshot(t, proj.StorageRoot)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		respondJSON(w, http.StatusOK, snapshot)

	case http.MethodPatch:
		if s.config.ReadOnly {
			respondError(w, http.StatusForbidden, fmt.Errorf("server is in read-only mode"))
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		var req taskUpdateRequest
		if err := json.Unmarshal(body, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		if req.Title != nil {
			t.SetTitle(*req.Title)
		}
		if req.Role != nil {
			t.Meta.Role = *req.Role
			t.MarkDirty()
		}
		if req.Priority != nil {
			t.Meta.Priority = task.NormalizePriority(*req.Priority)
			t.MarkDirty()
		}
		if req.Completed != nil {
			t.Meta.Completed = *req.Completed
			t.MarkDirty()
		}
		if req.Blockers != nil {
			t.Meta.Blockers = *req.Blockers
			t.MarkDirty()
		}
		if req.Blocks != nil {
			t.Meta.Blocks = *req.Blocks
			t.MarkDirty()
		}
		if req.Body != nil {
			t.SetBody(*req.Body)
		}

		if req.Parent != nil {
			oldParent := t.Meta.Parent
			newParent := *req.Parent
			if newParent != "" {
				resolved, err := db.ResolveID(newParent)
				if err != nil {
					respondError(w, http.StatusBadRequest, fmt.Errorf("invalid parent id: %w", err))
					return
				}
				newParent = resolved
			}

			if oldParent != newParent {
				if err := db.SetParent(t.ID, newParent); err != nil {
					respondError(w, http.StatusBadRequest, err)
					return
				}
				if oldParent != "" {
					db.UpdateParentTodos(oldParent)
				}
				if newParent != "" {
					db.UpdateParentTodos(newParent)
				}
			}
		}

		if _, err := db.SaveDirty(); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		snapshot, err := taskToSnapshot(t, proj.StorageRoot)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		respondJSON(w, http.StatusOK, snapshot)
	}
}

func (s *Server) handleTaskCreate(w http.ResponseWriter, r *http.Request, proj *ProjectInfo) {
	if s.config.ReadOnly {
		respondError(w, http.StatusForbidden, fmt.Errorf("server is in read-only mode"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	var req taskCreateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	// Create a bytes buffer to capture output
	var output strings.Builder

	// Build addOptions from the request
	opts := struct {
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
	}{
		ProjectName:       proj.Name,
		TemplateName:      req.TemplateName,
		Title:             req.Title,
		Role:              req.Role,
		Priority:          req.Priority,
		Parent:            req.Parent,
		Blockers:          req.Blockers,
		Blocks:            req.Blocks,
		Every:             []string{},
		RoleSpecified:     req.Role != "",
		PrioritySpecified: req.Priority != "",
		Body:              req.Body,
	}

	// We need to use the internal task creation logic
	// For now, we'll shell out to the add command logic
	// This is a simplified version - in production you'd refactor the add logic into a shared function
	if err := s.createTask(&output, opts, proj); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("%s: %w", output.String(), err))
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{
		"status":  "created",
		"message": output.String(),
	})
}

func taskToSnapshot(t *task.Task, storageRoot string) (*taskDetailResponse, error) {
	return &taskDetailResponse{
		ID:          t.ID,
		ShortID:     task.ShortID(t.ID),
		Title:       t.Title(),
		Role:        t.GetEffectiveRole(),
		Priority:    task.NormalizePriority(t.Meta.Priority),
		Completed:   t.Meta.Completed,
		Parent:      t.Meta.Parent,
		Blockers:    shortenIDs(t.Meta.Blockers),
		Blocks:      shortenIDs(t.Meta.Blocks),
		Path:        makeRelative(storageRoot, t.FilePath),
		DateCreated: t.Meta.DateCreated.Format(time.RFC3339),
		DateEdited:  t.Meta.DateEdited.Format(time.RFC3339),
		Body:        t.BodyContent,
	}, nil
}

func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	kind := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("kind")))
	entries, err := s.listFiles(proj, kind)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	respondJSON(w, http.StatusOK, entries)
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	proj, err := s.getProject(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	path := r.URL.Query().Get("path")
	if strings.TrimSpace(path) == "" {
		respondError(w, http.StatusBadRequest, fmt.Errorf("missing path"))
		return
	}

	resolved, err := s.resolvePath(proj, path)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		data, err := os.ReadFile(resolved)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		respondJSON(w, http.StatusOK, filePayload{Path: filepath.ToSlash(path), Content: string(data)})
	case http.MethodPut:
		if s.config.ReadOnly {
			respondError(w, http.StatusForbidden, fmt.Errorf("server is in read-only mode"))
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		var payload filePayload
		if err := json.Unmarshal(body, &payload); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		if payload.Content == "" {
			respondError(w, http.StatusBadRequest, fmt.Errorf("content cannot be empty"))
			return
		}
		if err := os.WriteFile(resolved, []byte(payload.Content), 0o644); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "saved"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		respondError(w, http.StatusInternalServerError, fmt.Errorf("streaming unsupported"))
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	updates := make(chan StreamUpdate, 32)
	s.broker.subscribe(updates)
	defer s.broker.unsubscribe(updates)

	writeSSE(w, "ready", map[string]string{"status": "connected"})
	flusher.Flush()

	keepalive := time.NewTicker(20 * time.Second)
	defer keepalive.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			writeSSE(w, "task", update)
			flusher.Flush()
		case <-keepalive.C:
			writeSSE(w, "ping", map[string]string{"time": time.Now().Format(time.RFC3339)})
			flusher.Flush()
		}
	}
}

func (s *Server) listTasks(proj *ProjectInfo) ([]taskListItem, error) {
	parser := task.NewParser()
	tasks, err := parser.LoadTasks(proj.TasksRoot)
	if err != nil {
		return nil, err
	}

	items := make([]taskListItem, 0, len(tasks))
	for _, t := range tasks {
		relPath := makeRelative(proj.StorageRoot, t.FilePath)
		items = append(items, taskListItem{
			ID:          t.ID,
			ShortID:     task.ShortID(t.ID),
			Title:       t.Title(),
			Role:        t.GetEffectiveRole(),
			Priority:    task.NormalizePriority(t.Meta.Priority),
			Completed:   t.Meta.Completed,
			Parent:      t.Meta.Parent,
			Blockers:    shortenIDs(t.Meta.Blockers),
			Blocks:      shortenIDs(t.Meta.Blocks),
			Path:        relPath,
			DateCreated: t.Meta.DateCreated.Format(time.RFC3339),
			DateEdited:  t.Meta.DateEdited.Format(time.RFC3339),
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority != items[j].Priority {
			return task.PriorityRank(items[i].Priority) < task.PriorityRank(items[j].Priority)
		}
		return items[i].ID < items[j].ID
	})

	return items, nil
}

func (s *Server) listFiles(proj *ProjectInfo, kind string) ([]fileEntry, error) {
	var root string
	switch kind {
	case "roles":
		root = proj.RolesRoot
	case "templates":
		root = proj.TemplatesRoot
	case "tasks":
		root = proj.TasksRoot
	default:
		return nil, fmt.Errorf("unknown kind: %s", kind)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	items := make([]fileEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			continue
		}
		rel := makeRelative(proj.StorageRoot, filepath.Join(root, name))
		items = append(items, fileEntry{
			Name: strings.TrimSuffix(name, filepath.Ext(name)),
			Path: rel,
			Kind: kind,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items, nil
}

func (s *Server) resolvePath(proj *ProjectInfo, rel string) (string, error) {
	cleaned := filepath.Clean(strings.TrimSpace(rel))
	cleaned = strings.TrimPrefix(cleaned, string(filepath.Separator))
	if cleaned == "." || cleaned == "" {
		return "", fmt.Errorf("invalid path")
	}

	abs := filepath.Join(proj.StorageRoot, cleaned)
	if !isWithin(abs, proj.StorageRoot) {
		return "", fmt.Errorf("path outside storage root")
	}
	if !isWithin(abs, proj.TasksRoot) && !isWithin(abs, proj.RolesRoot) && !isWithin(abs, proj.TemplatesRoot) {
		return "", fmt.Errorf("path not in tasks, roles, or templates")
	}

	if _, err := os.Stat(abs); err != nil {
		return "", err
	}

	return abs, nil
}

func isWithin(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".."
}

func makeRelative(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}

func shortenIDs(ids []string) []string {
	if len(ids) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, task.ShortID(id))
	}
	return out
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(payload)
}

func respondError(w http.ResponseWriter, status int, err error) {
	respondJSON(w, status, map[string]string{"error": err.Error()})
}

func writeSSE(w http.ResponseWriter, event string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "event: %s\n", event)
	fmt.Fprintf(w, "data: %s\n\n", data)
}

func (s *Server) createTask(w io.Writer, opts struct {
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
}, proj *ProjectInfo) error {
	db := task.NewTaskDB(proj.TasksRoot)
	if err := db.LoadAllIfEmpty(); err != nil {
		return err
	}

	tmplName := strings.TrimSpace(opts.TemplateName)
	if tmplName == "" {
		return fmt.Errorf("template name is required")
	}

	templates, err := template.LoadTemplates(proj.TemplatesRoot)
	if err != nil {
		return err
	}

	tmpl, ok := templates[tmplName]
	if !ok {
		return fmt.Errorf("unknown template %q", tmplName)
	}

	title := strings.TrimSpace(opts.Title)
	if title == "" {
		return fmt.Errorf("title is required")
	}

	// Reject placeholder titles
	invalidTitles := []string{"description", "task title", "new task", "title", "summary", "todo"}
	lowerTitle := strings.ToLower(title)
	for _, invalid := range invalidTitles {
		if lowerTitle == invalid {
			return fmt.Errorf("title %q looks like a placeholder; please provide a descriptive title", title)
		}
	}

	roleName := strings.TrimSpace(opts.Role)
	if !opts.RoleSpecified {
		roleName = strings.TrimSpace(tmpl.Meta.Role)
	}
	if roleName == "" {
		return fmt.Errorf("role is required")
	}

	roles, err := rPkg.LoadRoles(proj.RolesRoot)
	if err != nil {
		return err
	}

	if _, ok := roles[roleName]; !ok {
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
	if parent != "" {
		resolvedParent, err := db.ResolveID(parent)
		if err != nil {
			return fmt.Errorf("parent task %s does not exist: %w", parent, err)
		}
		parent = resolvedParent
		_, err = db.Get(parent)
		if err != nil {
			return fmt.Errorf("parent task %s does not exist: %w", parent, err)
		}
	}

	prefix := "T"
	if strings.Contains(strings.ToLower(tmplName), "epic") {
		prefix = "E"
	}

	id, err := idgen.GenerateID(prefix, title)
	if err != nil {
		return err
	}

	taskFile := filepath.Join(proj.TasksRoot, id+".md")
	if _, err := os.Stat(taskFile); err == nil {
		return fmt.Errorf("task file already exists: %s", taskFile)
	}

	blockers, err := db.ResolveIDs(normalizeTaskIDsWeb(opts.Blockers))
	if err != nil {
		return err
	}
	blocks, err := db.ResolveIDs(normalizeTaskIDsWeb(opts.Blocks))
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

	body := renderTemplateBodyWeb(tmpl.BodyContent, map[string]string{
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
	if err := writeTaskFileWeb(taskFile, meta, body); err != nil {
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

	if parent != "" {
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

	return nil
}

func normalizeTaskIDsWeb(items []string) []string {
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

func renderTemplateBodyWeb(body string, data map[string]string) string {
	out := body
	for key, value := range data {
		out = strings.ReplaceAll(out, "{{ ."+key+" }}", value)
	}
	return out
}

func writeTaskFileWeb(path string, meta task.Metadata, body string) error {
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
