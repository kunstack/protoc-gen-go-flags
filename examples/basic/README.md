# Basic Example

This example demonstrates basic usage of protoc-gen-go-flags with a simple server configuration.

## Overview

This example shows how to:
- Define basic protobuf messages with flag annotations
- Generate command-line flag bindings
- Use the generated code in a Go application
- Set default values and parse command-line arguments

## Project Structure

```
basic/
├── go.mod              # Go module definition
├── main.go             # Application entry point
├── buf.yaml            # Buf module configuration with flags dependency
├── buf.gen.yaml        # Code generation configuration
├── Makefile            # Build automation
├── README.md           # This file
└── proto/
    └── config.proto    # Server configuration definition
```

## Configuration

The example defines a `ServerConfig` message with three fields:

- `host` (string) - Server host address
- `port` (int32) - Server port number
- `debug` (bool) - Enable debug mode

Each field has flag annotations specifying:
- Flag name (e.g., `--host`)
- Short flag (e.g., `-H`)
- Usage text
- Default value

## Usage

### Build and Run

```bash
# Generate code and build
make build

# Run with defaults
./bin/server

# Run with custom flags
./bin/server --host 0.0.0.0 --port 9000 --debug

# Use short flags
./bin/server -H 127.0.0.1 -p 3000 -d

# Show help
./bin/server --help
```

### Expected Output

With defaults:
```
Starting server...
  Host: localhost
  Port: 8080
  Debug: false
```

With custom flags:
```
./bin/server --host 0.0.0.0 --port 9000 --debug

Starting server...
  Host: 0.0.0.0
  Port: 9000
  Debug: true
```

## Code Generation

The project uses Buf for dependency management and code generation:

```bash
# Update buf dependencies
buf mod update

# Generate code
buf generate
```

This generates:
- `proto/config.pb.go` - Standard protobuf Go code
- `proto/config.pb.flags.go` - Flag binding code with `AddFlags` and `SetDefaults` methods

## Key Concepts

### 1. Flag Annotations

In `proto/config.proto`, each field has a `(flags.value)` annotation:

```protobuf
string host = 1 [(flags.value).string = {
    name: "host"
    short: "H"
    usage: "Server host address"
    default: "localhost"
}];
```

### 2. Generated Methods

The plugin generates two methods:

- `AddFlags(fs *pflag.FlagSet, opts ...flags.Option)` - Registers flags
- `SetDefaults()` - Sets default values

### 3. Usage Pattern

In your application:

```go
config := &proto.ServerConfig{}
config.SetDefaults()  // Set defaults first

fs := pflag.NewFlagSet("server", pflag.ExitOnError)
config.AddFlags(fs)   // Register flags
fs.Parse(os.Args[1:]) // Parse arguments

// Use config.Host, config.Port, config.Debug...
```

## Next Steps

After understanding this basic example, check out:
- **nested/** - Learn about nested message configuration
- **advanced/** - Explore all supported types
- **hierarchical/** - Master flag organization with prefixes