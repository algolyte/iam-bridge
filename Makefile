.PHONY: all build run test clean docker-build docker-run lint help

# Variables
BINARY_NAME=goiam-bridge
DOCKER_IMAGE=goiam-bridge
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
GOLANGCI_LINT_VERSION=v1.54.2

# Colors for output
COLOR_RESET=\033[0m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m

all: help

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

build: ## Build the application
	@echo "$(COLOR_GREEN)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@go build -o bin/$(BINARY_NAME) cmd/main.go

run: ## Run the application
	@echo "$(COLOR_GREEN)Running $(BINARY_NAME)...$(COLOR_RESET)"
	@go run cmd/app/main.go

test: ## Run tests
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	@go test -v -race ./...

clean: ## Clean build artifacts
	@echo "$(COLOR_GREEN)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf bin/
	@go clean

docker-build: ## Build Docker image
	@echo "$(COLOR_GREEN)Building Docker image...$(COLOR_RESET)"
	@docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	@echo "$(COLOR_GREEN)Running Docker container...$(COLOR_RESET)"
	@docker-compose up -d

docker-stop: ## Stop Docker container
	@echo "$(COLOR_GREEN)Stopping Docker container...$(COLOR_RESET)"
	@docker-compose down

lint: install-lint ## Run linter
	@echo "$(COLOR_GREEN)Running linter...$(COLOR_RESET)"
	@golangci-lint run ./...

install-lint: ## Install golangci-lint
	@echo "$(COLOR_GREEN)Installing golangci-lint...$(COLOR_RESET)"
	@which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

dev: ## Run application with hot reload using Air
	@echo "$(COLOR_GREEN)Running with hot reload...$(COLOR_RESET)"
	@which air || go install github.com/cosmtrek/air@latest
	@air

tidy: ## Tidy and verify go modules
	@echo "$(COLOR_GREEN)Tidying up modules...$(COLOR_RESET)"
	@go mod tidy
	@go mod verify

coverage: ## Run tests with coverage
	@echo "$(COLOR_GREEN)Running tests with coverage...$(COLOR_RESET)"
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_YELLOW)Coverage report generated: coverage.html$(COLOR_RESET)"

fmt: ## Format code
	@echo "$(COLOR_GREEN)Formatting code...$(COLOR_RESET)"
	@gofmt -s -w .

vet: ## Run go vet
	@echo "$(COLOR_GREEN)Running go vet...$(COLOR_RESET)"
	@go vet ./...

generate: ## Run go generate
	@echo "$(COLOR_GREEN)Running go generate...$(COLOR_RESET)"
	@go generate ./...

init: ## Initialize development environment
	@echo "$(COLOR_GREEN)Initializing development environment...$(COLOR_RESET)"
	@go mod download
	@make install-lint
	@which air || go install github.com/cosmtrek/air@latest
