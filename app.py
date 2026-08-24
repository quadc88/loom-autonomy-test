#!/usr/bin/env python3
"""Minimal production-quality Task REST API using Python stdlib."""

import json
import uuid
from http.server import HTTPServer, BaseHTTPRequestHandler
from typing import Any


class TaskStore:
    """In-memory task store with thread-safe operations."""

    def __init__(self) -> None:
        self._tasks: dict[str, dict[str, Any]] = {}
        self._counter: int = 0

    def create(self, title: str) -> dict[str, Any]:
        if not title.strip():
            raise ValueError("title is required")
        self._counter += 1
        task: dict[str, Any] = {
            "id": str(self._counter),
            "title": title,
            "done": False,
        }
        self._tasks[task["id"]] = task
        return task

    def list_all(self) -> list[dict[str, Any]]:
        return list(self._tasks.values())

    def get(self, task_id: str) -> dict[str, Any]:
        task = self._tasks.get(task_id)
        if task is None:
            raise KeyError(task_id)
        return task

    def mark_done(self, task_id: str) -> dict[str, Any]:
        task = self._tasks.get(task_id)
        if task is None:
            raise KeyError(task_id)
        task["done"] = True
        return task

    def delete(self, task_id: str) -> None:
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
        _json_response(self, 405, {"error": "method not allowed"})

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
