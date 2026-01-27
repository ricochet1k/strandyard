# Epic Example: Scaffold CLI and Core Library

## Summary
Create the initial Go library and CLI scaffolding that will manage the filesystem-backed task database. This epic bootstraps the project with module init, command skeletons, and basic scan/sync functionality.

## Milestones
- Milestone 1: Initialize Go module and scaffold cobra CLI.
- Milestone 2: Implement `scan` to parse tasks and `sync` to update master lists deterministically.
- Milestone 3: Add unit tests and basic CI workflow.

## Tracks
- Track: CLI
  - Owner: dev-team/cli
  - Goals: Provide user-facing commands (`init`, `scan`, `sync`, `add`, `render`).

- Track: Core Library
  - Owner: dev-team/lib
  - Goals: Implement parsing, metadata model, and deterministic list maintenance.

- Track: Docs & Examples
  - Owner: dev-team/docs
  - Goals: Add example tasks, templates, and usage docs for agents and humans.

## Child tasks (examples)
- tasks/scaffold-cli/task.md — Scaffold cobra CLI and basic commands.
- tasks/implement-scan/task.md — Implement scanning and validation of task tree.
- tasks/implement-sync/task.md — Implement deterministic `root-tasks.md` and `free-tasks.md` generation.
- tasks/add-ci/task.md — Add a minimal CI workflow that runs `go test ./...`.

## TODOs
- [ ] (role: architect) Break milestones into child tasks and assign owners.
- [ ] (role: developer) Implement initial commands and library functions.
- [ ] (role: master-reviewer) Coordinate security and reliability reviews for CLI behaviours that touch filesystem or secrets.
