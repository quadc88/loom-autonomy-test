package store

import (
	"sync"
	"testing"
	"time"

	"loom-bootstrap-test-5/models"
)

func boolPtr(b bool) *bool { return &b }
func strPtr(s string) *string { return &s }

func TestNewInMemoryStore(t *testing.T) {
	s := NewInMemoryStore()
	if s == nil || s.tasks == nil {
		t.Fatal("Expected non-nil store")
	}
}

func TestStoreCreate_Valid(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	if err := s.Create(task); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	got, _ := s.Get(task.ID)
	if got.Title != task.Title || got.Description != task.Description {
		t.Error("Mismatch")
	}
	if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
		t.Error("Expected non-zero timestamps")
	}
}

func TestStoreCreate_InvalidTitleEmpty(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("", "Test")
	if err := s.Create(task); err == nil {
		t.Fatal("Expected error for empty title")
	}
}

func TestStoreCreate_TitleTooLong(t *testing.T) {
	s := NewInMemoryStore()
	task := &models.Task{ID: "test", Title: string(make([]byte, 256)), Description: "Test"}
	if err := s.Create(task); err == nil {
		t.Fatal("Expected error for long title")
	}
}

func TestStoreCreate_DescriptionTooLong(t *testing.T) {
	s := NewInMemoryStore()
	task := &models.Task{ID: "test", Title: "Valid", Description: string(make([]byte, 10001))}
	if err := s.Create(task); err == nil {
		t.Fatal("Expected error for long description")
	}
}

func TestStoreCreate_DuplicateTitleAllowed(t *testing.T) {
	s := NewInMemoryStore()
	t1 := models.NewTask("Same", "1")
	t2 := models.NewTask("Same", "2")
	s.Create(t1)
	s.Create(t2)
	g1, _ := s.Get(t1.ID)
	g2, _ := s.Get(t2.ID)
	if g1 == nil || g2 == nil || g1.ID == g2.ID {
		t.Error("Expected distinct tasks")
	}
}

func TestStoreCreate_OrderPreserved(t *testing.T) {
	s := NewInMemoryStore()
	tasks := make([]*models.Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = models.NewTask("T", "D")
		s.Create(tasks[i])
	}
	result, _ := s.List(nil)
	for i, task := range tasks {
		if result[i].ID != task.ID {
			t.Errorf("Wrong order at %d", i)
		}
	}
}

// 2.1.2 List

func TestStoreList_Empty(t *testing.T) {
	s := NewInMemoryStore()
	result, _ := s.List(nil)
	if len(result) != 0 {
		t.Errorf("Expected 0, got %d", len(result))
	}
}

func TestStoreList_NoFilter(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 3; i++ {
		s.Create(models.NewTask("T", "D"))
	}
	result, _ := s.List(nil)
	if len(result) != 3 {
		t.Errorf("Expected 3, got %d", len(result))
	}
}

func TestStoreList_FilterCompleted(t *testing.T) {
	s := NewInMemoryStore()
	tasks := make([]*models.Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = models.NewTask("T", "D")
		s.Create(tasks[i])
	}
	s.Update(tasks[2].ID, &models.TaskUpdate{Completed: boolPtr(true)})
	filtered, _ := s.List(&models.TaskFilter{Status: strPtr("completed")})
	if len(filtered) != 1 || filtered[0].ID != tasks[2].ID {
		t.Error("Expected only completed task")
	}
}

func TestStoreList_FilterActive(t *testing.T) {
	s := NewInMemoryStore()
	tasks := make([]*models.Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = models.NewTask("T", "D")
		s.Create(tasks[i])
	}
	s.Update(tasks[2].ID, &models.TaskUpdate{Completed: boolPtr(true)})
	filtered, _ := s.List(&models.TaskFilter{Status: strPtr("active")})
	if len(filtered) != 4 {
		t.Errorf("Expected 4 active, got %d", len(filtered))
	}
	for _, t := range filtered {
		if t.Completed {
			t.Error("Found completed in active")
		}
	}
}

func TestStoreList_Limit(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 5; i++ {
		s.Create(models.NewTask("T", "D"))
	}
	limited, _ := s.List(&models.TaskFilter{Limit: 2})
	if len(limited) != 2 {
		t.Errorf("Expected 2, got %d", len(limited))
	}
}

func TestStoreList_LimitExceedsTotal(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 3; i++ {
		s.Create(models.NewTask("T", "D"))
	}
	result, _ := s.List(&models.TaskFilter{Limit: 100})
	if len(result) != 3 {
		t.Errorf("Expected 3, got %d", len(result))
	}
}

func TestStoreList_FilterNone(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	result, _ := s.List(&models.TaskFilter{Status: nil})
	if len(result) != 1 {
		t.Errorf("Expected 1, got %d", len(result))
	}
}

// 2.1.3 Get

func TestStoreGet_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	if _, err := s.Get("nonexistent"); err == nil {
		t.Fatal("Expected error for missing task")
	}
}

func TestStoreGet_Existing(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	got, _ := s.Get(task.ID)
	if got.ID != task.ID || got.Title != task.Title {
		t.Error("Mismatch")
	}
}

func TestStoreGet_CompletedTask(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	s.Update(task.ID, &models.TaskUpdate{Completed: boolPtr(true)})
	got, _ := s.Get(task.ID)
	if !got.Completed {
		t.Error("Expected completed")
	}
}

// 2.1.4 Update

func TestStoreUpdate_Complete(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	before := time.Now()
	updated, _ := s.Update(task.ID, &models.TaskUpdate{Completed: boolPtr(true)})
	if !updated.Completed {
		t.Error("Expected completed")
	}
	if updated.UpdatedAt.Before(before) {
		t.Error("Expected recent UpdatedAt")
	}
	if updated.Title != task.Title {
		t.Error("Expected title unchanged")
	}
}

func TestStoreUpdate_Uncomplete(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	s.Update(task.ID, &models.TaskUpdate{Completed: boolPtr(true)})
	updated, _ := s.Update(task.ID, &models.TaskUpdate{Completed: boolPtr(false)})
	if updated.Completed {
		t.Error("Expected task to be uncompleted")
	}
}

func TestStoreUpdate_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	if _, err := s.Update("nonexistent", &models.TaskUpdate{Completed: boolPtr(true)}); err == nil {
		t.Fatal("Expected error for missing task")
	}
}

// 2.1.5 Delete

func TestStoreDelete_Existing(t *testing.T) {
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

func TestStoreDelete_NotFound(t *testing.T) {
	s := NewInMemoryStore()
	if err := s.Delete("nonexistent"); err == nil {
		t.Fatal("Expected error for missing task")
	}
}

func TestStoreDelete_EmptyStore(t *testing.T) {
	s := NewInMemoryStore()
	if err := s.Delete("any"); err == nil {
		t.Fatal("Expected error for missing task")
	}
}

func TestStoreDelete_AfterList(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	before, _ := s.List(nil)
	if len(before) != 1 {
		t.Fatal("Expected 1 task before delete")
	}
	s.Delete(task.ID)
	after, _ := s.List(nil)
	if len(after) != 0 {
		t.Errorf("Expected 0 tasks after delete, got %d", len(after))
	}
}

func TestStoreDelete_OtherTasksUnaffected(t *testing.T) {
	s := NewInMemoryStore()
	t1 := models.NewTask("Task1", "Desc")
	t2 := models.NewTask("Task2", "Desc")
	s.Create(t1)
	s.Create(t2)
	s.Delete(t1.ID)
	got, _ := s.Get(t2.ID)
	if got == nil || got.ID != t2.ID {
		t.Error("Expected second task to remain")
	}
}

// 2.1.6 Concurrent access safety

func TestStoreConcurrentCreate(t *testing.T) {
	s := NewInMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task := models.NewTask("Concurrent", "Desc")
			if err := s.Create(task); err != nil {
				t.Errorf("Create failed: %v", err)
			}
		}()
	}
	wg.Wait()
	result, err := s.List(nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 100 {
		t.Errorf("Expected 100 tasks, got %d", len(result))
	}
}

func TestStoreConcurrentReadWrite(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 10; i++ {
		s.Create(models.NewTask("T", "D"))
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				task := models.NewTask("RW", "Desc")
				s.Create(task)
			}
		}(i)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				s.List(nil)
			}
		}()
	}
	wg.Wait()
	result, err := s.List(nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 110 {
		t.Errorf("Expected 110 tasks, got %d", len(result))
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
