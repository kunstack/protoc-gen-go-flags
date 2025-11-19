.PHONY: help generate build run clean deps

# Default target
help:
	@echo "Available targets:"
	@echo "  deps     - Install required dependencies"
	@echo "  generate - Generate protobuf code"
	@echo "  build    - Build the application"
	@echo "  run      - Build and run the application"
	@echo "  clean    - Clean generated files and build artifacts"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go install github.com/kunstack/protoc-gen-go-flags@latest
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go mod download

# Generate protobuf code
generate:
	@echo "Updating buf dependencies..."
	@buf dep update
	@echo "Generating protobuf code..."
	@buf generate

# Build application
build: generate
	@echo "Building application..."
	@mkdir -p bin
	@go build -o bin/server .

# Run application
run: build
	@echo "Running application with defaults..."
	@./bin/server

# Clean generated files and build artifacts
clean:
	@echo "Cleaning generated files..."
	@rm -f proto/*.pb.go proto/*.pb.flags.go
	@rm -rf bin/
	@rm -f buf.lock
	@echo "Clean complete!"