# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Download migrate CLI tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
  | tar -xz -C /usr/local/bin

# Copy application source code
COPY . .

# Debug: List files to check structure
RUN echo "==== FILE STRUCTURE IN /app (BUILDER) ====" && ls -R /app

# Build the Go server
RUN go build -o server main.go

# Run stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies and dos2unix for CRLF fix
RUN apt-get update && apt-get install -y \
    netcat-openbsd curl make bash dos2unix && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Set GIN to release mode
ENV GIN_MODE=release

# Copy built Go server
COPY --from=builder /app/server .

# Copy migrate binary
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Copy environment and support files
COPY ./app.env .
COPY ./start.sh .
COPY ./Makefile .

# Copy migration files (assumes you have them in ./backend/db/migration)
COPY ./db/migration ./db/migration

# Fix line endings (if you're on Windows host)
RUN dos2unix start.sh app.env Makefile

# Debug: Check final file structure
RUN echo "==== FILE STRUCTURE IN /app (FINAL IMAGE) ====" && ls -R /app

# Make start.sh executable
RUN chmod +x start.sh

EXPOSE 8080

CMD ["./start.sh"]
