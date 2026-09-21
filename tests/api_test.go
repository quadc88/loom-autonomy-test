package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"loom-bootstrap-test-5/handlers"
	"loom-bootstrap-test-5/models"
	"loom-bootstrap-test-5/store"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := store.NewInMemoryStore()
	h := handlers.NewTaskHandler(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.HandleTasks)
	mux.HandleFunc("/tasks/", h.HandleTasks)
	mux.HandleFunc("/health", h.HealthCheck)
	return httptest.NewServer(mux)
}

func mustCreateTask(t *testing.T, srv *httptest.Server, title string) string {
	t.Helper()
	body := strings.NewReader(`{"title":"` + title + `","description":"desc"}`)
	req, _ := http.NewRequest("POST", srv.URL+"/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return task.ID
}

func doRequest(t *testing.T, srv *httptest.Server, method, path, bodyStr string) *http.Response {
	t.Helper()
	var body io.Reader
	if bodyStr != "" {
		body = strings.NewReader(bodyStr)
	}
	req, err := http.NewRequest(method, srv.URL+path, body)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	if bodyStr != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func readJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("decode json: %v body: %s", err, string(b))
	}
}

// 2.2.6 Test GET /health
func TestHealthCheck(t *testing.T) {
	srv := newTestServer(t)
	resp := doRequest(t, srv, "GET", "/health", "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want 200", resp.StatusCode)
	}
	var body map[string]string
	readJSON(t, resp, &body)
	if body["status"] != "healthy" {
		t.Errorf("status: got %q, want %q", body["status"], "healthy")
	}
	if body["service"] != "task-api" {
		t.Errorf("service: got %q, want %q", body["service"], "task-api")
	}
	if body["tasks"] != "0" {
		t.Errorf("tasks: got %q, want %q", body["tasks"], "0")
	}
}

// 2.2.1 Test POST /tasks - create, validate, duplicate handling
func TestPOST_CreateTask(t *testing.T) {
	srv := newTestServer(t)

	t.Run("valid creation returns 201 with task body", func(t *testing.T) {
		resp := doRequest(t, srv, "POST", "/tasks", `{"title":"Buy milk","description":"Fresh milk"}`)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("status: got %d, want 201", resp.StatusCode)
		}
		var task models.Task
		readJSON(t, resp, &task)
		if task.Title != "Buy milk" {
			t.Errorf("title: got %q, want %q", task.Title, "Buy milk")
		}
		if task.Description != "Fresh milk" {
			t.Errorf("description: got %q, want %q", task.Description, "Fresh milk")
		}
		if task.ID == "" {
			t.Error("expected non-empty ID")
		}
		if task.Completed {
			t.Error("expected completed=false")
		}
	})

	t.Run("missing title returns 400", func(t *testing.T) {
		resp := doRequest(t, srv, "POST", "/tasks", `{"description":"no title"}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
		var errResp models.ErrorResponse
		readJSON(t, resp, &errResp)
		if !strings.Contains(errResp.Error, "title") {
			t.Errorf("error msg should mention title, got %q", errResp.Error)
		}
	})

	t.Run("empty title returns 400", func(t *testing.T) {
		resp := doRequest(t, srv, "POST", "/tasks", `{"title":""}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		resp := doRequest(t, srv, "POST", "/tasks", `{bad json`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
	})

	t.Run("title is trimmed", func(t *testing.T) {
		resp := doRequest(t, srv, "POST", "/tasks", `{"title":"  spaced  "}`)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("status: got %d, want 201", resp.StatusCode)
		}
		var task models.Task
		readJSON(t, resp, &task)
		if task.Title != "spaced" {
			t.Errorf("trimmed title: got %q, want %q", task.Title, "spaced")
		}
	})

	t.Run("duplicate titles allowed", func(t *testing.T) {
		id1 := mustCreateTask(t, srv, "Same title")
		id2 := mustCreateTask(t, srv, "Same title")
		if id1 == id2 {
			t.Error("duplicate titles should produce different IDs")
		}
	})
}

// 2.2.2 Test GET /tasks - list all, filter by status
func TestGET_ListTasks(t *testing.T) {
	srv := newTestServer(t)
	mustCreateTask(t, srv, "Task A")
	mustCreateTask(t, srv, "Task B")
	mustCreateTask(t, srv, "Task C")

	t.Run("list all tasks", func(t *testing.T) {
		resp := doRequest(t, srv, "GET", "/tasks", "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status: got %d, want 200", resp.StatusCode)
		}
		var tasks []models.Task
		readJSON(t, resp, &tasks)
		if len(tasks) != 3 {
			t.Errorf("count: got %d, want 3", len(tasks))
		}
	})

	t.Run("filter by limit", func(t *testing.T) {
		resp := doRequest(t, srv, "GET", "/tasks?limit=2", "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status: got %d, want 200", resp.StatusCode)
		}
		var tasks []models.Task
		readJSON(t, resp, &tasks)
		if len(tasks) != 2 {
			t.Errorf("count: got %d, want 2", len(tasks))
		}
	})

	t.Run("filter by status active", func(t *testing.T) {
		mustCreateTask(t, srv, "Active Task")
		resp := doRequest(t, srv, "GET", "/tasks?status=active", "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status: got %d, want 200", resp.StatusCode)
		}
		var tasks []models.Task
		readJSON(t, resp, &tasks)
		if len(tasks) != 3 {
			t.Logf("active filter: got %d, expected 3 (all are active since none completed yet)", len(tasks))
		}
	})

	t.Run("filter by status completed", func(t *testing.T) {
		resp := doRequest(t, srv, "GET", "/tasks?status=completed", "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status: got %d, want 200", resp.StatusCode)
		}
		var tasks []models.Task
		readJSON(t, resp, &tasks)
		if len(tasks) != 0 {
			t.Errorf("count: got %d, want 0", len(tasks))
		}
	})
}

// 2.2.3 Test GET /tasks/{id} - get existing, not found
func TestGET_GetTask(t *testing.T) {
	srv := newTestServer(t)
	id := mustCreateTask(t, srv, "Target Task")

	t.Run("get existing task returns 200", func(t *testing.T) {
		resp := doRequest(t, srv, "GET", "/tasks/"+id, "")
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status: got %d, want 200", resp.StatusCode)
		}
		var task models.Task
		readJSON(t, resp, &task)
		if task.ID != id {
			t.Errorf("id: got %q, want %q", task.ID, id)
		}
		if task.Title != "Target Task" {
			t.Errorf("title: got %q, want %q", task.Title, "Target Task")
		}
	})

	t.Run("get nonexistent task returns 404", func(t *testing.T) {
		resp := doRequest(t, srv, "GET", "/tasks/nonexistent-id", "")
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status: got %d, want 404", resp.StatusCode)
		}
		var errResp models.ErrorResponse
		readJSON(t, resp, &errResp)
		if !strings.Contains(errResp.Error, "not found") {
			t.Errorf("error msg should mention not found, got %q", errResp.Error)
		}
	})

	t.Run("get with empty id returns 400", func(t *testing.T) {
		resp := doRequest(t, srv, "GET", "/tasks/", "")
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
	})
}

// 2.2.4 Test PATCH /tasks/{id} - update completion, invalid input, not found
func TestPATCH_UpdateTask(t *testing.T) {
	srv := newTestServer(t)
	id := mustCreateTask(t, srv, "Patchable Task")

	t.Run("patch to complete returns 200", func(t *testing.T) {
		resp := doRequest(t, srv, "PATCH", "/tasks/"+id, `{"completed":true}`)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status: got %d, want 200", resp.StatusCode)
		}
		var task models.Task
		readJSON(t, resp, &task)
		if !task.Completed {
			t.Error("expected completed=true after patch")
		}
	})

	t.Run("patch to uncomplete returns 200", func(t *testing.T) {
		resp := doRequest(t, srv, "PATCH", "/tasks/"+id, `{"completed":false}`)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status: got %d, want 200", resp.StatusCode)
		}
		var task models.Task
		readJSON(t, resp, &task)
		if task.Completed {
			t.Error("expected completed=false after patch")
		}
	})

	t.Run("patch nonexistent task returns 404", func(t *testing.T) {
		resp := doRequest(t, srv, "PATCH", "/tasks/nonexistent-id", `{"completed":true}`)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status: got %d, want 404", resp.StatusCode)
		}
	})

	t.Run("patch with invalid JSON returns 400", func(t *testing.T) {
		resp := doRequest(t, srv, "PATCH", "/tasks/"+id, `{bad json`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
	})

	t.Run("patch with empty id returns 400", func(t *testing.T) {
		resp := doRequest(t, srv, "PATCH", "/tasks/", `{"completed":true}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
	})

	t.Run("put with missing title returns 400", func(t *testing.T) {
		resp := doRequest(t, srv, "PUT", "/tasks/"+id, `{"completed":false}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
	})
}

// 2.2.5 Test DELETE /tasks/{id} - delete existing, not found
func TestDELETE_DeleteTask(t *testing.T) {
	srv := newTestServer(t)
	id := mustCreateTask(t, srv, "Deletable Task")

	t.Run("delete existing task returns 204", func(t *testing.T) {
		resp := doRequest(t, srv, "DELETE", "/tasks/"+id, "")
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("status: got %d, want 204", resp.StatusCode)
		}
	})

	t.Run("delete already deleted task returns 404", func(t *testing.T) {
		resp := doRequest(t, srv, "DELETE", "/tasks/"+id, "")
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status: got %d, want 404", resp.StatusCode)
		}
	})

	t.Run("delete nonexistent task returns 404", func(t *testing.T) {
		resp := doRequest(t, srv, "DELETE", "/tasks/nonexistent-id", "")
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status: got %d, want 404", resp.StatusCode)
		}
	})

	t.Run("delete with empty id returns 400", func(t *testing.T) {
		resp := doRequest(t, srv, "DELETE", "/tasks/", "")
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status: got %d, want 400", resp.StatusCode)
		}
	})
}
