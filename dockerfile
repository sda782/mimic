# ---------- Stage 1: Build ----------
FROM golang:1.25.1 AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download

# Copy backend source code
COPY backend/ ./backend/

# Copy webapp
COPY webapp/ ./webapp/

# Build the Go binary for Linux
WORKDIR /app/backend
RUN GOOS=linux GOARCH=amd64 go build -o ../server main.go

# ---------- Stage 2: Runtime ----------
FROM debian:bookworm-slim

WORKDIR /app

# Install ca-certificates (needed if Go server does HTTPS requests)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binary and webapp from builder
COPY --from=builder /app/server ./
COPY --from=builder /app/webapp ./webapp/

# Create runtime directories/files
RUN mkdir -p ./uploads && touch ./mimic.db


# Expose backend port
EXPOSE 6070

# Run the Go server
CMD ["./server"]