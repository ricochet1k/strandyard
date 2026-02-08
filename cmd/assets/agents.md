# Memmd Agent Instructions

## What is strand
strand is a cli task management tool intended for helping to direct AI agents to follow a specific workflow.
This workflow is defined as instructions in role documents and pre-filled todo lists in task templates.

## Purpose
These instructions define how agents should use strand to manage tasks.

## Core rules
- Use `strand next --claim` to select and claim work; respect role opt-in or ignore behavior.
- When asked to work on the next thing, run `strand next --claim`, follow the returned instructions, and report a brief task summary.
- Treat the role description returned by `strand next --claim` as canonical for how to execute the task.
- When a task is done (including planning-only), run `strand complete <task-id> "report of what was done"`.
- If blocked, record blockers with `strand block`.
- Use `strand add` for new tasks or issues; avoid ad-hoc task creation outside strand.
- Get the list of roles and task templates from `strand roles` and `strand templates`, add them to AGENTS.md and keep that part up to date as needed.
- If bugs or usability or missing features are discovered in attempting to use `strand`, file issues
  directly on the "strand" project with a command like `strand add issue --project strand "Issue title" <<EOF\n # Detailed markdown description \nEOF`
