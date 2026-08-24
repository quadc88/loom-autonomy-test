# Task REST API

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

## Local Development

### Prerequisites

- Python 3.11+

### Run Locally

```bash
# Start the server
python3 app.py

# In another terminal, test the API
curl http://localhost:8080/health
```

### Run Tests

```bash
python3 tests.py
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
```

## Project Structure

```
.
├── app.py          # Main application
├── tests.py        # Test suite
├── Dockerfile      # Container configuration
└── README.md       # This file
```
