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
	sourceDir := preset
	cleanup := func() {}

	if info, err := os.Stat(preset); err != nil || !info.IsDir() {
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
			return fmt.Errorf("failed to clone preset: %s", strings.TrimSpace(string(output)))
		}
		sourceDir = tempDir
	}
	defer cleanup()

	for _, name := range dirs {
		src := filepath.Join(sourceDir, name)
		dst := filepath.Join(baseDir, name)
		info, err := os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("preset is missing %s directory", name)
			}
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("preset %s is not a directory", src)
		}
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
