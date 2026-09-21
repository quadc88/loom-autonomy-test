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

func TestCreateTask(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id := createTask(t, server, "Test Task")
	if id == "" {
		t.Error("expected non-empty ID")
	}
}

func TestCreateTask_MissingTitle(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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

func TestHealthCheck(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL+"/health", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if result["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %s", result["status"])
	}
	if result["timestamp"] == "" {
		t.Error("expected non-empty timestamp")
	}
	if result["tasks"] != "0" {
		t.Errorf("expected tasks '0', got %s", result["tasks"])
	}
	if result["service"] != "task-api" {
		t.Errorf("expected service 'task-api', got %s", result["service"])
	}
}

func TestHealthCheck_WithTasks(t *testing.T) {
	h, server := newTestServer(t)
	defer server.Close()
	// Create a task
	createTask(t, server, "Health Check Test")
	// Verify health endpoint reflects the task count
	req, _ := http.NewRequest("GET", server.URL+"/health", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if result["tasks"] != "1" {
		t.Errorf("expected tasks '1', got %s", result["tasks"])
	}
	_ = h
}

func TestGetTask_NotFound(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL+"/tasks/999", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestUpdateTask(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id := createTask(t, server, "Update Test")
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
	req, _ = http.NewRequest("GET", server.URL+"/tasks/"+id, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	var task models.Task
	json.NewDecoder(resp.Body).Decode(&task)
	if !task.Completed {
		t.Error("expected task to be completed")
	}
}

func TestGetTask_Success(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id := createTask(t, server, "Get Task Test")
	req, _ := http.NewRequest("GET", server.URL+"/tasks/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var task models.Task
	json.NewDecoder(resp.Body).Decode(&task)
	if task.Title != "Get Task Test" {
		t.Errorf("expected title 'Get Task Test', got %s", task.Title)
	}
	if task.ID == "" {
		t.Error("expected non-empty ID")
	}
}



func TestCreateTask_WithDescription(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{"title": "Test", "description": "A description"})
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
	if task.Description != "A description" {
		t.Errorf("expected description 'A description', got %s", task.Description)
	}
	if task.ID == "" {
		t.Error("expected non-empty ID")
	}
	if task.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestDeleteTask(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id := createTask(t, server, "Delete Test")
	req, _ := http.NewRequest("DELETE", server.URL+"/tasks/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
	req, _ = http.NewRequest("GET", server.URL+"/tasks/"+id, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestListTasks_FilterByCompleted(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	// Create two tasks, complete one
	id1 := createTask(t, server, "Task 1")
	id2 := createTask(t, server, "Task 2")
	updateBody, _ := json.Marshal(map[string]bool{"completed": true})
	req, _ := http.NewRequest("PATCH", server.URL+"/tasks/"+id2, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)

	// Filter by completed
	req, _ = http.NewRequest("GET", server.URL+"/tasks?status=completed", nil)
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

	// Filter by active
	req, _ = http.NewRequest("GET", server.URL+"/tasks?status=active", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	var tasks2 []models.Task
	json.NewDecoder(resp.Body).Decode(&tasks2)
	if len(tasks2) != 1 {
		t.Errorf("expected 1 active task, got %d", len(tasks2))
	}
	if tasks2[0].ID != id1 {
		t.Errorf("expected id %s, got %s", id1, tasks2[0].ID)
	}
}

func TestListTasks_Limit(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	createTask(t, server, "Task 1")
	createTask(t, server, "Task 2")
	createTask(t, server, "Task 3")

	req, _ := http.NewRequest("GET", server.URL+"/tasks?limit=2", nil)
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
}

func TestUpdateTask_FullReplace(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id := createTask(t, server, "Original Title")

	newTask := map[string]interface{}{
		"id":          id,
		"title":       "Updated Title",
		"description": "Updated Description",
		"completed":   true,
	}
	updateBody, _ := json.Marshal(newTask)
	req, _ := http.NewRequest("PUT", server.URL+"/tasks/"+id, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("GET", server.URL+"/tasks/"+id, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	var task models.Task
	json.NewDecoder(resp.Body).Decode(&task)
	if task.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", task.Title)
	}
	if task.Description != "Updated Description" {
		t.Errorf("expected description 'Updated Description', got %s", task.Description)
	}
	if !task.Completed {
		t.Error("expected task to be completed")
	}
}

func TestUpdateTask_PutMissingTitle(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id := createTask(t, server, "Original Title")

	newTask := map[string]interface{}{
		"id":          id,
		"description": "No title",
	}
	updateBody, _ := json.Marshal(newTask)
	req, _ := http.NewRequest("PUT", server.URL+"/tasks/"+id, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestListTasks_FilterAndLimit(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	createTask(t, server, "Task 1")
	id2 := createTask(t, server, "Task 2")
	updateBody, _ := json.Marshal(map[string]bool{"completed": true})
	req, _ := http.NewRequest("PATCH", server.URL+"/tasks/"+id2, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
	createTask(t, server, "Task 3")

	req, _ = http.NewRequest("GET", server.URL+"/tasks?status=completed&limit=1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	var tasks []models.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != id2 {
		t.Errorf("expected id %s, got %s", id2, tasks[0].ID)
	}
}

// CreateTask validation tests
func TestCreateTask_SuccessWithDescription(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{
		"title":       "Test Task",
		"description": "This is a test description",
	})
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
	if task.Title != "Test Task" {
		t.Errorf("expected title 'Test Task', got %s", task.Title)
	}
	if task.Description != "This is a test description" {
		t.Errorf("expected description, got %s", task.Description)
	}
	if task.ID == "" {
		t.Error("expected non-empty ID")
	}
	if task.Completed {
		t.Error("expected task not completed")
	}
	if task.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestCreateTask_TitleTooLong(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	longTitle := string(make([]byte, 256))
	reqBody, _ := json.Marshal(map[string]string{"title": longTitle})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	var errResp models.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errResp)
	if errResp.Error == "" {
		t.Error("expected error message")
	}
}

func TestCreateTask_DescriptionTooLong(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	longDesc := string(make([]byte, 10001))
	reqBody, _ := json.Marshal(map[string]string{
		"title":       "Short Title",
		"description": longDesc,
	})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateTask_WhitespaceOnlyTitle(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	reqBody, _ := json.Marshal(map[string]string{"title": "   "})
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateTask_EmptyBody(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("POST", server.URL+"/tasks", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// Validation error handling tests
func TestGetTask_InvalidUUID(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL+"/tasks/not-a-valid-uuid", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_InvalidUUID(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	req, _ := http.NewRequest("DELETE", server.URL+"/tasks/not-a-valid-uuid", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_InvalidUUID(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	updateBody, _ := json.Marshal(map[string]bool{"completed": true})
	req, _ := http.NewRequest("PATCH", server.URL+"/tasks/not-a-valid-uuid", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

// Duplicate handling test
func TestCreateTask_DuplicateTitleAllowed(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id1 := createTask(t, server, "Same Title")
	id2 := createTask(t, server, "Same Title")
	if id1 == id2 {
		t.Error("expected different IDs for duplicate titles")
	}
}

// PATCH invalid input tests
func TestUpdateTask_PatchInvalidJSON(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	id := createTask(t, server, "Patch Test")
	req, _ := http.NewRequest("PATCH", server.URL+"/tasks/"+id, bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_PatchNotFound(t *testing.T) {
	_, server := newTestServer(t)
	defer server.Close()
	updateBody, _ := json.Marshal(map[string]bool{"completed": true})
	req, _ := http.NewRequest("PATCH", server.URL+"/tasks/nonexistent-id", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

// ListTasks filter tests
func TestListTasks_InvalidStatusFilter(t *testing.T) {
	h, server := newTestServer(t)
	defer server.Close()
	// Create two tasks
	id1 := createTask(t, server, "Task 1")
	id2 := createTask(t, server, "Task 2")
	// Mark one complete
	updateBody, _ := json.Marshal(map[string]bool{"completed": true})
	req, _ := http.NewRequest("PATCH", server.URL+"/tasks/"+id2, bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
	// Filter with invalid status should return all tasks (graceful fallback)
	req, _ = http.NewRequest("GET", server.URL+"/tasks?status=invalid", nil)
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

func TestListTasks_EmptyFilter(t *testing.T) {
	h, server := newTestServer(t)
	defer server.Close()
	createTask(t, server, "Task 1")
	createTask(t, server, "Task 2")
	// No filter should return all tasks
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
	_ = h
}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
