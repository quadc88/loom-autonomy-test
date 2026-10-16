package store

import (
	"testing"

	"loom-bootstrap-test-5/models"
)

func TestCreateTask(t *testing.T) {
	s := NewInMemoryStore()

	// Valid create
	task := models.NewTask("Test Task", "Description")
	err := s.Create(task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify task exists
	retrieved, err := s.Get(task.ID)
	if err != nil {
		t.Fatalf("expected task to exist, got error: %v", err)
	}
	if retrieved.Title != task.Title {
		t.Errorf("expected title %q, got %q", task.Title, retrieved.Title)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	s := NewInMemoryStore()

	// Empty title
	task := models.NewTask("", "Description")
	err := s.Create(task)
	if err == nil {
		t.Error("expected error for empty title")
	}

	// Title too long
	longTitle := string(make([]byte, 256))
	task = models.NewTask(longTitle, "Description")
	err = s.Create(task)
	if err == nil {
		t.Error("expected error for title > 255 chars")
	}

	// Description too long
	longDesc := string(make([]byte, 10001))
	task = models.NewTask("Valid Title", longDesc)
	err = s.Create(task)
	if err == nil {
		t.Error("expected error for description > 10000 chars")
	}
}

func TestReplaceTask(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Original", "Desc")
	s.Create(task)

	// Replace existing
	replaced := &models.Task{
		ID:          task.ID,
		Title:       "Replaced",n		Description: "New desc",
		Completed:   true,
	}
	result, err := s.Replace(task.ID, replaced)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Title != "Replaced" {
		t.Errorf("expected 'Replaced', got %q", result.Title)
	}
	if !result.Completed {
		t.Error("expected completed=true")
	}

	// Replace non-existing
	_, err = s.Replace("nonexistent", replaced)
	if err == nil {
		t.Error("expected error for non-existent task")
	}
}
}

func TestGetTask(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test Task", "Description")
	s.Create(task)

	// Get existing
	retrieved, err := s.Get(task.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrieved.ID != task.ID {
		t.Errorf("expected ID %q, got %q", task.ID, retrieved.ID)
	}

	// Get non-existing
	_, err = s.Get("non-existent-id")
	if err == nil {
		t.Error("expected error for non-existent task")
	}
}

func TestListTasks(t *testing.T) {
	s := NewInMemoryStore()

	// Create some tasks
	task1 := models.NewTask("Task 1", "Desc 1")
	task2 := models.NewTask("Task 2", "Desc 2")
	s.Create(task1)
	s.Create(task2)

	// Mark one as completed
	complete := true
	s.Update(task2.ID, &models.TaskUpdate{Completed: &complete})

	// List all
	all, err := s.List(nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(all))
	}

	// List active only
	status := "active"
	active, err := s.List(&models.TaskFilter{Status: &status})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(active) != 1 {
		t.Errorf("expected 1 active task, got %d", len(active))
	}
	if active[0].ID != task1.ID {
		t.Errorf("expected task1, got %s", active[0].ID)
	}

	// List completed only
	status = "completed"
	completed, err := s.List(&models.TaskFilter{Status: &status})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(completed) != 1 {
		t.Errorf("expected 1 completed task, got %d", len(completed))
	}
	if completed[0].ID != task2.ID {
		t.Errorf("expected task2, got %s", completed[0].ID)
	}

	// List with limit
	limit := 1
	limited, err := s.List(&models.TaskFilter{Limit: limit})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(limited) != 1 {
		t.Errorf("expected 1 task with limit=1, got %d", len(limited))
	}
}

func TestUpdateTask(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test Task", "Description")
	s.Create(task)

	// Update to completed
	complete := true
	updated, err := s.Update(task.ID, &models.TaskUpdate{Completed: &complete})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !updated.Completed {
		t.Error("expected task to be completed")
	}

	// Update non-existing
	_, err = s.Update("non-existent", &models.TaskUpdate{Completed: &complete})
	if err == nil {
		t.Error("expected error for non-existent task")
	}
}

func TestDeleteTask(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test Task", "Description")
	s.Create(task)

	// Delete existing
	err := s.Delete(task.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify deleted
	_, err = s.Get(task.ID)
	if err == nil {
		t.Error("expected task to be deleted")
	}

	// Delete non-existing
	err = s.Delete("non-existent")
	if err == nil {
		t.Error("expected error for non-existent task")
	}
}

func TestHealthCheck(t *testing.T) {
	s := NewInMemoryStore()

	// Empty store
	health := s.HealthCheck()
	if health["status"] != "healthy" {
		t.Errorf("expected healthy status, got %s", health["status"])
	}
	if health["tasks"] != "0" {
		t.Errorf("expected 0 tasks, got %s", health["tasks"])
	}

	// Add task and check
	task := models.NewTask("Test Task", "Description")
	s.Create(task)
	health = s.HealthCheck()
	if health["tasks"] != "1" {
		t.Errorf("expected 1 task, got %s", health["tasks"])
	}
}
