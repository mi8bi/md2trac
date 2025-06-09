# md2trac Dockerfile (Production)
# Multi-stage build for minimal final image

# Build stage
FROM golang:1.23-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY go.mod ./
COPY go.sum ./

# Build the application
WORKDIR /build/cmd/md2trac
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /build/md2trac

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests (if needed)
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 md2trac && \
    adduser -u 1000 -G md2trac -s /bin/sh -D md2trac

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/md2trac /usr/local/bin/md2trac

# Create workspace directory for file processing
RUN mkdir -p /workspace && \
    chown -R md2trac:md2trac /workspace

# Switch to non-root user
USER md2trac

# Set default workspace
WORKDIR /workspace

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD md2trac --help > /dev/null || exit 1

# Default command
ENTRYPOINT ["md2trac"]
CMD ["--help"]