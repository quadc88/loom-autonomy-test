package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Task represents a task in the system.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" validate:"required,max=255"`
	Description string    `json:"description,omitempty" validate:"max=10000"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskFilter represents filtering options for listing tasks.
type TaskFilter struct {
	Status *string
	Limit  int
}

// TaskUpdate represents partial updates for a task.
type TaskUpdate struct {
	Completed *bool
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewTask creates a new Task with a UUID v4 ID.
func NewTask(title, description string) *Task {
	now := time.Now().UTC()
	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ValidateCreate validates a task for creation.
func (t *Task) ValidateCreate() error {
	if t.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(t.Title) > 255 {
		return fmt.Errorf("title must be 255 characters or less")
	}
	if len(t.Description) > 10000 {
		return fmt.Errorf("description must be 10000 characters or less")
	}
	return nil
}

// ValidateUpdate validates a task update.
func (u *TaskUpdate) Validate() error {
	return nil
}

// ValidateReplace validates a task for replacement.
func (t *Task) ValidateReplace() error {
	return t.ValidateCreate()
}
