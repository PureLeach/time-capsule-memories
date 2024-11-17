build:
	go build -o bin/time_capsule_memories ./cmd/main.go


start:
	./bin/time_capsule_memories


run:
	go run -v ./cmd/main.go


docs_generation:
	swag init -g cmd/main.go


clean:
	go clean -cache && go clean -modcache


executable_scripts:
	chmod +x db_setup.sh migrate.sh

migrate:
	./migrate.sh