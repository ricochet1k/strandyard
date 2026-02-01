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

	if anchor == "HEAD" && log != nil && taskID != "" {
		if resolved, err := ResolveGitHash(repoPath, "HEAD"); err == nil {
			_ = log.WriteRecurrenceAnchorResolution(taskID, "HEAD", resolved)
		}
	}

	var cmd *exec.Cmd
	var output bytes.Buffer
	var stderr bytes.Buffer

	switch metricType {
	case "commits":
		cmd = exec.Command("git", "rev-list", "--count", fmt.Sprintf("%s..HEAD", anchor))
	case "lines_changed":
		cmd = exec.Command("git", "diff", "--numstat", fmt.Sprintf("%s..HEAD", anchor))
	default:
		return 0, fmt.Errorf("unsupported git metric type: %s", metricType)
	}

	cmd.Dir = repoPath
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	err := cmd.Run()
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

	activeLog, err := activity.Open(baseDir)
	if err != nil {
		return 0, fmt.Errorf("failed to open activity log: %w", err)
	}
	defer activeLog.Close()

	count, err := activeLog.CountCompletionsSince(anchorTime)
	if err != nil {
		return 0, fmt.Errorf("failed to count completions: %w", err)
	}

	return count, nil
}

// ResolveGitHash resolves a git reference to a commit hash.
func ResolveGitHash(repoPath, anchor string) (string, error) {
	cmd := exec.Command("git", "rev-parse", anchor)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// isHeadValid checks if the HEAD reference in a git repository is valid.
// It returns true if HEAD is valid (even detached HEAD pointing to a commit), false otherwise (unborn HEAD).
func isHeadValid(repoPath string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
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
