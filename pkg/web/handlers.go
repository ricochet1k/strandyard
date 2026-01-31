package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ricochet1k/strandyard/pkg/task"
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

type projectResponse struct {
	Name          string `json:"name"`
	StorageRoot   string `json:"storage_root"`
	TasksRoot     string `json:"tasks_root"`
	RolesRoot     string `json:"roles_root"`
	TemplatesRoot string `json:"templates_root"`
	GitRoot       string `json:"git_root"`
	Storage       string `json:"storage"`
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
			Parent:      task.ShortID(t.Meta.Parent),
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
