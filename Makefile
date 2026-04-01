.PHONY: help build run stop clean logs test

help: ## Show this help message
	@echo "Available commands:"
	@echo "  make build   - Build Docker containers"
	@echo "  make run     - Start all services"
	@echo "  make stop    - Stop all services"
	@echo "  make restart - Restart all services"
	@echo "  make logs    - View backend logs"
	@echo "  make clean   - Remove containers and volumes"
	@echo "  make test    - Run tests"
	@echo "  make dev     - Run locally without Docker"

build: ## Build Docker containers
	docker-compose build

run: ## Start all services
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "🌐 API: http://localhost:8080"
	@echo "🗄️ PostgreSQL: localhost:5432"

stop: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose down
	docker-compose up -d

logs: ## View backend logs
	docker-compose logs -f backend

clean: ## Remove containers and volumes
	docker-compose down -v
	@echo "✅ Cleaned up!"

test: ## Run tests
	cd backend && go test ./...

dev: ## Run locally without Docker
	cd backend && go run cmd/api/main.go

deps: ## Download Go dependencies
	cd backend && go mod download

tidy: ## Tidy Go dependencies
	cd backend && go mod tidy

psql: ## Connect to PostgreSQL
	docker exec -it mechanic_postgres psql -U postgres -d mechanic_db
