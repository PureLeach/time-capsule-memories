# DATABASE_URL for the migration targets, when a .env is present.
-include .env

.PHONY: help build run test lint fmt docs migrate-up migrate-down migrate-create clean

help: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  %-16s %s\n", $$1, $$2}'

build: ## Compile the server binary into bin/
	go build -trimpath -o bin/time-capsule-memories ./cmd

run: ## Run the server against the local .env
	go run ./cmd

test: ## Run the test suite with the race detector
	go test -race -count=1 ./...

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format and tidy
	go fmt ./...
	go mod tidy

docs: ## Regenerate the Swagger spec from handler annotations
	swag init -g cmd/main.go

migrate-up: ## Apply pending migrations
	goose -dir ./migrations postgres "$(DATABASE_URL)" up

migrate-down: ## Roll back the last migration
	goose -dir ./migrations postgres "$(DATABASE_URL)" down

migrate-create: ## Create a migration: make migrate-create name=add_foo
	goose -dir ./migrations create $(name) sql

clean: ## Remove build artifacts
	rm -rf bin/
