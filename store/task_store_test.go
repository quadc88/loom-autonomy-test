package store

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"loom-bootstrap-test-5/models"
)

func makeTask(title string) *models.Task {
	return &models.Task{
		ID:          fmt.Sprintf("task-%s", title),
		Title:       title,
		Description: "test description",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func makeTaskCompleted() *models.Task {
	t := makeTask("completed task")
	t.Completed = true
	return t
}

func TestCreate_Invalid(t *testing.T) {
	store := NewInMemoryStore()
	err := store.Create(&models.Task{Title: ""})
	if err == nil {
		t.Error("expected error for empty title")
	}
	err = store.Create(&models.Task{Title: "a" + string(make([]byte, 256))})
	if err == nil {
		t.Error("expected error for title > 255 chars")
	}
	err = store.Create(&models.Task{Title: "ok", Description: string(make([]byte, 10001))})
	if err == nil {
		t.Error("expected error for description > 10000 chars")
	}
}

func TestCreate_Valid(t *testing.T) {
	store := NewInMemoryStore()
	task := makeTask("hello")
	err := store.Create(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := store.Get(task.ID)
	if err != nil {
		t.Fatalf("expected to retrieve task: %v", err)
	}
	if got.Title != "hello" {
		t.Errorf("expected title hello, got %s", got.Title)
	}
}

func TestList_Empty(t *testing.T) {
	store := NewInMemoryStore()
	result, err := store.List(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(result))
	}
}

func TestList_WithTasks(t *testing.T) {
	store := NewInMemoryStore()
	store.Create(makeTask("task1"))
	store.Create(makeTask("task2"))
	result, err := store.List(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(result))
	}
}

func TestList_FilterCompleted(t *testing.T) {
	store := NewInMemoryStore()
	store.Create(makeTask("active"))
	store.Create(makeTaskCompleted())
	status := "completed"
	result, err := store.List(&models.TaskFilter{Status: &status})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 completed task, got %d", len(result))
	}
}

func TestList_Limit(t *testing.T) {
	store := NewInMemoryStore()
	store.Create(makeTask("task1"))
	store.Create(makeTask("task2"))
	store.Create(makeTask("task3"))
	result, err := store.List(&models.TaskFilter{Limit: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(result))
	}
}

func TestGet_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	_, err := store.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestGet_Success(t *testing.T) {
	store := NewInMemoryStore()
	task := makeTask("get test")
	store.Create(task)
	got, err := store.Get(task.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "get test" {
		t.Errorf("expected title get test, got %s", got.Title)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	flag := true
	_, err := store.Update("nonexistent", &models.TaskUpdate{Completed: &flag})
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestUpdate_Success(t *testing.T) {
	store := NewInMemoryStore()
	task := makeTask("update test")
	store.Create(task)
	flag := true
	updated, err := store.Update(task.ID, &models.TaskUpdate{Completed: &flag})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated.Completed {
		t.Error("expected task to be completed")
	}
}

func TestDelete_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	err := store.Delete("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestDelete_Success(t *testing.T) {
	store := NewInMemoryStore()
	task := makeTask("delete test")
	store.Create(task)
	err := store.Delete(task.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = store.Get(task.ID)
	if err == nil {
		t.Error("expected task to be deleted")
	}
}

func TestHealthCheck(t *testing.T) {
	store := NewInMemoryStore()
	stats := store.HealthCheck()
	if stats["status"] != "healthy" {
		t.Errorf("expected status healthy, got %s", stats["status"])
	}
	if stats["tasks"] != "0" {
		t.Errorf("expected tasks 0, got %s", stats["tasks"])
	}
	if stats["service"] != "task-api" {
		t.Errorf("expected service task-api, got %s", stats["service"])
	}
}

func TestHealthCheck_WithTasks(t *testing.T) {
	store := NewInMemoryStore()
	store.Create(makeTask("task1"))
	store.Create(makeTask("task2"))
	stats := store.HealthCheck()
	if stats["tasks"] != "2" {
		t.Errorf("expected tasks 2, got %s", stats["tasks"])
	}
}

func TestReplace_NotFound(t *testing.T) {
	store := NewInMemoryStore()
	task := makeTask("replace test")
	_, err := store.Replace("nonexistent", task)
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestReplace_Success(t *testing.T) {
	store := NewInMemoryStore()
	original := makeTask("original")
	store.Create(original)

	replaced := &models.Task{
		ID:          original.ID,
		Title:       "replaced",
		Description: "new description",
		Completed:   true,
		CreatedAt:   original.CreatedAt,
		UpdatedAt:   time.Now().UTC(),
	}
	result, err := store.Replace(original.ID, replaced)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Title != "replaced" {
		t.Errorf("expected title 'replaced', got %s", result.Title)
	}
	if !result.Completed {
		t.Error("expected task to be completed")
	}
	// Original creation time should be preserved
	if result.CreatedAt != original.CreatedAt {
		t.Error("expected CreatedAt to be preserved")
	}
}

func TestReplace_Invalid(t *testing.T) {
	store := NewInMemoryStore()
	task := makeTask("orig")
	store.Create(task)

	invalid := &models.Task{
		ID:          task.ID,
		Title:       "", // empty title is invalid
		Description: "test",
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	_, err := store.Replace(task.ID, invalid)
	if err == nil {
		t.Error("expected error for invalid replacement")
	}
}

func TestConcurrency_CreateRead(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup
	// Concurrent creates
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			task := &models.Task{
				ID:          fmt.Sprintf("cr-%d", n),
				Title:       fmt.Sprintf("task %d", n),
				Description: "test",
				Completed:   false,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			}
			store.Create(task)
		}(i)
	}
	wg.Wait()
	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			store.Get(fmt.Sprintf("cr-%d", n))
		}(i)
	}
	wg.Wait()
	result, err := store.List(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 100 {
		t.Errorf("expected 100 tasks, got %d", len(result))
	}
}

func TestConcurrency_MixedOperations(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("mix-%d", n)
			task := &models.Task{
				ID:          id,
				Title:       fmt.Sprintf("task %d", n),
				Description: "test",
				Completed:   false,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			}
			store.Create(task)
			if n%2 == 0 {
				flag := true
				store.Update(id, &models.TaskUpdate{Completed: &flag})
			}
			if n%3 == 0 {
				store.Delete(id)
			}
		}(i)
	}
	wg.Wait()
	result, err := store.List(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) > 50 {
		t.Errorf("expected at most 50 tasks, got %d", len(result))
	}
}
