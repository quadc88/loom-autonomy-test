package store

import (
	"strings"
	"testing"

	"loom-bootstrap-test-5/models"
)

func TestInMemoryStore_Create(t *testing.T) {
	s := NewInMemoryStore()
	t.Run("valid", func(t *testing.T) {
		task := &models.Task{ID: "c1", Title: "Test"}
		if err := s.Create(task); err != nil {
			t.Fatalf("unexpected: %v", err)
		}
		got, _ := s.Get("c1")
		if got.Title != "Test" {
			t.Errorf("got %q", got.Title)
		}
	})
	t.Run("empty title", func(t *testing.T) {
		task := &models.Task{ID: "c2", Title: ""}
		if err := s.Create(task); err == nil {
			t.Fatal("expected error")
		} else if !strings.Contains(err.Error(), "title") {
			t.Errorf("wrong err: %v", err)
		}
	})
	t.Run("long title", func(t *testing.T) {
		task := &models.Task{ID: "c3", Title: string(make([]byte, 256))}
		if s.Create(task) == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("long desc", func(t *testing.T) {
		task := &models.Task{ID: "c4", Title: "ok", Description: string(make([]byte, 10001))}
		if s.Create(task) == nil {
			t.Fatal("expected error")
		}
	})
}

func TestInMemoryStore_Get(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(&models.Task{ID: "g1", Title: "Found"})
	t.Run("exists", func(t *testing.T) {
		got, _ := s.Get("g1")
		if got.Title != "Found" {
			t.Errorf("got %q", got.Title)
		}
	})
	t.Run("not found", func(t *testing.T) {
		if _, err := s.Get("missing"); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestInMemoryStore_List(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(&models.Task{ID: "l1", Title: "A", Completed: false})
	s.Create(&models.Task{ID: "l2", Title: "B", Completed: true})
	s.Create(&models.Task{ID: "l3", Title: "C", Completed: false})
	t.Run("all", func(t *testing.T) {
		tasks, _ := s.List(nil)
		if len(tasks) != 3 {
			t.Errorf("count: got %d", len(tasks))
		}
	})
	t.Run("active", func(t *testing.T) {
		st := "active"
		tasks, _ := s.List(&models.TaskFilter{Status: &st})
		if len(tasks) != 2 {
			t.Errorf("count: got %d", len(tasks))
		}
	})
	t.Run("completed", func(t *testing.T) {
		st := "completed"
		tasks, _ := s.List(&models.TaskFilter{Status: &st})
		if len(tasks) != 1 || tasks[0].Title != "B" {
			t.Errorf("count=%d title=%s", len(tasks), tasks[0].Title)
		}
	})
	t.Run("limit", func(t *testing.T) {
		tasks, _ := s.List(&models.TaskFilter{Limit: 2})
		if len(tasks) != 2 {
			t.Errorf("count: got %d", len(tasks))
		}
	})
}

func TestInMemoryStore_Update(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(&models.Task{ID: "u1", Title: "Org"})
	t.Run("complete", func(t *testing.T) {
		d := true
		got, _ := s.Update("u1", &models.TaskUpdate{Completed: &d})
		if !got.Completed {
			t.Error("expected completed")
		}
	})
	t.Run("uncomplete", func(t *testing.T) {
		d := false
		got, _ := s.Update("u1", &models.TaskUpdate{Completed: &d})
		if got.Completed {
			t.Error("expected not completed")
		}
	})
	t.Run("missing", func(t *testing.T) {
		d := true
		if _, err := s.Update("missing", &models.TaskUpdate{Completed: &d}); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestInMemoryStore_Replace(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(&models.Task{ID: "r1", Title: "Old"})
	t.Run("replace", func(t *testing.T) {
		newT := &models.Task{ID: "r1", Title: "New"}
		got, _ := s.Replace("r1", newT)
		if got.Title != "New" {
			t.Errorf("got %q", got.Title)
		}
	})
	t.Run("missing", func(t *testing.T) {
		if _, err := s.Replace("missing", &models.Task{ID: "r1", Title: "X"}); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestInMemoryStore_Delete(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(&models.Task{ID: "d1", Title: "ToDelete"})
	t.Run("delete", func(t *testing.T) {
		if err := s.Delete("d1"); err != nil {
			t.Fatalf("unexpected: %v", err)
		}
		if _, err := s.Get("d1"); err == nil {
			t.Fatal("expected not found")
		}
	})
	t.Run("delete missing", func(t *testing.T) {
		if err := s.Delete("missing"); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestInMemoryStore_HealthCheck(t *testing.T) {
	s := NewInMemoryStore()
	s.Create(&models.Task{ID: "h1", Title: "A"})
	s.Create(&models.Task{ID: "h2", Title: "B"})
	stats := s.HealthCheck()
	if stats["status"] != "healthy" {
		t.Errorf("status: got %q", stats["status"])
	}
	if stats["tasks"] != "2" {
		t.Errorf("tasks: got %q", stats["tasks"])
	}
	if stats["service"] != "task-api" {
		t.Errorf("service: got %q", stats["service"])
	}
}
