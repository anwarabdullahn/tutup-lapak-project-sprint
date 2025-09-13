# ==========================
# Variables
# ==========================
MIGRATIONS_DIR=./migrations/db
# Read DATABASE_URL from .env if present; allow override via environment
DB_URL?=$(shell [ -f .env ] && grep -E '^DATABASE_URL=' .env | cut -d '=' -f2-)
MIGRATE_CMD=migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)"

.PHONY: check-db-url
check-db-url:
	@if [ -z "$(DB_URL)" ]; then \
		echo "❌ DATABASE_URL is empty."; \
		echo "   Create auth-service/.env with DATABASE_URL=... or pass inline:"; \
			echo "   DB_URL=postgres://user:pass@host:5432/db?sslmode=disable make migrate:up"; \
		exit 1; \
	fi


# ==========================
# Migration Commands
# ==========================
.PHONY: migrate\:create
migrate\:create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a migration name: make migrate-create name=create_table_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -format 20060102150405 $(name)

.PHONY: migrate\:up
migrate\:up: check-db-url
	$(MIGRATE_CMD) up

.PHONY: migrate\:up-1
migrate\:up-1: check-db-url
	$(MIGRATE_CMD) up 1

.PHONY: migrate\:down
migrate\:down: check-db-url
	$(MIGRATE_CMD) down

.PHONY: migrate\:down-1
migrate\:down-1: check-db-url
	$(MIGRATE_CMD) down 1

.PHONY: migrate\:version
migrate\:version: check-db-url
	$(MIGRATE_CMD) version

.PHONY: swagger-gen
swagger-gen:
	swag init --generalInfo app.go --output docs

# Backward compatible dash aliases
.PHONY: migrate-create migrate-up migrate-up-1 migrate-down migrate-down-1 migrate-version
migrate-create: migrate\:create
	@echo "⚠️  'migrate-create' is deprecated. Use 'make migrate:create'"
migrate-up: migrate\:up
	@echo "⚠️  'migrate-up' is deprecated. Use 'make migrate:up'"
migrate-up-1: migrate\:up-1
	@echo "⚠️  'migrate-up-1' is deprecated. Use 'make migrate:up-1'"
migrate-down: migrate\:down
	@echo "⚠️  'migrate-down' is deprecated. Use 'make migrate:down'"
migrate-down-1: migrate\:down-1
	@echo "⚠️  'migrate-down-1' is deprecated. Use 'make migrate:down-1'"
migrate-version: migrate\:version
	@echo "⚠️  'migrate-version' is deprecated. Use 'make migrate:version'"
# ==========================
# Help
# ==========================
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  migrate:create     Create a new migration (usage: make migrate:create name=create_table_users)"
	@echo "  migrate:up         Apply all migrations"
	@echo "  migrate:up-1       Apply one migration step"
	@echo "  migrate:down       Rollback all migrations"
	@echo "  migrate:down-1     Rollback one migration step"
	@echo "  migrate:version    Show current migration version"
