# Strand Web Dashboard

The `strand web` command provides a web-based dashboard for viewing and editing strand tasks, roles, and templates across all your projects.

## Features

- **Multi-project support**: Automatically discovers and watches all global and local strand projects
- **Real-time updates**: WebSocket-based live updates when files change
- **Project switcher**: Easy switching between projects in the UI
- **File editor**: Edit tasks, roles, and templates directly in the browser
- **Auto-open browser**: Opens your default browser to the current project automatically

## Usage

```bash
# Start the web server (auto-opens browser)
strand web

# Start without auto-opening browser
strand web --no-open

# Use a custom port
strand web --port 3000

# Require authentication token
strand web --auth-token your-secret-token

# Enable read-only mode (no file writes)
strand web --read-only

# Restrict CORS origins
strand web --allowed-origins "http://localhost:5173,http://localhost:3000"

# Filter to only global or local projects using environment variables
STRAND_STORAGE=global strand web
STRAND_STORAGE=local strand web
```

## Development

### Running with Vite Dev Server

For dashboard development with hot reload:

```bash
# Terminal 1: Start the API server
strand web --no-open

# Terminal 2: Start Vite dev server (proxies API to :8686)
cd apps/dashboard
npm run dev
# Visit http://localhost:5173
```

### Production Build

To build the complete production binary with embedded dashboard:

```bash
./scripts/build-web.sh
./strand web
```

This:
1. Builds the dashboard static files to `pkg/web/dist/`
2. Embeds them in the Go binary using `go:embed`
3. Creates a single `strand` binary with everything included

## Architecture

### Multi-Project Watching

The web server discovers projects from:
- **Global projects**: `~/.config/strand/projects/*/`
- **Local project**: `.strand/` in the current git repository

All projects are watched simultaneously, and the dashboard receives updates for all of them via a single SSE stream. The client filters updates by the currently selected project.

### API Endpoints

- `GET /api/health` - Health check
- `GET /api/projects` - List all available projects
- `GET /api/state?project=X` - Get project metadata
- `GET /api/tasks?project=X` - List all tasks for a project
- `GET /api/files?kind=roles&project=X` - List files (roles/templates)
- `GET /api/file?path=X&project=X` - Get file contents
- `PUT /api/file?path=X&project=X` - Save file contents
- `GET /api/stream` - Server-sent events stream for real-time updates

All endpoints (except `/api/projects` and `/api/health`) accept an optional `?project=X` query parameter to scope the request to a specific project.

## Security

### Authentication

The dashboard supports optional token-based authentication. When `--auth-token` is set, all API requests must include a valid token:

- **Bearer token in Authorization header**:
  ```
  Authorization: Bearer your-secret-token
  ```

- **Token in query string** (for convenience):
  ```
  /api/tasks?project=local&token=your-secret-token
  ```

**Note**: By default, no authentication is required. Always use `--auth-token` in production or shared environments.

### Read-Only Mode

Use `--read-only` to disable all file write operations. This prevents modifications while still allowing read access to tasks, roles, and templates.

### CORS

By default, CORS allows all origins (`*`). Use `--allowed-origins` to restrict access to specific origins:

```bash
strand web --allowed-origins "http://localhost:5173,https://yourdomain.com"
```

### File Write Safeguards

File writes are protected by:
- Path validation: Files can only be written to `tasks/`, `roles/`, or `templates/` directories
- Content validation: Empty content is rejected
- Read-only mode: When enabled, all write operations return HTTP 403 Forbidden

## Environment Variables

- `STRAND_ROOT` - Override git root detection (optional)
- `STRAND_STORAGE` - Filter projects: "global", "local", or empty for both (optional)
