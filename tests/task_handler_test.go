package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"loom-bootstrap-test-5/handlers"
	"loom-bootstrap-test-5/models"
	"loom-bootstrap-test-5/store"
)

func setupHandler() *handlers.TaskHandler {
	s := store.NewInMemoryStore()
	return handlers.NewTaskHandler(s)
}

func doRequest(h *handlers.TaskHandler, method, path string, body any) *httptest.ResponseRecorder {
	var buf io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewReader(b)
	} else {
		buf = http.NoBody
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, buf)
	r.Header.Set("Content-Type", "application/json")
	h.HandleTasks(w, r)
	return w
}

func expectStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Errorf("expected status %d, got %d: body=%s", expected, w.Code, w.Body.String())
	}
}

func parseJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var result map[string]any
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func TestHealthCheck(t *testing.T) {
	h := setupHandler()
	w := doRequest(h, http.MethodGet, "/health", nil)
	expectStatus(t, w, http.StatusOK)
}

func TestCreateTask_Success(t *testing.T) {
	h := setupHandler()
	body := models.CreateTaskRequest{Title: "Test Task"}
	w := doRequest(h, http.MethodPost, "/tasks", body)
	expectStatus(t, w, http.StatusCreated)
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	h := setupHandler()
	body := models.CreateTaskRequest{Title: ""}
	w := doRequest(h, http.MethodPost, "/tasks", body)
	expectStatus(t, w, http.StatusBadRequest)
}

func TestListTasks_Empty(t *testing.T) {
	h := setupHandler()
	w := doRequest(h, http.MethodGet, "/tasks", nil)
	expectStatus(t, w, http.StatusOK)
}

func TestGetTask_NotFound(t *testing.T) {
	h := setupHandler()
	w := doRequest(h, http.MethodGet, "/tasks/non-existent", nil)
	expectStatus(t, w, http.StatusNotFound)
}

func TestDeleteTask_NotFound(t *testing.T) {
	h := setupHandler()
	w := doRequest(h, http.MethodDelete, "/tasks/non-existent", nil)
	expectStatus(t, w, http.StatusNotFound)
}

func TestUpdateTask_NotFound(t *testing.T) {
	h := setupHandler()
	w := doRequest(h, http.MethodPatch, "/tasks/non-existent", map[string]bool{"completed": true})
	expectStatus(t, w, http.StatusNotFound)
}

func TestInvalidMethod(t *testing.T) {
	h := setupHandler()
	w := doRequest(h, http.MethodDelete, "/tasks", nil)
	expectStatus(t, w, http.StatusMethodNotAllowed)
}

func TestListWithFilter(t *testing.T) {
	h := setupHandler()
	body := models.CreateTaskRequest{Title: "Task 1"}
	doRequest(h, http.MethodPost, "/tasks", body)
	w := doRequest(h, http.MethodGet, "/tasks?status=active", nil)
	expectStatus(t, w, http.StatusOK)
}

func TestListWithLimit(t *testing.T) {
	h := setupHandler()
	for i := 0; i < 3; i++ {
		b := models.CreateTaskRequest{Title: "Task " + string(rune('1'+i))}
		doRequest(h, http.MethodPost, "/tasks", b)
	}
	w := doRequest(h, http.MethodGet, "/tasks?limit=2", nil)
	expectStatus(t, w, http.StatusOK)
}

func TestCreateAndGet(t *testing.T) {
	h := setupHandler()
	body := models.CreateTaskRequest{Title: "Get Me"}
	w1 := doRequest(h, http.MethodPost, "/tasks", body)
	result := parseJSON(t, w1)
	id, ok := result["id"].(string)
	if !ok || id == "" {
		t.Fatal("missing id")
	}
	w2 := doRequest(h, http.MethodGet, "/tasks/"+id, nil)
	expectStatus(t, w2, http.StatusOK)
	result = parseJSON(t, w2)
	if result["title"] != "Get Me" {
		t.Errorf("title mismatch")
	}
}

func TestCreateAndComplete(t *testing.T) {
	h := setupHandler()
	body := models.CreateTaskRequest{Title: "Complete Me"}
	w1 := doRequest(h, http.MethodPost, "/tasks", body)
	result := parseJSON(t, w1)
	id := result["id"].(string)
	w2 := doRequest(h, http.MethodPatch, "/tasks/"+id, map[string]bool{"completed": true})
	expectStatus(t, w2, http.StatusOK)
	result = parseJSON(t, w2)
	if !result["completed"].(bool) {
		t.Error("task not completed")
	}
}

func TestCreateAndDelete(t *testing.T) {
	h := setupHandler()
	body := models.CreateTaskRequest{Title: "Delete Me"}
	w1 := doRequest(h, http.MethodPost, "/tasks", body)
	result := parseJSON(t, w1)
	id := result["id"].(string)
	w2 := doRequest(h, http.MethodDelete, "/tasks/"+id, nil)
	expectStatus(t, w2, http.StatusNoContent)
	w3 := doRequest(h, http.MethodGet, "/tasks/"+id, nil)
	expectStatus(t, w3, http.StatusNotFound)
}
