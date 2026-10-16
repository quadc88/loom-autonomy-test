package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"loom-bootstrap-test-5/handlers"
	"loom-bootstrap-test-5/store"
)

// Unit tests for TaskHandler using ResponseRecorder directly
// These complement api_test.go which uses httptest.Server

func setupHandler() *handlers.TaskHandler {
	s := store.NewInMemoryStore()
	return handlers.NewTaskHandler(s)
}

func handleRequest(h *handlers.TaskHandler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)
	h.HandleTasks(w, r)
	return w
}

func TestHandler_NonExistentRoute(t *testing.T) {
	h := setupHandler()
	w := handleRequest(h, http.MethodGet, "/unknown")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// HealthCheck is routed via separate handler in api_test.go, not via HandleTasks
// This test verifies HandleTasks rejects /health as invalid task ID
func TestHandler_HealthViaHandleTasks(t *testing.T) {
	h := setupHandler()
	w := handleRequest(h, http.MethodGet, "/health")
	// /health routed via separate handler; HandleTasks treats it as unknown task -> 404
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandler_PostTaskList(t *testing.T) {
	h := setupHandler()
	w := handleRequest(h, http.MethodPost, "/tasks")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty body, got %d", w.Code)
	}
}
