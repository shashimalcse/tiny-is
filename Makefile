.PHONY: test build run

test:
	@echo "Running tests..."
	@go test -v ./... | grep -v "\[no test files\]" || (echo "Tests failed. Build aborted." && exit 1)
	@echo "Tests completed."	

build: test
	go build -o tinyis

run:
	./tinyis
