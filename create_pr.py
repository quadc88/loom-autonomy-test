#!/usr/bin/env python3
import os, json, urllib.request, urllib.error

token = os.environ.get('GITHUB_TOKEN', '')
if not token:
    print('GITHUB_TOKEN not set')
    exit(1)

pr_data = {
    'title': 'feat: add GET /tasks/{id} endpoint with test',
    'body': 'Implements the GetTask handler with proper UUID validation and 404 handling. Includes TestGetTask_Success unit test.',
    'head': 'feature/get-task-handler-test',
    'base': 'main'
}

url = 'https://api.github.com/repos/quadc88/loom-autonomy-test/pulls'
req = urllib.request.Request(url,
    data=json.dumps(pr_data).encode(),
    headers={'Authorization': f'token {token}', 'Accept': 'application/vnd.github.v3+json', 'Content-Type': 'application/json'})
try:
    resp = urllib.request.urlopen(req)
    print(f'Status: {resp.status}')
    print(resp.read().decode())
except urllib.error.HTTPError as e:
    print(f'HTTP Error: {e.code}')
    print(e.read().decode())
