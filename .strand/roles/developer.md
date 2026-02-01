---
description: "Implements tasks, writes code, and produces working software."
---

# Developer

## Role
Developer (human or AI) — implements tasks, writes code, and produces working software.

## Responsibilities
- Implement tasks assigned by the Architect
- Write clean, maintainable code following project conventions
- Add tests for new functionality
- Document code and update relevant documentation
- Fix bugs and address issues
- Ensure code passes validation and tests before marking tasks complete

## Deliverables
- Working code that meets acceptance criteria
- Tests covering the implemented functionality
- Updated documentation as needed
- Code that passes `go build` and `go test`

## Workflow
1. Read the assigned task and understand acceptance criteria
2. Implement the functionality described in the task
3. Write tests to verify the implementation
4. Run repair: `go build ./...`, `go test ./...`, `strand repair`
5. Update task status and mark as completed when done

## Constraints
- Follow existing code patterns and conventions in the codebase
- Ensure all changes are backward compatible unless explicitly noted
- Do not modify task metadata manually - use CLI commands
