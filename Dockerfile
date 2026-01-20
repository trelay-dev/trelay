# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build server
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/bin/trelay-server ./cmd/server

# Build CLI
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/bin/trelay ./cmd/trelay

# Runtime stage
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs tzdata

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/bin/trelay-server /app/trelay-server
COPY --from=builder /app/bin/trelay /usr/local/bin/trelay

# Create data directory
RUN mkdir -p /app/data

# Environment defaults
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV DB_DRIVER=sqlite3
ENV DB_PATH=/app/data/trelay.db

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Run server
CMD ["/app/trelay-server"]
