package store

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"loom-bootstrap-test-5/models"
)

func newTestStore() *InMemoryStore {
	return NewInMemoryStore()
}

func newTestTask(title string) *models.Task {
	return models.NewTask(title, "test description")
}

func TestCreate_ValidTask(t *testing.T) {
	store := newTestStore()
	task := newTestTask("Test Task")
	err := store.Create(task)
	require.NoError(t, err)
	retrieved, err := store.Get(task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.Title, retrieved.Title)
}

func TestCreate_EmptyTitle(t *testing.T) {
	store := newTestStore()
	task := &models.Task{ID: "test-id", Title: ""}
	err := store.Create(task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestCreate_TitleTooLong(t *testing.T) {
	store := newTestStore()
	task := &models.Task{ID: "test-id", Title: string(make([]byte, 256))}
	err := store.Create(task)
	assert.Error(t, err)
}

func TestList_NoFilter(t *testing.T) {
	store := newTestStore()
	for i := 0; i < 3; i++ {
		store.Create(newTestTask(string(rune('A' + i))))
	}
	tasks, err := store.List(nil)
	require.NoError(t, err)
	assert.Len(t, tasks, 3)
}

func TestList_Filters(t *testing.T) {
	store := newTestStore()
	task1 := newTestTask("Active")
	task2 := newTestTask("Done")
	store.Create(task1)
	store.Create(task2)
	b := true
	store.Update(task2.ID, &models.TaskUpdate{Completed: &b})
	filter := &models.TaskFilter{Status: strPtr("completed")}
	tasks, err := store.List(filter)
	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, task2.Title, tasks[0].Title)
}

func strPtr(s string) *string { return &s }

func TestList_Limit(t *testing.T) {
	store := newTestStore()
	for i := 0; i < 5; i++ {
		store.Create(newTestTask(string(rune('A' + i))))
	}
	filter := &models.TaskFilter{Limit: 3}
	tasks, err := store.List(filter)
	require.NoError(t, err)
	assert.Len(t, tasks, 3)
}

func TestGet_Found(t *testing.T) {
	store := newTestStore()
	task := newTestTask("Get Test")
	store.Create(task)
	retrieved, err := store.Get(task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.Title, retrieved.Title)
}

func TestGet_NotFound(t *testing.T) {
	store := newTestStore()
	_, err := store.Get("non-existent")
	assert.Error(t, err)
}

func TestUpdate_Found(t *testing.T) {
	store := newTestStore()
	task := newTestTask("Update Test")
	store.Create(task)
	b := true
	updated, err := store.Update(task.ID, &models.TaskUpdate{Completed: &b})
	require.NoError(t, err)
	assert.True(t, updated.Completed)
}

func TestUpdate_NotFound(t *testing.T) {
	store := newTestStore()
	b := true
	_, err := store.Update("non-existent", &models.TaskUpdate{Completed: &b})
	assert.Error(t, err)
}

func TestReplace_Found(t *testing.T) {
	store := newTestStore()
	task := newTestTask("Original")
	store.Create(task)
	replaced := models.NewTask("Replaced", "New desc")
	result, err := store.Replace(task.ID, replaced)
	require.NoError(t, err)
	assert.Equal(t, "Replaced", result.Title)
	assert.Equal(t, task.ID, result.ID)
}

func TestDelete_Found(t *testing.T) {
	store := newTestStore()
	task := newTestTask("Delete Me")
	store.Create(task)
	err := store.Delete(task.ID)
	require.NoError(t, err)
	_, err = store.Get(task.ID)
	assert.Error(t, err)
}

func TestConcurrentAccess(t *testing.T) {
	store := newTestStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			task := models.NewTask(string(rune('A' + id%26)), "desc")
			store.Create(task)
			store.Get(task.ID)
		}(i)
	}
	wg.Wait()
	tasks, _ := store.List(nil)
	assert.Len(t, tasks, 100)
}
