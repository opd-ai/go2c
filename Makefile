.PHONY: all build test clean install help

# Build variables
BINARY_NAME=go2c
LLVM2C_BINARY=llvm2c
BUILD_DIR=./bin
INSTALL_PATH=/usr/local/bin

all: build

# Build both binaries
build:
	@echo "Building go2c..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/go2c
	@echo "Building llvm2c..."
	@go build -o $(BUILD_DIR)/$(LLVM2C_BINARY) ./cmd/llvm2c
	@echo "Build complete!"

# Run tests
test:
	@echo "Running unit tests..."
	@go test ./... -v
	@echo "Running integration tests..."
	@go test -v ./integration_test.go
	@echo "All tests passed!"

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f go2c llvm2c
	@echo "Clean complete!"

# Install binaries to system
install: build
	@echo "Installing binaries to $(INSTALL_PATH)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@sudo cp $(BUILD_DIR)/$(LLVM2C_BINARY) $(INSTALL_PATH)/
	@echo "Installation complete!"

# Uninstall binaries from system
uninstall:
	@echo "Uninstalling binaries from $(INSTALL_PATH)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@sudo rm -f $(INSTALL_PATH)/$(LLVM2C_BINARY)
	@echo "Uninstall complete!"

# Run linter
lint:
	@echo "Running linter..."
	@go vet ./...
	@echo "Lint complete!"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete!"

# Run example
example: build
	@echo "Running hello world example..."
	@$(BUILD_DIR)/$(LLVM2C_BINARY) -input testdata/input/hello.ll -output /tmp/hello.c -verbose
	@echo ""
	@echo "Generated C code:"
	@cat /tmp/hello.c

# Show help
help:
	@echo "Available targets:"
	@echo "  make build     - Build all binaries"
	@echo "  make test      - Run all tests"
	@echo "  make coverage  - Run tests with coverage report"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make install   - Install binaries to $(INSTALL_PATH)"
	@echo "  make uninstall - Uninstall binaries from $(INSTALL_PATH)"
	@echo "  make lint      - Run linter"
	@echo "  make fmt       - Format code"
	@echo "  make example   - Run hello world example"
	@echo "  make help      - Show this help message"
