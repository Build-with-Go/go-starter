.PHONY: help build run test lint clean deps fmt vet check doctorreport

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
deps: ## Download dependencies
	go mod download
	go mod tidy

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

check: fmt vet lint ## Run all checks (fmt, vet, lint)

doctorreport: ## Run Doctor Report analysis
	@echo "🩺 Running Doctor Report analysis..."
	@echo "📝 Checking code formatting..."
	@go fmt ./...
	@echo "🔍 Running go vet analysis..."
	@go vet ./...
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "🔨 Verifying build..."
	@go build ./...
	@echo "✅ Doctor Report completed!"

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: ## Run tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build targets
build: ## Build the application
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/server cmd/server/main.go

build-debug: ## Build the application with debug info
	go build -o bin/server cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

run-config: ## Run with specific config file
	@if [ -z "$(CONFIG)" ]; then \
		echo "Usage: make run-config CONFIG=path/to/config.yaml"; \
		exit 1; \
	fi
	go run cmd/server/main.go -config $(CONFIG)

# Utility targets
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f go.sum

dev-setup: ## Set up development environment
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	go mod download

dev: ## Run with hot reload (requires air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not installed. Install with: make dev-setup"; \
		exit 1; \
	fi

# Docker targets
docker-build: ## Build Docker image
	docker build -t go-starter:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file configs/.env go-starter:latest

# CI/CD targets
ci: deps check doctorreport test ## Run CI pipeline locally

release: ## Build release binaries
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/server-linux-amd64 cmd/server/main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/server-darwin-amd64 cmd/server/main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/server-windows-amd64.exe cmd/server/main.go
