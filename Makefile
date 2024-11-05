.PHONY: all build run test clean lint wire docker-build docker-run

# Variables
BINARY_NAME=iam-wrapper
DOCKER_IMAGE=iam-wrapper-service
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

all: wire clean lint test build

# Generate wire_gen.go
wire:
	@echo "Generating wire_gen.go..."
	@wire ./internal/di

# Build the application
build: wire
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) ./cmd/api

# Run the application
run: wire
	@go run ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

# Clean build files
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Run with hot reload
dev: wire
	@air

# Generate documentation
docs:
	@echo "Generating documentation..."
	@swag init -g cmd/api/main.go

# Docker commands
docker-build: wire
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Docker Compose commands
docker-compose-up: wire
	@echo "Starting Docker Compose services..."
	@docker-compose up --build -d

docker-compose-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/google/wire/cmd/wire@latest

# Create necessary directories
init:
	@echo "Initializing project structure..."
	@mkdir -p cmd/api internal/auth internal/middleware internal/server pkg/logger pkg/errors api/v1 config docs scripts test build

# Help
help:
	@echo "Available commands:"
	@echo "  make wire            - Generate wire_gen.go"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build files"
	@echo "  make lint           - Run linter"
	@echo "  make dev            - Run with hot reload"
	@echo "  make docs           - Generate documentation"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"
	@echo "  make install-tools  - Install development tools"
