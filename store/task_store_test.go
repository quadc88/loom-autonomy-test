package store

import (
	"testing"

	"loom-bootstrap-test-5/models"
)

func TestNewInMemoryStore(t *testing.T) {
	s := NewInMemoryStore()
	if s == nil {
		t.Fatal("expected non-nil store")
	}
	if s.tasks == nil {
		t.Error("expected non-nil tasks map")
	}
}

func TestCreate(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test Task", "Test Description")

	err := s.Create(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	retrieved, err := s.Get(task.ID)
	if err != nil {
		t.Fatalf("failed to get task: %v", err)
	}
	if retrieved.Title != task.Title {
		t.Errorf("expected title %q, got %q", task.Title, retrieved.Title)
	}
	if retrieved.Completed {
		t.Error("expected task to be not completed")
	}
}

func TestCreateValidationErrors(t *testing.T) {
	s := NewInMemoryStore()
	tests := []struct {
		name    string
		task    *models.Task
		wantErr bool
	}{
		{
			name:    "empty title",
			task:    models.NewTask("", "desc"),
			wantErr: true,
		},
		{
			name:    "title too long",
			task:    models.NewTask(string(make([]byte, 300)), "desc"),
			wantErr: true,
		},
		{
			name:    "description too long",
			task:    models.NewTask("title", string(make([]byte, 10001))),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Create(tt.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestList(t *testing.T) {
	s := NewInMemoryStore()
	task1 := models.NewTask("Task 1", "Desc 1")
	task2 := models.NewTask("Task 2", "Desc 2")
	s.Create(task1)
	s.Create(task2)

	tasks, err := s.List(nil)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestListWithStatusFilter(t *testing.T) {
	s := NewInMemoryStore()
	task1 := models.NewTask("Active", "Desc")
	task2 := models.NewTask("Completed", "Desc")
	task2.Completed = true
	s.Create(task1)
	s.Create(task2)

	active := "active"
	filtered, err := s.List(&models.TaskFilter{Status: &active})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(filtered) != 1 {
		t.Errorf("expected 1 active task, got %d", len(filtered))
	}
	if filtered[0].Completed {
		t.Error("expected only active tasks")
	}

	completed := "completed"
	filtered, err = s.List(&models.TaskFilter{Status: &completed})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(filtered) != 1 {
		t.Errorf("expected 1 completed task, got %d", len(filtered))
	}
	if !filtered[0].Completed {
		t.Error("expected only completed tasks")
	}
}

func TestListWithLimit(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 5; i++ {
		task := models.NewTask("Task", "Desc")
		s.Create(task)
	}

	filtered, err := s.List(&models.TaskFilter{Limit: 2})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(filtered) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(filtered))
	}
}

func TestGet(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)

	retrieved, err := s.Get(task.ID)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if retrieved.ID != task.ID {
		t.Errorf("expected ID %q, got %q", task.ID, retrieved.ID)
	}

	_, err = s.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestUpdate(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)

	completed := true
	updated, err := s.Update(task.ID, &models.TaskUpdate{Completed: &completed})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if !updated.Completed {
		t.Error("expected task to be completed")
	}

	_, err = s.Update("nonexistent", &models.TaskUpdate{})
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestReplace(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)

	newTask := models.NewTask("Replaced", "New Desc")
	replaced, err := s.Replace(task.ID, newTask)
	if err != nil {
		t.Fatalf("Replace() error: %v", err)
	}
	if replaced.Title != "Replaced" {
		t.Errorf("expected title 'Replaced', got %q", replaced.Title)
	}

	_, err = s.Replace("nonexistent", newTask)
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestDelete(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)

	err := s.Delete(task.ID)
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	_, err = s.Get(task.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}

	err = s.Delete("nonexistent")
	if err == nil {
		t.Error("expected error for deleting nonexistent task")
	}
}

func TestHealthCheck(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)

	stats := s.HealthCheck()
	if stats["status"] != "healthy" {
		t.Errorf("expected status healthy, got %q", stats["status"])
	}
	if stats["tasks"] != "1" {
		t.Errorf("expected 1 task, got %q", stats["tasks"])
	}
	if stats["service"] != "task-api" {
		t.Errorf("expected service task-api, got %q", stats["service"])
	}
}
