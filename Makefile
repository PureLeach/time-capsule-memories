.PHONY: help up up-prod down logs build test lint fmt migrate-up migrate-down migrate-create

help: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  %-16s %s\n", $$1, $$2}'

up: ## Start the full stack, including the development tooling
	docker compose up --build

up-prod: ## Start without the development override (no pgAdmin, MailHog or dashboard)
	docker compose -f docker-compose.yml up --build

down: ## Stop and remove containers
	docker compose down

logs: ## Tail logs from all services
	docker compose logs -f

build: ## Build both container images
	docker compose build

test: ## Run backend tests
	$(MAKE) -C backend test

lint: ## Run backend and frontend linters
	$(MAKE) -C backend lint
	cd frontend && npm run lint && npm run format:check

fmt: ## Format both sides of the project
	$(MAKE) -C backend fmt
	cd frontend && npm run format

migrate-up: ## Apply pending DB migrations
	$(MAKE) -C backend migrate-up

migrate-down: ## Roll back the last DB migration
	$(MAKE) -C backend migrate-down

migrate-create: ## Create a migration: make migrate-create name=add_foo
	$(MAKE) -C backend migrate-create name=$(name)
