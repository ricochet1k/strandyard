# Date Parsing Options for --every Anchors

## Summary
Evaluate Go date parsing libraries for `strand recurring add --every` anchor inputs. The goal is to accept human-friendly dates while keeping parsing deterministic for CLI usage.

## Context
- design-docs/recurrence-anchor-flags-alternatives.md
- design-docs/recurrence-anchor-error-messages-alternatives.md
- design-docs/recurrence-metrics.md
- CLI.md (recurring add section)

## Requirements
- Deterministic results for identical input strings.
- Support for human-friendly dates (example: "Jan 28 2026 09:00").
- Optional support for relative expressions (example: "in 2 days", "next Tuesday").
- Explicit handling of time zone and locale defaults.
- Go module compatibility and a permissive license.

## Candidates

### Option A: `github.com/olebedev/when` (Apache-2.0)
- Description: Natural language date/time parser with pluggable rules. Parse is performed relative to a provided base time.
- Pros:
  - Rule-based, no locale auto-detection; behavior is controllable by selecting rule sets.
  - Supports relative expressions in English (for example: "next Tuesday", "in 2 days").
  - Deterministic when base time and location are fixed by the caller.
- Cons:
  - Limited locale support unless rules are explicitly added.
  - Partial matches are possible; needs extra validation to ensure full-string parsing.
  - Time zone parsing is limited; typically uses the base time location when unspecified.

### Option B: `github.com/araddon/dateparse` (MIT)
- Description: Fast parser for many absolute date formats without specifying a layout.
- Pros:
  - Broad absolute format coverage and good performance.
  - Provides `ParseStrict` to reject ambiguous numeric formats.
  - Simple dependency footprint.
- Cons:
  - No support for relative expressions.
  - Ambiguous numeric formats default to MM/DD unless `ParseStrict` is used.
  - Time zone defaults to `time.Local` when missing, which can be non-deterministic across environments unless controlled.

### Option C: `github.com/markusmobius/go-dateparser` (BSD-3)
- Description: Port of Python dateparser with locale-aware parsing, relative expressions, and timezone handling.
- Pros:
  - Supports relative and absolute dates, including natural language and timestamps.
  - Configurable `CurrentTime`, `DefaultTimezone`, `DateOrder`, and language preferences.
  - Works in many locales (can be restricted to a subset).
- Cons:
  - Large dependency and dataset footprint; heavier runtime cost.
  - Language detection and incomplete date defaults can introduce surprising results unless tightly configured.
  - Regex-heavy implementation; performance is lower than simple parsers.

## Determinism and Locale Notes
- Option A is deterministic if the caller sets a fixed base time and location and rejects partial matches.
- Option B is deterministic only if missing time zones are disallowed or `time.Local` is fixed (for example, set to UTC in the process).
- Option C is deterministic if the configuration pins `CurrentTime`, `DefaultTimezone`, `DefaultLanguages`, `DateOrder`, and uses strict parsing.

## Recommendation
Recommendation: Prefer Option A (`github.com/olebedev/when`) for `--every` anchors, with strict input validation and a fixed base time in UTC. This is the smallest dependency that can parse relative expressions while keeping behavior explicit and controllable by rule selection.

Constraints for CLI usage:
- Parse using a fixed base time in UTC (for example, `time.Now().UTC()` at command execution) and document this explicitly.
- Accept only full-string matches; if the matched substring does not cover the entire anchor input, return a parse error.
- Require unambiguous date formats for absolute anchors (prefer month names or ISO 8601 to avoid numeric ambiguity).
- Treat missing timezone as UTC; reject non-UTC offsets unless explicitly supported by the rule set.
- Limit rules to English (`en` + `common`) unless the CLI adds an explicit locale flag.

If broader locale support becomes a requirement, revisit Option C with a strict configuration that pins locale and time defaults.

## Decision
The project adopts **Option A (`github.com/olebedev/when`)**.

## Constraints for CLI usage
- Parse using a fixed base time in UTC (execution time).
- Require full-string matches; reject partial matches.
- Treat missing timezones as UTC.
- Limit rules to English (`en` + `common`).
