# protoc-gen-go-flags

[![Go Report Card](https://goreportcard.com/badge/github.com/kunstack/protoc-gen-go-flags)](https://goreportcard.com/report/github.com/kunstack/protoc-gen-go-flags)
[![Go Reference](https://pkg.go.dev/badge/github.com/kunstack/protoc-gen-go-flags.svg)](https://pkg.go.dev/github.com/kunstack/protoc-gen-go-flags)

[中文文档](README_zh.md) | English

protoc-gen-go-flags is a Go-based Protocol Buffer compiler plugin that automatically generates command-line flag bindings for protobuf messages. It generates `AddFlags` methods based on protobuf message definitions, seamlessly integrating with the `spf13/pflag` library to provide powerful command-line argument support for your protobuf messages.

## Why Use protoc-gen-go-flags

If your project meets any of the following criteria, protoc-gen-go-flags will greatly simplify your development workflow:

- ✅ Use Protocol Buffers to define configuration structures
- ✅ Need command-line argument support for CLI applications
- ✅ Want to avoid writing repetitive flag binding code manually
- ✅ Maintain consistency between configuration definitions and CLI interfaces
- ✅ Support complex nested configuration structures

**Traditional Approach vs protoc-gen-go-flags:**

The traditional approach requires manually writing flag bindings for each configuration field:
```go
// Manual approach: tedious and error-prone
fs.StringVar(&config.Host, "host", "localhost", "Server host")
fs.Int32Var(&config.Port, "port", 8080, "Server port")
fs.BoolVar(&config.Verbose, "verbose", false, "Enable verbose")
// ... repeat for every field
```

With protoc-gen-go-flags:
```go
// Auto-generated: concise and type-safe
config.AddFlags(fs)
```

## Table of Contents

- [Why Use protoc-gen-go-flags](#why-use-protoc-gen-go-flags)
- [Features](#features)
- [Quick Start](#quick-start)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Basic Usage](#basic-usage)
- [Complete Integration Tutorial](#complete-integration-tutorial)
  - [Step 1: Prepare Your Project](#step-1-prepare-your-project)
  - [Step 2: Add Flag Annotation Dependencies](#step-2-add-flag-annotation-dependencies)
  - [Step 3: Define Protobuf Messages](#step-3-define-protobuf-messages)
  - [Step 4: Generate Code](#step-4-generate-code)
  - [Step 5: Use in Your Application](#step-5-use-in-your-application)
- [Usage Examples](#usage-examples)
- [Supported Types](#supported-types)
- [Configuration Options](#configuration-options)
- [Hierarchical Flag Organization](#hierarchical-flag-organization)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

## Features

- 🚀 **Automated Code Generation**: Automatically generate command-line flag bindings from protobuf messages
- 🎯 **Complete Type Coverage**: Support all protobuf types (scalar types, enums, repeated, map, messages, etc.)
- 🔧 **Highly Configurable**: Support custom flag names, shortcuts, usage text, default values, and more
- 📦 **Nested Message Support**: Generate hierarchical flags for nested messages
- 🏗️ **Hierarchical Organization**: Support hierarchical flag naming through prefixes (dot, dash, underscore, colon separators)
- 🔒 **Best Practices**: Generate Go-idiomatic code with support for private/public methods
- 💾 **Default Value Support**: Provide default value settings for all types
- 🚦 **Deprecated Flags**: Support deprecated and hidden flags
- 🔄 **Package Aliasing**: Intelligently handle package name conflicts to avoid compilation errors

## Quick Start

### Prerequisites

Before getting started, ensure your development environment meets the following requirements:

- **Go 1.18+**: protoc-gen-go-flags requires Go 1.18 or higher
- **Protocol Buffers Compiler (protoc)**: Used to compile .proto files
  ```bash
  # macOS
  brew install protobuf

  # Ubuntu/Debian
  apt-get install protobuf-compiler

  # Or download from official releases: https://github.com/protocolbuffers/protobuf/releases
  ```
- **protoc-gen-go**: Go protobuf code generator
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  ```

### Installation

Install the protoc-gen-go-flags plugin:

```bash
go install github.com/kunstack/protoc-gen-go-flags@latest
```

Verify installation:
```bash
protoc-gen-go-flags --version
```

### Basic Usage

**1. Define a protobuf message with flag options:**

```protobuf
syntax = "proto3";

package example;

import "flags/annotations.proto";

option go_package = "github.com/example/project;example";

message Config {
    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server hostname"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port"
        default: 8080
    }];

    bool verbose = 3 [(flags.value).bool = {
        name: "verbose"
        short: "v"
        usage: "Enable verbose logging"
    }];
}
```

**2. Generate code:**

```bash
protoc -I. -I flags --go_out=paths=source_relative:. --flags_out=paths=source_relative:. config.proto
```

**3. Use in your application:**

```go
package main

import (
    "fmt"
    "os"

    pb "github.com/example/project"
    "github.com/spf13/pflag"
)

func main() {
    var config pb.Config

    // Create flag set and add flags
    fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
    config.AddFlags(fs)

    // Parse flags
    fs.Parse(os.Args[1:])

    // Use configuration (directly access fields)
    fmt.Printf("Server: %s:%d (verbose: %v)\n",
        config.Host, config.Port, config.Verbose)
}
```

### AddFlags vs SetDefaults

- **AddFlags method**: Registers configuration fields as command-line flags, allowing users to pass values via CLI arguments
- **SetDefaults method**: Sets default values for fields, used when no user-provided arguments are present

**Usage scenarios**:
- If you want default values to take effect before flag parsing, call `SetDefaults()` first
- If you only want to read configuration from command-line, you can use only `AddFlags()`
- Best practice is to combine both: provide defaults and allow user overrides

**Example calls**:

```go
var config pb.Config

// Method 1: Only use AddFlags (users must provide all values)
config.AddFlags(fs)

// Method 2: Combined usage (recommended)
config.SetDefaults()  // Set defaults first
config.AddFlags(fs)   // Then add flags for overrides

// Method 3: Use with custom flag set
customFS := pflag.NewFlagSet("custom", pflag.ExitOnError)
config.AddFlags(customFS)
```

## Complete Integration Tutorial

This section provides a complete step-by-step tutorial to help you integrate protoc-gen-go-flags into your own projects.

### Step 1: Prepare Your Project

Create a new Go project (or use an existing one):

```bash
mkdir myapp
cd myapp
go mod init github.com/yourname/myapp
```

Install necessary dependencies:

```bash
# Install pflag library
go get github.com/spf13/pflag

# Install protobuf runtime
go get google.golang.org/protobuf

# Install protoc-gen-go-flags runtime library
go get github.com/kunstack/protoc-gen-go-flags/flags
```

Create project structure:

```bash
myapp/
├── go.mod
├── main.go          # Application entry point
└── proto/
    └── config.proto # Protobuf definitions
```

### Step 2: Add Flag Annotation Dependencies

You need to add protoc-gen-go-flags annotation files to your project. There are two approaches:

#### Option 1: Use Buf Schema Registry (Recommended)

Add the dependency in your `buf.yaml`:

```yaml
version: v2
deps:
  - buf.build/kunstack/protoc-gen-go-flags
lint:
  use:
    - STANDARD
breaking:
  use:
    - FILE
```

Then configure code generation in `buf.gen.yaml`:

```yaml
version: v2
plugins:
  - remote: buf.build/protocolbuffers/go
    out: .
    opt: paths=source_relative
  - local: protoc-gen-go-flags
    out: .
    opt: paths=source_relative
```

Run buf commands to update dependencies and generate code:

```bash
buf mod update
buf generate
```

#### Option 2: Copy Files Directly

Download the `annotations.proto` file from the [protoc-gen-go-flags repository](https://github.com/kunstack/protoc-gen-go-flags/tree/main/flags) to your project:

```bash
mkdir -p proto/flags
curl -o proto/flags/annotations.proto \
  https://raw.githubusercontent.com/kunstack/protoc-gen-go-flags/main/flags/annotations.proto
```

Updated project structure:

```bash
myapp/
├── go.mod
├── main.go
└── proto/
    ├── config.proto
    └── flags/
        └── annotations.proto
```

### Step 3: Define Protobuf Messages

Define your configuration in `proto/config.proto`:

```protobuf
syntax = "proto3";

package myapp.config;

// Import flag annotations
import "flags/annotations.proto";

option go_package = "github.com/yourname/myapp/proto;config";

message ServerConfig {
    // Enable empty message generation
    option (flags.allow_empty) = true;

    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server host address"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port"
        default: 8080
    }];

    bool debug = 3 [(flags.value).bool = {
        name: "debug"
        short: "d"
        usage: "Enable debug mode"
    }];
}
```

### Step 4: Generate Code

Based on the option you chose in Step 2, use the appropriate command to generate code:

#### Using buf (if you chose Option 1)

If you chose the Buf Schema Registry approach in Step 2, the code was already generated when you ran `buf generate`.

Generated files:
- `proto/config.pb.go` - Standard protobuf Go code
- `proto/config.pb.flags.go` - Flag binding code

#### Using protoc (if you chose Option 2)

If you chose the direct file copy approach in Step 2, use the protoc command to generate:

```bash
protoc \
  -I./proto \
  -I./proto/flags \
  --go_out=. \
  --go_opt=paths=source_relative \
  --flags_out=. \
  --flags_opt=paths=source_relative \
  proto/config.proto
```

This generates two files:
- `proto/config.pb.go` - Standard protobuf Go code
- `proto/config.pb.flags.go` - Flag binding code

### Step 5: Use in Your Application

Use the generated code in `main.go`:

```go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/pflag"
    "github.com/yourname/myapp/proto"
)

func main() {
    // Create configuration instance
    config := &proto.ServerConfig{}

    // Set defaults (optional but recommended)
    config.SetDefaults()

    // Create flag set
    fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)

    // Add flags
    config.AddFlags(fs)

    // Parse command-line arguments
    if err := fs.Parse(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
        os.Exit(1)
    }

    // Use configuration
    fmt.Printf("Starting server...\n")
    fmt.Printf("  Host: %s\n", config.Host)
    fmt.Printf("  Port: %d\n", config.Port)
    fmt.Printf("  Debug: %v\n", config.Debug)

    // Start your application here...
}
```

### Step 6: Build and Run

Build the application:

```bash
go build -o myapp
```

Run and test command-line arguments:

```bash
# Use defaults
./myapp

# Output:
# Starting server...
#   Host: localhost
#   Port: 8080
#   Debug: false

# Custom parameters
./myapp --host 0.0.0.0 --port 3000 --debug

# Output:
# Starting server...
#   Host: 0.0.0.0
#   Port: 3000
#   Debug: true

# Use short options
./myapp -H 127.0.0.1 -p 9000 -d

# View help
./myapp --help
```

### Complete Project Example

Depending on which option you chose, the project structure will differ slightly:

#### Project Structure Using Buf Schema Registry

```bash
myapp/
├── go.mod
├── go.sum
├── main.go
├── buf.yaml              # Buf configuration
├── buf.gen.yaml          # Code generation config
├── buf.lock              # Dependency lock file (generated)
├── Makefile              # Optional: automation
└── proto/
    ├── config.proto
    ├── config.pb.go          # Generated
    └── config.pb.flags.go    # Generated
```

**Makefile example** (using buf):

```makefile
.PHONY: generate build run clean

# Generate protobuf code
generate:
	buf mod update
	buf generate

# Build application
build: generate
	go build -o bin/myapp .

# Run application
run: build
	./bin/myapp

# Clean generated files
clean:
	rm -f proto/*.pb.go proto/*.pb.flags.go
	rm -rf bin/
```

#### Project Structure Using Direct File Copy

```bash
myapp/
├── go.mod
├── go.sum
├── main.go
├── Makefile              # Optional: automation
└── proto/
    ├── config.proto
    ├── config.pb.go          # Generated
    ├── config.pb.flags.go    # Generated
    └── flags/
        └── annotations.proto
```

**Makefile example** (using protoc):

```makefile
.PHONY: generate build run clean

# Generate protobuf code
generate:
	protoc \
	  -I./proto \
	  -I./proto/flags \
	  --go_out=. \
	  --go_opt=paths=source_relative \
	  --flags_out=. \
	  --flags_opt=paths=source_relative \
	  proto/*.proto

# Build application
build: generate
	go build -o bin/myapp .

# Run application
run: build
	./bin/myapp

# Clean generated files
clean:
	rm -f proto/*.pb.go proto/*.pb.flags.go
	rm -rf bin/
```

Using the Makefile:

```bash
# Generate code
make generate

# Build
make build

# Run
make run

# Clean
make clean
```

### Advanced Integration: Nested Configuration

For complex applications, you may need nested configuration:

```protobuf
syntax = "proto3";

package myapp.config;

import "flags/annotations.proto";

option go_package = "github.com/yourname/myapp/proto;config";

message DatabaseConfig {
    string host = 1 [(flags.value).string = {
        name: "db-host"
        usage: "Database host"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "db-port"
        usage: "Database port"
        default: 5432
    }];
}

message AppConfig {
    option (flags.allow_empty) = true;

    string app_name = 1 [(flags.value).string = {
        name: "app-name"
        usage: "Application name"
        default: "MyApp"
    }];

    // Nested configuration
    DatabaseConfig database = 2 [(flags.value).message = {
        name: "db"
        nested: true
    }];
}
```

Using nested configuration:

```go
config := &proto.AppConfig{}
config.SetDefaults()

fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
config.AddFlags(fs)
fs.Parse(os.Args[1:])

fmt.Printf("App: %s\n", config.AppName)
fmt.Printf("DB: %s:%d\n", config.Database.Host, config.Database.Port)
```

Command-line usage:

```bash
./myapp --app-name "MyService" --db-db-host db.example.com --db-db-port 3306
```

### Troubleshooting

#### Issue 1: Cannot find annotations.proto

**Error message**:
```
proto/config.proto:3:1: Import "flags/annotations.proto" was not found.
```

**Solution**:

- **Using buf approach**: Ensure you ran `buf mod update` and that `buf.yaml` has the correct dependency:
  ```yaml
  deps:
    - buf.build/kunstack/protoc-gen-go-flags
  ```

- **Using protoc approach**: Ensure the protoc command includes the correct import path:
  ```bash
  protoc -I./proto -I./proto/flags ...
  ```

#### Issue 2: Generated code compilation errors

**Error message**:
```
undefined: flags.Option
undefined: types.Duration
```

**Solution**: Ensure you have installed the runtime libraries:
```bash
go get github.com/kunstack/protoc-gen-go-flags/flags
go get github.com/kunstack/protoc-gen-go-flags/types
go get github.com/kunstack/protoc-gen-go-flags/utils
```

The generated code will automatically import these packages; no manual import needed.

#### Issue 3: Flags not taking effect

**Symptom**: Command-line arguments are not being read; configuration uses zero values.

**Cause**: May have forgotten to call `SetDefaults()` or `AddFlags()`.

**Solution**: Call in the correct order:
```go
config := &proto.ServerConfig{}
config.SetDefaults()  // 1. Set defaults
config.AddFlags(fs)   // 2. Register flags
fs.Parse(os.Args[1:]) // 3. Parse arguments
```

#### Issue 4: buf generate failure

**Error message**:
```
Failure: plugin flags: not found
```

**Solution**: Ensure protoc-gen-go-flags is installed and in PATH:
```bash
# Install plugin
go install github.com/kunstack/protoc-gen-go-flags@latest

# Verify installation
which protoc-gen-go-flags
protoc-gen-go-flags --version
```

#### Issue 5: Package name conflicts

**Symptom**: Package name conflicts in generated code, such as using both `wrapperspb` and a custom `wrapperspb` package.

**Solution**: protoc-gen-go-flags automatically handles package name conflicts by generating aliases for conflicting packages. The generated code will automatically use aliased imports; no manual handling required.

#### Issue 6: Map format parsing errors

**Error message**:
```
invalid map format: ...
```

**Solution**: Ensure you use the correct format:
- **JSON format**: `--config='{"key": "value"}'`
- **STRING_TO_STRING**: `--labels="key1=value1,key2=value2"`
- **STRING_TO_INT**: `--limits="cpu=1000,memory=2048"`

Note: STRING_TO_INT only supports integer type values.

## Usage Examples

### Basic Configuration Example

```protobuf
syntax = "proto3";

package example;

import "flags/annotations.proto";

option go_package = "github.com/example/project;example";

message ServerConfig {
    option (flags.allow_empty) = true;

    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server host address"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port number"
        default: 8080
    }];

    bool https = 3 [(flags.value).bool = {
        name: "https"
        short: "s"
        usage: "Enable HTTPS"
    }];
}
```

### Hierarchical Flags (Using Prefixes)

```go
// Generate flags with prefix
fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
config.AddFlags(fs, flags.WithPrefix("server"))
fs.Parse(os.Args[1:])

// Result:
// --server.host
// --server.port
// --server.https
```

### Custom Delimiters

```go
// Use dash delimiter
fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
config.AddFlags(fs,
    flags.WithPrefix("server"),
    flags.WithDelimiter("-"))

// Result:
// --server-host
// --server-port
// --server-https
```

### Nested Messages

```protobuf
message DatabaseConfig {
    string url = 1 [(flags.value).string = {
        name: "database-url"
        usage: "Database connection URL"
    }];
}

message AppConfig {
    DatabaseConfig database = 1 [(flags.value).message = {
        name: "db"
        nested: true
    }];
}
```

Generated flags:
- `--db.database-url`

### Complete Configuration Example

```protobuf
syntax = "proto3";

package example;

import "flags/annotations.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/duration.proto";

option go_package = "github.com/example/project;example";

message Config {
    option (flags.allow_empty) = true;

    // Basic types
    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server hostname"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port"
        default: 8080
    }];

    // Special types
    google.protobuf.Duration timeout = 3 [(flags.value).duration = {
        name: "timeout"
        short: "t"
        usage: "Connection timeout"
        default: "30s"
    }];

    google.protobuf.Timestamp created = 4 [(flags.value).timestamp = {
        name: "created"
        usage: "Creation time"
        formats: ["RFC3339", "ISO8601"]
        default: "2024-01-01T00:00:00Z"
    }];

    // Repeated fields
    repeated string servers = 5 [(flags.value).repeated.string = {
        name: "servers"
        short: "s"
        usage: "Server addresses"
        default: ["localhost:8080"]
    }];

    // Map fields
    map<string, int32> limits = 6 [(flags.value).map = {
        name: "limits"
        usage: "Resource limits"
        format: MAP_FORMAT_TYPE_STRING_TO_INT
        default: "{\"cpu\": 1000, \"memory\": 2048}"
    }];

    // Nested messages
    DatabaseConfig database = 7 [(flags.value).message = {
        name: "database"
        nested: true
    }];
}
```

## Supported Types

protoc-gen-go-flags supports all Protocol Buffer types:

### Scalar Types

| Type | Go Type | Default Support | Repeated Support | Example |
|------|---------|-----------------|------------------|---------|
| `float` | `float32` | ✅ | ✅ | `3.14159` |
| `double` | `float64` | ✅ | ✅ | `2.71828` |
| `int32` | `int32` | ✅ | ✅ | `42` |
| `int64` | `int64` | ✅ | ✅ | `9223372036854775807` |
| `uint32` | `uint32` | ✅ | ✅ | `1000` |
| `uint64` | `uint64` | ✅ | ✅ | `18446744073709551615` |
| `sint32` | `int32` | ✅ | ✅ | `-42` |
| `sint64` | `int64` | ✅ | ✅ | `-9223372036854775808` |
| `fixed32` | `uint32` | ✅ | ✅ | `8080` |
| `fixed64` | `uint64` | ✅ | ✅ | `3000000000` |
| `sfixed32` | `int32` | ✅ | ✅ | `-1000` |
| `sfixed64` | `int64` | ✅ | ✅ | `-3000000000` |
| `bool` | `bool` | ✅ | ✅ | `true`, `false` |
| `string` | `string` | ✅ | ✅ | `"hello world"` |
| `bytes` | `[]byte` | ✅ | ✅ | `"aGVsbG8="` (base64) |

### Special Types

| Type | Go Type | Features | Example |
|------|---------|----------|---------|
| `enum` | Enum type | Default support, repeated fields | `1` (enum value) |
| `google.protobuf.Duration` | `*durationpb.Duration` | Default support, repeated fields | `"30s"`, `"1h"` |
| `google.protobuf.Timestamp` | `*timestamppb.Timestamp` | Multiple formats, default support, repeated fields | `"2024-01-01T00:00:00Z"` |
| `google.protobuf.StringValue` | `*wrapperspb.StringValue` | Default support, repeated fields | `"wrapper"` |
| `google.protobuf.Int32Value` | `*wrapperspb.Int32Value` | Default support, repeated fields | `42` |
| `google.protobuf.BoolValue` | `*wrapperspb.BoolValue` | Default support, repeated fields | `true` |

### Composite Types

| Type | Format Support | Default Support | Example |
|------|----------------|-----------------|---------|
| `repeated` (all scalar types) | - | ✅ | Slice types |
| `map<string, string>` | JSON, native | ✅ | `{"key": "value"}` |
| `map<string, int32>` | JSON, native | ✅ | `{"key": 123}` |
| `map<string, int64>` | JSON, native | ✅ | `{"key": 456}` |

### Nested Messages

Support generating hierarchical flags for nested messages, configured via the `message` flag type.

## Configuration Options

### Message-Level Options

Message-level options control flag generation behavior for the entire message:

```protobuf
message MyMessage {
  // Disable flag generation
  option (flags.disabled) = true;

  // Generate unexported flag methods (for custom wrapping)
  option (flags.unexported) = true;

  // Allow generating flag methods even without field configuration
  option (flags.allow_empty) = true;

  // Field definitions...
}
```

| Option | Type | Description |
|--------|------|-------------|
| `flags.disabled` | `bool` | Skip flag generation for this message |
| `flags.unexported` | `bool` | Generate unexported flag methods |
| `flags.allow_empty` | `bool` | Generate methods even without field configuration |

### Field-Level Options

Field-level options provide detailed configuration for individual fields:

```protobuf
string name = 1 [(flags.value).string = {
  name: "custom-name"           // Custom flag name
  short: "n"                    // Short flag (single character)
  usage: "Usage text"           // Usage description
  hidden: false                 // Hide flag (not shown in help)
  deprecated: true              // Mark as deprecated
  deprecated_usage: "Use --new-flag instead" // Deprecation message
  default: "default-value"      // Default value
}];
```

#### Common Field Options

All field types support the following options:

| Option | Type | Description |
|--------|------|-------------|
| `name` | `string` | Custom flag name (defaults to field name) |
| `short` | `string` | Short flag alias (single character) |
| `usage` | `string` | Help text (required) |
| `hidden` | `bool` | Hide flag |
| `deprecated` | `bool` | Deprecate flag |
| `deprecated_usage` | `string` | Deprecation message (required for deprecated flags) |

#### Bytes Type

Bytes type supports encoding format selection:

```protobuf
bytes data = 1 [(flags.value).bytes = {
  name: "data"
  usage: "Binary data"
  encoding: BYTES_ENCODING_TYPE_BASE64  // Or BYTES_ENCODING_TYPE_HEX
  default: "aGVsbG8="
}];
```

Supported encodings:
- `BYTES_ENCODING_TYPE_BASE64` - Standard base64 encoding (default)
- `BYTES_ENCODING_TYPE_HEX` - Hexadecimal encoding

#### Timestamp Type

Timestamp type supports multiple time formats:

```protobuf
google.protobuf.Timestamp created_at = 1 [(flags.value).timestamp = {
  name: "created-at"
  usage: "Creation timestamp"
  formats: ["RFC3339", "ISO8601"]  // Supported formats
  default: "2024-01-01T00:00:00Z"
}];
```

Supported time formats:

| Format Name | Go Time Format Constant | Example | Description |
|-------------|------------------------|---------|-------------|
| `RFC3339` | `time.RFC3339` | `2024-01-01T15:04:05Z` or `2024-01-01T15:04:05+08:00` | RFC 3339 standard format (recommended) |
| `RFC3339Nano` | `time.RFC3339Nano` | `2024-01-01T15:04:05.999999999Z` | RFC 3339 with nanosecond precision |
| `RFC822` | `time.RFC822` | `01 Jan 24 15:04 MST` | RFC 822 format |
| `RFC822Z` | `time.RFC822Z` | `01 Jan 24 15:04 -0700` | RFC 822 with numeric timezone |
| `RFC850` | `time.RFC850` | `Monday, 01-Jan-24 15:04:05 MST` | RFC 850 format |
| `RFC1123` | `time.RFC1123` | `Mon, 01 Jan 2024 15:04:05 MST` | RFC 1123 format |
| `RFC1123Z` | `time.RFC1123Z` | `Mon, 01 Jan 2024 15:04:05 -0700` | RFC 1123 with numeric timezone |
| `ISO8601` | Custom | `2024-01-01` | ISO 8601 date format |
| `ISO8601Time` | Custom | `2024-01-01T15:04:05` | ISO 8601 datetime format (no timezone) |
| `Kitchen` | `time.Kitchen` | `3:04PM` | Kitchen clock format |
| `Stamp` | `time.Stamp` | `Jan  1 15:04:05` | Timestamp format |
| `StampMilli` | `time.StampMilli` | `Jan  1 15:04:05.000` | Millisecond precision timestamp |
| `StampMicro` | `time.StampMicro` | `Jan  1 15:04:05.000000` | Microsecond precision timestamp |
| `StampNano` | `time.StampNano` | `Jan  1 15:04:05.000000000` | Nanosecond precision timestamp |
| `DateTime` | Custom | `2024-01-01 15:04:05` | Date-time format |
| `DateOnly` | Custom | `2024-01-01` | Date only format |
| `TimeOnly` | Custom | `15:04:05` | Time only format |

> **Note**:
> - Besides predefined formats, you can use any valid Go time format string
> - Format names support `RFC339` (typo) as an alias for `RFC3339` for backward compatibility
> - In the `formats` array, you can specify multiple formats; parsing will try them in order

#### Duration Type

```protobuf
google.protobuf.Duration timeout = 1 [(flags.value).duration = {
  name: "timeout"
  usage: "Timeout duration"
  default: "30s"
}];
```

Supported format: seconds + unit (e.g., "30s", "5m", "1h")

#### Map Type

Map types support multiple formats, with corresponding default values based on format type:

**1. JSON Format (Default)**

```protobuf
map<string, int32> config = 1 [(flags.value).map = {
  name: "config"
  usage: "Configuration key-value pairs"
  format: MAP_FORMAT_TYPE_JSON
  default: "{\"cpu\": 1000, \"memory\": 2048}"
}];
```

Command-line usage example:
```bash
# JSON format input
./myapp --config='{"cpu": 1000, "memory": 2048}'
```

**2. STRING_TO_STRING Format**

```protobuf
map<string, string> labels = 1 [(flags.value).map = {
  name: "labels"
  usage: "Key-value labels"
  format: MAP_FORMAT_TYPE_STRING_TO_STRING
  default: "env=production,region=us-west"
}];
```

Command-line usage example:
```bash
# Use comma-separated key-value pairs
./myapp --labels="env=production,region=us-west"

# Or specify multiple times (will overwrite)
./myapp --labels="env=staging" --labels="region=eu-central"
```

**3. STRING_TO_INT Format**

```protobuf
map<string, int32> limits = 1 [(flags.value).map = {
  name: "limits"
  usage: "Resource limits"
  format: MAP_FORMAT_TYPE_STRING_TO_INT
  default: "cpu=1000,memory=2048,disk=10000"
}];
```

Command-line usage example:
```bash
# Use comma-separated integer key-value pairs
./myapp --limits="cpu=1000,memory=2048,disk=10000"

# Single key-value pair
./myapp --limits="cpu=2000"

# Multiple specifications will merge
./myapp --limits="cpu=1000,memory=2048" --limits="disk=10000"
```

Supported formats:
- `MAP_FORMAT_TYPE_JSON` - JSON format (default)
  - Default value example: `"{\"key\": \"value\"}"`
- `MAP_FORMAT_TYPE_STRING_TO_STRING` - String key-value pair format
  - Default value example: `"key1=value1,key2=value2"`
  - Use commas to separate multiple key-value pairs, each connected with equals
- `MAP_FORMAT_TYPE_STRING_TO_INT` - String key to integer value format
  - Default value example: `"key1=123,key2=456"`
  - Use commas to separate multiple key-value pairs, values must be integers
  - **Supported integer types**: `int32`, `sint32`, `sfixed32`, `int64`, `sint64`, `sfixed64`, `uint32`, `fixed32`, `uint64`, `fixed64`

#### Repeated Fields

```protobuf
syntax = "proto3";
package tests;

import "flags/annotations.proto";

message Example {
  repeated string servers = 1 [(flags.value).repeated.string = {
    name: "servers"
    usage: "Server addresses (can be specified multiple times)"
    default: ["server1"]
  }];
}
```

### Nested Message Configuration

Nested messages use the `message` flag type:

```protobuf
message NestedConfig {
  string value = 1 [(flags.value).string = { name: "value" }];
}

message MainConfig {
  NestedConfig nested = 1 [(flags.value).message = {
    name: "nested"     // Prefix name for nested message
    nested: true       // Enable nested flag generation
  }];
}
```

| Option | Type | Description |
|--------|------|-------------|
| `name` | `string` | Prefix name for nested message (defaults to field name) |
| `nested` | `bool` | Whether to generate nested flags |

## Hierarchical Flag Organization

protoc-gen-go-flags supports hierarchical flag organization through `WithPrefix` and `WithDelimiter` options.

### Basic Prefix

```go
config.AddFlags(fs, flags.WithPrefix("server"))
```

Generates: `--server.host`, `--server.port`

### Multi-level Prefix

```go
config.AddFlags(fs, flags.WithPrefix("server", "database"))
```

Generates: `--server.database.host`, `--server.database.port`

### Custom Delimiters

```go
config.AddFlags(fs,
  flags.WithPrefix("server"),
  flags.WithDelimiter("-"))  // Dash
```

Generates: `--server-host`, `--server-port`

Supported delimiters:
- `flags.DelimiterDot` - Dot (default): `server.port`
- `flags.DelimiterDash` - Dash: `server-port`
- `flags.DelimiterUnderscore` - Underscore: `server_port`
- `flags.DelimiterColon` - Colon: `server:port`

### Custom Renaming Function

```go
config.AddFlags(fs,
  flags.WithPrefix("Server"),
  flags.WithRenamer(strings.ToLower))
```

Generates: `--server-host` (converted to lowercase)

## FAQ

### Q: How do I integrate protoc-gen-go-flags into an existing project?

**A:** Follow these steps:
1. Install the plugin: `go install github.com/kunstack/protoc-gen-go-flags@latest`
2. Copy `annotations.proto` to your project
3. Add flag annotations to your `.proto` files
4. Run `protoc` to generate code
5. Use the generated `AddFlags()` method in your application

For detailed steps, see the [Complete Integration Tutorial](#complete-integration-tutorial).

### Q: How do I handle complex nested configurations?

**A:** Use nested messages and the `message` flag type:

```protobuf
syntax = "proto3";
package tests;

import "flags/annotations.proto";

message DatabaseConfig {
    string url = 1 [(flags.value).string = {
        name: "url"
        usage: "Database connection URL"
    }];
}

message AppConfig {
    DatabaseConfig database = 1 [(flags.value).message = {
        name: "db"
        nested: true
    }];
}
```

This generates hierarchical flags like `--db-url`.

### Q: How do I customize flag naming (using prefixes or delimiters)?

**A:** Use options when calling `AddFlags`:

```go
// With prefix
config.AddFlags(fs, flags.WithPrefix("server"))
// Generates: --server.host

// Custom delimiter
config.AddFlags(fs,
    flags.WithPrefix("server"),
    flags.WithDelimiter("-"))
// Generates: --server-host
```

### Q: Generated code errors "undefined: flags.Option"

**A:** You need to install and import the runtime library:

```bash
go get github.com/kunstack/protoc-gen-go-flags/flags
```

```go
import "github.com/kunstack/protoc-gen-go-flags/flags"
```

### Q: How do I skip flag generation for specific fields?

**A:** Simply don't add flag annotations to that field. If you already added annotations, you can omit the field-level option:

```protobuf
string internal_field = 1;  // No flag annotation; no flag will be generated
```

### Q: How do I set default values for fields?

**A:** Use the `default` option in the flag annotation:

```protobuf
int32 port = 1 [(flags.value).int32 = {
    name: "port"
    usage: "Server port"
    default: 8080  // Set default value
}];
```

Then call `config.SetDefaults()` in your application to apply defaults.

### Q: Which protobuf types are supported?

**A:** protoc-gen-go-flags supports all standard protobuf types:
- Scalar types: string, int32, int64, bool, float, double, etc.
- Special types: google.protobuf.Duration, Timestamp
- Composite types: repeated (arrays), map (maps)
- Nested messages

For a detailed list, see the [Supported Types](#supported-types) section.

### Q: How do I use environment variables with flags?

**A:** protoc-gen-go-flags focuses on command-line flag binding. For environment variable support, combine with configuration management libraries like [viper](https://github.com/spf13/viper):

```go
import (
    "github.com/spf13/pflag"
    "github.com/spf13/viper"
)

config := &proto.Config{}
fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
config.AddFlags(fs)

// Bind to viper (supports environment variables)
viper.BindPFlags(fs)
viper.AutomaticEnv()

fs.Parse(os.Args[1:])
```

### Q: What is the naming convention for generated files?

**A:** For `.proto` files, corresponding `.pb.flags.go` files are generated:
- `config.proto` → `config.pb.go` + `config.pb.flags.go`
- `server.proto` → `server.pb.go` + `server.pb.flags.go`

### Q: Does it support gRPC?

**A:** protoc-gen-go-flags is fully compatible with gRPC. You can define both gRPC services and flag configurations in the same `.proto` file:

```bash
protoc \
    --go_out=. \
    --go-grpc_out=. \
    --flags_out=. \
    your_service.proto
```

## Contributing

Contributions are welcome! If you have suggestions or find issues, please:

- Submit an issue: [GitHub Issues](https://github.com/kunstack/protoc-gen-go-flags/issues)
- Submit a pull request: Fork the project and create a PR
- Improve documentation: Help enhance documentation and examples

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [protoc-gen-star](https://github.com/lyft/protoc-gen-star) - Code generation framework
- [spf13/pflag](https://github.com/spf13/pflag) - Command-line flag library
- [Google Protocol Buffers](https://protobuf.dev/) - Data serialization protocol
