# Tic Tac Toe Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=gofmt
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Binary names
BIN_DIR=bin
SERVER_BINARY=$(BIN_DIR)/server
BOTSERVICE_BINARY=$(BIN_DIR)/botservice
CLI_BINARY=$(BIN_DIR)/cli


# Default target
.PHONY: all
all: fmt lint test build

# Build all binaries
.PHONY: build
build: build-server build-botservice build-cli

# Build server
.PHONY: build-server
build-server:
	@echo "Building server..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(SERVER_BINARY) .

# Build bot service
.PHONY: build-botservice
build-botservice:
	@echo "Building botservice..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD)  -o $(BOTSERVICE_BINARY) ./botservice

# Build CLI
.PHONY: build-cli
build-cli:
	@echo "Building cli..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD)  -o $(CLI_BINARY) ./cli

# Run targets for development
.PHONY: run-server
run-server: build-server
	@echo "Starting server on :8080..."
	$(SERVER_BINARY) -port 8080 -store memory

.PHONY: run-botservice
run-botservice: build-botservice
	@echo "Starting botservice on :9090..."
	$(BOTSERVICE_BINARY) -botservice-port 9090

.PHONY: run-cli
run-cli: build-cli
	@echo "Starting CLI..."
	$(CLI_BINARY) -base-url http://localhost:8080 

# Test targets
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-cover
test-cover:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test-race
test-race:
	@echo "Running tests with race detector..."
	$(GOTEST) -race -v ./...

# Lint and format targets
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

.PHONY: lint
lint:
	@echo "Running linter..."
	$(GOLINT) run ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Dependency management
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

.PHONY: verify
verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify

# Clean targets
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

.PHONY: clean-data
clean-data:
	@echo "Cleaning game data files..."
	rm -f data/*.json

.PHONY: clean-all
clean-all: clean clean-data

# Install binaries to GOPATH/bin
.PHONY: install
install:
	@echo "Installing binaries..."
	$(GOCMD) install .
	$(GOCMD) install ./botservice
	$(GOCMD) install ./cli

# Development helpers
.PHONY: check
check: fmt vet lint test
	@echo "All checks passed!"
