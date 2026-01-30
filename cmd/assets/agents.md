# Memmd Agent Instructions

## What is memmd
memmd is a cli task management tool intended for helping to direct AI agents to follow a specific workflow.
This workflow is defined as instructions in role documents and pre-filled todo lists in task templates.

## Purpose
These instructions define how agents should use memmd to manage tasks.

## Core rules
- Use `memmd next` to select work; respect role opt-in or ignore behavior.
- When asked to work on the next thing, run `memmd next`, follow the returned instructions, and report a brief task summary.
- Treat the role description returned by `memmd next` as canonical for how to execute the task.
- When a task is done (including planning-only), run `memmd complete <task-id>`.
- If blocked, record blockers with `memmd block`.
- Use `memmd add` for new tasks or issues; avoid ad-hoc task creation outside memmd.
- Get the list of roles and task templates from `memmd roles` and `memmd templates`, add them to AGENTS.md and keep that part up to date as needed.
- If bugs or usability or missing features are discovered in attempting to use `memmd`, file issues
  directly on the "memmd" project with a command like `memmd add issue --project memmd "Issue title" <<EOF\n # Detailed markdown description \nEOF`