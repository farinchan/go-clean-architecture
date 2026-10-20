.PHONY: help build run dev test clean docker-build docker-up docker-down migrate-up migrate-down swagger lint fmt tidy

# Application
APP_NAME=go-clean-architecture
MAIN_PATH=./cmd/api
MIGRATE_PATH=./cmd/migrate
BINARY_NAME=main.exe
BINARY_PATH=./bin/$(BINARY_NAME)

# Docker
DOCKER_COMPOSE=docker-compose.yml

# Colors
GREEN=
NC=

## help: Show this help message
help:
	@echo 'Usage:'
	@echo ''
	@echo 'Development:'
	@echo '  make dev              - Run application with Air hot reload'
	@echo '  make run              - Run application without hot reload'
	@echo '  make build            - Build the application binary'
	@echo '  make clean            - Clean build artifacts'
	@echo ''
	@echo 'Testing:'
	@echo '  make test             - Run all tests'
	@echo '  make test-coverage    - Run tests with coverage report'
	@echo '  make test-verbose     - Run tests with verbose output'
	@echo ''
	@echo 'Code Quality:'
	@echo '  make lint             - Run linter (golangci-lint)'
	@echo '  make fmt              - Format code'
	@echo '  make tidy             - Tidy go.mod'
	@echo '  make vet              - Run go vet'
	@echo ''
	@echo 'Database:'
	@echo '  make migrate-up       - Run all migrations up'
	@echo '  make migrate-down     - Run all migrations down'
	@echo '  make migrate-create   - Create a new migration (NAME=migration_name)'
	@echo '  make migrate-force    - Force migration version (VERSION=1)'
	@echo '  make migrate-version  - Show current migration version'
	@echo '  make seed             - Run database seeder'
	@echo ''
	@echo 'Swagger:'
	@echo '  make swagger          - Generate Swagger documentation'
	@echo '  make swagger-fmt      - Format Swagger comments'
	@echo ''
	@echo 'Docker:'
	@echo '  make docker-up        - Start containers'
	@echo '  make docker-down      - Stop containers'

	@echo '  make docker-build     - Build Docker image'
	@echo '  make docker-logs      - View container logs'
	@echo ''
	@echo 'Dependencies:'
	@echo '  make deps             - Download dependencies'
	@echo '  make deps-upgrade     - Upgrade all dependencies'
	@echo '  make install-tools    - Install development tools'

## build: Build the application binary
build:
	@echo "Building application..."
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

## run: Run the application
run: build
	@echo "Running application..."
	@$(BINARY_PATH)

## dev: Run application with Air hot reload
dev:
	@echo "Starting development server with hot reload..."
	@air -c .air.toml

## clean: Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## test: Run all tests
test:
	@echo "Running tests..."
	@go test ./... -race

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@go test ./... -race -v

## lint: Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

## tidy: Tidy go modules
tidy:
	@echo "Tidying go.mod..."
	@go mod tidy

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download

## deps-upgrade: Upgrade all dependencies
deps-upgrade:
	@echo "Upgrading dependencies..."
	@go get -u ./...
	@go mod tidy

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed successfully"

## migrate-up: Run all migrations up
migrate-up:
	@echo "Running migrations up..."
	@go run $(MIGRATE_PATH)/main.go -direction up

## migrate-down: Run all migrations down
migrate-down:
	@echo "Running migrations down..."
	@go run $(MIGRATE_PATH)/main.go -direction down

## migrate-create: Create a new migration
migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Please provide NAME=migration_name"; exit 1; fi
	@echo "Creating migration: $(NAME)"
	@migrate create -ext sql -dir database/migrations -seq $(NAME)

## migrate-force: Force migration to a version
migrate-force:
	@if [ -z "$(VERSION)" ]; then echo "Please provide VERSION=1"; exit 1; fi
	@echo "Forcing migration to version $(VERSION)..."
	@go run $(MIGRATE_PATH)/main.go -direction force -force $(VERSION)

## migrate-version: Show current migration version
migrate-version:
	@echo "Getting current migration version..."
	@go run $(MIGRATE_PATH)/main.go -direction version

## seed: Run database seeder
seed:
	@echo "Running database seeder..."
	@go run cmd/seed/main.go

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger documentation generated"

## swagger-fmt: Format Swagger comments
swagger-fmt:
	@echo "Formatting Swagger comments..."
	@swag fmt

## docker-up: Start containers
docker-up:
	@echo "Starting containers..."
	@docker-compose -f $(DOCKER_COMPOSE) up -d
	@echo "Containers started"

## docker-down: Stop containers
docker-down:
	@echo "Stopping containers..."
	@docker-compose -f $(DOCKER_COMPOSE) down

## docker-prod: Alias for docker-up
docker-prod: docker-up

## docker-prod-down: Alias for docker-down
docker-prod-down: docker-down

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .
	@echo "Docker image built: $(APP_NAME):latest"

## docker-logs: View container logs
docker-logs:
	@docker-compose -f $(DOCKER_COMPOSE) logs -f

## docker-clean: Remove all containers and volumes
docker-clean:
	@echo "Removing all containers and volumes..."
	@docker-compose -f $(DOCKER_COMPOSE) down -v --remove-orphans
	@echo "Cleanup complete"
