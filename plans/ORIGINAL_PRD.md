# Product Requirements Document (PRD) - Task REST API

## 1. Executive Summary
Build a minimal, production-quality Task REST API with in-memory storage. Provide CRUD operations plus health check with automated tests, Dockerfile, and README.

## 2. User Personas
- Backend Developer: needs simple, predictable API
- QA Engineer: needs reliable API with test coverage
- DevOps Engineer: needs Docker support and health checks

## 3. Core Features
1. **Create Task**: POST /tasks with title (required) and description (optional)
   - 201 Created with task JSON
   - 400 Bad Request for invalid input
2. **List Tasks**: GET /tasks with optional ?status=filter
   - 200 OK with array of tasks
3. **Get Task**: GET /tasks/{id}
   - 200 OK with task
   - 404 Not Found
4. **Mark Completed**: PATCH /tasks/{id} with {"completed": true}
   - 200 OK with updated task
   - 404 Not Found
5. **Delete Task**: DELETE /tasks/{id}
   - 204 No Content
   - 404 Not Found
6. **Health Check**: GET /health
   - 200 OK with {"status": "healthy"}

## 4. Technical Requirements
- Language: Go (standard library net/http)
- Storage: In-memory map with sync.RWMutex
- Response format: JSON
- Status codes: RESTful conventions

## 5. Data Model
```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "completed": false,
  "created_at": "ISO8601",
  "updated_at": "ISO8601"
}
```

## 6. MVP Scope (P1)
- [x] All CRUD endpoints + health check
- [x] In-memory store with concurrency safety
- [x] Automated tests
- [x] Dockerfile
- [x] README

## 7. Nice to Have (P2)
- Filtering, pagination, structured logging

## 8. Future (P3)
- Database persistence, auth, webhooks, rate limiting

## 9. Non-Functional Requirements
- Performance: <100ms response time
- Security: input validation, CORS
- Scalability: stateless design
- Reliability: graceful shutdown, panic recovery
