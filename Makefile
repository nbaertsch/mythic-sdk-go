.PHONY: help test lint coverage integration-test bench install-tools clean

# Variables
GO := go
GOTEST := $(GO) test
GOVET := $(GO) vet
GOLINT := golangci-lint
COVERAGE_DIR := coverage

# Default target
help:
	@echo "Mythic Go SDK - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make test              - Run unit tests"
	@echo "  make integration-test  - Run integration tests (requires Docker)"
	@echo "  make lint              - Run linters"
	@echo "  make coverage          - Generate coverage report"
	@echo "  make bench             - Run benchmarks"
	@echo "  make install-tools     - Install development tools"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make all               - Run lint and test"

# Run unit tests
test:
	@echo "Running unit tests..."
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@echo "✓ Tests passed"

# Run integration tests
integration-test:
	@echo "Starting integration tests..."
	@./scripts/integration-test.sh

# Run linters
lint:
	@echo "Running linters..."
	$(GOLINT) run ./...
	$(GOVET) ./...
	@echo "✓ Linting passed"

# Generate coverage report
coverage: test
	@echo "Generating coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out
	@echo "✓ Coverage report generated: $(COVERAGE_DIR)/coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/golang/mock/mockgen@latest
	@echo "✓ Tools installed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -rf $(COVERAGE_DIR)/
	@rm -rf tmp/
	@echo "✓ Cleaned"

# Run all checks
all: lint test
	@echo "✓ All checks passed"

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "✓ Code formatted"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy
	@echo "✓ Dependencies tidied"

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "✓ Dependencies updated"

# Run security scan
security:
	@echo "Running security scan..."
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; $(GO) install github.com/securego/gosec/v2/cmd/gosec@latest; }
	gosec ./...
	@echo "✓ Security scan complete"

# Generate mocks
mocks:
	@echo "Generating mocks..."
	@command -v mockgen >/dev/null 2>&1 || { echo "mockgen not found. Run 'make install-tools'"; exit 1; }
	$(GO) generate ./...
	@echo "✓ Mocks generated"

# Initialize project
init: install-tools tidy
	@echo "Project initialized"
