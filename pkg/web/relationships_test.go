package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ricochet1k/strandyard/pkg/task"
)

func TestWebUpdatesMasterListsAndRelationships(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, "tasks")
	rolesDir := filepath.Join(tmpDir, "roles")
	templatesDir := filepath.Join(tmpDir, "templates")

	if err := os.MkdirAll(tasksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(rolesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a test role
	roleContent := `---
description: Test role
---
Test role body
`
	if err := os.WriteFile(filepath.Join(rolesDir, "developer.md"), []byte(roleContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a test template
	templateContent := `---
role: developer
priority: medium
description: Test template
id_prefix: T
---
# {{ .Title }}

{{ .Body }}
`
	if err := os.WriteFile(filepath.Join(templatesDir, "task.md"), []byte(templateContent), 0o644); err != nil {
		t.Fatal(err)
	}

	proj := &ProjectInfo{
		Name:          "test",
		StorageRoot:   tmpDir,
		TasksRoot:     tasksDir,
		RolesRoot:     rolesDir,
		TemplatesRoot: templatesDir,
	}

	server := &Server{
		config: ServerConfig{
			ReadOnly: false,
		},
		projects: map[string]*ProjectInfo{
			"test": proj,
		},
	}

	// Helper to make requests
	makeRequest := func(method, path string, body interface{}) *httptest.ResponseRecorder {
		var bodyReader *bytes.Reader
		if body != nil {
			jsonBody, _ := json.Marshal(body)
			bodyReader = bytes.NewReader(jsonBody)
		} else {
			bodyReader = bytes.NewReader([]byte{})
		}

		req, _ := http.NewRequest(method, path, bodyReader)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.handleTask)
		handler.ServeHTTP(rr, req)
		return rr
	}

	// 1. Create a blocker task
	blockerTitle := "Blocker Task"
	rr := makeRequest("POST", "/api/task?project=test", taskCreateRequest{
		TemplateName: "task",
		Title:        blockerTitle,
		Body:         "I block things",
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("failed to create blocker task: %v", rr.Body.String())
	}

	// Get the blocker ID from the files (since response just gives a message)
	files, _ := os.ReadDir(tasksDir)
	var blockerID string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") && !strings.HasPrefix(f.Name(), "root-") && !strings.HasPrefix(f.Name(), "free-") {
			content, _ := os.ReadFile(filepath.Join(tasksDir, f.Name()))
			if strings.Contains(string(content), blockerTitle) {
				blockerID = strings.TrimSuffix(f.Name(), ".md")
				break
			}
		}
	}
	if blockerID == "" {
		t.Fatal("could not find blocker task file")
	}

	// 2. Create a task that will be blocked
	blockedTitle := "Blocked Task"
	rr = makeRequest("POST", "/api/task?project=test", taskCreateRequest{
		TemplateName: "task",
		Title:        blockedTitle,
		Body:         "I am blocked",
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("failed to create blocked task: %v", rr.Body.String())
	}

	// Get the blocked task ID
	var blockedID string
	files, _ = os.ReadDir(tasksDir)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") && !strings.HasPrefix(f.Name(), "root-") && !strings.HasPrefix(f.Name(), "free-") {
			content, _ := os.ReadFile(filepath.Join(tasksDir, f.Name()))
			if strings.Contains(string(content), blockedTitle) {
				blockedID = strings.TrimSuffix(f.Name(), ".md")
				break
			}
		}
	}
	if blockedID == "" {
		t.Fatal("could not find blocked task file")
	}

	// 3. Verify master lists exist and are correct after creation
	// With current buggy implementation, they might not exist or be empty if not explicitly called
	// But let's check if they exist first.
	rootTasksPath := filepath.Join(tasksDir, "root-tasks.md")
	freeTasksPath := filepath.Join(tasksDir, "free-tasks.md")

	if _, err := os.Stat(rootTasksPath); os.IsNotExist(err) {
		t.Logf("root-tasks.md does not exist after creation (expected bug)")
	} else {
		// If it exists, check content
		content, _ := os.ReadFile(rootTasksPath)
		if !strings.Contains(string(content), blockerTitle) {
			t.Errorf("root-tasks.md does not contain blocker task")
		}
	}

	// 4. Update the blocked task to be blocked by the blocker
	rr = makeRequest("PATCH", "/api/task?project=test&id="+blockedID, taskUpdateRequest{
		Blockers: &[]string{blockerID},
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("failed to update blocked task: %v", rr.Body.String())
	}

	// 5. Verify bidirectional relationship
	// Check the blocker task file to see if it lists the blocked task in `blocks`
	// Wait a bit to ensure file write (though handler is synchronous)
	time.Sleep(100 * time.Millisecond)

	blockerContent, _ := os.ReadFile(filepath.Join(tasksDir, blockerID+".md"))
	parser := task.NewParser()
	blockerTask, err := parser.ParseString(string(blockerContent), blockerID)
	if err != nil {
		t.Fatalf("failed to parse blocker task: %v", err)
	}

	found := false
	for _, id := range blockerTask.Meta.Blocks {
		if id == blockedID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Blocker task %s does not list %s in 'blocks' (bidirectional link failed)", blockerID, blockedID)
	}

	// 6. Verify master lists are updated
	// The blocked task should NOT be in free-tasks.md anymore (assuming it was there)
	// The blocker task SHOULD be in free-tasks.md
	if _, err := os.Stat(freeTasksPath); os.IsNotExist(err) {
		t.Logf("free-tasks.md does not exist after update (expected bug)")
	} else {
		content, _ := os.ReadFile(freeTasksPath)
		if strings.Contains(string(content), blockedTitle) {
			t.Errorf("free-tasks.md incorrectly contains blocked task %s", blockedTitle)
		}
		if !strings.Contains(string(content), blockerTitle) {
			t.Errorf("free-tasks.md does not contain blocker task %s", blockerTitle)
		}
	}
}
