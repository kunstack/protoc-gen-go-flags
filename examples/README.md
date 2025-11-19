# Examples

This directory contains example projects demonstrating different use cases of protoc-gen-go-flags.

## Available Examples

### 1. Basic Example
**Directory**: `basic/`

A simple example showing basic flag generation for a server configuration with common field types (string, int, bool).

**Features**:
- Basic scalar types
- Short flags
- Default values
- Usage text

**Run**:
```bash
cd basic
make run
# Or with custom flags
./bin/server --host 0.0.0.0 --port 9000 --debug
```

### 2. Nested Configuration Example
**Directory**: `nested/`

Demonstrates hierarchical configuration with nested messages.

**Features**:
- Nested message configuration
- Hierarchical flag naming
- Multiple configuration sections (server, database, logging)

**Run**:
```bash
cd nested
make run
# Or with custom flags
./bin/app --server.host 0.0.0.0 --db.database-url postgres://localhost
```

### 3. Advanced Types Example
**Directory**: `advanced/`

Comprehensive example showing all supported protobuf types and features.

**Features**:
- All scalar types
- Duration and Timestamp types
- Repeated fields (arrays)
- Map fields with different formats
- Enum types
- Wrapper types

**Run**:
```bash
cd advanced
make run
```

### 4. Hierarchical Flags Example
**Directory**: `hierarchical/`

Shows how to use prefixes and custom delimiters for flag organization.

**Features**:
- Custom prefixes
- Different delimiter styles (dot, dash, underscore, colon)
- Multi-level prefixes

**Run**:
```bash
cd hierarchical
make run
```

## Prerequisites

All examples require:
- Go 1.18+
- [buf](https://buf.build/docs/installation) CLI tool
- protoc-gen-go-flags plugin installed

Install protoc-gen-go-flags:
```bash
go install github.com/kunstack/protoc-gen-go-flags@latest
```

## Building Examples

Each example has a Makefile with common targets:

```bash
# Generate protobuf code
make generate

# Build the example
make build

# Run the example
make run

# Clean generated files
make clean

# Show help
make help
```

## Project Structure

Each example follows this structure:

```
example-name/
├── go.mod              # Go module definition
├── main.go             # Application entry point
├── buf.yaml            # Buf module configuration
├── buf.gen.yaml        # Code generation config
├── Makefile            # Build automation
├── README.md           # Example-specific documentation
└── proto/
    └── config.proto    # Protobuf definitions
```

## Using Buf Schema Registry

All examples use the Buf Schema Registry approach, which is the recommended method:

1. **buf.yaml** declares the dependency on `buf.build/kunstack/protoc-gen-go-flags`
2. **buf.gen.yaml** configures code generation
3. Run `buf mod update` to fetch dependencies
4. Run `buf generate` to generate code

This approach eliminates the need to manually copy annotation files.

## Learning Path

We recommend exploring the examples in this order:

1. **basic** - Understand the fundamentals
2. **nested** - Learn about nested configurations
3. **advanced** - Explore all supported types
4. **hierarchical** - Master advanced flag organization