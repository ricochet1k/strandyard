package task

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ricochet1k/strandyard/pkg/activity"
)

// EvaluateGitMetric evaluates a git-based recurrence metric (commits or lines_changed).
// It returns the metric value and an error if one occurs.
// It treats unborn/invalid HEAD as a no-op (returns 0, nil).
// If a log and taskID are provided, it logs the resolution of "HEAD" anchors.
func EvaluateGitMetric(repoPath, metricType, anchor string, taskID string, log *activity.Log) (int, error) {
	// Check if HEAD is valid
	if !isHeadValid(repoPath) {
		return 0, nil // Unborn or invalid HEAD, treat as no-op
	}

	resolvedAnchor, err := ResolveGitHash(repoPath, anchor)
	if err != nil {
		// Treat unknown revision as no-op
		return 0, nil
	}

	if anchor == "HEAD" && log != nil && taskID != "" {
		_ = log.WriteRecurrenceAnchorResolution(taskID, "HEAD", resolvedAnchor)
	}

	var cmd *exec.Cmd
	var output bytes.Buffer
	var stderr bytes.Buffer

	switch metricType {
	case "commits":
		cmd = exec.Command("git", "rev-list", "--count", fmt.Sprintf("%s..HEAD", resolvedAnchor))
	case "lines_changed":
		cmd = exec.Command("git", "diff", "--numstat", fmt.Sprintf("%s..HEAD", resolvedAnchor))
	default:
		return 0, fmt.Errorf("unsupported git metric type: %s", metricType)
	}

	cmd.Dir = repoPath
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		errStr := stderr.String()
		// Treat "unknown revision" or "ambiguous argument" errors as no-op for recurrence metrics.
		if strings.Contains(errStr, "unknown revision") || strings.Contains(errStr, "ambiguous argument") {
			return 0, nil
		}
		return 0, fmt.Errorf("git command failed: %v, stderr: %s", err, errStr)
	}

	rawOutput := strings.TrimSpace(output.String())

	switch metricType {
	case "commits":
		count, err := strconv.Atoi(rawOutput)
		if err != nil {
			return 0, fmt.Errorf("failed to parse commit count: %w", err)
		}
		return count, nil
	case "lines_changed":
		// numstat output format: "additions\tdeletions\tfile"
		lines := strings.Split(rawOutput, "\n")
		totalLines := 0
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.Split(line, "\t")
			if len(parts) >= 2 {
				additions, err := strconv.Atoi(parts[0])
				if err != nil {
					return 0, fmt.Errorf("failed to parse additions from numstat: %w", err)
				}
				deletions, err := strconv.Atoi(parts[1])
				if err != nil {
					return 0, fmt.Errorf("failed to parse deletions from numstat: %w", err)
				}
				totalLines += additions + deletions
			}
		}
		return totalLines, nil
	}

	return 0, fmt.Errorf("unsupported metric type after execution: %s", metricType)
}

// EvaluateTasksCompletedMetric evaluates a tasks_completed recurrence metric.
// It queries the activity log to count task completions since the given anchor time.
// If a log and taskID are provided, it logs the resolution of "now" anchors.
func EvaluateTasksCompletedMetric(baseDir, anchor string, taskID string, log *activity.Log) (int, error) {
	var anchorTime time.Time
	var err error

	if anchor == "now" || anchor == "" {
		anchorTime = time.Now().UTC()
		if log != nil && taskID != "" {
			_ = log.WriteRecurrenceAnchorResolution(taskID, anchor, anchorTime.Format("Jan 2 2006 15:04 MST"))
		}
	} else {
		anchorTime, err = time.Parse("Jan 2 2006 15:04 MST", anchor)
		if err != nil {
			return 0, fmt.Errorf("invalid date anchor format: %w", err)
		}
	}

	activeLog := log
	if activeLog == nil {
		var err error
		activeLog, err = activity.Open(baseDir)
		if err != nil {
			return 0, fmt.Errorf("failed to open activity log: %w", err)
		}
		defer activeLog.Close()
	}

	count, err := activeLog.CountCompletionsSince(anchorTime)
	if err != nil {
		return 0, fmt.Errorf("failed to count completions: %w", err)
	}

	return count, nil
}

// ResolveGitHash resolves a git reference to a commit hash.
func ResolveGitHash(repoPath, anchor string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", "--end-of-options", anchor+"^{commit}")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// ValidateAnchor validates a recurrence anchor for a given metric.
func ValidateAnchor(metric, anchor, repoPath string, tasks map[string]*Task) error {
	if anchor == "" {
		return nil
	}

	switch metric {
	case "days", "weeks", "months":
		return ValidateDateAnchor(anchor)
	case "commits", "lines_changed":
		if anchor == "HEAD" {
			return nil
		}
		return ValidateCommitAnchor(repoPath, anchor)
	case "tasks_completed":
		if anchor == "now" {
			return nil
		}
		// If it's a date, validate as date
		if err := ValidateDateAnchor(anchor); err == nil {
			return nil
		}
		// Otherwise validate as task ID
		return ValidateTaskAnchor(tasks, anchor)
	default:
		return fmt.Errorf("unsupported metric: %s", metric)
	}
}

// ValidateDateAnchor checks if the anchor is a valid date (ISO 8601 or human-friendly).
func ValidateDateAnchor(anchor string) error {
	if anchor == "now" {
		return nil
	}
	// Try ISO 8601
	if _, err := time.Parse(time.RFC3339, anchor); err == nil {
		return nil
	}
	// Try human-friendly format
	if _, err := time.Parse("Jan 2 2006 15:04 MST", anchor); err == nil {
		return nil
	}
	return fmt.Errorf("invalid date format: %s (expected ISO 8601 or \"Jan 2 2006 15:04 MST\")", anchor)
}

// ValidateCommitAnchor checks if the anchor exists in the git repository.
func ValidateCommitAnchor(repoPath, anchor string) error {
	if anchor == "HEAD" {
		return nil
	}
	_, err := ResolveGitHash(repoPath, anchor)
	if err != nil {
		return fmt.Errorf("invalid commit anchor: %s (not found in repository)", anchor)
	}
	return nil
}

// ValidateTaskAnchor checks if the anchor is a valid task ID in the project.
func ValidateTaskAnchor(tasks map[string]*Task, anchor string) error {
	if tasks == nil {
		return nil // Can't validate if tasks are not loaded
	}
	_, err := ResolveTaskID(tasks, anchor)
	if err != nil {
		return fmt.Errorf("invalid task anchor: %s (task not found)", anchor)
	}
	return nil
}

// UpdateAnchor calculates the next anchor for a given metric and interval.
func UpdateAnchor(repoPath, baseDir, metric, currentAnchor string, interval int) (string, error) {
	switch metric {
	case "days", "weeks", "months":
		var anchorTime time.Time
		var err error
		if currentAnchor == "now" || currentAnchor == "" {
			anchorTime = time.Now().UTC()
		} else {
			anchorTime, err = time.Parse(time.RFC3339, currentAnchor)
			if err != nil {
				anchorTime, err = time.Parse("Jan 2 2006 15:04 MST", currentAnchor)
				if err != nil {
					return "", fmt.Errorf("invalid date anchor format: %w", err)
				}
			}
		}

		var nextDue time.Time
		switch metric {
		case "days":
			nextDue = anchorTime.AddDate(0, 0, interval)
		case "weeks":
			nextDue = anchorTime.AddDate(0, 0, interval*7)
		case "months":
			nextDue = anchorTime.AddDate(0, interval, 0)
		}

		// If nextDue is still in the past, skip ahead to avoid flooding
		now := time.Now().UTC()
		if nextDue.Before(now) {
			for nextDue.Before(now) {
				switch metric {
				case "days":
					nextDue = nextDue.AddDate(0, 0, interval)
				case "weeks":
					nextDue = nextDue.AddDate(0, 0, interval*7)
				case "months":
					nextDue = nextDue.AddDate(0, interval, 0)
				}
			}
		}

		return nextDue.Format(time.RFC3339), nil

	case "commits":
		newAnchor, err := GetCommitAtOffset(repoPath, currentAnchor, interval)
		if err != nil {
			// If anchor is missing (e.g. force push), fallback to HEAD
			if strings.Contains(err.Error(), "failed to list commits") ||
				strings.Contains(err.Error(), "unknown revision") ||
				strings.Contains(err.Error(), "failed to resolve anchor") {
				return ResolveGitHash(repoPath, "HEAD")
			}
			return "", err
		}
		return newAnchor, nil

	case "lines_changed":
		newAnchor, err := GetCommitCrossingLinesThreshold(repoPath, currentAnchor, interval)
		if err != nil {
			// If anchor is missing, fallback to HEAD
			if strings.Contains(err.Error(), "failed to list commits") ||
				strings.Contains(err.Error(), "unknown revision") ||
				strings.Contains(err.Error(), "failed to resolve anchor") {
				return ResolveGitHash(repoPath, "HEAD")
			}
			return "", err
		}
		return newAnchor, nil

	case "tasks_completed":
		var anchorTime time.Time
		var err error
		if currentAnchor == "now" || currentAnchor == "" {
			anchorTime = time.Now().UTC()
		} else {
			// Try task ID first
			// (Need tasks map for this, but for now we assume it's a date if it's not "now" or "T...")
			if strings.HasPrefix(currentAnchor, "T") || strings.HasPrefix(currentAnchor, "E") || strings.HasPrefix(currentAnchor, "I") {
				// TODO: Resolve task ID to completion time.
				// For now, let's just use date parsing as fallback.
				return "", fmt.Errorf("task ID anchors not yet supported in UpdateAnchor")
			}

			anchorTime, err = time.Parse("Jan 2 2006 15:04 MST", currentAnchor)
			if err != nil {
				anchorTime, err = time.Parse(time.RFC3339, currentAnchor)
				if err != nil {
					return "", fmt.Errorf("invalid date anchor format: %w", err)
				}
			}
		}

		log, err := activity.Open(baseDir)
		if err != nil {
			return "", err
		}
		defer log.Close()

		nextAnchorTime, err := log.GetCompletionTimestampAtOffset(anchorTime, interval)
		if err != nil {
			return "", err
		}

		return nextAnchorTime.Format(time.RFC3339), nil

	default:
		return "", fmt.Errorf("unsupported metric: %s", metric)
	}
}

// GetCommitAtOffset returns the commit hash that is exactly 'offset' commits after the 'anchor'.
func GetCommitAtOffset(repoPath, anchor string, offset int) (string, error) {
	resolvedAnchor, err := ResolveGitHash(repoPath, anchor)
	if err != nil {
		return "", fmt.Errorf("failed to resolve anchor %q: %w", anchor, err)
	}

	if offset <= 0 {
		return resolvedAnchor, nil
	}

	cmd := exec.Command("git", "rev-list", "--reverse", fmt.Sprintf("%s..HEAD", resolvedAnchor))
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list commits: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if offset > len(lines) {
		return "", fmt.Errorf("offset %d exceeds available commits (%d)", offset, len(lines))
	}

	return lines[offset-1], nil
}

// GetCommitCrossingLinesThreshold returns the first commit hash where the cumulative
// lines changed since 'anchor' meets or exceeds the 'threshold'.
func GetCommitCrossingLinesThreshold(repoPath, anchor string, threshold int) (string, error) {
	resolvedAnchor, err := ResolveGitHash(repoPath, anchor)
	if err != nil {
		return "", fmt.Errorf("failed to resolve anchor %q: %w", anchor, err)
	}

	if threshold <= 0 {
		return resolvedAnchor, nil
	}

	cmd := exec.Command("git", "rev-list", "--reverse", fmt.Sprintf("%s..HEAD", resolvedAnchor))
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list commits: %w", err)
	}

	commits := strings.Split(strings.TrimSpace(string(out)), "\n")
	cumulativeLines := 0
	for _, commit := range commits {
		if commit == "" {
			continue
		}

		// Get lines changed in this specific commit
		// We use git show --numstat <commit> and parse it
		lines, err := getLinesChangedInCommit(repoPath, commit)
		if err != nil {
			return "", fmt.Errorf("failed to get lines for commit %s: %w", commit, err)
		}

		cumulativeLines += lines
		if cumulativeLines >= threshold {
			return commit, nil
		}
	}

	return "", fmt.Errorf("threshold %d not reached (cumulative lines: %d)", threshold, cumulativeLines)
}

func getLinesChangedInCommit(repoPath, commit string) (int, error) {
	// commit should already be a resolved hash, but we use ResolveGitHash just in case
	// to ensure it doesn't look like a flag.
	resolvedCommit, err := ResolveGitHash(repoPath, commit)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve commit %q: %w", commit, err)
	}

	cmd := exec.Command("git", "show", "--numstat", "--format=", resolvedCommit)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	rawOutput := strings.TrimSpace(string(out))
	lines := strings.Split(rawOutput, "\n")
	totalLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			additions, err := strconv.Atoi(parts[0])
			if err != nil {
				continue // Skip binary files or other non-numeric entries
			}
			deletions, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			totalLines += additions + deletions
		}
	}
	return totalLines, nil
}

// isHeadValid checks if the HEAD reference in a git repository is valid.
// It returns true if HEAD is valid (even detached HEAD pointing to a commit), false otherwise (unborn HEAD).
func isHeadValid(repoPath string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", "--end-of-options", "HEAD")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		// git rev-parse --verify HEAD returns non-zero exit code if HEAD is invalid/unborn
		return false
	}
	return true
}

// tempGitRepo creates a temporary git repository for testing.
func tempGitRepo() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "git-test-repo-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to git init: %w", err)
	}

	return tmpDir, cleanup, nil
}
