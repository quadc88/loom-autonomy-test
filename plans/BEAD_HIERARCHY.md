# Bead Hierarchy for Task REST API

## Epic 1: Core API Implementation (P0)

### Story 1.1: In-Memory Task Store
**Type:** task | **Priority:** P0 | **Assignee:** backend-engineer

- [ ] Create `models/task.go` with Task struct and validation
- [ ] Create `store/task_store.go` with sync.RWMutex concurrency
- [ ] Generate UUID v4 for task IDs
- [ ] Implement Create, List, Get, Update, Delete, HealthCheck methods

**Acceptance Criteria:**
- Store is thread-safe with mutex protection
- All CRUD operations work correctly
- Validation enforces title required, max 255 chars, description max 10000 chars
- UUIDs are properly formatted v4

---

### Story 1.2: Task Creation Endpoint
**Type:** task | **Priority:** P0 | **Assignee:** backend-engineer

- [ ] Create `handlers/task_handler.go` with POST /tasks handler
- [ ] Parse JSON request body
- [ ] Validate input and return 400 on error
- [ ] Return 201 Created with task JSON on success

**Acceptance Criteria:**
- Returns 201 with created task
- Returns 400 with validation error message
- Content-Type is application/json

---

### Story 1.3: Task Listing Endpoint
**Type:** task | **Priority:** P0 | **Assignee:** backend-engineer

- [ ] Implement GET /tasks handler
- [ ] Support optional ?status=active|completed query parameter
- [ ] Return array of tasks

**Acceptance Criteria:**
- Returns 200 with task array
- Filtering by status works correctly
- Empty list returns empty JSON array []

---

### Story 1.4: Task Retrieval by ID
**Type:** task | **Priority:** P0 | **Assignee:** backend-engineer

- [ ] Implement GET /tasks/{id} handler
- [ ] Parse UUID from URL path
- [ ] Return 404 if task not found

**Acceptance Criteria:**
- Returns 200 with task JSON when found
- Returns 404 with error message when not found
- Handles invalid UUID format gracefully

---

### Story 1.5: Task Completion Update
**Type:** task | **Priority:** P0 | **Assignee:** backend-engineer

- [ ] Implement PATCH /tasks/{id} handler
- [ ] Parse JSON body with completed field
- [ ] Update task and return updated object

**Acceptance Criteria:**
- Returns 200 with updated task
- Returns 400 if request body is invalid
- Returns 404 if task not found
- Updates updated_at timestamp

---

### Story 1.6: Task Deletion
**Type:** task | **Priority:** P0 | **Assignee:** backend-engineer

- [ ] Implement DELETE /tasks/{id} handler
- [ ] Remove task from store
- [ ] Return 204 No Content on success

**Acceptance Criteria:**
- Returns 204 on successful deletion
- Returns 404 if task not found
- Task is actually removed from store

---

### Story 1.7: Health Check Endpoint
**Type:** task | **Priority:** P0 | **Assignee:** backend-engineer

- [ ] Implement GET /health handler
- [ ] Return healthy status

**Acceptance Criteria:**
- Returns 200 with {"status": "healthy"}
- Response time < 10ms

---

## Epic 2: Testing & Quality (P1)

### Story 2.1: Unit Tests for Store
**Type:** task | **Priority:** P1 | **Assignee:** qa-engineer

- [ ] Write tests for Create operation
- [ ] Write tests for List operation (with filters)
- [ ] Write tests for Get operation (found/not found)
- [ ] Write tests for Update operation
- [ ] Write tests for Delete operation
- [ ] Write tests for concurrent access

**Acceptance Criteria:**
- All store operations tested
- Concurrency tests verify thread safety
- 80%+ code coverage

---

### Story 2.2: Integration Tests for API Endpoints
**Type:** task | **Priority:** P1 | **Assignee:** qa-engineer

- [ ] Write integration tests for POST /tasks
- [ ] Write integration tests for GET /tasks
- [ ] Write integration tests for GET /tasks/{id}
- [ ] Write integration tests for PATCH /tasks/{id}
- [ ] Write integration tests for DELETE /tasks/{id}
- [ ] Write integration tests for GET /health

**Acceptance Criteria:**
- All endpoints tested with valid and invalid inputs
- Status codes verified
- Response bodies validated
- Error cases covered

---

### Story 2.3: Validation Error Handling Tests
**Type:** task | **Priority:** P1 | **Assignee:** qa-engineer

- [ ] Test missing title returns 400
- [ ] Test empty title returns 400
- [ ] Test title > 255 chars returns 400
- [ ] Test description > 10000 chars returns 400
- [ ] Test invalid JSON returns 400

**Acceptance Criteria:**
- All validation rules tested
- Error messages are clear and helpful
- Consistent error response format

---

## Epic 3: Deployment & Documentation (P2)

### Story 3.1: Dockerfile with Multi-Stage Build
**Type:** task | **Priority:** P2 | **Assignee:** devops-engineer

- [ ] Create multi-stage Dockerfile
- [ ] Build stage: golang:1.21-alpine
- [ ] Runtime stage: alpine:3.19
- [ ] Expose port 8080
- [ ] Add healthcheck instruction

**Acceptance Criteria:**
- Docker image builds successfully
- Image size is optimized
- Container runs and responds to /health

---

### Story 3.2: README Documentation
**Type:** task | **Priority:** P2 | **Assignee:** technical-writer

- [ ] Write project overview
- [ ] Document local setup and run instructions
- [ ] Document Docker build and run
- [ ] Document API endpoints with examples
- [ ] Document test running

**Acceptance Criteria:**
- README is comprehensive and clear
- Includes all required sections
- API usage examples work

---

### Story 3.3: CORS and Middleware
**Type:** task | **Priority:** P2 | **Assignee:** backend-engineer

- [ ] Create middleware package
- [ ] Implement CORS middleware
- [ ] Implement panic recovery middleware
- [ ] Implement logging middleware
- [ ] Wire middleware into main.go

**Acceptance Criteria:**
- CORS headers present in responses
- Panics are recovered gracefully
- Requests are logged

---

## File Structure

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
    ├── BEAD_HIERARCHY.md
    ├── BOOTSTRAP.md
    ├── ORIGINAL_PRD.md
    └── SRD.md
```
