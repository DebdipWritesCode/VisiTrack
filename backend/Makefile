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

# Generate Protobuf files
.PHONY: proto
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=visitrack \
	--experimental_allow_proto3_optional \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc -f

# Run all migrations (up)
.PHONY: migrateup
migrateup:
	migrate -path ./db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose up

# Rollback all migrations (down)
.PHONY: migratedown
migratedown:
	migrate -path db/migration -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose down

# Run the Go server
.PHONY: run
run:
	go run main.go

# Run Evans gRPC client
.PHONY: evans
evans:
	evans --proto service_visitrack.proto --path proto --host localhost --port 9000

# Install dependencies
.PHONY: install
install:
	go mod tidy
	go get github.com/gin-gonic/gin
	go get github.com/jackc/pgx/v4
	go get github.com/joho/godotenv
	go get github.com/sqlc-dev/sqlc
