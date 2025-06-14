# go-locate Makefile

# Project-specific variables
BINARY_NAME := glocate
OUTPUT_DIR := bin
CMD_DIR := cmd/glocate

TAG_NAME ?= $(shell head -n 1 .release-version 2>/dev/null || echo "v0.0.0")
VERSION_RAW ?= $(shell tail -n 1 .release-version 2>/dev/null || echo "dev")
VERSION ?= $(VERSION_RAW)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GO_VERSION := $(shell cat .go-version 2>/dev/null || echo "1.24.2")
GO_FILES := $(wildcard $(CMD_DIR)/*.go internal/**/*.go)
GOPATH ?= $(shell go env GOPATH)
GOLANGCI_LINT = $(GOPATH)/bin/golangci-lint
STATICCHECK = $(GOPATH)/bin/staticcheck
GOIMPORTS = $(GOPATH)/bin/goimports

# Ensure the output directory exists
$(OUTPUT_DIR):
	@mkdir -p $(OUTPUT_DIR)

# Default target
.PHONY: default
default: fmt vet imports lint staticcheck build quicktest

# Display help information
.PHONY: help
help:
	@echo "go-locate - Modern File Search Tool"
	@echo ""
	@echo "Available targets:"
	@echo "  Building and Running:"
	@echo "  ===================="
	@echo "  default         - Run formatting, vetting, linting, staticcheck, build, and quick tests"
	@echo "  run             - Run the application locally"
	@echo "  dev             - Run in development mode with hot reload"
	@echo "  build           - Build the application for the current OS/architecture"
	@echo "  build-debug     - Build debug version with debug symbols"
	@echo "  build-cross     - Build binaries for multiple platforms (Linux, macOS, Windows)"

	@echo "  install         - Install binary to /usr/local/bin"
	@echo "  uninstall       - Remove binary from /usr/local/bin"
	@echo ""
	@echo "  Testing and Validation:"
	@echo "  ======================"
	@echo "  test            - Run all tests with standard coverage"
	@echo "  test-with-race  - Run all tests with race detection and coverage"

	@echo "  quicktest       - Run quick tests without additional checks"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-race       - Run tests with race detection"
	@echo "  test-integration- Run integration tests"
	@echo "  test-all        - Run all tests and benchmarks"
	@echo ""
	@echo "  Benchmarking:"
	@echo "  ============="
	@echo "  benchmark       - Run basic benchmarks"
	@echo "  benchmark-long  - Run comprehensive benchmarks with longer duration"
	@echo "  benchmark-search- Run file search benchmarks"

	@echo "  benchmark-report- Generate a markdown report of all benchmarks"
	@echo ""
	@echo "  Code Quality:"
	@echo "  ============"
	@echo "  fmt             - Check and format Go code"
	@echo "  vet             - Analyze code with go vet"
	@echo "  imports         - Format imports with goimports"
	@echo "  lint            - Run golangci-lint"
	@echo "  lint-fix        - Run linters with auto-fix"
	@echo "  staticcheck     - Run staticcheck static analyzer"
	@echo "  check-all       - Run all code quality checks"
	@echo ""
	@echo "  Dependencies:"
	@echo "  ============="
	@echo "  deps            - Install project dependencies"
	@echo "  install-deps    - Install project dependencies (alias for deps)"
	@echo "  upgrade-deps    - Upgrade all dependencies to latest versions"
	@echo "  clean-deps      - Clean up dependencies"
	@echo "  install-tools   - Install development tools"
	@echo ""
	@echo "  Version Management:"
	@echo "  =================="
	@echo "  version         - Show current version information"
	@echo "  bump-patch      - Bump patch version"
	@echo "  bump-minor      - Bump minor version"
	@echo "  bump-major      - Bump major version"
	@echo "  release         - Build release version with all optimizations"
	@echo ""
	@echo "  Cleanup:"
	@echo "  ========"
	@echo "  clean           - Clean build artifacts"
	@echo "  clean-coverage  - Clean coverage and benchmark files"
	@echo "  clean-all       - Clean everything including dependencies"
	@echo ""
	@echo "Examples:"
	@echo "  make build              - Build the binary"
	@echo "  make test               - Run all tests"
	@echo "  make build-cross        - Build for multiple platforms"
	@echo "  make run ARGS=\"*.go\"    - Run with arguments"
	@echo ""
	@echo "For CLI usage instructions, run: ./bin/glocate --help"

# Build and run the application locally
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	go run ./$(CMD_DIR) $(ARGS)

# Dependencies
.PHONY: deps install-deps upgrade-deps clean-deps install-tools
deps: install-deps

install-deps:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully"

upgrade-deps:
	@echo "Upgrading all dependencies to latest versions..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies upgraded. Please test thoroughly before committing!"

clean-deps:
	@echo "Cleaning up dependencies..."
	rm -rf vendor

install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Development tools installed successfully"

# Build targets
.PHONY: build build-debug build-cross

build: $(OUTPUT_DIR)
	@echo "Building $(BINARY_NAME) with version $(VERSION)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build \
		-ldflags="-X 'main.Version=$(VERSION)'" \
		-o $(OUTPUT_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

build-debug: $(OUTPUT_DIR)
	@echo "Building debug version..."
	CGO_ENABLED=0 go build \
		-gcflags="all=-N -l" \
		-ldflags="-X 'main.Version=$(VERSION)'" \
		-o $(OUTPUT_DIR)/$(BINARY_NAME)-debug ./$(CMD_DIR)

build-cross: $(OUTPUT_DIR)
	@echo "Building cross-platform binaries..."
	GOOS=linux   GOARCH=amd64   CGO_ENABLED=0 go build -ldflags="-X 'main.Version=$(VERSION)'" -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=darwin  GOARCH=arm64   CGO_ENABLED=0 go build -ldflags="-X 'main.Version=$(VERSION)'" -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=darwin  GOARCH=amd64   CGO_ENABLED=0 go build -ldflags="-X 'main.Version=$(VERSION)'" -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64   CGO_ENABLED=0 go build -ldflags="-X 'main.Version=$(VERSION)'" -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "Cross-platform binaries are available in $(OUTPUT_DIR):"
	@ls -1 $(OUTPUT_DIR)



# Development targets
.PHONY: dev run-built

dev:
	@echo "Running in development mode..."
	go run ./$(CMD_DIR) $(ARGS)

run-built: build
	./$(OUTPUT_DIR)/$(BINARY_NAME) $(ARGS)

# Testing
.PHONY: test test-with-race quicktest test-coverage test-race test-integration test-all

test:
	@echo "Running Go tests..."
	go test -v ./... -cover

test-with-race:
	@echo "Running all tests with race detection and coverage..."
	go test -v -race -cover ./...



quicktest:
	@echo "Running quick tests..."
	go test ./...

test-coverage:
	@echo "Running tests with coverage report..."
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./...

test-integration: build
	@echo "Running integration tests..."
	# TODO: Add integration tests when implemented
	@echo "Integration tests not yet implemented"

test-all: test-coverage test-race benchmark
	@echo "All tests and benchmarks completed"

# Benchmark targets
.PHONY: benchmark benchmark-long benchmark-search benchmark-report

benchmark:
	@echo "Running benchmarks..."
	go test -v -bench=. -benchmem ./...

benchmark-long:
	@echo "Running comprehensive benchmarks (longer duration)..."
	go test -v -bench=. -benchmem -benchtime=5s ./...

benchmark-search: build
	@echo "Running file search benchmarks..."
	# TODO: Add search-specific benchmarks when implemented
	@echo "Search benchmarks not yet implemented"



benchmark-report:
	@echo "Generating benchmark report..."
	@echo "# Benchmark Results" > benchmark-report.md
	@echo "\nGenerated on \`$$(date)\`\n" >> benchmark-report.md
	@echo "## Go Benchmarks" >> benchmark-report.md
	@echo "| Test | Iterations | Time/op | Memory/op | Allocs/op |" >> benchmark-report.md
	@echo "|------|------------|---------|-----------|-----------|" >> benchmark-report.md
	@go test -bench=. -benchmem ./... 2>/dev/null | grep "Benchmark" | awk '{print "| " $$1 " | " $$2 " | " $$3 " | " $$5 " | " $$7 " |"}' >> benchmark-report.md
	@echo "Benchmark report generated: benchmark-report.md"

# Code quality
.PHONY: fmt vet imports lint lint-fix staticcheck check-all

fmt:
	@echo "Checking and formatting code..."
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "Code formatting completed"

vet:
	@echo "Running go vet..."
	go vet ./...

# Run goimports
.PHONY: imports
imports:
	@if command -v $(GOIMPORTS) >/dev/null 2>&1; then \
		echo "Running goimports..."; \
		$(GOIMPORTS) -local github.com/Gosayram/go-locate -w $(GO_FILES); \
		echo "Imports formatting completed!"; \
	else \
		echo "goimports is not installed. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		echo "Running goimports..."; \
		$(GOIMPORTS) -local github.com/Gosayram/go-locate -w $(GO_FILES); \
		echo "Imports formatting completed!"; \
	fi

# Run linter
.PHONY: lint
lint:
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "Running linter..."; \
		$(GOLANGCI_LINT) run; \
		echo "Linter ended!"; \
	else \
		echo "golangci-lint is not installed. Skipping linter. Run 'make install-lint' to install."; \
	fi

# Run staticcheck tool
.PHONY: staticcheck
staticcheck:
	@if command -v $(STATICCHECK) >/dev/null 2>&1; then \
		echo "Running staticcheck..."; \
		$(STATICCHECK) ./...; \
		echo "Staticcheck ended!"; \
	else \
		echo "staticcheck is not installed. Skipping staticcheck. Run 'make install-staticcheck' to install."; \
	fi

lint-fix:
	@echo "Running linters with auto-fix..."
	@$(GOLANGCI_LINT) run --fix
	@echo "Auto-fix completed"

check-all: imports lint staticcheck
	@echo "All code quality checks completed"

# Release and installation
.PHONY: release install uninstall

release: test lint staticcheck
	@echo "Building release version $(VERSION)..."
	@mkdir -p $(OUTPUT_DIR)
	cd $(RUST_DIR) && cargo build --release
	CGO_ENABLED=1 go build \
		-ldflags="-X 'main.Version=$(VERSION)' -s -w" \
		-o $(OUTPUT_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@strip $(OUTPUT_DIR)/$(BINARY_NAME) 2>/dev/null || true
	@echo "Release build completed: $(OUTPUT_DIR)/$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(OUTPUT_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation completed"

uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation completed"

# Cleanup
.PHONY: clean clean-coverage clean-all

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(OUTPUT_DIR)
	rm -f coverage.out coverage.html benchmark-report.md
	go clean -cache
	@echo "Cleanup completed"

clean-coverage:
	@echo "Cleaning coverage and benchmark files..."
	rm -f coverage.out coverage.html benchmark-report.md
	@echo "Coverage files cleaned"

clean-all: clean clean-deps
	@echo "Deep cleaning everything including dependencies..."
	go clean -modcache
	@echo "Deep cleanup completed"

# Version management
.PHONY: version bump-patch bump-minor bump-major

version:
	@echo "Project: go-locate"
	@echo "Go version: $(GO_VERSION)"
	@echo "Release version: $(VERSION)"
	@echo "Tag name: $(TAG_NAME)"
	@echo "Build target: $(GOOS)/$(GOARCH)"

bump-patch:
	@current=$$(cat .release-version); \
	new=$$(echo $$current | awk -F. '{$$3=$$3+1; print $$1"."$$2"."$$3}'); \
	echo $$new > .release-version; \
	echo "Version bumped from $$current to $$new"

bump-minor:
	@current=$$(cat .release-version); \
	new=$$(echo $$current | awk -F. '{$$2=$$2+1; $$3=0; print $$1"."$$2"."$$3}'); \
	echo $$new > .release-version; \
	echo "Version bumped from $$current to $$new"

bump-major:
	@current=$$(cat .release-version); \
	new=$$(echo $$current | awk -F. '{$$1=$$1+1; $$2=0; $$3=0; print $$1"."$$2"."$$3}'); \
	echo $$new > .release-version; \
	echo "Version bumped from $$current to $$new"

# Docker support (for future use)
.PHONY: docker-build docker-run

docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(TAG_NAME) .
	@echo "Docker image built: $(BINARY_NAME):$(TAG_NAME)"

docker-run:
	@echo "Running Docker image..."
	docker run -it --rm $(BINARY_NAME):$(TAG_NAME) --version

# Documentation
.PHONY: docs docs-api

docs:
	@echo "Generating documentation..."
	@mkdir -p docs
	@echo "# go-locate API Documentation" > docs/api.md
	@echo "\nGenerated on \`$$(date)\`\n" >> docs/api.md
	@go doc -all ./... >> docs/api.md
	@echo "Documentation generated in docs/api.md"

docs-api:
	@echo "Generating API documentation..."
	@mkdir -p docs
	go doc -all ./... > docs/api.md
	@echo "API documentation generated" 