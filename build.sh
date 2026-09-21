#!/bin/bash
# Build and test script using Docker
cd /app/data/projects/loom-bootstrap-test-5/main
docker run --rm \
  -v "$(pwd):/go/src/loom-bootstrap-test-5" \
  -w /go/src/loom-bootstrap-test-5 \
  golang:1.21-alpine \
  sh -c 'go test ./...' 2>&1
