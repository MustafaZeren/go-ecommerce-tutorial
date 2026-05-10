.PHONY: help build run dev lint migrate-up migrate-down

DB_URL=postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make dev         - Run the application in development mode"
	@echo "  make lint        - Run linter on the codebase"
	@echo "  make format      - Format the code and re-arrange imports"
	@echo "  make migrate-up  - Apply database migrations"
	@echo "  make migrate-down- Rollback database migrations"

build:
	go build -o bin/app ./cmd/api

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint:
	golangci-lint run ./...

migrate-up:
	docker run --rm -v $(shell pwd)/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "$(DB_URL)" up

migrate-down:
	docker run --rm -v $(shell pwd)/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "$(DB_URL)" down 1

docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down

docker-tools-up:
	docker compose -f docker/docker-compose.tools.yml up -d

docker-tools-down:
	docker compose -f docker/docker-compose.tools.yml down

docker-all-up: docker-tools-up docker-up
	@echo "All systems are up and running!"

docker-all-down: docker-down docker-tools-down
	@echo "All systems are stopped."