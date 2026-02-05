# Anchor Help Text and Examples — Design Alternatives

## Summary
Propose concise `--help` text and CLI.md documentation for `strand add --every`, covering time units (days, weeks, months), git units (commits, lines_changed), and tasks_completed metrics. Focus on clarity for users who skim help output while maintaining precision.

## Context
- `design-docs/recurrence-anchor-flags-alternatives.md` (Alternative D adopted)
- `design-docs/recurrence-anchor-error-messages.md`
- `CLI.md` (recurring add section — outdated, uses `--unit`/`--anchor` format)
- `cmd/strand` (current `--help` output: `--every strings recurrence rule (e.g., "10 days", "50 commits from HEAD")`)

## Project Principles
- Help text should be scannable (short summary + clear examples).
- Anchor type must be unambiguous per metric to prevent user errors.
- Examples should use deterministic strings for tests and automation.
- Defaults should be explicit (not implied) to reduce cognitive load.

## Per-Unit Anchor Mappings (Authoritative)

| Metric | Expected Anchor Type | Default Anchor | Format Example | Full Example |
|--------|---------------------|----------------|----------------|--------------|
| `days` | date | now (implicit) | ISO 8601 or human-friendly date | `--every "10 days from Jan 28 2026 09:00 UTC"` |
| `weeks` | date | now (implicit) | ISO 8601 or human-friendly date | `--every "2 weeks from Jan 28 2026 09:00 UTC"` |
| `months` | date | now (implicit) | ISO 8601 or human-friendly date | `--every "3 months from Jan 28 2026 09:00 UTC"` |
| `commits` | commit | HEAD | 40-char hex or `HEAD` | `--every "50 commits from HEAD"` |
| `lines_changed` | commit | HEAD | 40-char hex or `HEAD` | `--every "500 lines_changed from HEAD"` |
| `tasks_completed` | date or task ID | now (implicit) | ISO 8601, human date, or task ID | `--every "20 tasks_completed from T1a1a"` |

**Notes**:
- Date anchors accept ISO 8601 (`2026-01-28T09:00:00Z`) or human-friendly formats parsed by `when` library (e.g., `Jan 28 2026 09:00 UTC`).
- Commit anchors require exact 40-character hex string or the literal `HEAD`.
- `from <anchor>` is optional; omission uses the default anchor.
- Use `from now` to trigger immediate run then recur; use `after now` to schedule first run at next interval (when implemented).

## Alternatives

### Alternative A — Compact one-line help with examples

**`--help` text**:
```
--every strings   recurrence rule: "<amount> <metric> [from <anchor>]" (repeatable)
                  metrics: days, weeks, months, commits, lines_changed, tasks_completed
                  examples: "10 days", "50 commits from HEAD", "20 tasks_completed from T1a1a"
```

**CLI.md section**:
```markdown
### `add` — Create recurring task definitions

Creates a recurring task definition using the `--every` flag.

```bash
strand add [title] --every "<amount> <metric> [from <anchor>]" [flags]
```

**Metrics and anchors**:

| Metric | Anchor Type | Default Anchor | Example |
|--------|-------------|----------------|---------|
| `days` | date | now (implicit) | `--every "10 days"` |
| `weeks` | date | now (implicit) | `--every "2 weeks from Jan 28 2026 09:00 UTC"` |
| `months` | date | now (implicit) | `--every "3 months"` |
| `commits` | commit | HEAD | `--every "50 commits from HEAD"` |
| `lines_changed` | commit | HEAD | `--every "500 lines_changed"` |
| `tasks_completed` | date or task ID | now (implicit) | `--every "20 tasks_completed from T1a1a"` |

**Multiple metrics**: repeat `--every` to combine triggers (when implemented).

**Example**:
```bash
strand add "Weekly sync" --every "1 week"
```

**Definition metadata** (excerpt):
```markdown
---
recurrence_triggers:
  - metric: days
    interval: 1
    anchor: now
---
```

**Flags**:
- `--every`: repeatable recurrence rule (`<amount> <metric> [from <anchor>]`)
- `-t, --title`: definition title (optional, can be passed as positional argument)
- `-r, --role`: role responsible for generated tasks (defaults from template)
- `-p, --parent`: parent task ID
- `--priority`: priority: high, medium, or low
- `--blocker`: blocker task ID(s)
- `--no-repair`: skip repair and master list updates
```

**Pros**:
- Minimal screen space; fits in typical `--help` width (80 chars).
- Examples cover all metric types with defaults and explicit anchors.
- Table in CLI.md is scannable for reference.

**Cons**:
- Assumes users read the full help line; long metrics list may wrap.
- No explicit "time vs git" grouping; users must infer from table.

**Risks**:
- Users on narrow terminals may see line breaks that reduce clarity.

**Effort**: Small (help text update + CLI.md section rewrite).

---

### Alternative B — Grouped help with section headers

**`--help` text**:
```
--every strings   recurrence rule: "<amount> <metric> [from <anchor>]" (repeatable)
                  time metrics (days, weeks, months, tasks_completed):
                    "10 days", "2 weeks from Jan 28 2026 09:00 UTC", "20 tasks_completed"
                  git metrics (commits, lines_changed):
                    "50 commits from HEAD", "500 lines_changed"
```

**CLI.md section** (similar to Alternative A, with grouped sections):
```markdown
**Time-based metrics**:
- `days`: date anchor, default `now`
- `weeks`: date anchor, default `now`
- `months`: date anchor, default `now`
- `tasks_completed`: date or task ID anchor, default `now`

**Git-based metrics**:
- `commits`: commit anchor (hash or `HEAD`), default `HEAD`
- `lines_changed`: commit anchor (hash or `HEAD`), default `HEAD`
```

**Pros**:
- Explicit grouping helps users identify anchor types quickly.
- Scannable for "time" vs "git" intent.

**Cons**:
- Longer help text; may require scrolling on small screens.
- Grouping adds cognitive load if users only need one example.

**Risks**:
- Formatting may look dense on narrow terminals.

**Effort**: Small to Medium (help text update + CLI.md section rewrite with groups).

---

### Alternative C — Minimal help with reference URL

**`--help` text**:
```
--every strings   recurrence rule: "<amount> <metric> [from <anchor>]" (repeatable)
                  metrics: days, weeks, months, commits, lines_changed, tasks_completed
                  see CLI.md for examples and anchor formats
```

**CLI.md section**: Comprehensive with full examples (same as Alternative A).

**Pros**:
- Shortest help text; maximizes scannability.
- Encourages users to read documentation for details.

**Cons**:
- Requires users to open CLI.md to see examples.
- Breaks flow for quick-reference users.

**Risks**:
- Users may not read CLI.md; more likely to make formatting errors.

**Effort**: Small (help text update + CLI.md section rewrite).

---

## Recommendation

**Adopt Alternative A**. The compact one-line help with examples balances scannability and completeness. The examples cover all metric types, show both default and explicit anchor usage, and fit within typical terminal widths. The CLI.md table provides a quick reference for anchor types and defaults.

**Rationale**:
- Users who skim help see: `<amount> <metric> [from <anchor>]` + three examples covering time, git, and tasks_completed.
- CLI.md provides a comprehensive reference table with scannable columns.
- Alternative B adds grouping without significant benefit; Alternative A's examples already distinguish time vs git.
- Alternative C pushes too much to documentation; users should not need to open CLI.md for basic usage.

## Review Requests
- Request review from: `master-reviewer`, `reviewer-usability`, `reviewer-reliability`.
