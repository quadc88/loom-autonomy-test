package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"loom-bootstrap-test-5/handlers"
	"loom-bootstrap-test-5/store"
)

func setupHandler() *handlers.TaskHandler {
	s := store.NewInMemoryStore()
	return handlers.NewTaskHandler(s)
}

func TestHealthCheck(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.HandleTasks(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("got Content-Type %q, want application/json", ct)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if body["status"] != "healthy" {
		t.Errorf("got status %q, want %q", body["status"], "healthy")
	}
}

func TestCreateTask(t *testing.T) {
	h := setupHandler()
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"valid", `{"title":"Test","description":"A test"}`, http.StatusCreated},
		{"missing title", `{"description":"No title"}`, http.StatusBadRequest},
		{"empty title", `{"title":""}`, http.StatusBadRequest},
		{"invalid json", `{bad`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.HandleTasks(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	h := setupHandler()
	for _, title := range []string{"A", "B", "C"} {
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"title":`+title+`}`))
		rec := httptest.NewRecorder()
		h.HandleTasks(rec, req)
	}
	tests := []struct {
		url       string
		wantCount int
	}{
		{"/tasks", 3},
		{"/tasks?limit=2", 2},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			h.HandleTasks(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("got %d, want 200", rec.Code)
			}
			var result []map[string]interface{}
			json.NewDecoder(rec.Body).Decode(&result)
			if len(result) != tt.wantCount {
				t.Errorf("got %d, want %d", len(result), tt.wantCount)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"title":"X"}`))
	rec := httptest.NewRecorder()
	h.HandleTasks(rec, req)
	var created struct{ ID string `json:"id"` }
	json.NewDecoder(rec.Body).Decode(&created)

	tests := []struct {
		path       string
		wantStatus int
	}{
		{"/tasks/" + created.ID, http.StatusOK},
		{"/tasks/nonexistent", http.StatusNotFound},
		{"/tasks/health", http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			h.HandleTasks(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"title":"Y"}`))
	rec := httptest.NewRecorder()
	h.HandleTasks(rec, req)
	var created struct{ ID string `json:"id"` }
	json.NewDecoder(rec.Body).Decode(&created)

	tests := []struct {
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"PATCH", "/tasks/" + created.ID, `{"completed":true}`, http.StatusOK},
		{"PATCH", "/tasks/nonexistent", `{"completed":true}`, http.StatusNotFound},
		{"PUT", "/tasks/" + created.ID, `{"id":"`+created.ID+`","title":"Updated","completed":true}`, http.StatusOK},
		{"PUT", "/tasks/" + created.ID, `{"id":"`+created.ID+`","completed":false}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.method+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()
			h.HandleTasks(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"title":"Z"}`))
	rec := httptest.NewRecorder()
	h.HandleTasks(rec, req)
	var created struct{ ID string `json:"id"` }
	json.NewDecoder(rec.Body).Decode(&created)

	tests := []struct {
		path       string
		wantStatus int
	}{
		{"/tasks/" + created.ID, http.StatusNoContent},
		{"/tasks/nonexistent", http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			rec := httptest.NewRecorder()
			h.HandleTasks(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
