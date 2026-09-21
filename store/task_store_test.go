package store

import (
	"sync"
	"testing"
	"time"

	"loom-bootstrap-test-5/models"
)

// helpers
func makeTask(title string) *models.Task {
	return &models.Task{
		ID:          "test-1",
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

// --- Create ---

func TestCreate_Invalid(t *testing.T) {
	store := NewInMemoryStore()

	// empty title
	err := store.Create(&models.Task{Title: ""})
	if err == nil {
		t.Error("expected error for empty title")
	}

	// title too long
	err = store.Create(&models.Task{Title: "a" + string(make([]byte, 256))})
	if err == nil {
		t.Error("expected error for title > 255 chars")
	}

	// description too long
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

// --- List ---

func TestList_Empty(t *testing.T) {
	store := NewInMemoryStore()
	result, err := store.List(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 tasks, got %d",
