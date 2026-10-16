package store
import (
	"sync"
	"testing"
	"loom-bootstrap-test-5/models"
)
func TestNewInMemoryStore(t *testing.T) {
	s := NewInMemoryStore()
	if s == nil { t.Fatal("nil store") }
	if s.tasks == nil { t.Error("nil tasks map") }
}
func TestCreate(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	if err := s.Create(task); err != nil { t.Fatalf("unexpected error: %v", err) }
	retrieved, err := s.Get(task.ID)
	if err != nil { t.Fatalf("failed to get task: %v", err) }
	if retrieved.Title != task.Title { t.Errorf("expected %q, got %q", task.Title, retrieved.Title) }
	if retrieved.Completed { t.Error("expected not completed") }
}
func TestCreateValidationErrors(t *testing.T) {
	s := NewInMemoryStore()
	tests := []struct{name string; task *models.Task; wantErr bool}{
		{"empty title", models.NewTask("", "desc"), true},
		{"title too long", models.NewTask(string(make([]byte, 300)), "desc"), true},
		{"desc too long", models.NewTask("title", string(make([]byte, 10001))), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Create(tt.task); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestList(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(models.NewTask("Task 1", "D1"))
	s.Create(models.NewTask("Task 2", "D2"))
	tasks, _ := s.List(nil)
	if len(tasks) != 2 { t.Errorf("expected 2, got %d", len(tasks)) }
}
func TestListWithStatusFilter(t *testing.T) {
	s := NewInMemoryStore()
	t1 := models.NewTask("Active", "D")
	t2 := models.NewTask("Done", "D")
	t2.Completed = true
	s.Create(t1); s.Create(t2)
	active := "active"
	filtered, _ := s.List(&models.TaskFilter{Status: &active})
	if len(filtered) != 1 || filtered[0].Completed { t.Error("expected active") }
	completed := "completed"
	filtered, _ = s.List(&models.TaskFilter{Status: &completed})
	if len(filtered) != 1 || !filtered[0].Completed { t.Error("expected completed") }
}
func TestListWithLimit(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 5; i++ { s.Create(models.NewTask("T", "D")) }
	filtered, _ := s.List(&models.TaskFilter{Limit: 2})
	if len(filtered) != 2 { t.Errorf("expected 2, got %d", len(filtered)) }
}
func TestListEmptyStore(t *testing.T) {
	s := NewInMemoryStore()
	tasks, _ := s.List(nil)
	if len(tasks) != 0 { t.Errorf("expected 0, got %d", len(tasks)) }
}
func TestGet(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	retrieved, err := s.Get(task.ID)
	if err != nil { t.Fatalf("Get() error: %v", err) }
	if retrieved.ID != task.ID { t.Errorf("ID mismatch") }
	_, err = s.Get("nonexistent")
	if err == nil { t.Error("expected error") }
}
func TestUpdate(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	c := true
	updated, err := s.Update(task.ID, &models.TaskUpdate{Completed: &c})
	if err != nil { t.Fatalf("Update() error: %v", err) }
	if !updated.Completed { t.Error("expected completed") }
	_, err = s.Update("nonexistent", &models.TaskUpdate{})
	if err == nil { t.Error("expected error") }
}
func TestReplace(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	newTask := models.NewTask("Replaced", "New")
	replaced, err := s.Replace(task.ID, newTask)
	if err != nil { t.Fatalf("Replace() error: %v", err) }
	if replaced.Title != "Replaced" { t.Errorf("expected Replaced, got %s", replaced.Title) }
	_, err = s.Replace("nonexistent", newTask)
	if err == nil { t.Error("expected error") }
}
func TestDelete(t *testing.T) {
	s := NewInMemoryStore()
	task := models.NewTask("Test", "Desc")
	s.Create(task)
	if err := s.Delete(task.ID); err != nil { t.Fatalf("Delete() error: %v", err) }
	_, err := s.Get(task.ID)
	if err == nil { t.Error("expected error after delete") }
	if err := s.Delete("nonexistent"); err == nil { t.Error("expected error") }
}
func TestHealthCheck(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(models.NewTask("Test", "Desc"))
	stats := s.HealthCheck()
	if stats["status"] != "healthy" { t.Errorf("expected healthy, got %s", stats["status"]) }
	if stats["tasks"] != "1" { t.Errorf("expected 1, got %s", stats["tasks"]) }
	if stats["service"] != "task-api" { t.Errorf("expected task-api, got %s", stats["service"]) }
}
func TestConcurrentAccess(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 10; i++ { s.Create(models.NewTask("T", "D")) }
	done := make(chan bool, 3)
	go func() {
		for i := 0; i < 100; i++ { s.List(nil) }
		done <- true
	}()
	go func() {
		for i := 0; i < 50; i++ { s.Create(models.NewTask("Conc", "Test")) }
		done <- true
	}()
	go func() {
		tasks, _ := s.List(nil)
		for i := 0; i < 50; i++ {
			if len(tasks) > 0 {
				c := true
				s.Update(tasks[i%len(tasks)].ID, &models.TaskUpdate{Completed: &c})
			}
		}
		done <- true
	}()
	for i := 0; i < 3; i++ { <-done }
	tasks, err := s.List(nil)
	if err != nil { t.Fatalf("concurrent access error: %v", err) }
	if len(tasks) == 0 { t.Error("expected tasks after concurrent access") }
}
