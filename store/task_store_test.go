package store

import (
	"sync"
	"testing"

	"loom-bootstrap-test-5/models"
)

func TestNewInMemoryStore(t *testing.T) {
	s := NewInMemoryStore()
	if s == nil {
		t.Fatal("nil store")
	}
	if s.tasks == nil {
		t.Error("nil tasks map")
	}
}

func TestCreate_Success(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	if err := s.Create(task); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.Get(task.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Title != task.Title {
		t.Errorf("title mismatch")
	}
	if got.Completed {
		t.Error("expected not completed")
	}
}

func TestCreate_Multiple(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 5; i++ {
		s.Create(models.NewTask("T", "D"))
	}
	tasks, _ := s.List(nil)
	if len(tasks) != 5 {
		t.Errorf("expected 5, got %d", len(tasks))
	}
}

func TestCreate_EmptyTitle(t *testing.T) {
	s := NewInMemoryStore()
	if err := s.Create(models.NewTask("", "desc")); err == nil {
		t.Error("expected error for empty title")
	}
}

func TestCreate_TitleTooLong(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask(string(make([]byte, 300)), "desc")
	if err := s.Create(task); err == nil {
		t.Error("expected error for long title")
	}
}

func TestCreate_DescTooLong(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("title", string(make([]byte, 10001)))
	if err := s.Create(task); err == nil {
		t.Error("expected error for long description")
	}
}

func TestCreate_ValidBoundary(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask(string(make([]byte, 255)), string(make([]byte, 10000)))
	task.Title = "A"
	if err := s.Create(task); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreate_EmptyDescription(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("T", "")
	if err := s.Create(task); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := s.Get(task.ID)
	if got.Description != "" {
		t.Error("expected empty description")
	}
}

func TestList_Empty(t *testing.T) {
	s := NewInMemoryStore()
	tasks, _ := s.List(nil)
	if len(tasks) != 0 {
		t.Errorf("expected 0, got %d", len(tasks))
	}
}

func TestList_All(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(models.NewTask("A", "D"))
	s.Create(models.NewTask("B", "D"))
	s.Create(models.NewTask("C", "D"))
	tasks, _ := s.List(nil)
	if len(tasks) != 3 {
		t.Errorf("expected 3, got %d", len(tasks))
	}
}

func TestList_FilterActive(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(models.NewTask("A", "D"))
	s.Create(models.NewTask("B", "D"))
	c := models.NewTask("C", "D")
	c.Completed = true
	s.Create(c)
	active := "active"
	tasks, _ := s.List(&models.TaskFilter{Status: &active})
	if len(tasks) != 2 {
		t.Errorf("expected 2, got %d", len(tasks))
	}
	for _, tc := range tasks {
		if tc.Completed {
			t.Error("unexpected completed in active")
		}
	}
}

func TestList_FilterCompleted(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(models.NewTask("A", "D"))
	c := models.NewTask("B", "D")
	c.Completed = true
	s.Create(c)
	completed := "completed"
	tasks, _ := s.List(&models.TaskFilter{Status: &completed})
	if len(tasks) != 1 {
		t.Errorf("expected 1, got %d", len(tasks))
	}
	if !tasks[0].Completed {
		t.Error("expected completed")
	}
}

func TestList_Limit(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 5; i++ {
		s.Create(models.NewTask("T", "D"))
	}
	tasks, _ := s.List(&models.TaskFilter{Limit: 2})
	if len(tasks) != 2 {
		t.Errorf("expected 2, got %d", len(tasks))
	}
}

func TestList_FilterLimitCombo(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(models.NewTask("A", "D"))
	c := models.NewTask("B", "D")
	c.Completed = true
	s.Create(c)
	active := "active"
	tasks, _ := s.List(&models.TaskFilter{Status: &active, Limit: 1})
	if len(tasks) != 1 {
		t.Errorf("expected 1, got %d", len(tasks))
	}
}

func TestGet_Existing(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("T", "D")
	s.Create(task)
	got, err := s.Get(task.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.ID != task.ID {
		t.Errorf("ID mismatch")
	}
	if got.Title != task.Title {
		t.Errorf("title mismatch")
	}
}

func TestGet_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Error("expected error")
	}
}

func TestGet_EmptyID(t *testing.T) {
	s := NewInMemoryStore()
	_, err := s.Get("")
	if err == nil {
		t.Error("expected error for empty ID")
	}
}

func TestUpdate_Complete(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("T", "D")
	s.Create(task)
	c := true
	got, err := s.Update(task.ID, &models.TaskUpdate{Completed: &c})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if !got.Completed {
		t.Error("expected completed")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	c := true
	_, err := s.Update("nonexistent", &models.TaskUpdate{Completed: &c})
	if err == nil {
		t.Error("expected error")
	}
}

func TestReplace_Success(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("T", "D")
	s.Create(task)
	newTask := models.NewTask("Replaced", "New")
func TestUpdate_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	c := true
	_, err := s.Update("nonexistent", &models.TaskUpdate{Completed: &c})
	if err == nil {
		t.Error("expected error")
	}
}

func TestReplace_Success(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("T", "D")
	s.Create(task)
	newTask := models.NewTask("Replaced", "New")
	replaced, err := s.Replace(task.ID, newTask)
	if err != nil {
		t.Fatalf("Replace failed: %v", err)
	}
	if replaced.Title != "Replaced" {
		t.Errorf("expected Replaced, got %s", replaced.Title)
	}
}

func TestReplace_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	newTask := models.NewTask("T", "D")
	_, err := s.Replace("nonexistent", newTask)
	if err == nil {
		t.Error("expected error")
	}
}

func TestDelete_Success(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("T", "D")
	s.Create(task)
	if err := s.Delete(task.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err := s.Get(task.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	if err := s.Delete("nonexistent"); err == nil {
		t.Error("expected error")
	}
}

func TestHealthCheck(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(models.NewTask("T", "D"))
	stats := s.HealthCheck()
	if stats["status"] != "healthy" {
		t.Errorf("expected healthy, got %s", stats["status"]) 
	}
	if stats["tasks"] != "1" {
		t.Errorf("expected 1, got %s", stats["tasks"]) 
	}
	if stats["service"] != "task-api" {
		t.Errorf("expected task-api, got %s", stats["service"]) 
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := NewInMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				s.Create(models.NewTask("T", "D"))
			}
		}()
	}
	wg.Wait()
	tasks, err := s.List(nil)
	if err != nil {
		t.Fatalf("concurrent list error: %v", err)
	}
	if len(tasks) == 0 {
		t.Error("expected tasks after concurrent access")
	}
}
