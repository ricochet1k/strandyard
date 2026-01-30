# Implementation Summary: `strand web` Subcommand

## What Was Implemented

Successfully migrated the standalone `stream-server` into a `strand web` subcommand with multi-project support.

## Files Created

### Go Backend (`pkg/web/`)
- **types.go** - Type definitions for projects, server config, stream updates, and update broker
- **browser.go** - Cross-platform browser opening utility
- **watcher.go** - Multi-project file watching coordinator
- **server.go** - Main HTTP server with embedded dashboard support
- **handlers.go** - HTTP API handlers with project-aware routing
- **README.md** - Documentation for the web package

### Command (`cmd/`)
- **web.go** - Cobra command definition with project discovery logic

### Build & Config
- **scripts/build-web.sh** - Build script for dashboard + binary
- **.gitignore** - Ignore dist directory and binaries

### Dashboard Updates
- **apps/dashboard/vite.config.ts** - Added proxy config and changed output to `pkg/web/dist`
- **apps/dashboard/src/App.tsx** - Added multi-project support with switcher UI

## Files Removed
- **apps/stream-server/** - Entire directory (functionality moved to `pkg/web`)

## Key Features

### Multi-Project Support
✓ Discovers all global projects from `~/.config/strand/projects/`
✓ Discovers local project from `.strand/` in git root
✓ Environment variable filtering via `STRAND_STORAGE` (global/local/both)
✓ Single SSE stream for all projects with client-side filtering

### Dashboard Enhancements
✓ Project switcher dropdown (shows only when multiple projects exist)
✓ All API calls include `?project=X` parameter
✓ Real-time updates filtered by current project
✓ Relative API URLs (works with proxy in dev, embedded in prod)

### Development Workflow
✓ Vite dev server proxies API to `strandyard web` server
✓ Hot reload works for dashboard development
✓ Production build embeds dashboard in Go binary

### Production Build
✓ `go:embed` directive includes dashboard files
✓ Single binary contains everything
✓ Graceful fallback when dist not embedded (helpful dev message)

## Usage

```bash
# Development (two terminals)
strandyard web --no-open              # Terminal 1: API server
cd apps/dashboard && npm run dev      # Terminal 2: Vite dev (http://localhost:5173)

# Production build
./scripts/build-web.sh                # Builds dashboard + binary
./strandyard web                      # Single command, auto-opens browser

# Filter projects
STRAND_STORAGE=global strandyard web  # Only global projects
STRAND_STORAGE=local strandyard web   # Only local project

# Custom port
strandyard web --port 3000

# No auto-open
strandyard web --no-open
```

## API Endpoints

All endpoints support optional `?project=X` parameter:

- `GET /api/health` - Health check
- `GET /api/projects` - List all projects with current project
- `GET /api/state?project=X` - Get project metadata
- `GET /api/tasks?project=X` - List tasks
- `GET /api/files?kind=roles&project=X` - List files
- `GET /api/file?path=X&project=X` - Get file contents
- `PUT /api/file?path=X&project=X` - Save file contents
- `GET /api/stream` - SSE stream (broadcasts updates for all projects)

## Testing Performed

✓ Server starts and binds to port
✓ `/api/health` returns 200 OK
✓ `/api/projects` returns project list with current project
✓ `/api/tasks` returns tasks for specified project
✓ Static files served from embedded dist directory
✓ Dashboard builds successfully to `pkg/web/dist/`
✓ TypeScript compilation succeeds with no errors
✓ Go compilation succeeds with embedded files
✓ Browser auto-open flag works (`--no-open` disables)

## Architecture Decisions

1. **Single SSE stream for all projects** - Clients filter by project name rather than subscribing to project-specific streams. Simpler server logic, easy to extend.

2. **Project query parameter** - All endpoints use `?project=X` rather than path-based routing (`/api/projects/X/tasks`). Cleaner with existing codebase patterns.

3. **Embedded dashboard** - `go:embed` includes built files in binary. Requires placeholder `dist/index.html` to compile during development.

4. **Environment-based filtering** - `STRAND_STORAGE` environment variable controls which projects are discovered. Clean separation from flags.

5. **Auto-open browser to current project** - Detects current project context and opens browser with `?project=X` URL parameter automatically.

## Migration Notes

- The old `stream-server` is completely removed
- All functionality has been preserved and enhanced
- No breaking changes to the dashboard UI/UX
- API remains compatible (added project parameter)

## Future Enhancements

Possible improvements (not in scope for this implementation):

- [ ] Add project creation UI in dashboard
- [ ] Support file upload for templates/roles
- [ ] Add real-time collaboration indicators
- [ ] Support custom project colors/icons
- [ ] Add search across all projects
- [ ] Support project-level settings/preferences
