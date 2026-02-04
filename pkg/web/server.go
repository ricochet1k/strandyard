package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
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
	mux.HandleFunc("/api/health", server.withAuth(server.handleHealth))
	mux.HandleFunc("/api/projects", server.withAuth(server.handleProjects))
	mux.HandleFunc("/api/state", server.withAuth(server.handleState))
	mux.HandleFunc("/api/tasks", server.withAuth(server.handleTasks))
	mux.HandleFunc("/api/task", server.withAuth(server.handleTask))
	mux.HandleFunc("/api/files", server.withAuth(server.handleFiles))
	mux.HandleFunc("/api/file", server.withAuth(server.handleFile))
	mux.HandleFunc("/api/stream", server.withAuth(server.handleStream))

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

	handler := server.withCORS(mux)

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

func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.config.AuthToken != "" {
			authHeader := r.Header.Get("Authorization")
			token := r.URL.Query().Get("token")

			valid := false
			if strings.HasPrefix(authHeader, "Bearer ") {
				valid = authHeader[7:] == s.config.AuthToken
			} else if token != "" {
				valid = token == s.config.AuthToken
			}

			if !valid {
				w.Header().Set("WWW-Authenticate", "Bearer")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigin := "*"
		for _, allowed := range s.config.AllowedOrigins {
			if allowed == "*" || allowed == origin {
				allowedOrigin = allowed
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
