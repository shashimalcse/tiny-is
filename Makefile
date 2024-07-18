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

all: coverage build run
