.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down swagger lint fmt

# Variables
APP_NAME=golang-ddd-template
MAIN_PATH=cmd/api/main.go
BINARY_NAME=bin/$(APP_NAME)
DOCKER_COMPOSE=docker-compose

# Colors for terminal output
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

.PHONY: run
run:
	@echo "Generating docs..."
	swag init -g cmd/api/main.go
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "Running application..."
	@go run $(MAIN_PATH)

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit: ## Run unit tests only
	@go test -v -short ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

docker-up: ## Start all Docker containers
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Containers started. Waiting for services to be ready..."
	@sleep 5
	@$(DOCKER_COMPOSE) ps

docker-down: ## Stop all Docker containers
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE) down
	@echo "Containers stopped"

docker-logs: ## View Docker container logs
	@$(DOCKER_COMPOSE) logs -f

docker-rebuild: ## Rebuild and restart Docker containers
	@echo "Rebuilding containers..."
	@$(DOCKER_COMPOSE) up -d --build
	@echo "Rebuild complete"

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide migration name. Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir migrations -seq $(name)
	@echo "Migration files created in migrations/"

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/ddd_template?sslmode=disable" up
	@echo "Migrations complete"

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	@migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/ddd_template?sslmode=disable" down 1
	@echo "Rollback complete"

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "Error: Please provide version. Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	@migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/ddd_template?sslmode=disable" force $(version)

seed: ## Seed database with initial data
	@echo "Seeding database..."
	@go run scripts/seed.go
	@echo "Seeding complete"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o docs
	@echo "Swagger docs generated in docs/"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .
	@echo "Code formatted"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Tools installed"

dev: docker-up migrate-up run ## Start development environment

setup: deps install-tools docker-up migrate-up swagger ## Setup project for first time

.DEFAULT_GOAL := help
