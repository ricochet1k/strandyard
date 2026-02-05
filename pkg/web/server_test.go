package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandleHealth(t *testing.T) {
	server := &Server{}
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleHealth)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status": "ok"}`
	var gotMap, wantMap map[string]string
	json.Unmarshal(rr.Body.Bytes(), &gotMap)
	json.Unmarshal([]byte(expected), &wantMap)

	if gotMap["status"] != wantMap["status"] {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleProjects(t *testing.T) {
	projects := []ProjectInfo{
		{Name: "test-proj", StorageRoot: "/tmp/test"},
	}
	server := &Server{
		config: ServerConfig{
			Projects:       projects,
			CurrentProject: "test-proj",
		},
	}

	req, err := http.NewRequest("GET", "/api/projects", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleProjects)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response["current"] != "test-proj" {
		t.Errorf("expected current project test-proj, got %v", response["current"])
	}

	projectsResp := response["projects"].([]any)
	if len(projectsResp) != 1 {
		t.Errorf("expected 1 project, got %v", len(projectsResp))
	}

	proj := projectsResp[0].(map[string]any)
	if proj["name"] != "test-proj" {
		t.Errorf("expected project name test-proj, got %v", proj["name"])
	}
}

func TestHandleStream(t *testing.T) {
	broker := newUpdateBroker()
	server := &Server{
		broker: broker,
	}

	req, err := http.NewRequest("GET", "/api/stream", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Use a context with timeout to close the connection
	ctx, cancel := context.WithCancel(context.Background())
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Start the handler in a goroutine because it's a blocking loop
	done := make(chan bool)
	go func() {
		server.handleStream(rr, req)
		done <- true
	}()

	// Wait for "ready" event
	time.Sleep(100 * time.Millisecond)

	// Broadcast an update
	update := StreamUpdate{
		Event:   "test-event",
		Project: "test-proj",
		Path:    "tasks/T1-task.md",
	}
	broker.broadcast(update)

	// Wait for update to be processed
	time.Sleep(100 * time.Millisecond)

	// Cancel context to stop the stream
	cancel()
	<-done

	// Verify headers
	if rr.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected text/event-stream, got %v", rr.Header().Get("Content-Type"))
	}

	// Verify body contains the update
	body := rr.Body.String()
	if !strings.Contains(body, "event: ready") {
		t.Errorf("expected ready event, not found in: %s", body)
	}
	if !strings.Contains(body, "event: task") {
		t.Errorf("expected task event, not found in: %s", body)
	}
	if !strings.Contains(body, "test-event") {
		t.Errorf("expected test-event data, not found in: %s", body)
	}
}

func TestWithAuth(t *testing.T) {
	server := &Server{
		config: ServerConfig{
			AuthToken: "secret-token",
		},
	}

	handler := server.withAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test 1: No auth
	req, _ := http.NewRequest("GET", "/api/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %v", rr.Code)
	}

	// Test 2: Valid Bearer token
	req, _ = http.NewRequest("GET", "/api/health", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 with Bearer, got %v", rr.Code)
	}

	// Test 3: Valid query token
	req, _ = http.NewRequest("GET", "/api/health?token=secret-token", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 with query token, got %v", rr.Code)
	}

	// Test 4: Invalid token
	req, _ = http.NewRequest("GET", "/api/health", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with wrong token, got %v", rr.Code)
	}
}
