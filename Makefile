# Makefile for TRexT

.PHONY: sqlc goose-up goose-down goose-status goose-create generate build run

# Database URL - set this in your environment or .env file
# Example: export DATABASE_URL="postgres://user:password@localhost:5432/trext?sslmode=disable"

# Generate Go code from SQL queries using sqlc
sqlc:
	sqlc generate

# Run all goose migrations up
goose-up:
	goose -dir sql/schemas postgres "$(DATABASE_URL)" up

# Rollback the last goose migration
goose-down:
	goose -dir sql/schemas postgres "$(DATABASE_URL)" down

# Show goose migration status
goose-status:
	goose -dir sql/schemas postgres "$(DATABASE_URL)" status

# Create a new goose migration
goose-create:
	@read -p "Enter migration name: " name; \
	goose -dir sql/schemas create $$name sql

# Generate all code (sqlc)
generate: sqlc

# Build the application
build:
	go build -o trext ./cmd/trext

# Run the application (requires DATABASE_URL to be set)
run: build
	./trext
