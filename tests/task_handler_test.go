package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestHandler_DeleteTask_Success(t *testing.T) {
	h := setupHandler()
	// Create a task first
	reqBody := `{"title":"Test Task","description":"A test"}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(reqBody))
	r.Header.Set("Content-Type", "application/json")
	h.HandleTasks(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 on create, got %d", w.Code)
	}
	// Extract ID from response
	var created map[string]string
	json.NewDecoder(w.Body).Decode(&created)
	id := created["id"]

	// Delete the task
	w = handleRequest(h, http.MethodDelete, "/tasks/"+id)
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204 on delete, got %d", w.Code)
	}

	// Verify task is gone via GET
	w = handleRequest(h, http.MethodGet, "/tasks/"+id)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", w.Code)
	}
}

func TestHandler_DeleteTask_NotFound(t *testing.T) {
	h := setupHandler()
	w := handleRequest(h, http.MethodDelete, "/tasks/nonexistent-id")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent task, got %d", w.Code)
	}
}

func TestCreateTask_MissingTitle_Returns400(t *testing.T) {
	h := setupHandler()
	reqBody := `{"description":"no title"}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(reqBody))
	r.Header.Set("Content-Type", "application/json")
	h.HandleTasks(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing title, got %d", w.Code)
	}
	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["error"] != "title is required" {
		t.Errorf("expected 'title is required' error, got %s", errResp["error"])
	}
}

func TestCreateTask_EmptyTitle_Returns400(t *testing.T) {
	h := setupHandler()
	reqBody := `{"title":""}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(reqBody))
	r.Header.Set("Content-Type", "application/json")
	h.HandleTasks(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty title, got %d", w.Code)
	}
	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["error"] != "title is required" {
		t.Errorf("expected 'title is required' error, got %s", errResp["error"])
	}
}

func TestCreateTask_TitleTooLong_Returns400(t *testing.T) {
	h := setupHandler()
	longTitle := strings.Repeat("a", 256)
	reqBody := `{"title":"` + longTitle + `"}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(reqBody))
	r.Header.Set("Content-Type", "application/json")
	h.HandleTasks(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for title > 255 chars, got %d", w.Code)
	}
	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["error"] != "title must be 255 characters or less" {
		t.Errorf("expected 'title must be 255 characters or less' error, got %s", errResp["error"])
	}
}

func TestCreateTask_DescriptionTooLong_Returns400(t *testing.T) {
	h := setupHandler()
	longDesc := strings.Repeat("a", 10001)
	reqBody := `{"title":"Test","description":"` + longDesc + `"}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(reqBody))
	r.Header.Set("Content-Type", "application/json")
	h.HandleTasks(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for description > 10000 chars, got %d", w.Code)
	}
	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["error"] != "description must be 10000 characters or less" {
		t.Errorf("expected 'description must be 10000 characters or less' error, got %s", errResp["error"])
	}
}

func TestCreateTask_InvalidJSON_Returns400(t *testing.T) {
	h := setupHandler()
	reqBody := `{invalid json}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(reqBody))
	r.Header.Set("Content-Type", "application/json")
	h.HandleTasks(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", w.Code)
	}
	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["error"] != "invalid request body" {
		t.Errorf("expected 'invalid request body' error, got %s", errResp["error"])
	}
}

func TestGetTask_MalformedUUID_Returns400(t *testing.T) {
	h := setupHandler()
	w := handleRequest(h, http.MethodGet, "/tasks/not-a-uuid")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed UUID, got %d", w.Code)
	}
}
	var errResp map[string]string
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["error"] != "task not found" {
		t.Errorf("expected 'task not found' error, got %s", errResp["error"])
	}
}
