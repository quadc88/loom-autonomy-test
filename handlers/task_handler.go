package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"loom-bootstrap-test-5/models"
	"loom-bootstrap-test-5/store"
)

type TaskHandler struct {
	store store.TaskStore
}

func NewTaskHandler(s store.TaskStore) *TaskHandler {
	return &TaskHandler{store: s}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg})
}

func extractID(path string) string {
	return strings.TrimPrefix(path, "/tasks/")
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch r.Method {
	case http.MethodGet:
		if path == "/tasks" || path == "/tasks/" {
			h.ListTasks(w, r)
			return
		}
		id := extractID(path)
		if id != "" {
			h.GetTask(w, r)
			return
		}
		writeError(w, http.StatusBadRequest, "invalid path")
	case http.MethodPost:
		if path == "/tasks" || path == "/tasks/" {
			h.CreateTask(w, r)
			return
		}
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	case http.MethodPatch:
		id := extractID(path)
		if id != "" {
			h.UpdateTask(w, r)
			return
		}
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	case http.MethodDelete:
		id := extractID(path)
		if id != "" {
			h.DeleteTask(w, r)
			return
		}
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task := &models.Task{
		Title:       input.Title,
		Description: input.Description,
	}
	if err := h.store.Create(task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	filter := &models.TaskFilter{}
	status := r.URL.Query().Get("status")
	if status != "" {
		filter.Status = &status
	}
	tasks, err := h.store.List(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "task ID is required")
		return
	}
	task, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "task ID is required")
		return
	}

	var input models.TaskUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.store.Update(id, &input)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "task ID is required")
		return
	}

	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}
