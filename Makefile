.PHONY: up down logs test lint fmt migrate-up migrate-down migrate-create help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  %-15s %s\n", $$1, $$2}'

up: ## Start the full stack
	docker compose up

down: ## Stop and remove containers
	docker compose down

logs: ## Tail logs from all services
	docker compose logs -f

test: ## Run backend tests
	cd backend && $(MAKE) test

lint: ## Run backend and frontend linters
	cd backend && $(MAKE) lint
	cd frontend && npm run lint

fmt: ## Format the frontend (the backend is auto-formatted by golangci-lint)
	cd frontend && npm run format

migrate-up: ## Apply pending DB migrations
	cd backend && $(MAKE) migrate_up

migrate-down: ## Roll back the last DB migration
	cd backend && $(MAKE) migrate_down

migrate-create: ## Create a new migration: make migrate-create name=add_foo
	cd backend && $(MAKE) create_migration name=$(name)
