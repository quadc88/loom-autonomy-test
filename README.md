# Task REST API

<<<<<<< HEAD
A minimal production-quality Task REST API built with Go. Supports creating, listing, retrieving, updating, and deleting tasks with an in-memory store.

## Features

- ✅ Create tasks with title and optional description
- ✅ List tasks with optional status filter (active/completed)
- ✅ Get task by ID
- ✅ Update task (mark as completed)
- ✅ Delete task
- ✅ Health check endpoint
- ✅ Thread-safe in-memory store with UUID v4 IDs
- ✅ JSON responses with proper HTTP status codes
- ✅ CORS support
- ✅ Recovery and logging middleware
- ✅ Docker support

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |
| GET | `/tasks` | List all tasks (optional: `?status=active|completed&limit=N`) |
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
GET /tasks?status=active&limit=10
```

**Update a task:**
```bash
PATCH /tasks/{id}
Content-Type: application/json

{
  "completed": true
}
```
=======
A minimal production-quality Task REST API built with Python stdlib.

## Features

- Create tasks
- List all tasks
- Get task by ID
- Mark task as completed
- Delete task
- Health check endpoint
- JSON responses with proper HTTP status codes
- In-memory store (restart clears data)

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /tasks | List all tasks |
| POST | /tasks | Create a task |
| GET | /tasks/:id | Get task by ID |
| PUT | /tasks/:id | Mark task as done |
| DELETE | /tasks/:id | Delete a task |
>>>>>>> origin/main

## Local Development

### Prerequisites

<<<<<<< HEAD
- Go 1.21 or later
=======
- Python 3.11+
>>>>>>> origin/main

### Run Locally

```bash
<<<<<<< HEAD
go run .
```

The server starts on port `8080` by default. Use the `PORT` environment variable to change it:

```bash
PORT=3000 go run .
=======
# Start the server
python3 app.py

# In another terminal, test the API
curl http://localhost:8080/health
>>>>>>> origin/main
```

### Run Tests

```bash
<<<<<<< HEAD
go test ./...
```

### Run Tests with Verbose Output

```bash
go test ./... -v
```

### Build Binary

```bash
go build -o task-api .
=======
python3 tests.py
>>>>>>> origin/main
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

<<<<<<< HEAD
With custom port:

```bash
docker run -p 8080:8080 -e PORT=3000 task-api
```

### Docker Health Check

```bash
docker inspect --format='{{.State.Health.Status}}' <container_id>
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
  -d '{"title": "Task title", "description": "Optional description"}'
```

### List Tasks
```bash
curl http://localhost:8080/tasks
curl http://localhost:8080/tasks?status=active
curl http://localhost:8080/tasks?status=completed&limit=10
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
=======
### Test with Docker

```bash
curl http://localhost:8080/health
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"title":"My Task"}'
curl http://localhost:8080/tasks
```

## Example Requests

```bash
# Create a task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Python"}'

# List all tasks
curl http://localhost:8080/tasks

# Get task by ID
curl http://localhost:8080/tasks/1

# Mark task as done
curl -X PUT http://localhost:8080/tasks/1

# Delete task
curl -X DELETE http://localhost:8080/tasks/1

# Health check
curl http://localhost:8080/health
```

## Response Format

### Success Response
```json
{
  "id": "1",
  "title": "Learn Python",
  "done": false
}
```

### Error Response
```json
{
  "error": "task not found"
}
>>>>>>> origin/main
```

## Project Structure

```
.
<<<<<<< HEAD
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksum
├── Dockerfile           # Multi-stage Docker build
├── README.md            # This file
├── handlers/
│   └── task_handler.go  # HTTP request handlers
├── models/
│   └── task.go          # Task data model and validation
├── store/
│   ├── task_store.go    # In-memory task store
│   └── task_store_test.go # Store unit tests
├── middleware/
│   └── middleware.go    # CORS, logging, recovery middleware
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

## License

MIT
=======
├── app.py          # Main application
├── tests.py        # Test suite
├── Dockerfile      # Container configuration
└── README.md       # This file
```
>>>>>>> origin/main
