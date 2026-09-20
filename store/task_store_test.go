package store

import (
	"testing"
	"time"

	"loom-bootstrap-test-5/models"
)

func boolPtr(b bool) *bool { return &b }
func strPtr(s string) *string { return &s }

func TestNewInMemoryStore(t *testing.T) {
	s := NewInMemoryStore()
	if s == nil || s.tasks == nil {
		t.Fatal("Expected non-nil store with initialized map")
	}
}

func TestStoreCreate(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test Task", "Test Desc")
	if err := s.Create(task); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	retrieved, err := s.Get(task.ID)
	if err != nil {
		t.Fatalf("Expected to find task: %v", err)
	}
	if retrieved.Title != task.Title || retrieved.Description != task.Description {
		t.Error("Mismatch in retrieved task")
	}
	// NewTask sets timestamps to current time, so they should not be zero.
	if retrieved.CreatedAt.IsZero() || retrieved.UpdatedAt.IsZero() {
		t.Error("Expected non-zero timestamps after create")
	}
}

func TestStoreCreate_InvalidTitle(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("", "Test")
	if err := s.Create(task); err == nil {
		t.Fatal("Expected error for empty title")
	} else if err.Error() != "title is required" {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestStoreCreate_TitleTooLong(t *testing.T) {
	s := NewInMemoryStore()
	task := &models.Task{ID: "test", Title: string(make([]byte, 256)), Description: "Test"}
	if err := s.Create(task); err == nil {
		t.Fatal("Expected error for long title")
	} else if err.Error() != "title must be 255 characters or less" {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestStoreGet_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	if _, err := s.Get("nonexistent"); err == nil {
		t.Fatal("Expected error for missing task")
	} else if err.Error() != "task not found" {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestStoreList_FilterAndLimit(t *testing.T) {
	s := NewInMemoryStore()
	tasks := make([]*models.Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = models.NewTask("Task", "Desc")
		s.Create(tasks[i])
	}
	updates := &models.TaskUpdate{Completed: boolPtr(true)}
	s.Update(tasks[2].ID, updates)

	filtered, err := s.List(&models.TaskFilter{Status: strPtr("completed")})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(filtered) != 1 || filtered[0].ID != tasks[2].ID {
		t.Error("Expected only completed task")
	}

	limited, err := s.List(&models.TaskFilter{Limit: 2})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(limited) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(limited))
	}
}

func TestStoreUpdate(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	before := time.Now()
	updates := &models.TaskUpdate{Completed: boolPtr(true)}
	updated, err := s.Update(task.ID, updates)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !updated.Completed {
		t.Error("Expected task to be completed")
	}
	if updated.UpdatedAt.Before(before) {
		t.Error("Expected UpdatedAt to be recent")
	}
	if updated.Title != task.Title {
		t.Error("Expected title unchanged")
	}
}

func TestStoreDelete(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	if err := s.Delete(task.ID); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, err := s.Get(task.ID); err == nil {
		t.Error("Expected task to be deleted")
	}
}

func TestStoreReplace(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Original", "Desc")
	s.Create(task)
	replacement := &models.Task{ID: task.ID, Title: "Replaced", Description: "New Desc", Completed: true}
	replaced, err := s.Replace(task.ID, replacement)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if replaced.Title != "Replaced" || !replaced.Completed {
		t.Error("Expected replaced task")
	}
}

func TestStoreHealthCheck(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	health := s.HealthCheck()
	if health["status"] != "healthy" || health["tasks"] != "1" {
		t.Errorf("Unexpected health check: %v", health)
	}
}
