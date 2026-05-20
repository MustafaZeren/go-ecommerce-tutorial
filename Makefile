.PHONY: \
	help \
	build \
	run \
	dev \
	format \
	check-format \
	lint \
	test \
	ci \
	migrate-up \
	migrate-down \
	docker-up \
	docker-down \
	docker-tools-up \
	docker-tools-down \
	docker-all-up \
	docker-all-down

DB_URL=postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable

BINARY_NAME=app
CMD_PATH=./cmd/api

help:
	@echo "Available commands:"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application"
	@echo "  make dev              - Run the application in development mode"
	@echo "  make format           - Format code and organize imports"
	@echo "  make check-format     - Verify code formatting"
	@echo "  make lint             - Run golangci-lint"
	@echo "  make test             - Run test suite"
	@echo "  make ci               - Run CI quality checks"
	@echo "  make migrate-up       - Apply database migrations"
	@echo "  make migrate-down     - Rollback database migration"
	@echo "  make docker-up        - Start application containers"
	@echo "  make docker-down      - Stop application containers"
	@echo "  make docker-tools-up  - Start development tools"
	@echo "  make docker-tools-down- Stop development tools"

build:
	go build -o bin/$(BINARY_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

dev:
	go run $(CMD_PATH)

format:
	@gofmt -s -w .
	@goimports -w .

check-format:
	@test -z "$$(gofmt -l .)" || (echo "gofmt check failed"; exit 1)
	@goimports -l . | grep -q . && (echo "goimports check failed"; exit 1) || true

lint:
	golangci-lint run ./...

test:
	go test ./...

ci: check-format lint test build
	@echo "CI quality checks passed successfully!"

migrate-up:
	docker run --rm \
		-v $(shell pwd)/db/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations/ \
		-database "$(DB_URL)" up

migrate-down:
	docker run --rm \
		-v $(shell pwd)/db/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations/ \
		-database "$(DB_URL)" down 1

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