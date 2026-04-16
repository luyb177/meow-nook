.PHONY: proto build test lint tidy help

# ──────────────────────────────────────────────
# Proto generation
# ──────────────────────────────────────────────
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=module=github.com/luyb177/meow-nook \
	       --go-grpc_out=. --go-grpc_opt=module=github.com/luyb177/meow-nook \
	       service/user/pb/user.proto
	protoc --go_out=. --go_opt=module=github.com/luyb177/meow-nook \
	       --go-grpc_out=. --go-grpc_opt=module=github.com/luyb177/meow-nook \
	       service/cat/pb/cat.proto
	protoc --go_out=. --go_opt=module=github.com/luyb177/meow-nook \
	       --go-grpc_out=. --go-grpc_opt=module=github.com/luyb177/meow-nook \
	       service/task/pb/task.proto
	protoc --go_out=. --go_opt=module=github.com/luyb177/meow-nook \
	       --go-grpc_out=. --go-grpc_opt=module=github.com/luyb177/meow-nook \
	       service/adoption/pb/adoption.proto
	protoc --go_out=. --go_opt=module=github.com/luyb177/meow-nook \
	       --go-grpc_out=. --go-grpc_opt=module=github.com/luyb177/meow-nook \
	       service/post/pb/post.proto
	@echo "Done."

# ──────────────────────────────────────────────
# Build all services
# ──────────────────────────────────────────────
build:
	@echo "Building all services..."
	go build ./...

# ──────────────────────────────────────────────
# Tests
# ──────────────────────────────────────────────
test:
	go test ./... -count=1

# ──────────────────────────────────────────────
# Linting
# ──────────────────────────────────────────────
lint:
	golangci-lint run ./...

# ──────────────────────────────────────────────
# Tidy
# ──────────────────────────────────────────────
tidy:
	go mod tidy

# ──────────────────────────────────────────────
# Docker Compose (development)
# ──────────────────────────────────────────────
up:
	docker compose -f deploy/docker-compose.yml up -d

down:
	docker compose -f deploy/docker-compose.yml down

# ──────────────────────────────────────────────
# Help
# ──────────────────────────────────────────────
help:
	@echo "Available targets:"
	@echo "  proto  - regenerate protobuf Go code"
	@echo "  build  - build all services"
	@echo "  test   - run all tests"
	@echo "  lint   - run golangci-lint"
	@echo "  tidy   - run go mod tidy"
	@echo "  up     - start all services via docker compose"
	@echo "  down   - stop all services via docker compose"
