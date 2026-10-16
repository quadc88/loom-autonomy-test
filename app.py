#!/usr/bin/env python3
"""Minimal production-quality Task REST API using Python stdlib."""

import json
import threading
from datetime import datetime, timezone
from http.server import HTTPServer, BaseHTTPRequestHandler
from typing import Any


class TaskStore:
    """In-memory task store with thread-safe operations."""

    def __init__(self) -> None:
        self._tasks: dict[str, dict[str, Any]] = {}
        self._counter: int = 0
        self._lock = threading.Lock()

    def create(self, title: str) -> dict[str, Any]:
        if not title.strip():
            raise ValueError("title is required")
        with self._lock:
            self._counter += 1
            now = datetime.now(timezone.utc).isoformat()
            task: dict[str, Any] = {
                "id": str(self._counter),
                "title": title,
                "done": False,
                "created_at": now,
                "updated_at": now,
            }
            self._tasks[task["id"]] = task
            return task

    def list_all(self) -> list[dict[str, Any]]:
        with self._lock:
            return list(self._tasks.values())

    def get(self, task_id: str) -> dict[str, Any]:
        with self._lock:
            task = self._tasks.get(task_id)
            if task is None:
                raise KeyError(task_id)
            return task

    def mark_done(self, task_id: str) -> dict[str, Any]:
        with self._lock:
            task = self._tasks.get(task_id)
            if task is None:
                raise KeyError(task_id)
            task["done"] = True
            task["updated_at"] = datetime.now(timezone.utc).isoformat()
            return task

    def patch(self, task_id: str, data: dict[str, Any]) -> dict[str, Any]:
        with self._lock:
            task = self._tasks.get(task_id)
            if task is None:
                raise KeyError(task_id)
            if "completed" in data:
                if not isinstance(data["completed"], bool):
                    raise ValueError("completed must be a boolean")
                task["done"] = data["completed"]
            task["updated_at"] = datetime.now(timezone.utc).isoformat()
            return task

    def delete(self, task_id: str) -> None:
        with self._lock:
            if task_id not in self._tasks:
                raise KeyError(task_id)
            del self._tasks[task_id]


store = TaskStore()


def _json_response(handler: BaseHTTPRequestHandler, status: int, body: Any) -> None:
    data = json.dumps(body).encode("utf-8") if body is not None else b""
    handler.send_response(status)
    handler.send_header("Content-Type", "application/json")
    handler.send_header("Content-Length", str(len(data)))
    handler.end_headers()
    handler.wfile.write(data)


class TaskHandler(BaseHTTPRequestHandler):
    """HTTP request handler for the Task API."""

    def _parse_path(self) -> tuple[str, str | None]:
        path = self.path.rstrip("/")
        parts = path.split("/")
        # /health
        if len(parts) == 2 and parts[1] == "health":
            return "health", None
        # /tasks or /tasks/<id>
        if len(parts) == 2 and parts[1] == "tasks":
            return "tasks", None
        if len(parts) == 3 and parts[1] == "tasks":
            return "task", parts[2]
        return "unknown", None

    def do_GET(self) -> None:
        resource, task_id = self._parse_path()
        try:
            if resource == "tasks":
                _json_response(self, 200, store.list_all())
            elif resource == "task" and task_id:
                _json_response(self, 200, store.get(task_id))
            elif resource == "health":
                _json_response(self, 200, {"status": "ok"})
            else:
                _json_response(self, 404, {"error": "not found"})
        except KeyError:
            _json_response(self, 404, {"error": "task not found"})

    def do_POST(self) -> None:
        resource, _ = self._parse_path()
        if resource != "tasks":
            _json_response(self, 404, {"error": "not found"})
            return
        try:
            length = int(self.headers.get("Content-Length", 0))
            body = json.loads(self.rfile.read(length)) if length > 0 else {}
            title = body.get("title", "")
            task = store.create(title)
            _json_response(self, 201, task)
        except (ValueError, KeyError) as e:
            _json_response(self, 400, {"error": str(e)})
        except json.JSONDecodeError:
            _json_response(self, 400, {"error": "invalid JSON"})

    def do_PUT(self) -> None:
        resource, task_id = self._parse_path()
        if resource != "task" or not task_id:
            _json_response(self, 404, {"error": "not found"})
            return
        try:
            task = store.mark_done(task_id)
            _json_response(self, 200, task)
        except KeyError:
            _json_response(self, 404, {"error": "task not found"})

    def do_DELETE(self) -> None:
        resource, task_id = self._parse_path()
        if resource != "task" or not task_id:
            _json_response(self, 404, {"error": "not found"})
            return
        try:
            store.delete(task_id)
            _json_response(self, 204, None)
        except KeyError:
            _json_response(self, 404, {"error": "task not found"})

    def do_PATCH(self) -> None:
        resource, task_id = self._parse_path()
        if resource != "task" or not task_id:
            _json_response(self, 404, {"error": "not found"})
            return
        try:
            length = int(self.headers.get("Content-Length", 0))
            if length == 0:
                _json_response(self, 400, {"error": "request body is required"})
                return
            body = json.loads(self.rfile.read(length))
            if not isinstance(body, dict):
                _json_response(self, 400, {"error": "request body must be JSON object"})
                return
            if "completed" not in body:
                _json_response(self, 400, {"error": "completed field is required"})
                return
            task = store.patch(task_id, body)
            _json_response(self, 200, task)
        except json.JSONDecodeError:
            _json_response(self, 400, {"error": "invalid JSON"})
        except ValueError as e:
            _json_response(self, 400, {"error": str(e)})
        except KeyError:
            _json_response(self, 404, {"error": "task not found"})

    def log_message(self, format: str, *args: Any) -> None:  # type: ignore[override]
        print(f"{self.address_string()} - {format % args}")


def main() -> None:
    host = "0.0.0.0"
    port = 8080
    server = HTTPServer((host, port), TaskHandler)
    print(f"Task API running on http://{host}:{port}")
    server.serve_forever()


if __name__ == "__main__":
    main()
