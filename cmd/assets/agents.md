# Strand Agent Instructions

`strand` is a task workflow CLI for humans and agents.
Tasks live as Markdown files with YAML frontmatter, and role docs define how each task should be done.

## Core workflow
- Pull the next task and claim it in one step with `strand next --claim`.
- If you already know the task ID, claim it directly with `strand claim <task-id>`.
- Treat the role document printed by `strand next` as part of the assignment context.
- Finish work with `strand complete <task-id> "report of what was done"`.

## Common commands
- Create work: `strand add <type> "title"` (or `strand add issue "title"`).
- Update metadata/content: `strand edit <task-id>`.
- Reassign role: `strand assign <task-id> <role>`.
- Inspect tasks: `strand list`, `strand search <query>`, `strand show <task-id>`.
- Manage status: `strand claim <task-id>`, `strand cancel <task-id> [reason]`, `strand mark-duplicate <task-id> <duplicate-of>`.
- Explore templates and roles: `strand templates`, `strand roles`.

## Rules for agents
- Prefer `strand` commands over manual task file edits whenever possible.
- Use short IDs or full IDs; both are accepted by task-ID commands.
- Run `strand repair` after manual edits to task files or when task indexes look stale.
- If CLI behavior is missing/awkward, file an issue task via `strand add issue` instead of silently working around it.
