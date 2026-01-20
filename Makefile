.PHONY: all build build-server build-cli run test lint clean migrate dev

# Binary names
SERVER_BINARY=trelay-server
CLI_BINARY=trelay

# Build directories
BIN_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

all: build

build: build-server build-cli

build-server:
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(SERVER_BINARY) ./cmd/server

build-cli:
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(CLI_BINARY) ./cmd/trelay

run: build-server
	./$(BIN_DIR)/$(SERVER_BINARY)

dev:
	$(GOCMD) run ./cmd/server

test:
	$(GOTEST) -v -race -cover ./...

test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

deps:
	$(GOMOD) download
	$(GOMOD) tidy

migrate-up:
	$(GOCMD) run ./cmd/server migrate up

migrate-down:
	$(GOCMD) run ./cmd/server migrate down

migrate-create:
	@read -p "Migration name: " name; \
	goose -dir migrations create $$name sql
