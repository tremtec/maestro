.PHONY: build build-all install test clean lint release-local

APP_NAME = maestro
BIN_DIR = bin
DIST_DIR = dist
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

# Default build for current platform
build:
	@mkdir -p $(BIN_DIR)
	@go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME) .
	@echo "Built $(BIN_DIR)/$(APP_NAME) (version: $(VERSION))"

# Install locally
install:
	@go install $(LDFLAGS) .
	@echo "Installed $(APP_NAME) to $(shell go env GOPATH)/bin/$(APP_NAME)"

# Run tests
test:
	@go test -v -race ./...

# Run linter
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, running go vet..."; \
		go vet ./...; \
	fi

# Clean build artifacts
clean:
	@rm -rf $(BIN_DIR) $(DIST_DIR)
	@echo "Cleaned build artifacts"

# Build for all platforms (local testing)
build-all: clean
	@mkdir -p $(DIST_DIR)
	@echo "Building for all platforms..."
	
	# Linux builds
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 .
	@GOOS=linux GOARCH=386 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-386 .
	
	# macOS builds
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .
	
	# Windows builds
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .
	@GOOS=windows GOARCH=386 CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows-386.exe .
	
	@echo "All builds complete:"
	@ls -lh $(DIST_DIR)/

# Create release archives locally
release-local: build-all
	@mkdir -p $(DIST_DIR)/release
	@echo "Creating release archives..."
	
	@for file in $(DIST_DIR)/$(APP_NAME)-*; do \
		if [ -f "$$file" ]; then \
			filename=$$(basename "$$file"); \
			platform=$${filename#$(APP_NAME)-}; \
			if [[ "$$filename" == *.exe ]]; then \
				tar -czf "$(DIST_DIR)/release/$(APP_NAME)-$(VERSION)-$${platform%.exe}.tar.gz" -C $(DIST_DIR) "$$filename"; \
			else \
				tar -czf "$(DIST_DIR)/release/$(APP_NAME)-$(VERSION)-$$platform.tar.gz" -C $(DIST_DIR) "$$filename"; \
			fi; \
		fi; \
	done
	
	@echo "Generating checksums..."
	@cd $(DIST_DIR)/release && sha256sum *.tar.gz > SHA256SUMS.txt
	
	@echo "Release artifacts ready in $(DIST_DIR)/release/:"
	@ls -lh $(DIST_DIR)/release/

# Show version
version:
	@echo "$(VERSION)"

# Run the CLI
run: build
	@./$(BIN_DIR)/$(APP_NAME)

# Development mode - watch and rebuild
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

# Format code
fmt:
	@go fmt ./...

# Check module dependencies
check-deps:
	@go mod verify
	@go mod tidy
	@echo "Dependencies verified"

# CI check (runs all quality checks)
ci: fmt lint test build
	@echo "All CI checks passed"

# Test release workflow locally (simulates GitHub Actions)
test-release:
	@./scripts/test-release-locally.sh

# Test release with specific version
test-release-version:
	@./scripts/test-release-locally.sh $(VERSION)

# Quick build test (builds current platform only)
test-build:
	@echo "Quick build test..."
	@go build -ldflags="-s -w -X main.version=test" -o /tmp/maestro-test .
	@/tmp/maestro-test --help > /dev/null && echo "✓ Build successful" || echo "✗ Build failed"
	@rm -f /tmp/maestro-test

# Run all tests (build, test, lint)
test-all: test-build test lint
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════"
	@echo "  All tests passed!                                          "
	@echo "═══════════════════════════════════════════════════════════════"
