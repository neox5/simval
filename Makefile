.PHONY: help test test-short test-race

# Default target - show available commands
help:
	@echo "Available targets:"
	@echo "  make test       - Run all tests including stress tests"
	@echo "  make test-short - Run only non-stress tests (matches CI)"
	@echo "  make test-race  - Run tests with race detector"

# Run all tests including stress tests
test:
	go test ./...

# Run only non-stress tests (matches CI behavior)
test-short:
	go test -short ./...

# Run tests with race detector
test-race:
	go test -race ./...
