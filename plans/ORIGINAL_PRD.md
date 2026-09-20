# Product Requirements Document (PRD)
## Task REST API

**Version:** 1.0
**Status:** Draft
**Last Updated:** 2026-09-20

---

## 1. Executive Summary

### Project Vision
Build a minimal, production-quality REST API for managing tasks. This service provides a simple, reliable backend for task management operations with an in-memory data store, containerized deployment, and comprehensive documentation.

### Problem Statement
Teams need a lightweight, self-contained task management API that can be deployed quickly without external dependencies. This MVP establishes the foundation for future task management features.

### Success Metrics
- All CRUD operations return correct HTTP status codes
- 95%+ test coverage on core endpoints
- Container builds successfully and runs without errors
- API documentation is complete and accurate

---

## 2. User Personas

### Primary Persona: Developer
- **Name:** Alex Chen
- **Role:** Full-stack developer
- **Needs:** Quick integration, clear API contracts, Docker support
- **Goals:** Add task management to existing projects

### Secondary Persona: DevOps Engineer
- **Name:** Jordan Rivera
- **Role:** Infrastructure engineer
- **Needs:** Easy deployment, health checks, logging
- **Goals:** Deploy and monitor in production environments

---

## 3. Core Features & User Stories

### Feature 1: Task Creation
**Story:** As a user, I want to create a new task so I can track work items.

**Acceptance Criteria:**
- POST /tasks accepts JSON body with title (required) and optional description
- Returns 201 Created with the created task
- Returns 400 Bad Request for invalid input
- Task has auto-generated unique ID

### Feature 2: List Tasks
**Story:** As a user, I want to list all tasks so I can see my work queue.

**Acceptance Criteria:**
- GET /tasks returns array of all tasks
- Supports optional pagination (not required for MVP)
- Returns 200 OK with JSON array

### Feature 3: Get Task by ID
**Story:** As a user, I want to retrieve a specific task so I can view its details.

**Acceptance Criteria:**
- GET /tasks/{id} returns the task with matching ID
- Returns 404 Not Found if task doesn't exist
- Returns 400 Bad Request for invalid ID format

### Feature 4: Mark Task Completed
**Story:** As a user, I want to mark a task as completed so I can track progress.

**Acceptance Criteria:**
- PATCH /tasks/{id}/complete marks task as completed
- Returns updated task with completed=true
- Returns 404 if task not found
- Returns 400 if task is already completed

### Feature 5: Delete Task
**Story:** As a user, I want to delete a task so I can remove completed or irrelevant items.

**Acceptance Criteria:**
- DELETE /tasks/{id} removes the task
- Returns 204 No Content on success
- Returns 404 if task not found
- Returns 400 for invalid ID format

### Feature 6: Health Check
**Story:** As an operator, I want to check if the service is healthy so I can monitor availability.

**Acceptance Criteria:**
- GET /health returns 200 OK
- Response includes service status and timestamp
- No authentication required

---

## 4. Technical Requirements

### Stack Decisions
- **Language:** Go (Golang)
- **Web Framework:** Standard library net/http (minimal, no external dependencies)
- **Data Store:** In-memory map (sync.RWMutex for concurrency)
- **Containerization:** Docker with multi-stage build
- **Testing:** Go testing package

### API Specifications
- Base URL: `http://localhost:8080`
- Content-Type: application/json
- Request/Response format: JSON
- Error responses: `{"error": "message"}`

### Data Model
```go
type Task struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    Completed   bool   `json:"completed"`
    CreatedAt   string `json:"created_at"`
}
```

---

## 5. Architecture Considerations

### In-Memory Store
- Single process with concurrent-safe access
- No persistence across restarts (MVP limitation)
- Thread-safe using sync.RWMutex

### API Design
- RESTful resource naming
- Proper HTTP status codes
- JSON validation with clear error messages
- Content-Type headers on all responses

### Deployment
- Dockerfile with multi-stage build
- Slim base image (alpine)
- Expose port 8080
- Health check endpoint for orchestration

---

## 6. MVP Scope Definition

### P1 - Must Have (MVP)
- [ ] Create task endpoint (POST /tasks)
- [ ] List tasks endpoint (GET /tasks)
- [ ] Get task by ID endpoint (GET /tasks/{id})
- [ ] Mark task completed endpoint (PATCH /tasks/{id}/complete)
- [ ] Delete task endpoint (DELETE /tasks/{id})
- [ ] Health check endpoint (GET /health)
- [ ] In-memory data store
- [ ] JSON request/response validation
- [ ] Automated tests for all endpoints
- [ ] Dockerfile
- [ ] README with setup and usage

### P2 - Nice to Have
- [ ] Pagination for list endpoint
- [ ] Task filtering by completed status
- [ ] Input sanitization
- [ ] Request ID headers for tracing
- [ ] Structured logging

### P3 - Future
- [ ] Persistent storage (SQLite/PostgreSQL)
- [ ] Authentication/authorization
- [ ] Task priorities and due dates
- [ ] REST API versioning
- [ ] OpenAPI/Swagger documentation
- [ ] Rate limiting

---

## 7. Non-Functional Requirements

### Performance
- Response time < 100ms for single task operations
- Support concurrent requests (tested with 10+ simultaneous clients)
- No memory leaks under normal operation

### Security
- Input validation to prevent injection attacks
- No sensitive data in responses
- Proper HTTP status codes (no information leakage)

### Reliability
- Graceful handling of malformed requests
- Consistent error response format
- Service restarts without data corruption concerns (in-memory MVP)

### Maintainability
- Clean separation of concerns
- Commented code for key logic
- Test coverage > 80%

---

## 8. Acceptance Criteria Summary

| Feature | Status Code | Response Format | Test Coverage |
|---------|-------------|-----------------|---------------|
| POST /tasks | 201, 400 | JSON task object | Required |
| GET /tasks | 200 | JSON array | Required |
| GET /tasks/{id} | 200, 404, 400 | JSON task object | Required |
| PATCH /tasks/{id}/complete | 200, 404, 400 | JSON task object | Required |
| DELETE /tasks/{id} | 204, 404, 400 | Empty/JSON error | Required |
| GET /health | 200 | JSON status | Required |

---

## 9. Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| In-memory data loss on restart | Low (MVP) | Document limitation; plan for persistence in P3 |
| Concurrency bugs | Medium | Use sync.RWMutex; test with concurrent requests |
| Resource exhaustion | Low | In-memory limits scope; add monitoring if needed |

---

**Approved by:** Engineering Team
**Next Steps:** Create System Requirements Document (SRD)
