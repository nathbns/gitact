# GitHub Activity CLI - Makefile

# Variables
BINARY_NAME=gitact
VERSION=1.0.0
BUILD_DIR=build
DIST_DIR=dist

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -s -w"
BUILD_FLAGS=-trimpath

# Platform targets
PLATFORMS=windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: all build clean test deps run help install uninstall release

# Default target
all: clean deps test build

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "✅ Clean completed"

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...
	@echo "✅ Tests completed"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies updated"

# Run the application
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Run with example user
demo: build
	@echo "🎮 Running demo with karpathy profile..."
	./$(BUILD_DIR)/$(BINARY_NAME) karpathy

# Install to system PATH
install: build
	@echo "📥 Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Installation completed"

# Uninstall from system
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstallation completed"

# Cross-compile for all platforms
release: clean deps test
	@echo "🏗️  Building release binaries for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		platform_split=($$(echo $$platform | tr '/' ' ')); \
		GOOS=$${platform_split[0]}; \
		GOARCH=$${platform_split[1]}; \
		output_name=$(BINARY_NAME)-$(VERSION)-$$GOOS-$$GOARCH; \
		if [ $$GOOS = "windows" ]; then output_name=$$output_name.exe; fi; \
		echo "Building $$output_name..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$$output_name .; \
		if [ $$? -ne 0 ]; then \
			echo "❌ Failed to build $$output_name"; \
			exit 1; \
		fi; \
	done
	@echo "📦 Creating release archives..."
	@cd $(DIST_DIR) && for binary in *; do \
		if [[ $$binary == *.exe ]]; then \
			zip $${binary%.*}.zip $$binary; \
		else \
			tar -czf $${binary}.tar.gz $$binary; \
		fi; \
	done
	@echo "✅ Release build completed: $(DIST_DIR)/"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@$(GOCMD) fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Vet code
vet:
	@echo "🔍 Vetting code..."
	@$(GOCMD) vet ./...
	@echo "✅ Code vetted"

# Security scan
security:
	@echo "🔒 Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Check for updates
update:
	@echo "🔄 Checking for dependency updates..."
	@$(GOCMD) list -u -m all

# Development mode with live reload
dev:
	@echo "🔥 Starting development mode..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "🔄 Falling back to regular run..."; \
		make run ARGS="$(ARGS)"; \
	fi

# Benchmark
benchmark:
	@echo "⚡ Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./...

# Profile application
profile: build
	@echo "📊 Profiling application..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -cpuprofile=cpu.prof -memprofile=mem.prof $(ARGS)
	@echo "📈 Profile data saved: cpu.prof, mem.prof"

# Docker build
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "✅ Docker image built: $(BINARY_NAME):$(VERSION)"

# Docker run
docker-run: docker-build
	@echo "🐳 Running Docker container..."
	@docker run --rm -it $(BINARY_NAME):latest $(ARGS)

# Check code quality
quality: fmt vet lint test
	@echo "✅ Code quality checks completed"

# Full CI pipeline
ci: deps quality build release
	@echo "🎉 CI pipeline completed successfully"

# Show help
help:
	@echo "🐙 GitHub Activity CLI - Make Commands"
	@echo ""
	@echo "📋 Available targets:"
	@echo "  build      - Build the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  deps       - Download dependencies"
	@echo "  run        - Build and run (use ARGS='username' for arguments)"
	@echo "  demo       - Run demo with karpathy profile"
	@echo "  install    - Install to /usr/local/bin"
	@echo "  uninstall  - Remove from system"
	@echo "  release    - Cross-compile for all platforms"
	@echo ""
	@echo "🔧 Development:"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code (requires golangci-lint)"
	@echo "  vet        - Vet code"
	@echo "  security   - Security scan (requires gosec)"
	@echo "  dev        - Development mode with live reload (requires air)"
	@echo "  quality    - Run all code quality checks"
	@echo ""
	@echo "📊 Analysis:"
	@echo "  benchmark  - Run benchmarks"
	@echo "  profile    - Profile application"
	@echo "  update     - Check for dependency updates"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run in Docker container"
	@echo ""
	@echo "🎯 CI/CD:"
	@echo "  ci         - Full CI pipeline"
	@echo ""
	@echo "💡 Examples:"
	@echo "  make run ARGS='karpathy'        - Run with karpathy profile"
	@echo "  make run ARGS='--repos torvalds' - List torvalds repositories"
	@echo "  make dev ARGS='octocat'         - Development mode with octocat"
