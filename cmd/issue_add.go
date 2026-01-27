/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/ricochet1k/memmd/pkg/idgen"
	"github.com/ricochet1k/memmd/pkg/task"
	"github.com/spf13/cobra"
)

const issuePrefix = "I"

type issueTemplateData struct {
	Kind        string
	Role        string
	Priority    string
	DateCreated string
	DateEdited  string
	Title       string
}

var (
	issueAddTitle      string
	issueAddPriority   string
	issueAddNoValidate bool
)

var issueAddCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Create a new issue task",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runIssueAdd(args)
	},
}

func init() {
	issueCmd.AddCommand(issueAddCmd)
	issueAddCmd.Flags().StringVarP(&issueAddTitle, "title", "t", "", "issue title")
	issueAddCmd.Flags().StringVar(&issueAddPriority, "priority", "medium", "priority: high, medium, or low")
	issueAddCmd.Flags().BoolVar(&issueAddNoValidate, "no-validate", false, "skip validation and master list updates")
}

func runIssueAdd(args []string) error {
	title := strings.TrimSpace(issueAddTitle)
	if title == "" && len(args) > 0 {
		title = strings.TrimSpace(strings.Join(args, " "))
	}
	if title == "" {
		return fmt.Errorf("title is required (use --title or provide it as an argument)")
	}

	priority := task.NormalizePriority(issueAddPriority)
	if !task.IsValidPriority(priority) {
		return fmt.Errorf("invalid priority: %s", issueAddPriority)
	}

	parentDir := "tasks"

	id, err := idgen.GenerateID(issuePrefix, title)
	if err != nil {
		return err
	}

	taskDir := filepath.Join(parentDir, id)
	if _, err := os.Stat(taskDir); err == nil {
		return fmt.Errorf("task directory already exists: %s", taskDir)
	}
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		return fmt.Errorf("failed to create task directory: %w", err)
	}

	taskFile := filepath.Join(taskDir, id+".md")

	now := time.Now().UTC().Format(time.RFC3339)
	data := issueTemplateData{
		Kind:        "issue",
		Role:        "triage",
		Priority:    priority,
		DateCreated: now,
		DateEdited:  now,
		Title:       title,
	}

	if err := renderTemplateToFile("templates/issue.md", taskFile, data); err != nil {
		return err
	}

	fmt.Printf("✓ Issue created: %s\n", filepath.ToSlash(taskFile))

	if !issueAddNoValidate {
		if err := runValidate("tasks", "tasks/root-tasks.md", "tasks/free-tasks.md", "text"); err != nil {
			return err
		}
	}

	return nil
}

func renderTemplateToFile(templatePath, outputPath string, data issueTemplateData) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer out.Close()

	if err := tmpl.Execute(out, data); err != nil {
		return fmt.Errorf("failed to render template %s: %w", templatePath, err)
	}

	return nil
}
