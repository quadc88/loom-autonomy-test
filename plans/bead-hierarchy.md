# Bead Hierarchy - Task REST API

Generated from: plans/SRD.md
Created: 2026-09-21

---

## Epic 1: Core API Implementation (P0)

### Story 1.1: In-Memory Task Store
**Priority:** P0
**Type:** task
**Assignee:** backend-engineer

**Tasks:**
- [ ] Implement `store/task_store.go` with sync.RWMutex concurrency control
- [ ] Add UUID v4 generation for task IDs
- [ ] Implement Create, List, Get, Update, Delete methods per TaskStore interface
- [ ] Acceptance: All store operations pass unit tests with concurrent access

### Story 1.2: Task Creation Endpoint
**Priority:** P0
**Type:** task
**Assignee:** backend-engineer

**Tasks:**
- [ ] Implement POST /tasks handler in `handlers/task_handler.go`
- [ ] Add input validation (title required, max 255 chars)
- [ ] Return 201 Created with Task JSON on success
- [ ] Acceptance: Returns 201 with valid input, 400 with invalid input

### Story 1.3: Task Listing Endpoint
**Priority:** P0
**Type:** task
**Assignee:** backend-engineer

**Tasks:**
- [ ] Implement GET /tasks handler with optional status filter
- [ ] Support query params: ?status=active|completed
- [ ] Return 200 OK with Task array
- [ ] Acceptance: Returns all tasks or filtered results correctly

### Story 1.4: Task Retrieval by ID
**Priority:** P0
**Type:** task
**Assignee:** backend-engineer

**Tasks:**
- [ ] Implement GET /tasks/{id} route handler
- [ ] Parse URL parameter for task ID
- [ ] Return 200 OK with Task JSON or 404 if not found
- [ ] Acceptance: Returns task by ID or 404 error

### Story 1.5: Task Completion Update
**Priority:** P0
**Type:** task
**Assignee:** backend-engineer

**Tasks:**
- [ ] Implement PATCH /tasks/{id} handler
- [ ] Accept request body with completed boolean
- [ ] Update timestamps on modification
- [ ] Return 200 OK with updated Task or appropriate errors
- [ ] Acceptance: Updates task status and returns updated entity

### Story 1.6: Task Deletion
**Priority:** P0
**Type:** task
**Assignee:** backend-engineer

**Tasks:**
- [ ] Implement DELETE /tasks/{id} handler
- [ ] Return 204 No Content on successful deletion
- [ ] Return 404 if task not found
- [ ] Acceptance: Deletes task and returns correct status code

### Story 1.7: Health Check Endpoint
**Priority:** P0
**Type:** task
**Assignee:** backend-engineer

**Tasks:**
- [ ] Implement GET /health handler
- [ ] Return 200 OK with {"status": "healthy"}
- [ ] Ensure response time < 10ms
- [ ] Acceptance: Health endpoint responds correctly

---

## Epic 2: Testing & Quality (P1)

### Story 2.1: Unit Tests for In-Memory Store
**Priority:** P1
**Type:** task
**Assignee:** qa-engineer

**Tasks:**
- [ ] Write tests for Create operation (success and validation errors)
- [ ] Write tests for List operation (with and without filters)
- [ ] Write tests for Get operation (found and not found)
- [ ] Write tests for Update operation (success and not found)
- [ ] Write tests for Delete operation (success and not found)
- [ ] Test concurrent access safety with goroutines
- [ ] Acceptance: 80%+ code coverage on store package

### Story 2.2: Integration Tests for API Endpoints
**Priority:** P1
**Type:** task
**Assignee:** qa-engineer

**Tasks:**
- [ ] Write integration test for POST /tasks (success and error cases)
- [ ] Write integration test for GET /tasks (with filters)
- [ ] Write integration test for GET /tasks/{id} (found and 404)
- [ ] Write integration test for PATCH /tasks/{id} (update and not found)
- [ ] Write integration test for DELETE /tasks/{id} (delete and 404)
- [ ] Write integration test for GET /health
- [ ] Acceptance: All API endpoints covered with table-driven tests

### Story 2.3: Validation Error Handling Tests
**Priority:** P1
**Type:** task
**Assignee:** qa-engineer

**Tasks:**
- [ ] Test empty title validation
- [ ] Test title exceeding 255 characters
- [ ] Test description exceeding 10000 characters
- [ ] Test invalid UUID format handling
- [ ] Test malformed JSON request bodies
- [ ] Acceptance: All validation paths return 400 with error messages

---

## Epic 3: Deployment & Documentation (P2)

### Story 3.1: Dockerfile with Multi-Stage Build
**Priority:** P2
**Type:** task
**Assignee:** devops-engineer

**Tasks:**
- [ ] Create Dockerfile with golang:1.21-alpine build stage
- [ ] Create alpine:3.19 runtime stage
- [ ] Copy binary and set up non-root user
- [ ] Expose port 8080
- [ ] Add healthcheck directive
- [ ] Acceptance: Docker image builds and runs successfully

### Story 3.2: README Documentation
**Priority:** P2
**Type:** task
**Assignee:** tech-writer

**Tasks:**
- [ ] Write setup instructions (local and Docker)
- [ ] Document all API endpoints with examples
- [ ] Include test commands and coverage info
- [ ] Add troubleshooting section
- [ ] Acceptance: README covers all required sections from SRD

### Story 3.3: CORS and Middleware Configuration
**Priority:** P2
**Type:** task
**Assignee:** backend-engineer

**Tasks:**
- [ ] Implement CORS middleware in `middleware/middleware.go`
- [ ] Add panic recovery middleware
- [ ] Add request logging middleware
- [ ] Wire middleware into main.go handler registration
- [ ] Acceptance: Middleware stack functions correctly

---

## Summary

| Epic | Priority | Stories | Total Tasks |
|------|----------|---------|-------------|
| Core API Implementation | P0 | 7 | 21 |
| Testing & Quality | P1 | 3 | 14 |
| Deployment & Documentation | P2 | 3 | 9 |
| **Total** | | **13** | **44** |
