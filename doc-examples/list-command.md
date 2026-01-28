# memmd list Examples

Examples of listing tasks with filters and output formats.

```bash
# List all tasks
memmd list

# Root tasks only
memmd list --scope root

# Free tasks grouped by priority in Markdown
memmd list --scope free --format md --group priority

# Children of a parent task
memmd list --parent E2k7x-metadata-format

# Tasks under a subtree path
memmd list --path tasks/E2k7x-metadata-format

# Filter by role and priority, sorted by created date
memmd list --role developer --priority high --sort created --order desc
```
