# Frontend build stage
FROM oven/bun:1 AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/bun.lock* ./
RUN bun install --frozen-lockfile

COPY frontend/ ./
RUN bun run build

# Backend build stage
FROM golang:1.22-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/bin/trelay-server ./cmd/server
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/bin/trelay ./cmd/trelay

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates sqlite-libs tzdata

WORKDIR /app

COPY --from=backend-builder /app/bin/trelay-server /app/trelay-server
COPY --from=backend-builder /app/bin/trelay /usr/local/bin/trelay
COPY --from=frontend-builder /app/frontend/build /app/static

RUN mkdir -p /app/data

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV DB_DRIVER=sqlite3
ENV DB_PATH=/app/data/trelay.db
ENV STATIC_DIR=/app/static

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["/app/trelay-server"]
