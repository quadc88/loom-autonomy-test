package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"loom-bootstrap-test-5/handlers"
	"loom-bootstrap-test-5/models"
	"loom-bootstrap-test-5/store"
)

func newTestServer(t *testing.T) (*handlers.TaskHandler, *httptest.Server) {
	s := store.NewInMemoryStore()
	h := handlers.NewTaskHandler(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.HandleTasks)
	mux.HandleFunc("/tasks/", h.HandleTasks)
	mux.HandleFunc("/health", h.HealthCheck)
	server := httptest.NewServer(mux)
	return h, server
}

func TestCreateTask(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{"title": "Test Task"})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var task models.Task
	json.NewDecoder(resp.Body).Decode(&task)
	if task.Title != "Test Task" { t.Errorf("expected title 'Test Task', got %s", task.Title) }
	if task.ID == "" { t.Error("expected non-empty ID") }
}

func TestCreateTask_MissingTitle(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestListTasks(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL+"/tasks", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var tasks []models.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 0 { t.Errorf("expected 0 tasks, got %d", len(tasks)) }
}

func TestHealthCheck(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL+"/health", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL+"/tasks/999", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("DELETE", server.URL+"/tasks/999", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestUpdateTask(t *testing.T) {
	h, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{"title": "Test Task"})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	h.HandleTasks(nil, req)
	req, _ = http.NewRequest("GET", server.URL+"/tasks/1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	var task models.Task
	json.NewDecoder(resp.Body).Decode(&task)
	updateBody, _ := json.Marshal(map[string]bool{"completed": true})
	req, _ = http.NewRequest("PATCH", server.URL+"/tasks/"+task.ID, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	h.HandleTasks(nil, req)
	req, _ = http.NewRequest("GET", server.URL+"/tasks/"+task.ID, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&task)
	if !task.Completed { t.Error("expected task to be completed") }
}

func TestGetTask_Success(t *testing.T) {
	h, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{"title": "Get Task Test"})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	h.HandleTasks(nil, req)

	req, _ = http.NewRequest("GET", server.URL+"/tasks/1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var task models.Task
	json.NewDecoder(resp.Body).Decode(&task)
	if task.Title != "Get Task Test" { t.Errorf("expected title 'Get Task Test', got %s", task.Title) }
	if task.ID == "" { t.Error("expected non-empty ID") }
}

func TestDeleteTask(t *testing.T) {
	h, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{"title": "Test Task"})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	h.HandleTasks(nil, req)
	req, _ = http.NewRequest("DELETE", server.URL+"/tasks/1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
	req, _ = http.NewRequest("GET", server.URL+"/tasks/1", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
