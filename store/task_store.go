package store

import (
	"fmt"
	"sync"
	"time"

	"loom-bootstrap-test-5/models"
	"github.com/google/uuid"
)

type TaskStore interface {
	Create(task *models.Task) error
	List(filter *models.TaskFilter) ([]*models.Task, error)
	Get(id string) (*models.Task, error)
	Update(id string, updates *models.TaskUpdate) (*models.Task, error)
	Delete(id string) error
	HealthCheck() map[string]string
}

type InMemoryStore struct {
	mu     sync.RWMutex
	tasks  map[string]*models.Task
	order  []string
	nextID uint64
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		tasks:  make(map[string]*models.Task),
		order:  make([]string, 0),
		nextID: 1,
	}
}

func (s *InMemoryStore) Create(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := task.ValidateCreate(); err != nil {
		return err
	}

	task.ID = uuid.New().String()
	ts := time.Now().UTC()
	task.CreatedAt = ts
	task.UpdatedAt = ts

	s.tasks[task.ID] = task
	s.order = append(s.order, task.ID)
	s.nextID++
	return nil
}

func (s *InMemoryStore) List(filter *models.TaskFilter) ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.Task, 0)
	for _, id := range s.order {
		t, ok := s.tasks[id]
		if !ok {
			continue
		}
		if filter != nil && filter.Status != nil {
			if *filter.Status == "active" && t.Completed {
				continue
			}
			if *filter.Status == "completed" && !t.Completed {
				continue
			}
		}
		result = append(result, t)
	}
	return result, nil
}

func (s *InMemoryStore) Get(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	return t, nil
}

func (s *InMemoryStore) Update(id string, updates *models.TaskUpdate) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}

	if updates.Completed != nil {
		t.Completed = *updates.Completed
	}
	t.UpdatedAt = time.Now().UTC()
	return t, nil
}

func (s *InMemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found")
	}
	delete(s.tasks, id)
	for i, oid := range s.order {
		if oid == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
	return nil
}

func (s *InMemoryStore) HealthCheck() map[string]string {
	s.mu.RLock()
	count := len(s.tasks)
	s.mu.RUnlock()
	return map[string]string{
		"status":  "healthy",
		"tasks":   fmt.Sprintf("%d", count),
		"service": "task-api",
	}
}
