package models

import (
	"fmt"
	"strings"
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

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type TaskUpdate struct {
	Completed *bool
}

type TaskFilter struct {
	Status *string
	Limit  int
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var idCounter uint64

func generateID() string {
	idCounter++
	temp := idCounter
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = chars[temp%uint64(len(chars))]
		temp /= uint64(len(chars))
	}
	return string(buf)
}

func NewTask(title, description string) *Task {
	now := time.Now().UTC()
	return &Task{
		ID:          generateID(),
		Title:       strings.TrimSpace(title),
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) ValidateCreate() error {
	t.Title = strings.TrimSpace(t.Title)
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

func (t *Task) ValidateReplace() error {
	return t.ValidateCreate()
}
