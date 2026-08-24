# System Requirements Document (SRD) - Task REST API

**Version:** 1.0
**Status:** Draft
**Date:** 2026-07-10
**Author:** Agnes (Sapiens AI)
**Based on:** plans/ORIGINAL_PRD.md

---

## 1. System Architecture

### 1.1 High-Level Architecture
Client Applications -> API Server -> In-Memory Store

Components:
- main.go: Application entry point and server initialization
- handlers/: HTTP request handlers for all endpoints
- models/: Data structures (Task struct, request/response types)
- store/: In-memory data store with concurrency control
- middleware/: CORS, panic recovery, logging middleware

### 1.2 Technology Stack
- Language: Go 1.21+
- Framework: net/http standard library
- Concurrency: sync.RWMutex for thread-safe store access
- Serialization: encoding/json
- Testing: testing package with httptest

### 1.3 Deployment
- Local: go run . (listens on :8080)
- Docker: Multi-stage build, Alpine Linux base image
- Health Check: GET /health endpoint

---

## 2. Data Models

### 2.1 Task Entity
```json
{
  "id": "string (UUID v4)",
  "title": "string (required, max 255 chars)",
  "description": "string (optional, max 10000 chars)",
  "completed": "boolean (default: false)",
  "created_at": "ISO 8601 UTC timestamp",
  "updated_at": "ISO 8601 UTC timestamp"
}
```

### 2.2 Validation Rules
- Title: required, non-empty, max 255 characters
- Description: optional, max 10000 characters
- ID: auto-generated UUID v4
- Timestamps: auto-generated on create/update

### 2.3 Store Interface
```go
type TaskStore interface {
    Create(task *Task) error
    List(filter *TaskFilter) ([]*Task, error)
    Get(id string) (*Task, error)
    Update(id string, updates *TaskUpdate) (*Task, error)
    Delete(id string) error
    HealthCheck() map[string]string
}
```

---

## 3. API Endpoints

### 3.1 POST /tasks - Create Task
- Auth: None
- Request: {"title": "string", "description": "string?"}
- Success: 201 Created + Task JSON
- Error: 400 Bad Request + {"error": "message"}

### 3.2 GET /tasks - List Tasks
- Auth: None
- Query Params: ?status=active|completed (optional)
- Success: 200 OK + [Task] array
- Error: N/A

### 3.3 GET /tasks/{id} - Get Task
- Auth: None
- Success: 200 OK + Task JSON
- Error: 404 Not Found + {"error": "task not found"}

### 3.4 PATCH /tasks/{id} - Update Task
- Auth: None
- Request: {"completed": boolean}
- Success: 200 OK + Updated Task JSON
- Error: 400 Bad Request, 404 Not Found

### 3.5 DELETE /tasks/{id} - Delete Task
- Auth: None
- Success: 204 No Content
- Error: 404 Not Found + {"error": "task not found"}

### 3.6 GET /health - Health Check
- Auth: None
- Success: 200 OK + {"status": "healthy"}
- Response Time: < 10ms

---

## 4. Integration Points

### 4.1 Internal
- Handlers call Store interface methods
- Middleware wraps handlers
- Models used across all components

### 4.2 Future (Out of Scope)
- PostgreSQL database driver
- Redis for caching
- Prometheus metrics endpoint
- JWT authentication
- GraphQL API

---

## 5. Infrastructure Requirements

### 5.1 Containerization
- Dockerfile with multi-stage build
- Base image: golang:1.21-alpine (build), alpine:3.19 (runtime)
- Expose port 8080
- Health check via GET /health

### 5.2 Configuration
- PORT environment variable (default: 8080)
- No external configuration files for MVP

### 5.3 Observability
- Standard output logging
- Panic recovery middleware
- Structured JSON responses

---

## 6. Security Requirements

### 6.1 Input Validation
- All user input validated before processing
- SQL injection prevention (not applicable - in-memory)
- XSS prevention (output encoding)

### 6.2 Headers
- Content-Type: application/json
- CORS headers (optional, configurable)
- No sensitive data in responses

---

## 7. Implementation Constraints

### 7.1 Must Implement (P1)
- All 6 API endpoints
- In-memory store with mutex
- Automated tests (80%+ coverage)
- Dockerfile
- README.md

### 7.2 Out of Scope
- Database persistence
- Authentication/Authorization
- Rate limiting
- Pagination
- Webhooks

### 7.3 Technical Decisions
- Use UUID v4 for task IDs
- ISO 8601 timestamps with UTC timezone
- Standard library only (no third-party dependencies)
- Table-driven tests

---

## 8. Test Requirements

### 8.1 Unit Tests
- Store operations (Create, List, Get, Update, Delete)
- Validation logic
- Error handling

### 8.2 Integration Tests
- All API endpoints
- Status code verification
- Response body validation
- Error case testing

### 8.3 Coverage Target
- Minimum 80% code coverage
- All public functions tested

---

## 9. Acceptance Criteria

### 9.1 Functional
- [ ] POST /tasks creates task and returns 201
- [ ] GET /tasks returns all tasks
- [ ] GET /tasks/{id} returns specific task or 404
- [ ] PATCH /tasks/{id} updates task status
- [ ] DELETE /tasks/{id} removes task or returns 404
- [ ] GET /health returns healthy status
- [ ] All validation errors return 400

### 9.2 Non-Functional
- [ ] Docker image builds successfully
- [ ] All tests pass with go test ./...
- [ ] README includes setup instructions
- [ ] Response time < 100ms for all operations

---

## 10. File Structure
```
├── main.go
├── go.mod
├── go.sum
├── handlers/
│   └── task_handler.go
├── models/
│   └── task.go
├── store/
│   └── task_store.go
├── middleware/
│   └── middleware.go
├── tests/
│   └── api_test.go
├── Dockerfile
├── README.md
└── plans/
    ├── BOOTSTRAP.md
    ├── ORIGINAL_PRD.md
    └── SRD.md
```
