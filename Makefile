# Makefile for dotenv Go package

.PHONY: all test install fmt vet

# Default target
all: fmt vet test

# Run tests
test:
	go test ./...

# Install the package
install:
	go install ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...
