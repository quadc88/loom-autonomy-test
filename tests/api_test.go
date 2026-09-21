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

func createTask(t *testing.T, server *httptest.Server, title string) string {
	reqBody, _ := json.Marshal(map[string]string{"title": title})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var task models.Task
	json.NewDecoder(resp.Body).Decode(&task)
	return task.ID
}

func completeTask(t *testing.T, server *httptest.Server, id string) {
	updateBody, _ := json.Marshal(map[string]bool{"completed": true})
	req, _ := http.NewRequest("PATCH", server.URL+"/tasks/"+id, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestListTasks_Empty(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL+"/tasks", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var tasks []models.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestListTasks_All(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id1 := createTask(t, server, "Task 1")
	id2 := createTask(t, server, "Task 2")

	req, _ := http.NewRequest("GET", server.URL+"/tasks", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var tasks []models.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
	ids := make(map[string]bool)
	for _, task := range tasks {
		ids[task.ID] = true
	}
	if !ids[id1] || !ids[id2] {
		t.Error("expected both tasks in response")
	}
}

func TestListTasks_StatusActive(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id1 := createTask(t, server, "Active Task")
	id2 := createTask(t, server, "Completed Task")
	completeTask(t, server, id2)

	req, _ := http.NewRequest("GET", server.URL+"/tasks?status=active", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var tasks []models.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("expected 1 active task, got %d", len(tasks))
	}
	if tasks[0].ID != id1 {
		t.Errorf("expected id %s, got %s", id1, tasks[0].ID)
	}
}

func TestListTasks_StatusCompleted(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id1 := createTask(t, server, "Active Task")
	id2 := createTask(t, server, "Completed Task")
	completeTask(t, server, id2)

	req, _ := http.NewRequest("GET", server.URL+"/tasks?status=completed", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var tasks []models.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("expected 1 completed task, got %d", len(tasks))
	}
	if tasks[0].ID != id2 {
		t.Errorf("expected id %s, got %s", id2, tasks[0].ID)
	}
}

func TestListTasks_InvalidStatusFilter(t *testing.T) {
	h, server := newTestServer(t)
	defer server.Close()
	createTask(t, server, "Task 1")
	createTask(t, server, "Task 2")

	req, _ := http.NewRequest("GET", server.URL+"/tasks?status=invalid", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var tasks []models.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks (fallback), got %d", len(tasks))
	}
	_ = h
}
