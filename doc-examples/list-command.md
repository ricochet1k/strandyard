# StrandYard list Examples

Examples of listing tasks with filters and output formats.

```bash
# List all tasks
strand list

# Root tasks only
strand list --scope root

# Free tasks grouped by priority in Markdown
strand list --scope free --format md --group priority

# Children of a parent task
strand list --parent E2k7x-metadata-format

# Tasks under a subtree path
strand list --path tasks/E2k7x-metadata-format

# Filter by role and priority, sorted by created date
strand list --role developer --priority high --sort created --order desc
```
