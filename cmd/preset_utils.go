package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func applyPreset(w io.Writer, baseDir, preset string, dirs []string) error {
	if preset == "" {
		return fmt.Errorf("preset path or URL cannot be empty")
	}

	sourceDir := preset
	cleanup := func() {}

	// Check if preset is a local directory
	if info, err := os.Stat(preset); err != nil || !info.IsDir() {
		// Not a local directory - try git clone
		if strings.TrimSpace(preset) == "" {
			return fmt.Errorf("preset cannot be empty or whitespace")
		}

		fmt.Fprintf(w, "Cloning preset from %s...\n", preset)
		tempDir, err := os.MkdirTemp("", "strand-preset-")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}
		cleanup = func() {
			_ = os.RemoveAll(tempDir)
		}
		cmd := exec.Command("git", "clone", "--depth", "1", "--", preset, tempDir)
		if output, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			outputStr := strings.TrimSpace(string(output))
			if strings.Contains(outputStr, "not found") || strings.Contains(outputStr, "could not read") {
				return fmt.Errorf("failed to clone preset: repository not found or inaccessible\n  URL: %s\n  Hint: check the URL is correct and accessible", preset)
			}
			if strings.Contains(outputStr, "Authentication failed") || strings.Contains(outputStr, "authentication required") {
				return fmt.Errorf("failed to clone preset: authentication required\n  URL: %s\n  Hint: use HTTPS URLs for public repos or configure SSH keys for private repos", preset)
			}
			return fmt.Errorf("failed to clone preset from %s:\n  %s", preset, outputStr)
		}
		sourceDir = tempDir
		fmt.Fprintf(w, "✓ Cloned preset to temporary directory\n")
	} else {
		fmt.Fprintf(w, "Using local preset directory: %s\n", preset)
	}
	defer cleanup()

	// Validate preset structure before copying
	fmt.Fprintf(w, "Validating preset structure...\n")
	missingDirs := []string{}
	for _, name := range dirs {
		src := filepath.Join(sourceDir, name)
		info, err := os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				missingDirs = append(missingDirs, name)
				continue
			}
			return fmt.Errorf("failed to check %s directory: %w", name, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("preset %s exists but is not a directory", name)
		}
	}

	if len(missingDirs) > 0 {
		return fmt.Errorf("preset is missing required directories: %s\n  Expected: %s\n  Location: %s\n  Hint: a valid preset must contain %s subdirectories",
			strings.Join(missingDirs, ", "),
			strings.Join(dirs, ", "),
			sourceDir,
			strings.Join(dirs, " and "))
	}
	fmt.Fprintf(w, "✓ Preset structure validated\n\n")

	// Copy each directory
	for _, name := range dirs {
		src := filepath.Join(sourceDir, name)
		dst := filepath.Join(baseDir, name)
		fmt.Fprintf(w, "Refreshing %s/...\n", name)
		if err := copyDir(w, src, dst, name); err != nil {
			return fmt.Errorf("failed to copy %s: %w", name, err)
		}
	}

	return nil
}

func copyDir(w io.Writer, src, dst, logPrefix string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			if rel == "." {
				return nil
			}
			return os.MkdirAll(target, 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if w != nil {
			fmt.Fprintf(w, "  Refreshing %s\n", filepath.Join(logPrefix, rel))
		}
		return os.WriteFile(target, data, info.Mode())
	})
}
