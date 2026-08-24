# Task REST API

A minimal production-quality Task REST API built with Go. Supports creating, listing, retrieving, updating, and deleting tasks with an in-memory store.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |
| GET | `/tasks` | List all tasks (optional: `?status=completed|pending`) |
| POST | `/tasks` | Create a new task |
| GET | `/tasks/{id}` | Get a task by ID |
| PATCH | `/tasks/{id}` | Update a task (mark as completed) |
| DELETE | `/tasks/{id}` | Delete a task |

### Request/Response Examples

**Create a task:**
```bash
POST /tasks
Content-Type: application/json

{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**List tasks:**
```bash
GET /tasks
GET /tasks?status=completed
```

**Update a task:**
```bash
PATCH /tasks/{id}
Content-Type: application/json

{
  "completed": true
}
```

## Local Development

### Prerequisites

- Go 1.21 or later

### Run Locally

```bash
go run .
```

The server starts on port `8080` by default. Use the `PORT` environment variable to change it:

```bash
PORT=3000 go run .
```

### Run Tests

```bash
go test ./tests/...
```

### Build Binary

```bash
go build -o task-api .
```

## Docker

### Build Image

```bash
docker build -t task-api .
```

### Run Container

```bash
docker run -p 8080:8080 task-api
```

With custom port:

```bash
docker run -p 8080:8080 -e PORT=3000 task-api
```

## API Usage

### Health Check
```bash
curl http://localhost:8080/health
```

### Create Task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Task title"}'
```

### List Tasks
```bash
curl http://localhost:8080/tasks
```

### Get Task by ID
```bash
curl http://localhost:8080/tasks/{id}
```

### Mark Task Completed
```bash
curl -X PATCH http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

### Delete Task
```bash
curl -X DELETE http://localhost:8080/tasks/{id}
```

## Project Structure

```
.
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── Dockerfile           # Multi-stage Docker build
├── README.md            # This file
├── handlers/
│   └── task_handler.go  # HTTP request handlers
├── models/
│   └── task.go          # Task data model and validation
├── store/
│   └── task_store.go    # In-memory task store
└── tests/
    └── api_test.go      # Integration tests
```

## Validation Rules

- `title` is required for task creation
- `title` must be 255 characters or less
- `description` is optional, must be 10000 characters or less
- `id` must be a valid UUID format

## Error Responses

All errors return JSON in the format:
```json
{
  "error": "error message"
}
```

Common status codes:
- `400 Bad Request` - Invalid input or missing required fields
- `404 Not Found` - Task not found
- `405 Method Not Allowed` - Invalid HTTP method for endpoint
- `500 Internal Server Error` - Server-side error
