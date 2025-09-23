# Makefile for Billetera Backend

# Load environment variables from .env file
-include .env

# Ensure that the required variables are exported to be available for the shell commands
export USER_DB PASSWORD_DB HOST_DB PORT_DB DATABASE_DB

# Construct the database URL from environment variables
DB_URL := postgres://$(USER_DB):$(PASSWORD_DB)@$(HOST_DB):$(PORT_DB)/$(DATABASE_DB)?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create build

## Run all pending up migrations
migrate-up:
	@echo "Running up migrations..."
	@go run ./cmd/migrate/main.go up

## Revert the last applied migration
migrate-down:
	@echo "Running down migrations..."
	@go run ./cmd/migrate/main.go down

build:
	@go build -o billetera_backend .

## Create a new migration file. Requires a 'name' argument.
## Example: make migrate-create name=add_user_table
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)..."
	@migrate create -ext sql -dir db/migrations $(name)