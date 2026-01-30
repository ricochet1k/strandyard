package web

import (
	"sync"

	"github.com/ricochet1k/streamyard/pkg/task"
)

type ProjectInfo struct {
	Name          string
	StorageRoot   string
	TasksRoot     string
	RolesRoot     string
	TemplatesRoot string
	GitRoot       string
	Storage       string // "global" or "local"
}

type ServerConfig struct {
	Port           int
	Projects       []ProjectInfo
	CurrentProject string
	AutoOpen       bool
}

type StreamUpdate struct {
	Event   string             `json:"event"`
	Path    string             `json:"path"`
	Project string             `json:"project"`
	Task    *task.TaskSnapshot `json:"task,omitempty"`
}

type updateBroker struct {
	mu      sync.Mutex
	clients map[chan StreamUpdate]struct{}
}

func newUpdateBroker() *updateBroker {
	return &updateBroker{clients: make(map[chan StreamUpdate]struct{})}
}

func (b *updateBroker) subscribe(ch chan StreamUpdate) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[ch] = struct{}{}
}

func (b *updateBroker) unsubscribe(ch chan StreamUpdate) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, ch)
	close(ch)
}

func (b *updateBroker) broadcast(update StreamUpdate) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		select {
		case ch <- update:
		default: // Skip slow clients
		}
	}
}
