.PHONY: help run build test clean docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run cmd/server/main.go

build: ## Build the application
	go build -o bin/server cmd/server/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

docker-build: ## Build Docker image
	docker build -t url-shortener:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env url-shortener:latest

install: ## Install dependencies
	go mod download
	go mod tidy

dev: ## Run in development mode with hot reload (requires air)
	air
