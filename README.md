# Strand

A task management system with real-time web dashboard.

## Installation

```bash
go install github.com/ricochet1k/strandyard/cmd/strand@latest
```

This will install the `strand` binary to your `$GOPATH/bin` (or `$GOBIN` if set).

## Quick Start

```bash
# Initialize a new project
strand init

# Start the web dashboard
strand web

# View all tasks
strand list

# Add a new task
strand add "Implement feature X"
```

## Web Dashboard

The `strand web` command starts a web server that watches all your strand projects and provides a real-time dashboard for viewing and editing tasks, roles, and templates.

Features:
- Multi-project support (watches global and local projects)
- Real-time updates via WebSocket
- Project switcher UI
- In-browser file editor
- Auto-opens browser to current project

See [pkg/web/README.md](pkg/web/README.md) for detailed documentation.

## Development

### Building from Source

```bash
# Build the dashboard and binary
./scripts/build-web.sh

# Or build just the binary
go build -o strand ./cmd/strand
```

### Development Workflow

For dashboard development with hot reload:

```bash
# Terminal 1: Start the API server
strand web --no-open

# Terminal 2: Start Vite dev server
cd apps/dashboard
npm run dev
# Visit http://localhost:5173
```

## Environment Variables

- `STRAND_ROOT` - Override git root detection (optional)
- `STRAND_STORAGE` - Filter projects: "global", "local", or empty for both (optional)

## Project Structure

- `cmd/strand/` - Main binary entry point
- `cmd/` - Cobra command implementations
- `pkg/task/` - Task parsing and management
- `pkg/web/` - Web server and dashboard
- `apps/dashboard/` - SolidJS web UI

## License

[Add license information]
