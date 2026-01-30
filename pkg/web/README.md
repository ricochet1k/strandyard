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
strandyard web

# Start without auto-opening browser
strandyard web --no-open

# Use a custom port
strandyard web --port 3000

# Filter to only global or local projects using environment variables
STRAND_STORAGE=global strandyard web
STRAND_STORAGE=local strandyard web
```

## Development

### Running with Vite Dev Server

For dashboard development with hot reload:

```bash
# Terminal 1: Start the API server
strandyard web --no-open

# Terminal 2: Start Vite dev server (proxies API to :8686)
cd apps/dashboard
npm run dev
# Visit http://localhost:5173
```

### Production Build

To build the complete production binary with embedded dashboard:

```bash
./scripts/build-web.sh
./strandyard web
```

This:
1. Builds the dashboard static files to `pkg/web/dist/`
2. Embeds them in the Go binary using `go:embed`
3. Creates a single `strandyard` binary with everything included

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

## Environment Variables

- `STRAND_ROOT` - Override git root detection (optional)
- `STRAND_STORAGE` - Filter projects: "global", "local", or empty for both (optional)
