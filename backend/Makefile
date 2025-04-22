# Load environment variables from app.env
include app.env
export $(shell sed 's/=.*//' app.env)

# Default target
.PHONY: all
all: generate migrateup run

# Run SQLC code generation
.PHONY: generate
generate:
	sqlc generate

# Run all migrations (up)
.PHONY: migrateup
migrateup:
	migrate -path ./db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=require" -verbose up

# Rollback all migrations (down)
.PHONY: migratedown
migratedown:
	migrate -path db/migration -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose down

# Run the Go server with Gin
.PHONY: run
run:
	go run main.go

# Install dependencies
.PHONY: install
install:
	go mod tidy
	go get github.com/gin-gonic/gin
	go get github.com/jackc/pgx/v4
	go get github.com/joho/godotenv
	go get github.com/sqlc-dev/sqlc
