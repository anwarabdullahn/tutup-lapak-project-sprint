APP_NAME=product-service
PORT?=3003
MIGRATIONS_DIR=./migrations/db
DB_URL?=$(shell [ -f .env ] && grep -E '^DATABASE_URL=' .env | cut -d '=' -f2-)
MIGRATE_CMD=migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)"

.PHONY: run build docker-build tidy lint check-db-url migrate\:create migrate\:up migrate\:up-1 migrate\:down migrate\:down-1 migrate\:version help

run:
	go run ./cmd/server

build:
	go build -o bin/$(APP_NAME) ./cmd/server

docker-build:
	docker build -t $(APP_NAME):local .

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

check-db-url:
	@if [ -z "$(DB_URL)" ]; then \
		echo "❌ DATABASE_URL is empty. Provide via .env or inline DB_URL=... make migrate:up"; \
		exit 1; \
	fi

migrate\:create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a migration name: make migrate:create name=create_table_products"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -format 20060102150405 $(name)

migrate\:up: check-db-url
	$(MIGRATE_CMD) up

migrate\:up-1: check-db-url
	$(MIGRATE_CMD) up 1

migrate\:down: check-db-url
	$(MIGRATE_CMD) down

migrate\:down-1: check-db-url
	$(MIGRATE_CMD) down 1

migrate\:version: check-db-url
	$(MIGRATE_CMD) version

# Backward compatible simple aliases
migrate-create: migrate\:create
migrate-up: migrate\:up
migrate-up-1: migrate\:up-1
migrate-down: migrate\:down
migrate-down-1: migrate\:down-1
migrate-version: migrate\:version

help:
	@echo "Commands:"
	@echo "  run                Run service"
	@echo "  build              Build binary"
	@echo "  docker-build       Build docker image"
	@echo "  migrate:create name=...  Create new migration"
	@echo "  migrate:up         Apply all migrations"
	@echo "  migrate:down       Rollback all migrations"
	@echo "  migrate:version    Show migration version"
