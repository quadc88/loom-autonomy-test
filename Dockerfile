# Multi-stage Dockerfile for Task REST API

# Build stage
FROM golang:1.21-alpine AS builder

# Install ca-certificates for HTTPS access to module proxy
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Copy dependency files first for better layer caching
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o task-api .

# Run tests during build (optional, uncomment if needed)
# RUN go test ./...

# Runtime stage
FROM alpine:3.19

# Install ca-certificates and required runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set timezone
ENV TZ=UTC

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/task-api .

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose the default port
EXPOSE 8080

# Environment variables
ENV PORT=8080
ENV GIN_MODE=release

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./task-api"]
