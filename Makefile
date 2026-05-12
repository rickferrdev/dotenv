# Makefile for dotenv Go package

.PHONY: all build test clean install fmt vet

# Default target
all: fmt vet test build

# Build the package
build:
	go build ./...

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	go clean ./...

# Install the package
install:
	go install ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...