package task

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// TaskEvent describes the kind of file change observed for a task.
type TaskEvent string

const (
	TaskCreated  TaskEvent = "created"
	TaskModified TaskEvent = "modified"
	TaskRemoved  TaskEvent = "removed"
	TaskRenamed  TaskEvent = "renamed"
)

// TaskSnapshot is a JSON-marshable representation of a task at a point in time.
type TaskSnapshot struct {
	ID       string     `json:"id"`
	Dir      string     `json:"dir"`
	FilePath string     `json:"file_path"`
	Meta     Metadata   `json:"meta"`
	Content  string     `json:"content"`
	Title    string     `json:"title"`
	Todos    []TaskItem `json:"todos,omitempty"`
	Subtasks []TaskItem `json:"subtasks,omitempty"`
}

// TaskUpdate represents a task change event.
type TaskUpdate struct {
	Task  *TaskSnapshot `json:"task,omitempty"`
	Event TaskEvent     `json:"event"`
	Path  string        `json:"path"`
}

// WatchTasks starts watching task files under tasksRoot. It returns a channel
// of updates and a channel of errors. Cancel the context to stop the watcher.
func WatchTasks(ctx context.Context, tasksRoot string) (<-chan TaskUpdate, <-chan error, error) {
	if tasksRoot == "" {
		return nil, nil, fmt.Errorf("tasks root is required")
	}

	info, err := os.Stat(tasksRoot)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to stat tasks root %s: %w", tasksRoot, err)
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("tasks root is not a directory: %s", tasksRoot)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	if err := addWatchDirs(watcher, tasksRoot); err != nil {
		_ = watcher.Close()
		return nil, nil, err
	}

	updates := make(chan TaskUpdate, 32)
	errors := make(chan error, 8)
	parser := NewParser()

	go func() {
		defer close(updates)
		defer close(errors)
		defer watcher.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-watcher.Errors:
				if err != nil {
					errors <- err
				}
			case event := <-watcher.Events:
				if event.Name == "" {
					continue
				}

				if event.Op&fsnotify.Create != 0 {
					if addDirErr := watchIfDir(watcher, event.Name); addDirErr != nil {
						errors <- addDirErr
					}
				}

				if !isTaskFilePath(event.Name) {
					continue
				}

				switch {
				case event.Op&fsnotify.Write != 0:
					sendTaskUpdate(parser, updates, errors, event.Name, TaskModified)
				case event.Op&fsnotify.Create != 0:
					sendTaskUpdate(parser, updates, errors, event.Name, TaskCreated)
				case event.Op&fsnotify.Remove != 0:
					updates <- TaskUpdate{Event: TaskRemoved, Path: event.Name}
				case event.Op&fsnotify.Rename != 0:
					updates <- TaskUpdate{Event: TaskRenamed, Path: event.Name}
				}
			}
		}
	}()

	return updates, errors, nil
}

func sendTaskUpdate(parser *Parser, updates chan<- TaskUpdate, errors chan<- error, path string, event TaskEvent) {
	task, err := parser.ParseFile(path)
	if err != nil {
		errors <- err
		return
	}

	snapshot, err := snapshotFromTask(task)
	if err != nil {
		errors <- err
		return
	}

	updates <- TaskUpdate{
		Task:  snapshot,
		Event: event,
		Path:  path,
	}
}

func snapshotFromTask(task *Task) (*TaskSnapshot, error) {
	if task == nil {
		return nil, nil
	}

	content := task.Content()

	return &TaskSnapshot{
		ID:       task.ID,
		Dir:      task.Dir,
		FilePath: task.FilePath,
		Meta:     task.Meta,
		Content:  content,
		Title:    task.Title(),
		Todos:    task.TodoItems,
		Subtasks: task.SubsItems,
	}, nil
}

func isTaskFilePath(path string) bool {
	base := filepath.Base(path)
	if base == "task.md" || base == "README.md" {
		return true
	}
	dir := filepath.Base(filepath.Dir(path))
	return base == dir+".md"
}

func addWatchDirs(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if err := watcher.Add(path); err != nil {
			return fmt.Errorf("failed to watch %s: %w", path, err)
		}
		return nil
	})
}

func watchIfDir(watcher *fsnotify.Watcher, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if !info.IsDir() {
		return nil
	}
	return addWatchDirs(watcher, path)
}
