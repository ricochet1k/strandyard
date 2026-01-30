package web

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/ricochet1k/streamyard/pkg/task"
)

func (s *Server) startWatchers(ctx context.Context) error {
	for _, proj := range s.config.Projects {
		updates, errs, err := task.WatchTasks(ctx, proj.TasksRoot)
		if err != nil {
			return err
		}

		go s.relayUpdates(ctx, proj.Name, proj.StorageRoot, updates, errs)
	}

	return nil
}

func (s *Server) relayUpdates(ctx context.Context, projectName, storageRoot string, updates <-chan task.TaskUpdate, errs <-chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-errs:
			if !ok {
				return
			}
			if err != nil {
				s.logger.Printf("[%s] watcher error: %v", projectName, err)
			}
		case update, ok := <-updates:
			if !ok {
				return
			}
			// Add project context
			relPath := strings.TrimPrefix(update.Path, storageRoot+string(filepath.Separator))
			enriched := StreamUpdate{
				Event:   string(update.Event),
				Path:    relPath,
				Project: projectName,
				Task:    update.Task,
			}
			s.broker.broadcast(enriched)
		}
	}
}
