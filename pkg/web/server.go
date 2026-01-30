package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed dist
var distFS embed.FS

type Server struct {
	config   ServerConfig
	broker   *updateBroker
	logger   *log.Logger
	projects map[string]*ProjectInfo
}

func Serve(ctx context.Context, cfg ServerConfig) error {
	logger := log.New(os.Stdout, "strand-web ", log.LstdFlags)

	// Build project map
	projectsMap := make(map[string]*ProjectInfo)
	for i := range cfg.Projects {
		projectsMap[cfg.Projects[i].Name] = &cfg.Projects[i]
	}

	server := &Server{
		config:   cfg,
		broker:   newUpdateBroker(),
		logger:   logger,
		projects: projectsMap,
	}

	// Start watchers for all projects
	if err := server.startWatchers(ctx); err != nil {
		return fmt.Errorf("failed to start watchers: %w", err)
	}

	// Setup routes
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/health", server.handleHealth)
	mux.HandleFunc("/api/projects", server.handleProjects)
	mux.HandleFunc("/api/state", server.handleState)
	mux.HandleFunc("/api/tasks", server.handleTasks)
	mux.HandleFunc("/api/files", server.handleFiles)
	mux.HandleFunc("/api/file", server.handleFile)
	mux.HandleFunc("/api/stream", server.handleStream)

	// Static files (embedded dashboard)
	stripped, err := fs.Sub(distFS, "dist")
	if err != nil {
		// If dist doesn't exist, serve a helpful message
		logger.Printf("Warning: embedded dashboard not found (run build-web.sh)")
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Strand Web</title></head>
<body>
<h1>Strand Web API Server</h1>
<p>Dashboard files not embedded. For development, run the Vite dev server:</p>
<pre>cd apps/dashboard && npm run dev</pre>
<p>Then visit <a href="http://localhost:5173">http://localhost:5173</a></p>
<p>API is available at <code>/api/*</code> endpoints.</p>
</body>
</html>`)
		})
	} else {
		mux.Handle("/", http.FileServer(http.FS(stripped)))
	}

	handler := withCORS(mux)

	addr := fmt.Sprintf(":%d", cfg.Port)
	url := fmt.Sprintf("http://localhost:%d", cfg.Port)

	logger.Printf("Watching %d projects", len(cfg.Projects))
	for _, proj := range cfg.Projects {
		logger.Printf("  - %s (%s)", proj.Name, proj.StorageRoot)
	}
	logger.Printf("Server listening on %s", url)

	// Auto-open browser
	if cfg.AutoOpen {
		projectURL := url
		if cfg.CurrentProject != "" {
			projectURL = fmt.Sprintf("%s?project=%s", url, cfg.CurrentProject)
		}
		go func() {
			time.Sleep(500 * time.Millisecond)
			if err := openBrowser(projectURL); err != nil {
				logger.Printf("Failed to open browser: %v", err)
			}
		}()
	}

	// Start server with graceful shutdown
	httpServer := &http.Server{Addr: addr, Handler: handler}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(shutdownCtx)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
