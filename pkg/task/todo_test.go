package task

import (
	"testing"
)

func TestTaskDB_TodoOperations(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	id := "T1aaa-task"
	createTaskFile(t, tasksRoot, id, "Test Task")

	// Test AddTodo
	if err := db.AddTodo(id, "Todo 1"); err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}
	if err := db.AddTodo(id, "(role: developer) Todo 2"); err != nil {
		t.Fatalf("AddTodo 2 failed: %v", err)
	}

	task, _ := db.Get(id)
	if len(task.TodoItems) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(task.TodoItems))
	}
	if task.TodoItems[0].Text != "Todo 1" {
		t.Errorf("expected Todo 1, got %q", task.TodoItems[0].Text)
	}
	if task.TodoItems[1].Text != "Todo 2" || task.TodoItems[1].Role != "developer" {
		t.Errorf("expected Todo 2 with role developer, got %q role %q", task.TodoItems[1].Text, task.TodoItems[1].Role)
	}

	// Test EditTodo
	if err := db.EditTodo(id, 1, "Updated Todo 1"); err != nil {
		t.Fatalf("EditTodo failed: %v", err)
	}
	task, _ = db.Get(id)
	if task.TodoItems[0].Text != "Updated Todo 1" {
		t.Errorf("expected Updated Todo 1, got %q", task.TodoItems[0].Text)
	}

	// Test CheckTodo (via CompleteTodo)
	if _, err := db.CompleteTodo(id, 1, "Report 1"); err != nil {
		t.Fatalf("CompleteTodo failed: %v", err)
	}
	task, _ = db.Get(id)
	if !task.TodoItems[0].Checked || task.TodoItems[0].Report != "Report 1" {
		t.Errorf("todo should be checked with report")
	}

	// Test UncheckTodo
	if err := db.UncheckTodo(id, 1); err != nil {
		t.Fatalf("UncheckTodo failed: %v", err)
	}
	task, _ = db.Get(id)
	if task.TodoItems[0].Checked {
		t.Error("todo should be unchecked")
	}

	// Test RemoveTodo
	if err := db.RemoveTodo(id, 1); err != nil {
		t.Fatalf("RemoveTodo failed: %v", err)
	}
	task, _ = db.Get(id)
	if len(task.TodoItems) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(task.TodoItems))
	}
	if task.TodoItems[0].Text != "Todo 2" {
		t.Errorf("expected remaining todo to be Todo 2, got %q", task.TodoItems[0].Text)
	}

	// Test ReorderTodo
	db.AddTodo(id, "Todo 3")
	db.AddTodo(id, "Todo 4")
	// Current state: [Todo 2, Todo 3, Todo 4]
	if err := db.ReorderTodo(id, 3, 1); err != nil {
		t.Fatalf("ReorderTodo failed: %v", err)
	}
	// Expected state: [Todo 4, Todo 2, Todo 3]
	task, _ = db.Get(id)
	if task.TodoItems[0].Text != "Todo 4" || task.TodoItems[1].Text != "Todo 2" || task.TodoItems[2].Text != "Todo 3" {
		t.Errorf("reorder failed, got: %v, %v, %v", task.TodoItems[0].Text, task.TodoItems[1].Text, task.TodoItems[2].Text)
	}
}

func TestTaskDB_TodoCompletionLogic(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	id := "T1aaa-task"
	createTaskFile(t, tasksRoot, id, "Test Task")

	db.AddTodo(id, "Todo 1")
	db.AddTodo(id, "Todo 2")

	// Complete all todos
	db.CompleteTodo(id, 1, "")
	result, _ := db.CompleteTodo(id, 2, "")

	if !result.TaskCompleted {
		t.Error("task should be marked as completed after last todo")
	}

	task, _ := db.Get(id)
	if !task.Meta.Completed || task.Meta.Status != "done" {
		t.Errorf("task metadata should be completed and status done, got %v status %q", task.Meta.Completed, task.Meta.Status)
	}

	// Add new todo - should mark task as incomplete
	db.AddTodo(id, "Todo 3")
	task, _ = db.Get(id)
	if task.Meta.Completed || task.Meta.Status == "done" {
		t.Error("task should be incomplete after adding new todo")
	}

	// Complete again
	db.CompleteTodo(id, 3, "")
	task, _ = db.Get(id)
	if !task.Meta.Completed {
		t.Error("task should be completed again")
	}

	// Uncheck one - should mark task as incomplete
	db.UncheckTodo(id, 3)
	task, _ = db.Get(id)
	if task.Meta.Completed {
		t.Error("task should be incomplete after unchecking a todo")
	}
}
