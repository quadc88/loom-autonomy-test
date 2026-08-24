package models

import (
	"fmt"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskFilter struct {
	Status *string
}

type TaskUpdate struct {
	Completed *bool
}

type ErrorResponse struct {
	Error string `json:"error"`
}

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

func (u *TaskUpdate) Validate() error {
	return nil
}
