DB_NAME=databases/tinyis.db
SCHEMA_FILE=scripts/sqlite3.sql
JWT_ALGO ?= ed25519
JWT_KEY_DIR = resources/crypto/jwt
SERVER_KEY_DIR = resources/crypto/server

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

generate_jwt_key:
	mkdir -p $(JWT_KEY_DIR)
	openssl genpkey -algorithm $(JWT_ALGO) -out $(JWT_KEY_DIR)/eddsa.pem

generate_server_keypair:
	mkdir -p $(SERVER_KEY_DIR)
	openssl req -newkey ed25519 -keyout $(SERVER_KEY_DIR)/server-key.pem -out $(SERVER_KEY_DIR)/server-cert.pem -x509 -nodes -days 365

all: coverage build run
