DB_NAME=databases/tinyis.db
SCHEMA_FILE=scripts/sqlite3.sql

.PHONY: test build run

test:
	@echo "Running tests..."
	@go test -v ./... | grep -v "\[no test files\]" || (echo "Tests failed. Build aborted." && exit 1)
	@echo "Tests completed."

coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=test/results/coverage.out ./...
	@go tool cover -func=test/results/coverage.out
	@echo "Opening coverage report in browser..."
	@go tool cover -html=test/results/coverage.out

build: test
	@echo "Building tinyis..."
	@go build -o tinyis

run:
	@echo "Running tinyis..."
	@./tinyis

create_db:
	@if ! command -v sqlite3 &> /dev/null; then \
		echo "sqlite3 could not be found. Please install sqlite3 to proceed."; \
		exit 1; \
	fi
	@echo "Creating SQLite3 database..."
	@sqlite3 $(DB_NAME) < $(SCHEMA_FILE)
	@echo "Database created at $(DB_NAME)"	

all: coverage build run
