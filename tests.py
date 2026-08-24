#!/usr/bin/env python3
"""Tests for the Task REST API."""

import json
import threading
import unittest
import urllib.request
import urllib.error
from http.server import HTTPServer

from app import TaskStore, TaskHandler


BASE_URL = "http://localhost:9002"


class TestTaskAPI(unittest.TestCase):
    """Integration tests against a running server."""

    @classmethod
    def setUpClass(cls):
        cls.server = HTTPServer(("127.0.0.1", 9002), TaskHandler)
        cls.thread = threading.Thread(target=cls.server.serve_forever)
        cls.thread.daemon = True
        cls.thread.start()

    @classmethod
    def tearDownClass(cls):
        cls.server.shutdown()
        cls.server.server_close()

    def _request(self, method, path, data=None):
        url = BASE_URL + path
        body = json.dumps(data).encode("utf-8") if data else None
        req = urllib.request.Request(url, data=body, method=method)
        if body:
            req.add_header("Content-Type", "application/json")
        try:
            with urllib.request.urlopen(req) as resp:
                resp_data = resp.read().decode("utf-8")
                return resp.status, json.loads(resp_data) if resp_data else None
        except urllib.error.HTTPError as e:
            resp_data = e.read().decode("utf-8")
            if resp_data:
                return e.code, json.loads(resp_data)
            return e.code, None

    def test_health(self):
        status, body = self._request("GET", "/health")
        self.assertEqual(status, 200)
        self.assertEqual(body["status"], "ok")

    def test_create_task(self):
        status, body = self._request("POST", "/tasks", {"title": "Test task"})
        self.assertEqual(status, 201)
        self.assertEqual(body["title"], "Test task")
        self.assertFalse(body["done"])
        self.assertIn("id", body)

    def test_create_task_empty_title(self):
        status, body = self._request("POST", "/tasks", {"title": ""})
        self.assertEqual(status, 400)
        self.assertIn("error", body)

    def test_list_tasks(self):
        self._request("POST", "/tasks", {"title": "Task 1"})
        status, body = self._request("GET", "/tasks")
        self.assertEqual(status, 200)
        self.assertIsInstance(body, list)
        self.assertGreaterEqual(len(body), 1)

    def test_get_task(self):
        status, created = self._request("POST", "/tasks", {"title": "Get me"})
        task_id = created["id"]
        status, body = self._request("GET", f"/tasks/{task_id}")
        self.assertEqual(status, 200)
        self.assertEqual(body["title"], "Get me")

    def test_get_task_not_found(self):
        status, body = self._request("GET", "/tasks/999")
        self.assertEqual(status, 404)
        self.assertIn("error", body)

    def test_mark_done(self):
        status, created = self._request("POST", "/tasks", {"title": "Do me"})
        task_id = created["id"]
        status, body = self._request("PUT", f"/tasks/{task_id}")
        self.assertEqual(status, 200)
        self.assertTrue(body["done"])

    def test_mark_done_not_found(self):
        status, body = self._request("PUT", "/tasks/999")
        self.assertEqual(status, 404)
        self.assertIn("error", body)

    def test_delete_task(self):
        status, created = self._request("POST", "/tasks", {"title": "Delete me"})
        task_id = created["id"]
        status, _ = self._request("DELETE", f"/tasks/{task_id}")
        self.assertEqual(status, 204)
        status, body = self._request("GET", f"/tasks/{task_id}")
        self.assertEqual(status, 404)
        self.assertIn("error", body)

    def test_delete_task_not_found(self):
        status, body = self._request("DELETE", "/tasks/999")
        self.assertEqual(status, 404)
        self.assertIn("error", body)

    def test_method_not_allowed(self):
        status, body = self._request("PATCH", "/tasks")
        self.assertEqual(status, 405)
        self.assertIn("error", body)


if __name__ == "__main__":
    unittest.main()
